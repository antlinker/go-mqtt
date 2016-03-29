package client

import (
	"time"

	"github.com/antlinker/go-mqtt/event"
)

const (
	EVENT_CONN int = iota
	EVENT_BULISH
	EVENT_RECVBULISH
	EVENT_SUBSCRIBE
	EVENT_UNSUBSCRIBE
	EVENT_PACKET
	EVENT_DISCONNECT
)

//基础mqtt事件
type MqttEvent struct {
	event.BaseEvent
	client    MqttClienter
	occurtime time.Time
}

func createMqttEvent(client MqttClienter, mtype int) *MqttEvent {
	event := &MqttEvent{}
	event.Init(mtype, client)
	return event
}
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

type ConnStatus int

const (
	STATUS_CONN_START ConnStatus = iota
	STATUS_CONN_SUCCESS
	STATUS_CONN_FAILURE
)

//mqtt连接事件
type MqttConnEvent struct {
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
	event.Init(EVENT_CONN, client)
	event.status = status
	event.curRecnt = curRecnt
	event.totalRecnt = totalRecnt
	return event

}

//统计发布数据类型
type CntType int

const (
	PubCnt_QoS0 CntType = iota
	PubCnt_QoS1
	PubCnt_QoS2
	PubCnt_TOTAL
)

//创建连接事件
func createMqttPubEvent(client MqttClienter, status PubStatus, pubcnt map[CntType]*pubCnt) *MqttPubEvent {
	event := &MqttPubEvent{}
	event.Init(EVENT_BULISH, client)
	event.status = status
	event.pubcnt = pubcnt
	return event

}

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
	event.Init(EVENT_RECVBULISH, client)
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
