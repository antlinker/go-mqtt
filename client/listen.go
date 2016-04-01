package client

import (
	"sync/atomic"

	"github.com/antlinker/go-mqtt/event"
	"github.com/antlinker/go-mqtt/packet"
)

type pubCnt struct {
	send    int64
	success int64
	final   int64
}

// MqttConnListener 连接事件
type MqttConnListener interface {
	event.Listener
	//开始建立连接
	OnConnStart(event *MqttConnEvent)
	//连接成功
	OnConnSuccess(event *MqttConnEvent)
	//连接失败
	OnConnFailure(event *MqttConnEvent, returncode int, err error)
}

type mqttListen struct {
	event.Generator
	baseClientStatus
}

//注册连接监听
func (cl *mqttListen) AddConnListener(listener MqttConnListener) {
	cl.AddListener(EventConn, listener)
}

//移除连接监听
func (cl *mqttListen) RemoveConnListener(listener MqttConnListener) {
	cl.RemoveListener(EventConn, listener)
}

func (cl *mqttListen) fireOnConnStart(client MqttClienter) {
	//Mlog.Debug("fireOnConnStart")
	cl.FireListener(EventConn, "OnConnStart", createMqttConnEvent(client, StatusConnStart, atomic.AddInt64(&cl.totalRecnt, 1), atomic.AddInt64(&cl.curRecnt, 1)))
}
func (cl *mqttListen) fireOnConnSuccess(client MqttClienter) {
	cl.curRecnt = 0
	cl.FireListener(EventConn, "OnConnSuccess", createMqttConnEvent(client, StatusConnSuccess, cl.curRecnt, cl.totalRecnt))
}
func (cl *mqttListen) fireOnConnFailure(client MqttClienter, returncode int, err error) {
	cl.FireListener(EventConn, "OnConnFailure", createMqttConnEvent(client, StatusConnFailure, cl.curRecnt, cl.totalRecnt), returncode, err)
}

// MqttPubListener 发布消息事件监听
type MqttPubListener interface {
	event.Listener
	//准备发送消息
	OnPubReady(event *MqttPubEvent, packet *MqttPacket)
	//消息发布成功
	OnPubSuccess(event *MqttPubEvent, packet *MqttPacket)
	//消息发布完成
	OnPubFinal(event *MqttPubEvent, packet *MqttPacket)
}

//注册发布消息监听
func (cl *mqttListen) AddPubListener(listener MqttPubListener) {
	cl.AddListener(EventPublish, listener)
}

//移除发布消息监听
func (cl *mqttListen) RemovePubListener(listener MqttPubListener) {
	cl.RemoveListener(EventPublish, listener)
}

func (cl *mqttListen) fireOnPubReady(client MqttClienter, mpacket *MqttPacket) {
	atomic.AddInt64(&(cl.pubcnt[PubCntTOTAL].send), 1)
	switch mpacket.Packet.(*packet.Publish).GetQos() {
	case 0:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS0].send, 1)
	case 1:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS1].send, 1)
	case 2:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS2].send, 1)
	}
	cl.FireListener(EventPublish, "OnPubReady", createMqttPubEvent(client, PubStatusReady, cl.pubcnt), mpacket)
}
func (cl *mqttListen) fireOnPubSuccess(client MqttClienter, mpacket *MqttPacket) {
	atomic.AddInt64(&cl.pubcnt[PubCntTOTAL].success, 1)
	switch mpacket.Packet.(*packet.Publish).GetQos() {
	case 0:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS0].success, 1)
	case 1:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS1].success, 1)
	case 2:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS2].success, 1)
	}
	cl.FireListener(EventPublish, "OnPubSuccess", createMqttPubEvent(client, PubStatusSuccess, cl.pubcnt), mpacket)
}
func (cl *mqttListen) fireOnPubFinal(client MqttClienter, mpacket *MqttPacket) {
	atomic.AddInt64(&cl.pubcnt[PubCntTOTAL].final, 1)
	switch mpacket.Packet.(*packet.Publish).GetQos() {
	case 0:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS0].final, 1)
	case 1:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS1].final, 1)
	case 2:
		atomic.AddInt64(&cl.pubcnt[PubCntQoS2].final, 1)
	}
	cl.FireListener(EventPublish, "OnPubFinal", createMqttPubEvent(client, PubStatusfinal, cl.pubcnt), mpacket)
}

// MqttRecvPubListener 接收消息事件监听
type MqttRecvPubListener interface {
	event.Listener
	//接收到了消息
	OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS)
}

//注册接收发布消息监听
func (cl *mqttListen) AddRecvPubListener(listener MqttRecvPubListener) {
	cl.AddListener(EventRecvPub, listener)
}

//移除接收发布消息监听
func (cl *mqttListen) RemoveRecvPubListener(listener MqttRecvPubListener) {
	cl.RemoveListener(EventRecvPub, listener)
}

func (cl *mqttListen) fireOnRecvPublish(client MqttClienter, topic string, payload []byte, qos QoS) {
	recvPubCnt := atomic.AddInt64(&cl.recvPubCnt, 1)

	cl.FireListener(EventRecvPub, "OnRecvPublish", createMqttRecvPubEvent(client, recvPubCnt), topic, payload, qos)
}

// MqttSubListener 订阅事件监听
type MqttSubListener interface {
	event.Listener
	OnSubStart(event *MqttEvent, sub []SubFilter)
	OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS)
}

//注册订阅消息监听
func (cl *mqttListen) AddSubListener(listener MqttSubListener) {
	cl.AddListener(EventSubscribe, listener)
}

//移除订阅消息监听
func (cl *mqttListen) RemoveSubListener(listener MqttSubListener) {
	cl.RemoveListener(EventSubscribe, listener)
}
func (cl *mqttListen) fireOnSubStart(client MqttClienter, sub []SubFilter) {
	cl.FireListener(EventSubscribe, "OnSubStart", createMqttEvent(client, EventSubscribe), sub)
}
func (cl *mqttListen) fireOnSubSuccess(client MqttClienter, sub []SubFilter, result []QoS) {
	cl.FireListener(EventSubscribe, "OnSubSuccess", createMqttEvent(client, EventSubscribe), sub, result)
}

// MqttUnSubListener 取消订阅事件监听
type MqttUnSubListener interface {
	event.Listener
	OnUnSubStart(event *MqttEvent, filter []string)
	OnUnSubSuccess(event *MqttEvent, filter []string)
}

//注册取消订阅消息监听
func (cl *mqttListen) AddUnSubListener(listener MqttUnSubListener) {
	cl.AddListener(EventUnSub, listener)
}

//移除取消订阅消息监听
func (cl *mqttListen) RemoveUnSubListener(listener MqttUnSubListener) {
	cl.RemoveListener(EventUnSub, listener)
}

func (cl *mqttListen) fireOnUnSubStart(client MqttClienter, filter []string) {
	cl.FireListener(EventUnSub, "OnUnSubStart", createMqttEvent(client, EventUnSub), filter)
}
func (cl *mqttListen) fireOnUnSubSuccess(client MqttClienter, filter []string) {
	cl.FireListener(EventUnSub, "OnUnSubSuccess", createMqttEvent(client, EventUnSub), filter)
}

// MqttPacketListener 发送接收报文接口
type MqttPacketListener interface {
	event.Listener
	//接收到报文
	OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64)
	//发送报文
	OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error)
}

//注册取消订阅消息监听
func (cl *mqttListen) AddPacketListener(listener MqttPacketListener) {
	cl.AddListener(EventPacket, listener)
}

//移除取消订阅消息监听
func (cl *mqttListen) RemovePacketListener(listener MqttPacketListener) {
	cl.RemoveListener(EventPacket, listener)
}

func (cl *mqttListen) fireOnRecvPacket(client MqttClienter, packet packet.MessagePacket) {
	recvPacketCnt := atomic.AddInt64(&cl.recvPacketCnt, 1)
	cl.FireListener(EventPacket, "OnRecvPacket", createMqttEvent(client, EventPacket), packet, recvPacketCnt)
}
func (cl *mqttListen) fireOnSendPacket(client MqttClienter, packet packet.MessagePacket, err error) {
	sendPacketCnt := atomic.AddInt64(&cl.sendPacketCnt, 1)
	cl.FireListener(EventPacket, "OnSendPacket", createMqttEvent(client, EventPacket), packet, sendPacketCnt, err)
}

// MqttDisConnListener 断开连接事件
type MqttDisConnListener interface {
	event.Listener
	OnDisconning(event *MqttEvent)
	OnDisconned(event *MqttEvent)
	OnLostConn(event *MqttEvent, err error)
}

//注册取消订阅消息监听
func (cl *mqttListen) AddDisConnListener(listener MqttDisConnListener) {
	cl.AddListener(EventDisconn, listener)
}

//移除取消订阅消息监听
func (cl *mqttListen) RemoveDisConnListener(listener MqttDisConnListener) {
	cl.RemoveListener(EventDisconn, listener)
}

func (cl *mqttListen) fireOnDisconning(client MqttClienter) {
	cl.FireListener(EventDisconn, "OnDisconning", createMqttEvent(client, EventDisconn))
}
func (cl *mqttListen) fireOnDisconned(client MqttClienter) {
	cl.FireListener(EventDisconn, "OnDisconned", createMqttEvent(client, EventDisconn))
}
func (cl *mqttListen) fireOnLostConn(client MqttClienter, err error) {
	cl.FireListener(EventDisconn, "OnLostConn", createMqttEvent(client, EventDisconn), err)
}

// DefaultListener 默认全部事件实现，不做任何事情
type DefaultListener struct {
	DefaultConnListen
	DefaultSubscribeListen
	DefaultPacketListen
	DefaultDisConnListen
	DefaultPubListen
	DefaultRecvPubListen
}

// DefaultConnListen 默认连接事件监听
type DefaultConnListen struct {
}

// OnConnStart 连接开始
func (*DefaultConnListen) OnConnStart(event *MqttConnEvent) {
}

// OnConnSuccess 连接开成功
func (*DefaultConnListen) OnConnSuccess(event *MqttConnEvent) {

}

// OnConnFailure 连接失败
func (*DefaultConnListen) OnConnFailure(event *MqttConnEvent, returncode int, err error) {

}

// DefaultSubscribeListen 默认订阅事件监听
type DefaultSubscribeListen struct {
}

// OnSubStart 订阅开始
func (*DefaultSubscribeListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
}

// OnSubSuccess 订阅成功
func (*DefaultSubscribeListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
}

// DefaultUnSubListen 取消订阅事件监听
type DefaultUnSubListen struct {
}

// OnUnSubStart 取消订阅开始
func (*DefaultUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {

}

// OnUnSubSuccess 取消订阅成功
func (*DefaultUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {

}

// DefaultPubListen 默认发布消息事件监听
type DefaultPubListen struct {
}

// OnPubReady 发布准备
func (*DefaultPubListen) OnPubReady(event *MqttPubEvent, mp *MqttPacket) {

}

// OnPubSuccess 发布成功
func (*DefaultPubListen) OnPubSuccess(event *MqttPubEvent, mp *MqttPacket) {

}

// OnPubFinal 发布完成
func (*DefaultPubListen) OnPubFinal(event *MqttPubEvent, mp *MqttPacket) {

}

// DefaultRecvPubListen 接收消息事件监听
type DefaultRecvPubListen struct {
}

// OnRecvPublish 接收到消息
// topic 主题
// payload 有效载荷
// qos 服务质量
func (l *DefaultRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {

}

// DefaultPacketListen 发送接收报文事件监听
type DefaultPacketListen struct {
}

// OnRecvPacket 接收报文事件
// event 事件
// msg 报文
// recvPacketCnt 接收报文数量
func (l *DefaultPacketListen) OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64) {

}

// OnSendPacket 发送报文事件
// event 事件
// msg 报文
// sendPacketCnt 发送报文数量
func (l *DefaultPacketListen) OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error) {

}

// DefaultDisConnListen 断开连接事件监听
type DefaultDisConnListen struct {
}

// OnDisconning 关闭连接中
func (l *DefaultDisConnListen) OnDisconning(event *MqttEvent) {

}

// OnDisconned 已经关闭连接
func (l *DefaultDisConnListen) OnDisconned(event *MqttEvent) {

}

// OnLostConn 失去连接事件
func (l *DefaultDisConnListen) OnLostConn(event *MqttEvent, err error) {

}

// DefaultPrintListener mqtt全部事件,可以输出事件信息
type DefaultPrintListener struct {
	DefaultPrintConnListen
	DefaultPrintSubscribeListen
	DefaultPrintPacketListen
	DefaultPrintDisConnListen
	DefaultPrintPubListen
	DefaultPrintRecvPubListen
}

// DefaultPrintConnListen 默认连接事件监听
type DefaultPrintConnListen struct {
}

// OnConnStart 连接开始
func (*DefaultPrintConnListen) OnConnStart(event *MqttConnEvent) {
	Mlog.Debugf("OnConnStart:共连接:%d,本次连接:%d", event.GetTotalRecnt(), event.GetCurRecnt())
}

// OnConnSuccess 连接开成功
func (*DefaultPrintConnListen) OnConnSuccess(event *MqttConnEvent) {
	Mlog.Debugf("OnConnSuccess:共连接:%d,本次连接:%d", event.GetTotalRecnt(), event.GetCurRecnt())
}

// OnConnFailure 连接失败
func (*DefaultPrintConnListen) OnConnFailure(event *MqttConnEvent, returncode int, err error) {
	Mlog.Debugf("OnConnFailure(%d):%v", returncode, err)
}

// DefaultPrintSubscribeListen 默认订阅事件监听
type DefaultPrintSubscribeListen struct {
}

// OnSubStart 订阅开始
func (*DefaultPrintSubscribeListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
	Mlog.Debugf("OnSubStart:%v", sub)
}

// OnSubSuccess 订阅成功
func (*DefaultPrintSubscribeListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	Mlog.Debugf("OnSubSuccess:%v:%v", sub, result)

}

// DefaultPrintUnSubListen 取消订阅事件监听
type DefaultPrintUnSubListen struct {
}

// OnUnSubStart 取消订阅开始
func (*DefaultPrintUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {
	Mlog.Debugf("OnUnSubStart:%v", filter)
}

// OnUnSubSuccess 取消订阅成功
func (*DefaultPrintUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {
	Mlog.Debugf("OnUnSubSuccess:%v", filter)
}

// DefaultPrintPubListen 默认发布消息事件监听
type DefaultPrintPubListen struct {
}

// OnPubReady 发布准备
func (*DefaultPrintPubListen) OnPubReady(event *MqttPubEvent, mp *MqttPacket) {
	Mlog.Debugf("OnPubReady:%v", event.GetSendCnt(PubCntTOTAL))
}

// OnPubSuccess 发布成功
func (*DefaultPrintPubListen) OnPubSuccess(event *MqttPubEvent, mp *MqttPacket) {
	Mlog.Debugf("OnPubSuccess:%v", mp.Packet)

}

// OnPubFinal 发布完成
func (*DefaultPrintPubListen) OnPubFinal(event *MqttPubEvent, mp *MqttPacket) {
	Mlog.Debugf("OnPubFinal:%v", mp.Packet)
}

// DefaultPrintRecvPubListen 接收消息事件监听
type DefaultPrintRecvPubListen struct {
}

// OnRecvPublish 接收到消息
// topic 主题
// payload 有效载荷
// qos 服务质量
func (l *DefaultPrintRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {
	Mlog.Debugf("OnRecvPublish:%s(%d) :%s", topic, qos, string(payload))
}

// DefaultPrintPacketListen 发送接收报文事件监听
type DefaultPrintPacketListen struct {
}

// OnRecvPacket 接收报文事件
// event 事件
// msg 报文
// recvPacketCnt 接收报文数量
func (l *DefaultPrintPacketListen) OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64) {
	Mlog.Debugf("OnRecvPacket:(%d) :%v", recvPacketCnt, msg)
}

// OnSendPacket 发送报文事件
// event 事件
// msg 报文
// sendPacketCnt 发送报文数量
func (l *DefaultPrintPacketListen) OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error) {
	Mlog.Debugf("OnSendPacket:(%d) :%v\n %v", sendPacketCnt, msg, err)
}

// DefaultPrintDisConnListen 断开连接事件监听
type DefaultPrintDisConnListen struct {
}

// OnLostConn 失去连接事件
func (l *DefaultPrintDisConnListen) OnLostConn(event *MqttEvent, err error) {
	Mlog.Debugf("OnLostConn :%v", err)
}

// OnDisconning 关闭连接中
func (l *DefaultPrintDisConnListen) OnDisconning(event *MqttEvent) {
	Mlog.Debugf("OnDisconning")
}

// OnDisconned 已经关闭连接
func (l *DefaultPrintDisConnListen) OnDisconned(event *MqttEvent) {
	Mlog.Debugf("OnDisconned")
}
