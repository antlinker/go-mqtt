package client

import (
	"time"

	"github.com/antlinker/go-mqtt/packet"
)

// PacketManager 报文发送接收管理
// 实现 subscribe unsubscrie publish报文发送时的报文标识分配
// 完成qos1 qos2报文的重发
// 记录 subscribe unsubscrie publish 发送时记录,完成时移除
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
	//qos2存储id,重复id只保留一个
	AddReceivePacket(*MqttPacket)
	RemoveReceivePacket(id uint16) *MqttPacket
}

const (
	// DirectSend 发送
	DirectSend int = iota
	// DirectRecive 接收
	DirectRecive
)

// CreateMqttPacket 创建发送报文
func CreateMqttPacket(direct int, msg packet.PacketIdMessage) *MqttPacket {
	return &MqttPacket{
		Direct:  direct,
		Packet:  msg,
		Lastime: time.Now(),
		wait:    make(chan struct{}),
	}
}

// MqttPacket 发送时
//发送的订阅取消订阅或发布消息报文
//订阅、取消订阅报文发送失败不会重新发送
//订阅报文直到收到Suback报文，Wait()才会返回，触发OnSubscribeSuccess事件
//取消订阅报文直到收到UnSuback报文，Wait()才会返回，触发OnUnSubscribeSuccess事件
//发布消息报文，qos0发送不是不会重新发布 发送后Wait()返回，触发OnPubFinal事件
//		qos1 收到Puback报文后Wait()返回，触发OnPubFinal事件
//		qos2收到Pubcomp报文后Wait()才会返回，触发OnPubFinal事件
//Pulish报文发送后触发OnPubSuccess事件
//接收时，只有少量属性生效，大部分属性不生效
type MqttPacket struct {
	//发送或接收
	Direct int
	//发送的消息报文
	Packet packet.PacketIdMessage
	//最后一次更新事件
	Lastime time.Time
	wait    chan struct{}
	err     error
	rectime time.Time
	sending bool
}

// IsSending 是否正在发送中
func (m *MqttPacket) IsSending() bool {
	return m.sending
}

// Rectime qos2的Pubrec报文接收事件
func (m *MqttPacket) Rectime() time.Time {
	return m.rectime
}

// Err 错误，如果不为nil发送失败
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

// WaitTimeout 等待发送完成、失败，或超时
//返回　true等待了内部返回　false未等到返回超时
//注：并不意味着返回false就是发送失败，必须判断Err()==nil则发送失败
func (m *MqttPacket) WaitTimeout(timeout time.Duration) bool {
	select {
	case <-m.wait:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Wait 等待发送完成或发送失败
func (m *MqttPacket) Wait() {
	<-m.wait
}
