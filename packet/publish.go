/**
作者：ybq
时间:2015-09-08
**/
package packet

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/antlinker/go-mqtt/util"
)

/**
mqtt　Publish类型报文

作者：ybq
时间:2015-09-08
**/
type Publish struct {
	//固定头
	FixedHeader

	//可变头
	topic []byte //主题名

	packetId uint16 //报文标识符

	cache      bool      //司否启用缓存 不启用缓存使用packetdata 使用缓存 使用datain
	payload    []byte    //有效载荷
	payloadin  io.Reader //数据源
	datainsize int       //数据大小

	ReciveStartTime time.Time
}

//构造函数
func NewPublish() Publish {
	msg := Publish{}
	msg.hdata = byte(TYPE_PUBLISH << 4)
	return msg
}

func NewPublishAll(topic string, payload []byte, qos QoS, retain bool) *Publish {
	msg := &Publish{topic: []byte(topic), payload: payload}
	if retain {
		msg.hdata = (TYPE_PUBLISH << 4) | (uint8(qos) << 1) | 1
	} else {
		msg.hdata = (TYPE_PUBLISH << 4) | (uint8(qos) << 1)
	}

	return msg
}

//构造一个保留消息
func NewPublishRetain(topic string, payload []byte, qos QoS) *Publish {
	msg := &Publish{topic: []byte(topic), payload: payload}
	msg.hdata = (TYPE_PUBLISH << 4) | (uint8(qos) << 1) | 1
	return msg
}
func (p *Publish) CloneByQoS(qos QoS) *Publish {
	msg := *p
	msg.SetQos(qos)
	return &msg
}

//设置控制标志
//参数
//	dup 重发标志
//	qos 服务质量等级 常量 QOS_0 QOS_1 QOS_2
//	retain 保留标志
func (c *Publish) SetControlFlag(dup bool, qos QoS, retain bool) {
	c.hdata = (c.hdata & 0xf0)
	if dup {
		c.hdata |= 8
	}
	c.hdata = c.hdata | (uint8(qos) << 1)
	if retain {
		c.hdata |= 1
	}
	//if c.ispacket {
	//c.ispacket = false
	//}
}

//设置重发标志
//true　表示这可能是一个早前报文请求的重发
//false 表示这是客户端或服务端第一次请求发送这个 PUBLISH 报文
func (c *Publish) SetDupFlag(flag bool) {
	c.hdata = (c.hdata & 0xf7)
	if flag {
		c.hdata |= 8
	}
	// if c.ispacket {
	// 	c.packetbuff[0] = c.hdata
	// }
	//c.ispacket = false
}

//获取重发标志
func (c *Publish) GetDupFlag() bool {
	return c.hdata&0x8 == 0x8
}

//设置Qos 服务质量等级
func (c *Publish) SetQos(qos QoS) {
	c.hdata = (c.hdata & 0xf9) | (uint8(qos) << 1)
	// if c.ispacket {
	// 	c.packetbuff[0] = c.hdata
	// }
	//c.ispacket = false
}

//获取Qos 服务质量等级
func (c *Publish) GetQos() QoS {
	return QoS(c.hdata >> 1 & 0x3)
}

//设置保留标志
//true　表示这可能是一个早前报文请求的重发
//false 表示这是客户端或服务端第一次请求发送这个 PUBLISH 报文
func (c *Publish) SetRetain(flag bool) {
	c.hdata = (c.hdata & 0xfe)
	if flag {
		c.hdata |= 1
	}
	// if c.ispacket {
	// 	c.packetbuff[0] = c.hdata
	// }
	//c.ispacket = false
}

//获取保留标志
func (c *Publish) GetRetain() bool {
	return c.hdata&0x1 == 0x1
}

//设置主题名
func (c *Publish) SetTopic(topic []byte) {
	c.topic = topic
	//c.ispacket = false
}

//设置主题名
func (c *Publish) SetTopicByString(topic string) {
	c.topic = []byte(topic)
	//c.ispacket = false
}

//获取主题名
func (c *Publish) GetTopic() []byte {
	return c.topic
}

//获取主题名
func (c *Publish) GetTopicByString() string {
	return string(c.topic)
}

//设置报文标识符
func (c *Publish) SetPacketId(packetId uint16) {
	c.packetId = packetId
	//c.ispacket = false
}

//获取报文标识符
func (c *Publish) GetPacketId() uint16 {
	return c.packetId
}

//设置有效载荷
func (c *Publish) SetPayload(payload []byte) {
	c.payload = payload
	//c.ispacket = false
}

//获取有效载荷
func (c *Publish) GetPayloadIn() io.Reader {
	return c.payloadin
}

//设置有效载荷
func (c *Publish) SetPayloadIn(in io.Reader, size int) {

	c.payloadin = in
	c.datainsize = size
	//c.ispacket = false
}

//获取有效载荷
func (c *Publish) GetPayload() []byte {
	return c.payload
}

//计算剩余长度
func (c *Publish) totalRemain() int {
	if c.GetQos() == 0 {
		return len(c.topic) + 2 + len(c.payload)
	}
	return len(c.topic) + 4 + len(c.payload)
}
func formatPayload(payload []byte) string {
	n := len(payload)
	if n > 128 {
		n = 128
	}
	return fmt.Sprintf("(%d)%X", len(payload), payload[0:n])
}
func (c *Publish) String() string {
	return fmt.Sprintf("PUBLISH重发标志:%t 服务质量等级(Qos):%d 保留标志:%t 主题名:%s 报文标志:%d  接收时间:%v有效载荷:%s", c.GetDupFlag(), c.GetQos(), c.GetRetain(), c.topic, c.packetId, c.ReciveStartTime, formatPayload(c.payload))
}

//解包
func (c *Publish) UnPacket(header byte, msg []byte) error {
	c.cache = false

	c.remlen = len(msg)
	// c.SetRetain(header&0x1 == 0x1)
	// c.SetQos(QoS(header >> 1 & 0x3))
	// c.SetDupFlag(header&0x8 == 0x8)
	c.hdata = header
	curindex := 0
	c.topic, curindex = BytesRBString(msg, curindex)

	if c.GetQos() > QOS_0 {
		c.packetId, curindex = BytesRUint16(msg, curindex)
	}
	if curindex < len(msg) {
		c.payload = BytesRBytesEnd(msg, curindex)
	}

	if !util.VerifyZeroByMqtt([]rune(c.GetTopicByString())) {
		return errors.New("主题包含U+0000字符关闭连接")
	}
	//c.setBuffByRemdata(msg)
	return nil
}

//封包
func (c *Publish) Packet() []byte {
	// if c.ispacket {
	// 	return c.packetbuff.Bytes()
	// }
	//固定报头
	c.remlen = c.totalRemain()
	//fmt.Printf("剩余字节%d %X\n", c.remlen, c.remlen)
	remlenbyte := Remlen2Bytes(int32(c.remlen))
	//fmt.Printf("剩余字节 %X\n", remlenbyte)
	data := make([]byte, c.remlen+1+len(remlenbyte))
	curind := 0
	curind = BytesWByte(data, c.hdata, curind)
	curind = BytesWBytes(data, remlenbyte, curind)
	curind = BytesWBString(data, c.topic, curind)
	if c.GetQos() > QOS_0 {
		curind = BytesWUint16(data, c.packetId, curind)
	}
	if len(c.payload) > 0 {
		curind = BytesWBytes(data, c.payload, curind)
	}
	c.totalen = len(data)
	//c.packetbuff.Reset()
	//c.packetbuff.Write(data)
	//c.packetbuff = data
	//c.ispacket = true
	return data
}

//将信息发送到一个Writer中
func (c *Publish) WriteTo(out io.Writer) (int64, error) {
	num, err := out.Write(c.Packet())
	return int64(num), err
}

//从文件解码 主要为PUBLIST类型报文使用，大数据传输缓存中读取数据其他报文不需要重载
//还未实现
func (c *Publish) UnPacketFile(header byte, varheader []byte, in io.Reader) error {

	return nil
}

//打包到，大数据传输缓存中写入数据 除PUBLIST外其他报文不需要重载
func (c *Publish) PacketTo(out io.Writer, size int) error {
	c.remlen = len(c.topic) + 4 + c.datainsize
	//fmt.Printf("剩余字节%d %X\n", c.remlen, c.remlen)
	remlenbyte := Remlen2Bytes(int32(c.remlen))
	//fmt.Printf("剩余字节 %X\n", remlenbyte)
	data := make([]byte, 1+len(remlenbyte)+len(c.topic)+4)
	curind := 0
	curind = BytesWByte(data, c.hdata, curind)
	curind = BytesWBytes(data, remlenbyte, curind)
	curind = BytesWBString(data, c.topic, curind)
	curind = BytesWUint16(data, c.packetId, curind)
	out.Write(data)
	buff := [1024]byte{}
	for i := 0; i < size; {
		cnt, err := c.payloadin.Read(buff[0:])
		if err != nil {
			return err
		}
		_, err = out.Write(buff[0:size])
		if err != nil {
			return err
		}
		i += cnt
	}
	return nil
}
