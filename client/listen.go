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

//连接事件
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
	event.EventGenerator
	baseClientStatus
}

//注册连接监听
func (cl *mqttListen) AddConnListener(listener MqttConnListener) {
	cl.AddListener(EVENT_CONN, listener)
}

//移除连接监听
func (cl *mqttListen) RemoveConnListener(listener MqttConnListener) {
	cl.RemoveListener(EVENT_CONN, listener)
}

func (cl *mqttListen) fireOnConnStart(client MqttClienter) {
	//Mlog.Debug("fireOnConnStart")
	cl.FireListener(EVENT_CONN, "OnConnStart", createMqttConnEvent(client, STATUS_CONN_START, atomic.AddInt64(&cl.totalRecnt, 1), atomic.AddInt64(&cl.curRecnt, 1)))
}
func (cl *mqttListen) fireOnConnSuccess(client MqttClienter) {
	cl.curRecnt = 0
	cl.FireListener(EVENT_CONN, "OnConnSuccess", createMqttConnEvent(client, STATUS_CONN_SUCCESS, cl.curRecnt, cl.totalRecnt))
}
func (cl *mqttListen) fireOnConnFailure(client MqttClienter, returncode int, err error) {
	cl.FireListener(EVENT_CONN, "OnConnFailure", createMqttConnEvent(client, STATUS_CONN_FAILURE, cl.curRecnt, cl.totalRecnt), returncode, err)
}

//发布消息事件监听
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
	cl.AddListener(EVENT_BULISH, listener)
}

//移除发布消息监听
func (cl *mqttListen) RemovePubListener(listener MqttPubListener) {
	cl.RemoveListener(EVENT_BULISH, listener)
}

func (cl *mqttListen) fireOnPubReady(client MqttClienter, mpacket *MqttPacket) {
	atomic.AddInt64(&(cl.pubcnt[PubCnt_TOTAL].send), 1)
	switch mpacket.Packet.(*packet.Publish).GetQos() {
	case 0:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS0].send, 1)
	case 1:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS1].send, 1)
	case 2:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS2].send, 1)
	}
	cl.FireListener(EVENT_BULISH, "OnPubReady", createMqttPubEvent(client, PubStatus_Ready, cl.pubcnt), mpacket)
}
func (cl *mqttListen) fireOnPubSuccess(client MqttClienter, mpacket *MqttPacket) {
	atomic.AddInt64(&cl.pubcnt[PubCnt_TOTAL].success, 1)
	switch mpacket.Packet.(*packet.Publish).GetQos() {
	case 0:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS0].success, 1)
	case 1:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS1].success, 1)
	case 2:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS2].success, 1)
	}
	cl.FireListener(EVENT_BULISH, "OnPubSuccess", createMqttPubEvent(client, PubStatus_Success, cl.pubcnt), mpacket)
}
func (cl *mqttListen) fireOnPubFinal(client MqttClienter, mpacket *MqttPacket) {
	atomic.AddInt64(&cl.pubcnt[PubCnt_TOTAL].final, 1)
	switch mpacket.Packet.(*packet.Publish).GetQos() {
	case 0:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS0].final, 1)
	case 1:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS1].final, 1)
	case 2:
		atomic.AddInt64(&cl.pubcnt[PubCnt_QoS2].final, 1)
	}
	cl.FireListener(EVENT_BULISH, "OnPubFinal", createMqttPubEvent(client, PubStatus_final, cl.pubcnt), mpacket)
}

//接收消息事件监听
type MqttRecvPubListener interface {
	event.Listener
	//接收到了消息
	OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS)
}

//注册接收发布消息监听
func (cl *mqttListen) AddRecvPubListener(listener MqttRecvPubListener) {
	cl.AddListener(EVENT_RECVBULISH, listener)
}

//移除接收发布消息监听
func (cl *mqttListen) RemoveRecvPubListener(listener MqttRecvPubListener) {
	cl.RemoveListener(EVENT_RECVBULISH, listener)
}

func (cl *mqttListen) fireOnRecvPublish(client MqttClienter, topic string, payload []byte, qos QoS) {
	recvPubCnt := atomic.AddInt64(&cl.recvPubCnt, 1)

	cl.FireListener(EVENT_RECVBULISH, "OnRecvPublish", createMqttRecvPubEvent(client, recvPubCnt), topic, payload, qos)
}

type MqttSubListener interface {
	event.Listener
	OnSubStart(event *MqttEvent, sub []SubFilter)
	OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS)
}

//注册订阅消息监听
func (cl *mqttListen) AddSubListener(listener MqttSubListener) {
	cl.AddListener(EVENT_SUBSCRIBE, listener)
}

//移除订阅消息监听
func (cl *mqttListen) RemoveSubListener(listener MqttSubListener) {
	cl.RemoveListener(EVENT_SUBSCRIBE, listener)
}
func (cl *mqttListen) fireOnSubStart(client MqttClienter, sub []SubFilter) {
	cl.FireListener(EVENT_SUBSCRIBE, "OnSubStart", createMqttEvent(client, EVENT_SUBSCRIBE), sub)
}
func (cl *mqttListen) fireOnSubSuccess(client MqttClienter, sub []SubFilter, result []QoS) {
	cl.FireListener(EVENT_SUBSCRIBE, "OnSubSuccess", createMqttEvent(client, EVENT_SUBSCRIBE), sub, result)
}

type MqttUnSubListener interface {
	event.Listener
	OnUnSubStart(event *MqttEvent, filter []string)
	OnUnSubSuccess(event *MqttEvent, filter []string)
}

//注册取消订阅消息监听
func (cl *mqttListen) AddUnSubListener(listener MqttUnSubListener) {
	cl.AddListener(EVENT_UNSUBSCRIBE, listener)
}

//移除取消订阅消息监听
func (cl *mqttListen) RemoveUnSubListener(listener MqttUnSubListener) {
	cl.RemoveListener(EVENT_UNSUBSCRIBE, listener)
}

func (cl *mqttListen) fireOnUnSubStart(client MqttClienter, filter []string) {
	cl.FireListener(EVENT_UNSUBSCRIBE, "OnUnSubStart", createMqttEvent(client, EVENT_UNSUBSCRIBE), filter)
}
func (cl *mqttListen) fireOnUnSubSuccess(client MqttClienter, filter []string) {
	cl.FireListener(EVENT_UNSUBSCRIBE, "OnUnSubSuccess", createMqttEvent(client, EVENT_UNSUBSCRIBE), filter)
}

//发送接收报文接口
type MqttPacketListener interface {
	event.Listener
	//接收到报文
	OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64)
	//发送报文
	OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error)
}

//注册取消订阅消息监听
func (cl *mqttListen) AddPacketListener(listener MqttPacketListener) {
	cl.AddListener(EVENT_PACKET, listener)
}

//移除取消订阅消息监听
func (cl *mqttListen) RemovePacketListener(listener MqttPacketListener) {
	cl.RemoveListener(EVENT_PACKET, listener)
}

func (cl *mqttListen) fireOnRecvPacket(client MqttClienter, packet packet.MessagePacket) {
	recvPacketCnt := atomic.AddInt64(&cl.recvPacketCnt, 1)
	cl.FireListener(EVENT_PACKET, "OnRecvPacket", createMqttEvent(client, EVENT_PACKET), packet, recvPacketCnt)
}
func (cl *mqttListen) fireOnSendPacket(client MqttClienter, packet packet.MessagePacket, err error) {
	sendPacketCnt := atomic.AddInt64(&cl.sendPacketCnt, 1)
	cl.FireListener(EVENT_PACKET, "OnSendPacket", createMqttEvent(client, EVENT_PACKET), packet, sendPacketCnt, err)
}

type MqttDisConnListener interface {
	event.Listener
	OnDisconning(event *MqttEvent)
	OnDisconned(event *MqttEvent)
	OnLostConn(event *MqttEvent, err error)
}

//注册取消订阅消息监听
func (cl *mqttListen) AddDisConnListener(listener MqttDisConnListener) {
	cl.AddListener(EVENT_DISCONNECT, listener)
}

//移除取消订阅消息监听
func (cl *mqttListen) RemoveDisConnListener(listener MqttDisConnListener) {
	cl.RemoveListener(EVENT_DISCONNECT, listener)
}

func (cl *mqttListen) fireOnDisconning(client MqttClienter) {
	cl.FireListener(EVENT_DISCONNECT, "OnDisconning", createMqttEvent(client, EVENT_DISCONNECT))
}
func (cl *mqttListen) fireOnDisconned(client MqttClienter) {
	cl.FireListener(EVENT_DISCONNECT, "OnDisconned", createMqttEvent(client, EVENT_DISCONNECT))
}
func (cl *mqttListen) fireOnLostConn(client MqttClienter, err error) {
	cl.FireListener(EVENT_DISCONNECT, "OnLostConn", createMqttEvent(client, EVENT_DISCONNECT), err)
}

type DefaultListener struct {
	DefaultConnListen
	DefaultSubscribeListen
	DefaultPacketListen
	DefaultDisConnListen
	DefaultPubListen
	DefaultRecvPubListen
}
type DefaultConnListen struct {
}

func (*DefaultConnListen) OnConnStart(event *MqttConnEvent) {
}
func (*DefaultConnListen) OnConnSuccess(event *MqttConnEvent) {

}
func (*DefaultConnListen) OnConnFailure(event *MqttConnEvent, returncode int, err error) {

}

type DefaultSubscribeListen struct {
}

func (*DefaultSubscribeListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
}
func (*DefaultSubscribeListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
}

type DefaultUnSubListen struct {
}

func (*DefaultUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {

}
func (*DefaultUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {

}

type DefaultPubListen struct {
}

func (*DefaultPubListen) OnPubReady(event *MqttPubEvent, mp *MqttPacket) {

}
func (*DefaultPubListen) OnPubSuccess(event *MqttPubEvent, mp *MqttPacket) {

}
func (*DefaultPubListen) OnPubFinal(event *MqttPubEvent, mp *MqttPacket) {

}

type DefaultRecvPubListen struct {
}

func (l *DefaultRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {

}

type DefaultPacketListen struct {
}

func (l *DefaultPacketListen) OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64) {

}
func (l *DefaultPacketListen) OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error) {

}

type DefaultDisConnListen struct {
}

func (l *DefaultDisConnListen) OnDisconning(event *MqttEvent) {

}
func (l *DefaultDisConnListen) OnDisconned(event *MqttEvent) {

}
func (l *DefaultDisConnListen) OnLostConn(event *MqttEvent, err error) {

}

type DefaultPrintListener struct {
	DefaultPrintConnListen
	DefaultPrintSubscribeListen
	DefaultPrintPacketListen
	DefaultPrintDisConnListen
	DefaultPrintPubListen
	DefaultPrintRecvPubListen
}
type DefaultPrintConnListen struct {
}

func (*DefaultPrintConnListen) OnConnStart(event *MqttConnEvent) {
	Mlog.Debugf("OnConnStart:共连接:%d,本次连接:%d", event.GetTotalRecnt(), event.GetCurRecnt())
}
func (*DefaultPrintConnListen) OnConnSuccess(event *MqttConnEvent) {
	Mlog.Debugf("OnConnSuccess:共连接:%d,本次连接:%d", event.GetTotalRecnt(), event.GetCurRecnt())
}
func (*DefaultPrintConnListen) OnConnFailure(event *MqttConnEvent, returncode int, err error) {
	Mlog.Debugf("OnConnFailure(%d):%v", returncode, err)
}

type DefaultPrintSubscribeListen struct {
}

func (*DefaultPrintSubscribeListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
	Mlog.Debugf("OnSubStart:%v", sub)
}
func (*DefaultPrintSubscribeListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	Mlog.Debugf("OnSubSuccess:%v:%v", sub, result)

}

type DefaultPrintUnSubListen struct {
}

func (*DefaultPrintUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {
	Mlog.Debugf("OnUnSubStart:%v", filter)
}
func (*DefaultPrintUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {
	Mlog.Debugf("OnUnSubSuccess:%v", filter)
}

type DefaultPrintPubListen struct {
}

func (*DefaultPrintPubListen) OnPubReady(event *MqttPubEvent, mp *MqttPacket) {
	Mlog.Debugf("OnPubReady:%v", event.GetSendCnt(PubCnt_TOTAL))
}
func (*DefaultPrintPubListen) OnPubSuccess(event *MqttPubEvent, mp *MqttPacket) {
	Mlog.Debugf("OnPubSuccess:%v", mp.Packet)

}
func (*DefaultPrintPubListen) OnPubFinal(event *MqttPubEvent, mp *MqttPacket) {
	Mlog.Debugf("OnPubFinal:%v", mp.Packet)
}

type DefaultPrintRecvPubListen struct {
}

func (l *DefaultPrintRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {
	Mlog.Debugf("OnRecvPublish:%s(%d) :%s", topic, qos, string(payload))
}

type DefaultPrintPacketListen struct {
}

func (l *DefaultPrintPacketListen) OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64) {
	Mlog.Debugf("OnRecvPacket:(%d) :%v", recvPacketCnt, msg)
}
func (l *DefaultPrintPacketListen) OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error) {
	Mlog.Debugf("OnSendPacket:(%d) :%v\n %v", sendPacketCnt, msg, err)
}

type DefaultPrintDisConnListen struct {
}

func (l *DefaultPrintDisConnListen) OnLostConn(event *MqttEvent, err error) {
	Mlog.Debugf("OnLostConn :%v", err)
}

func (l *DefaultPrintDisConnListen) OnDisconning(event *MqttEvent) {
	Mlog.Debugf("OnDisconning")
}
func (l *DefaultPrintDisConnListen) OnDisconned(event *MqttEvent) {
	Mlog.Debugf("OnDisconned")
}
