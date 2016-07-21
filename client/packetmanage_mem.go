package client

import (
	"container/list"
	"sync"
	"time"

	"github.com/antlinker/go-mqtt/util"

	"github.com/antlinker/go-cmap"
	"github.com/antlinker/go-mqtt/packet"
)

// NewMemPacketManager 内存版报文管理
func NewMemPacketManager(client *antClient) PacketManager {
	m := new(MemPacketManager)
	m.client = client
	m.reSendInterval = 10 * time.Second
	m.init()

	return m
}

// BasePacketManager 报文管理基础类
type BasePacketManager struct {
	client         *antClient
	reSendInterval time.Duration
}

// Send 发送报文 主要用来重发pub rel
func (b *BasePacketManager) Send(msg packet.MessagePacket) error {
	return b.client._send(msg)
}

// MemPacketManager 内存报文管理
type MemPacketManager struct {
	memSendPacketer
	memSendUnfinaler
	memSendQos2Unfinaler
	memReceiveQos2
	BasePacketManager
	stop            chan struct{}
	packetIDFactory *util.PacketIdFactory
}

func (m *MemPacketManager) init() {
	m.packetIDFactory = util.NewPacketIdFactory()

	m.memSendUnfinaler.packetIDFactory = m.packetIDFactory
	m.memSendUnfinaler.storePacketMap = cmap.NewConcurrencyMap()
	m.memSendQos2Unfinaler.sendQsos2Map = cmap.NewConcurrencyMap()
	m.memReceiveQos2.recvQsos2Map = cmap.NewConcurrencyMap()

	m.memSendPacketer.start(m.packetIDFactory)
}

// PopSend 弹出第一个需要发送的报文
func (m *MemPacketManager) PopSend() *MqttPacket {
	mp := m.memSendPacketer.PopSend()
	if mp != nil {
		if mp.Packet.GetHeaderType() == packet.TYPE_PUBLISH {
			pub := mp.Packet.(*packet.Publish)
			if pub.GetQos() == packet.QOS_0 {
				return mp
			}
		}
		m.memSendUnfinaler.MoveSendUnfinal(mp)
		return mp
	}
	return nil
}

// Start 开始报文管理,创建重发检测go程
func (m *MemPacketManager) Start() {
	m.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(m.reSendInterval):
			case <-m.stop:
				return
			}
			m.unfinal2send()

		}
	}()
}

// Stop 停止报文管理
func (m *MemPacketManager) Stop() {

	m.memSendPacketer._stop()
	select {
	case <-m.stop:
	default:
		close(m.stop)
	}
}
func (m *MemPacketManager) unfinal2send() {
	defer func() {
		if err := recover(); err != nil {
			Mlog.Warn("重发出错", err)
		}
	}()
	//检测消息
	smap := m.storePacketMap.ToMap()
	for _, v := range smap {
		mp := v.(*MqttPacket)
		msg := mp.Packet
		if !mp.sending && mp.Lastime.Add(m.reSendInterval).Sub(time.Now()) <= 0 {
			if msg.GetHeaderType() == packet.TYPE_PUBLISH {

				nmp, _ := m.storePacketMap.Remove(msg.GetPacketId())
				if nmp != nil {
					msg.(*packet.Publish).SetDupFlag(true)
					m.AddSendPacket(mp)
				}
			} else {
				nmp := m.RemoveSendUnfinal(msg.GetPacketId())
				if nmp != nil {
					nmp.finalish(ErrTimeout)
				}

			}
		}

	}

	//qos2收到rec 未收到rel 重发rel
	sqos2map := m.sendQsos2Map.ToMap()
	for _, q := range sqos2map {
		mp := q.(*MqttPacket)
		msg := mp.Packet

		if mp.Rectime().Add(m.reSendInterval).Sub(time.Now()) <= 0 {
			pubrel := packet.NewPubrel()
			pubrel.SetPacketId(msg.GetPacketId())
			m.Send(pubrel)
		}
	}
}

type memSendPacketer struct {
	send            *list.List
	packetIDFactory *util.PacketIdFactory
	rwlock          sync.RWMutex

	sendcond *sync.Cond
	stop     bool
	stopwait sync.WaitGroup
}

func (m *memSendPacketer) start(packetIDFactory *util.PacketIdFactory) {
	m.send = list.New()
	m.sendcond = sync.NewCond(&m.rwlock)
	m.packetIDFactory = packetIDFactory
	m.stop = false
}
func (m *memSendPacketer) AddSendPacket(pkt *MqttPacket) error {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	msg := pkt.Packet
	pid := msg.GetPacketId()
	if pid == 0 {
		if msg.GetHeaderType() == packet.TYPE_PUBLISH {
			pub := msg.(*packet.Publish)
			if pub.GetQos() != packet.QOS_0 {
				pid := m.packetIDFactory.CreateId()
				msg.SetPacketId(pid)

			}

		} else {
			pid := m.packetIDFactory.CreateId()
			msg.SetPacketId(pid)
		}
	}
	m.send.PushFront(pkt)
	m.sendcond.Signal()
	//Mlog.Debug("AddSendPacket：", pkt)
	return nil
}

// //获取一个要发送的报文
// func (m *memSendPacketer) NextSend() *MqttPacket {
// 	m.rwlock.RLock()
// 	elem := m.send.Front()
// 	m.rwlock.RUnlock()
// 	if elem != nil {
// 		//	Mlog.Debug("NextSend", elem)

// 		return elem.Value.(*MqttPacket)
// 	}
// 	return nil

// }

//弹出一个发送报文
func (m *memSendPacketer) PopSend() *MqttPacket {
	//defer Mlog.Debug("PopSend end")
	m.sendcond.L.Lock()
	defer m.sendcond.L.Unlock()

	//Mlog.Debug("PopSend start")
	for {

		elem := m.send.Back()
		if elem != nil {
			mp := elem.Value.(*MqttPacket)
			m.send.Remove(elem)
			mp.sending = true
			return mp
		} else if m.stop {
			return nil
		} else {
			m.stopwait.Add(1)
			m.sendcond.Wait()
			m.stopwait.Done()
		}
	}

}
func (m *memSendPacketer) _stop() {
	m.sendcond.L.Lock()

	m.stop = true
	m.sendcond.Broadcast()
	m.sendcond.L.Unlock()
	m.stopwait.Wait()
}

type memSendUnfinaler struct {
	storePacketMap  cmap.ConcurrencyMap
	packetIDFactory *util.PacketIdFactory
}

//移动到未完成
func (m *memSendUnfinaler) MoveSendUnfinal(mp *MqttPacket) {
	m.storePacketMap.Set(mp.Packet.GetPacketId(), mp)
}

//获取未完成信息
func (m *memSendUnfinaler) GetSendUnfinal(id uint16) *MqttPacket {
	e, _ := m.storePacketMap.Get(id)
	if e != nil {
		return e.(*MqttPacket)

	}
	return nil
}

//移除未完成发送信息
func (m *memSendUnfinaler) RemoveSendUnfinal(id uint16) *MqttPacket {
	m.packetIDFactory.ReleaseId(id)
	e, _ := m.storePacketMap.Remove(id)

	if e != nil {
		return e.(*MqttPacket)

	}
	return nil
}

type memSendQos2Unfinaler struct {
	sendQsos2Map cmap.ConcurrencyMap
}

func (m *memSendQos2Unfinaler) AddSendQos2Pakcet(mp *MqttPacket) {
	m.sendQsos2Map.Set(mp.Packet.GetPacketId(), mp)
}
func (m *memSendQos2Unfinaler) RemoveSendQos2Pakcet(id uint16) *MqttPacket {
	e, _ := m.sendQsos2Map.Remove(id)
	if e != nil {
		return e.(*MqttPacket)

	}
	return nil
}

type memReceiveQos2 struct {
	recvQsos2Map cmap.ConcurrencyMap
}

//qos2存储id
func (m *memReceiveQos2) AddReceivePacket(mp *MqttPacket) {

	m.recvQsos2Map.SetIfAbsent(mp.Packet.GetPacketId(), mp)

}

func (m *memReceiveQos2) RemoveReceivePacket(id uint16) *MqttPacket {
	e, _ := m.recvQsos2Map.Remove(id)
	if e != nil {
		return e.(*MqttPacket)

	}
	return nil
}
