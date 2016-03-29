package client

/**
存储重发队列
连接参数
订阅发布
取消订阅
消息发布
重连

事件机制

	连接开始
	连接成功
	连接失败

	正常断开连接
	异常断开连接

	发布消息准备成功
	发布消息成功
	发布消息完成

	收到消息


	发出订阅
	订阅返回

	发去取消订阅
	取消订阅返回


**/

type QoS uint8

const (
	QoS0 QoS = iota
	QoS1
	QoS2
	Refused = 0x80
)

func CreateSubFilter(filter string, qos QoS) SubFilter {
	return SubFilter{
		filter: filter,
		qos:    qos,
	}
}

type SubFilter struct {
	filter string
	qos    QoS
}

type MqttClienter interface {
	MqttConner
	MqttPublisher
	MqttSubscriber
	MqttUnSubscriber
	MqttPackerListener
	MqttDisConner
}

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

type MqttPublisher interface {
	//发布消息
	Publish(topic string, qos QoS, retain bool, payload interface{}) (*MqttPacket, error)
	//注册发布消息监听
	AddPublishListener(listener MqttPublishListener)
	//移除发布消息监听
	RemovePublishListener(listener MqttPublishListener)

	//注册接收消息监听
	AddRecvPubListener(listener MqttRecvPubListener)
	//移除发布消息监听
	RemoveRecvPubListener(listener MqttRecvPubListener)
}

type MqttSubscriber interface {
	//订阅
	Subscribe(filter string, qos QoS) (*MqttPacket, error)
	//批量订阅
	Subscribes(filters ...SubFilter) (*MqttPacket, error)
	AddSubListener(listener MqttSubListener)
	RemoveSubListener(listener MqttSubListener)
}
type MqttUnSubscriber interface {
	//取消订阅
	UnSubscribe(filter string) (*MqttPacket, error)
	//批量取消订阅
	UnSubscribes(filters ...string) (*MqttPacket, error)
	AddUnSubListener(listener MqttUnSubListener)
	RemoveUnSubListener(listener MqttUnSubListener)
}

type MqttPackerListener interface {
	//增加报文接收发送监听
	AddPacketListener(listener MqttPacketListener)
	//移除报文接收发送监听
	RemovePacketListener(listener MqttPacketListener)
}
type MqttDisConner interface {
	//断开连接
	Disconnect()
	//增加报文接收发送监听
	AddDisConnListener(listener MqttDisConnListener)
	//移除报文接收发送监听
	RemoveDisConnListener(listener MqttDisConnListener)
}

type BaseClientStatus struct {
	//当前重连次数统计
	curRecnt int64
	//总重连次数统计
	totalRecnt int64
	//接收发布消息报文统计
	recvPubCnt int64
	//发布消息次数统计
	pubcnt map[CntType]*pubCnt

	//接收报文数量
	recvPacketCnt int64
	//发送报文数量
	sendPacketCnt int64
}

func (c *BaseClientStatus) init() {
	c.pubcnt = make(map[CntType]*pubCnt)
	c.pubcnt[PubCnt_QoS0] = &pubCnt{}
	c.pubcnt[PubCnt_QoS1] = &pubCnt{}
	c.pubcnt[PubCnt_QoS2] = &pubCnt{}
	c.pubcnt[PubCnt_TOTAL] = &pubCnt{}
}
