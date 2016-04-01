package client

import (
	"time"

	"github.com/antlinker/go-mqtt/event"
)

const (
	// EventConn 连接事件
	EventConn int = iota
	// EventPublish 发布消息事件
	EventPublish
	// EventRecvPub 接收消息事件
	EventRecvPub
	// EventSubscribe 订阅主题事件
	EventSubscribe
	// EventUnSub 取消订阅事件
	EventUnSub
	// EventPacket 消息报文事件
	EventPacket
	// EventDisconn 断开连接事件
	EventDisconn
)

// MqttEvent 基础mqtt事件
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

// Init 初始化事件
//etype 事件类型
//client 客户端
func (e *MqttEvent) Init(etype int, client MqttClienter) {
	e.BaseEvent.Init(etype, client)
	e.client = client
	e.occurtime = time.Now()
}

// GetClient 获取mqtt客户端
func (e *MqttEvent) GetClient() MqttClienter {
	return e.client
}

// GetOccurtime 获取事件发生事件
func (e *MqttEvent) GetOccurtime() time.Time {
	return e.occurtime
}

// ConnStatus 连接状态
type ConnStatus int

const (
	// StatusConnStart 连接开始事件
	StatusConnStart ConnStatus = iota
	// StatusConnSuccess 连接成功
	StatusConnSuccess
	// StatusConnFailure 连接失败
	StatusConnFailure
)

// MqttConnEvent mqtt连接事件
type MqttConnEvent struct {
	//事件基础实现
	MqttEvent
	status ConnStatus
	//总连接次数
	totalRecnt int64
	//本次连接次数
	curRecnt int64
}

// GetTotalRecnt 获取连接状态
func (e *MqttConnEvent) GetTotalRecnt() int64 {
	return e.totalRecnt
}

// GetCurRecnt 获取连接状态
func (e *MqttConnEvent) GetCurRecnt() int64 {
	return e.curRecnt
}

// GetStatus 获取连接状态
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

// CntType 统计发布数据类型
type CntType int

const (
	// PubCntQoS0 QoS0统计
	PubCntQoS0 CntType = iota
	// PubCntQoS1 QoS1统计
	PubCntQoS1
	//PubCntQoS2 QoS2统计
	PubCntQoS2
	//PubCntTOTAL 总发布消息统计
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

// PubStatus 发布状态
type PubStatus int

const (
	// PubStatusReady 发布开始
	PubStatusReady PubStatus = iota
	// PubStatusSuccess 发布成功，已经调用conn.Write发送完成
	PubStatusSuccess
	// PubStatusfinal 发送完成，qos1 qos2按协议报文发送完成
	PubStatusfinal
)

// MqttPubEvent mqtt连接事件
type MqttPubEvent struct {
	MqttEvent
	status PubStatus
	pubcnt map[CntType]*pubCnt
}

// GetFinalCnt 获取完成报文统计数据
func (e *MqttPubEvent) GetFinalCnt(cntType CntType) int64 {
	return e.pubcnt[cntType].final
}

// GetSendCnt 获取发送报文统计数据
func (e *MqttPubEvent) GetSendCnt(cntType CntType) int64 {
	return e.pubcnt[cntType].send
}

// GetSuccessCnt 获取成功发送报文统计数据
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

// MqttRecvPubEvent mqtt连接事件
type MqttRecvPubEvent struct {
	MqttEvent
	recvCnt int64
}

// GetRecvCnt 获取完成报文统计数据
func (e *MqttRecvPubEvent) GetRecvCnt() int64 {
	return e.recvCnt
}
