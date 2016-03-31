package client

import (
	"time"

	"github.com/antlinker/go-mqtt/event"
)

const (
	//连接事件
	EventConn int = iota
	//发布消息事件
	EventPublish
	//接收消息事件
	EventRecvPub
	//订阅主题事件
	EventSubscribe
	//取消订阅事件
	EventUnSub
	//消息报文事件
	EventPacket
	//断开连接事件
	EventDisconn
)

//基础mqtt事件
type MqttEvent struct {
	//事件基础实现
	event.BaseEvent
	client    MqttClienter
	occurtime time.Time
}

func createMqttEvent(client MqttClienter, mtype int) *MqttEvent {
	event := &MqttEvent{}
	event.Init(mtype, client)
	return event
}

//初始化事件
//etype 事件类型
//client 客户端
func (e *MqttEvent) Init(etype int, client MqttClienter) {
	e.BaseEvent.Init(etype, client)
	e.client = client
	e.occurtime = time.Now()
}

//获取mqtt客户端
func (e *MqttEvent) GetClient() MqttClienter {
	return e.client
}

//获取事件发生事件
func (e *MqttEvent) GetOccurtime() time.Time {
	return e.occurtime
}

//连接状态
type ConnStatus int

const (
	//连接开始事件
	StatusConnStart ConnStatus = iota
	//连接成功
	StatusConnSuccess
	//连接失败
	StatusConnFailure
)

//mqtt连接事件
type MqttConnEvent struct {
	//事件基础实现
	MqttEvent
	status ConnStatus
	//总连接次数
	totalRecnt int64
	//本次连接次数
	curRecnt int64
}

//获取连接状态
func (e *MqttConnEvent) GetTotalRecnt() int64 {
	return e.totalRecnt
}

//获取连接状态
func (e *MqttConnEvent) GetCurRecnt() int64 {
	return e.curRecnt
}

//获取连接状态
func (e *MqttConnEvent) GetStatus() ConnStatus {
	return e.status
}

//创建连接事件
func createMqttConnEvent(client MqttClienter, status ConnStatus, curRecnt, totalRecnt int64) *MqttConnEvent {
	event := &MqttConnEvent{}
	event.Init(EventConn, client)
	event.status = status
	event.curRecnt = curRecnt
	event.totalRecnt = totalRecnt
	return event

}

//统计发布数据类型
type CntType int

const (
	PubCntQoS0 CntType = iota
	PubCntQoS1
	PubCntQoS2
	PubCntTOTAL
)

//创建连接事件
func createMqttPubEvent(client MqttClienter, status PubStatus, pubcnt map[CntType]*pubCnt) *MqttPubEvent {
	event := &MqttPubEvent{}
	event.Init(EventPublish, client)
	event.status = status
	event.pubcnt = pubcnt
	return event

}

//发布状态
type PubStatus int

const (
	PubStatus_Ready PubStatus = iota
	PubStatus_Success
	PubStatus_final
)

//mqtt连接事件
type MqttPubEvent struct {
	MqttEvent
	status PubStatus
	pubcnt map[CntType]*pubCnt
}

//获取完成报文统计数据
func (e *MqttPubEvent) GetFinalCnt(cntType CntType) int64 {
	return e.pubcnt[cntType].final
}

//获取发送报文统计数据
func (e *MqttPubEvent) GetSendCnt(cntType CntType) int64 {
	return e.pubcnt[cntType].send
}

//获取成功发送报文统计数据
func (e *MqttPubEvent) GetSuccessCnt(cntType CntType) int64 {
	return e.pubcnt[cntType].success
}

//创建连接事件
func createMqttRecvPubEvent(client MqttClienter, recvCnt int64) *MqttRecvPubEvent {
	event := &MqttRecvPubEvent{}
	event.Init(EventRecvPub, client)
	event.recvCnt = recvCnt
	return event

}

//mqtt连接事件
type MqttRecvPubEvent struct {
	MqttEvent
	recvCnt int64
}

//获取完成报文统计数据
func (e *MqttRecvPubEvent) GetRecvCnt() int64 {
	return e.recvCnt
}
