package client

// QoS 服务质量
type QoS uint8

const (
	// QoS0 服务质量0
	QoS0 QoS = iota
	// QoS1 服务质量1
	QoS1
	// QoS2 服务质量2
	QoS2
	// Refused 服务器返回拒绝订阅
	Refused = 0x80
)

// CreateSubFilter 创建订阅对象
// filter　订阅规则
// qos 服务质量
func CreateSubFilter(filter string, qos QoS) SubFilter {
	return SubFilter{
		filter: filter,
		qos:    qos,
	}
}

// SubFilter 订阅结构体
type SubFilter struct {
	filter string
	qos    QoS
}

func (f SubFilter) Filter() string {
	return f.filter
}
func (f SubFilter) QoS() QoS {
	return f.qos
}

// MqttClienter MQTT客户端
type MqttClienter interface {
	MqttConner
	MqttPublisher
	MqttSubscriber
	MqttUnSubscriber
	MqttPackerListener
	MqttDisConner
	GetOpt() MqttOption
}

// MqttConner MQTT连接器
type MqttConner interface {
	//开始连接
	Connect() error
	//是否已连接
	IsConnect() bool
	//注册连接监听
	AddConnListener(listener MqttConnListener)
	//移除连接监听
	RemoveConnListener(listener MqttConnListener)
}

// MqttPublisher MQTT消息发布消息接收
type MqttPublisher interface {
	//发布消息
	Publish(topic string, qos QoS, retain bool, payload interface{}) (*MqttPacket, error)
	//注册发布消息监听
	AddPubListener(listener MqttPubListener)
	//移除发布消息监听
	RemovePubListener(listener MqttPubListener)

	//注册接收消息监听
	AddRecvPubListener(listener MqttRecvPubListener)
	//移除发布消息监听
	RemoveRecvPubListener(listener MqttRecvPubListener)
}

// MqttSubscriber MQTT订阅
type MqttSubscriber interface {
	//订阅
	Subscribe(filter string, qos QoS) (*MqttPacket, error)
	//批量订阅
	Subscribes(filters ...SubFilter) (*MqttPacket, error)
	AddSubListener(listener MqttSubListener)
	RemoveSubListener(listener MqttSubListener)
}

// MqttUnSubscriber MQTT取消订阅
type MqttUnSubscriber interface {
	//取消订阅
	UnSubscribe(filter string) (*MqttPacket, error)
	//批量取消订阅
	UnSubscribes(filters ...string) (*MqttPacket, error)
	AddUnSubListener(listener MqttUnSubListener)
	RemoveUnSubListener(listener MqttUnSubListener)
}

// MqttPackerListener MQTT报文发送接收
type MqttPackerListener interface {
	//增加报文接收发送监听
	AddPacketListener(listener MqttPacketListener)
	//移除报文接收发送监听
	RemovePacketListener(listener MqttPacketListener)
}

// MqttDisConner MQTT连接断开
type MqttDisConner interface {
	//断开连接
	Disconnect()
	//增加报文接收发送监听
	AddDisConnListener(listener MqttDisConnListener)
	//移除报文接收发送监听
	RemoveDisConnListener(listener MqttDisConnListener)
}
