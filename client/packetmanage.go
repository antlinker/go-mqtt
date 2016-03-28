package client

import (
	"time"

	"github.com/antlinker/go-mqtt/packet"
)

type PacketManager interface {
	sendPacketer
	sendUnfinaler
	sendQos2Unfinaler
	receiveQos2

	Start()
	Stop()
}

//publish subscribe unsbuscribe发送时管理
type sendPacketer interface {
	//加入发送的报文同时分配id并锁定id
	AddSendPacket(*MqttPacket) error
	//弹出一个发送报文
	PopSend() *MqttPacket
}

//publish（qos1 qos2） subscribe unsbuscribe发送后管理
type sendUnfinaler interface {
	//移动到未完成
	MoveSendUnfinal(*MqttPacket)
	//获取未完成信息
	GetSendUnfinal(id uint16) *MqttPacket
	//移除发送未完成信息同时释放id
	RemoveSendUnfinal(id uint16) *MqttPacket
}

//publis(qos2)收到pubrec后管理
type sendQos2Unfinaler interface {
	AddSendQos2Pakcet(*MqttPacket)
	RemoveSendQos2Pakcet(id uint16) *MqttPacket
}

type receiveQos2 interface {
	//qos2存储id
	AddReceivePacket(*MqttPacket)
	//检测qos2
	ExistsReveivePacket(id uint16) bool
	RemoveReceivePacket(id uint16) *MqttPacket
}

const (
	//发送
	Direct_Send int = iota
	//接收
	Direct_Recive
)

func CreateMqttPacket(direct int, msg packet.PacketIdMessage) *MqttPacket {
	return &MqttPacket{
		Direct:  direct,
		Packet:  msg,
		Lastime: time.Now(),
		wait:    make(chan struct{}),
	}
}

type MqttPacket struct {
	Direct  int
	Packet  packet.PacketIdMessage
	Lastime time.Time
	wait    chan struct{}
	err     error
	rectime time.Time
	sending bool
}

func (m *MqttPacket) IsSending() bool {
	return m.sending
}

func (m *MqttPacket) Rectime() time.Time {
	return m.rectime
}

func (m *MqttPacket) Err() error {
	return m.err
}
func (m *MqttPacket) finalish(err error) {
	m.err = err
	select {
	case <-m.wait:
	default:
		close(m.wait)
	}
}
func (m *MqttPacket) WaitTimeout(timeout time.Duration) bool {
	select {
	case <-m.wait:
		return true
	case <-time.After(timeout):
		return false
	}
	return false
}
func (m *MqttPacket) Wait() {
	<-m.wait
}
