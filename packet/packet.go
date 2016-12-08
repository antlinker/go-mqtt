package packet

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// QoS 服务质量
type QoS uint8

// PacketIdMessage 带有报文标识的报文
type PacketIdMessage interface {
	MessagePacket
	GetPacketId() uint16
	SetPacketId(packetID uint16)
}

// PublishBackSuccessMessage qos1 qos2报文成功发送后的回执
type PublishBackSuccessMessage interface {
	GetPacketId() uint16
}

// Value 服务质量
func (qos QoS) Value() uint8 {
	return uint8(qos)
}
func (qos QoS) String() string {
	return strconv.Itoa(int(qos))
}

type Topic struct {
	value     string
	topicnode []string
	len       int
}

func NewTopic(topic string) *Topic {
	t := &Topic{value: topic}
	t.topicnode = splitTopic(topic)
	t.len = len([]byte(topic))
	return t
}

// IsValidTopic 是否是有效主题
func (t *Topic) IsValidTopic() bool {
	if strings.Index(t.value, "+") >= 0 {
		return false
	}
	if strings.Index(t.value, "#") >= 0 {
		return false
	}
	return true
}

func (t *Topic) Len() int {
	return t.len
}
func (t *Topic) BytesValue() []byte {
	return []byte(t.value)
}
func (t *Topic) String() string {
	return t.value
}

func (t *Topic) GetTopicnode() []string {
	return t.topicnode
}

//比较主题和主题过滤器
func _compare(topics []string, filters []string) bool {

	ti := 0
	fi := 0
	if topics[0] == "$" && filters[0] != "$" {
		return false
	}
	for ti < len(topics) && fi < len(filters) {
		if topics[ti] == filters[fi] || (filters[fi] == "+" && !(ti == 0 && topics[ti] == "$")) {
			ti++
			fi++
		} else if filters[fi] == "#" {
			if fi == len(filters)-1 {
				return true
			}

			for ; ti < len(topics); ti++ {
				if _compare(topics[ti:], filters[fi+1:]) {
					return true
				}
			}
			return false
		} else {
			return false
		}

	}

	return len(topics) == len(filters)
}

// CompareFilter 比较主题核主题是否被过滤器所匹配
// 返回 true 匹配 false 不匹配
func (t *Topic) CompareFilter(filter *Topic) bool {
	return _compare(t.topicnode, filter.topicnode)
}

//分割主题过滤器和主题
func splitTopic(topic string) []string {
	out := make([]string, 0, 4)
	data := topic
	if topic[0] == '$' {
		out = append(out, "$")

		data = string([]byte(topic)[1:])

	}
	tmp := strings.Split(data, "/")

	out = append(out, tmp...)
	return out
}

//只分析了固定头还未解包的消息报文
type MessageUnPacket struct {
	FixedHeader
	remdata []byte
}

func NewMessageUnPacket(head byte, remlen int, remdata []byte) *MessageUnPacket {
	mup := &MessageUnPacket{}
	mup.hdata = head
	mup.remlen = remlen
	mup.remdata = remdata
	return mup
}
func (m *MessageUnPacket) String() string {
	return fmt.Sprintf("报文类型:%s (%X)(%X)(%d):%X", GetTypeName(m.GetHeaderType()), m.GetHeaderType(), m.GetHeaderTypeFlag(), m.GetRemlen(), m.remdata)
}
func (m *MessageUnPacket) GetRemdata() []byte {
	return m.remdata
}
func (m *MessageUnPacket) Parse() (MessagePacket, error) {
	msg, err := createNewMessagePacket(m.hdata)
	if err != nil {
		return nil, err
	}
	err = msg.UnPacket(m.hdata, m.remdata)
	return msg, err
}
func (m *MessageUnPacket) UnPacket(head byte, data []byte) error {
	m.hdata = head
	m.remlen = len(data)
	m.remdata = data
	return nil
}
func (m *MessageUnPacket) Packet() []byte {

	//fmt.Printf("剩余字节%d %X\n", c.remlen, c.remlen)
	remlenbyte := Remlen2Bytes(int32(m.remlen))
	//fmt.Printf("剩余字节 %X\n", remlenbyte)
	data := make([]byte, int(m.remlen)+1+len(remlenbyte))
	curind := 0
	curind = BytesWByte(data, m.hdata, curind)
	curind = BytesWBytes(data, remlenbyte, curind)
	curind = BytesWBytes(data, m.remdata, curind)
	return data
}

//消息接口
type MessagePacket interface {
	String() string
	GetHeaderType() uint8
	UnPacket(head byte, data []byte) error
	Packet() []byte
	WriteTo(out io.Writer) (int64, error)
	// //从文件解码 主要为PUBLIST类型报文使用，大数据传输缓存中读取数据
	// UnPacketFile(head byte, vheader []byte, in io.Reader) error
	//打包到
	PacketTo(out io.Writer, size int) error
	//获取数据总的字节长度
	GetDataLen() (sum int)
	setDataLen(sum int)
}

func createNewMessagePacket(type_flag byte) (msg MessagePacket, err error) {
	switch uint8(type_flag) {
	case TYPE_FLAG_CONNECT:
		return NewConnect(), nil
	case TYPE_FLAG_PUBACK:
		return NewPuback(), nil
	case TYPE_FLAG_PUBCOMP:
		return NewPubcomp(), nil
	case TYPE_FLAG_PUBREC:
		return NewPubrec(), nil
	case TYPE_FLAG_PUBREL:
		return NewPubrel(), nil
	case TYPE_FLAG_SUBSCRIBE:
		return NewSubscribe(), nil
	case TYPE_FLAG_UNSUBSCRIBE:
		return NewUnSubscribe(), nil
	case TYPE_FLAG_PINGREQ:
		return pingreq, nil
	case TYPE_FLAG_DISCONNECT:
		return disconnect, nil
	case TYPE_FLAG_CONNACK:
		return NewConnbak(), nil
	case TYPE_FLAG_PINGRESP:
		return PingrespObj, nil
	case TYPE_FLAG_SUBACK:
		return NewSuback(), nil
	case TYPE_FLAG_UNSUBACK:
		return NewUnSuback(), nil
	default:
		//判断是否是PUBLISSH报文
		if (type_flag >> 4) != TYPE_PUBLISH {

			return nil, fmt.Errorf("读取数据失败，无匹配类型报文放弃该连接:%X\n", type_flag)
		}
		pub := NewPublish()
		return &pub, nil
	}
}
