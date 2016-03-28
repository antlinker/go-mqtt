/**
作者：ybq
时间:2015-09-08
**/
package packet

import (
	"errors"
	"fmt"

	"github.com/antlinker/go-mqtt/util"

	"strconv"
)

//主题过滤器
type TopicFilter struct {
	filter *Topic
	qos    QoS
}

func NewTopicFilter(filter *Topic, qos QoS) *TopicFilter {
	return &TopicFilter{filter, qos}
}

//获取主题过滤器
func (t *TopicFilter) GetFilter() *Topic {
	return t.filter
}
func (t *TopicFilter) GetQos() QoS {
	return t.qos
}
func (t *TopicFilter) SetQos(qos QoS) {
	t.qos = qos
}
func (t *TopicFilter) SetFilter(filter *Topic) {
	t.filter = filter
}

//是否时有效主题过滤器
func (t *TopicFilter) IsValidTopicFilter() bool {

	//	pos := strings.Index(tmpv, "#")
	//	if pos >= 0 && pos+1 != len(tmpv) {
	//		//fmt.Println(t.filter.value, pos, len(t.filter.value))
	//		return false
	//	}
	if t.filter.value == "$#" {
		return true
	}
	//	if pos > 0 {

	//		if tmpv[pos-1] != '/' {
	//			//fmt.Println(t.filter.value, pos, len(t.filter.value), t.filter.value[pos-1])
	//			return false
	//		}
	//	}
	tmpv := []rune(t.filter.value)
	for i, v := range tmpv {
		if v == '+' {
			if i+1 < len(tmpv) && tmpv[i+1] != '/' {
				return false

			}
			if i > 0 && tmpv[i-1] != '/' {
				return false
			}
		} else if v == '#' {
			if i+1 < len(tmpv) {
				return false

			}
			if i > 0 && tmpv[i-1] != '/' {
				return false
			}
		} else if v == 0 {
			return false
		}
	}
	return true
}

//主题过滤器列表
type TopicFilterList struct {
	filters []*TopicFilter
}

//增加主题过滤器
func (t *TopicFilterList) AddFilter(filter *Topic, qos QoS) *TopicFilterList {
	t.filters = append(t.filters, &TopicFilter{filter, qos})
	return t
}
func (t *TopicFilterList) GetFilters() []*TopicFilter {
	return t.filters
}
func (t *TopicFilterList) GetFilter(ind int) (*TopicFilter, error) {
	if ind >= len(t.filters) {
		return nil, fmt.Errorf("参数越界，最大为%d 当前为%d", len(t.filters), ind)
	}
	return t.filters[ind], nil
}

/**
mqtt　SUBSCRIBE类型报文

作者：ybq
时间:2015-09-08
**/
type Subscribe struct {
	//固定头
	FixedHeader
	//可变头
	packetId uint16 //报文标识符
	TopicFilterList
}

//构造函数
func NewSubscribe(args ...int) *Subscribe {

	msg := &Subscribe{}
	size := 1
	if len(args) == 1 {
		size = args[0]
	}
	msg.hdata = TYPE_FLAG_SUBSCRIBE
	msg.filters = make([]*TopicFilter, 0, size)

	return msg
}

//设置报文标识符
func (c *Subscribe) SetPacketId(packetId uint16) {
	c.packetId = packetId
}

//获取报文标识符
func (c *Subscribe) GetPacketId() uint16 {
	return c.packetId
}

//计算剩余长度
func (c *Subscribe) totalRemain() int {
	sum := 2
	var fs = c.filters
	for i := range fs {
		sum += fs[i].filter.len + 3
	}
	return sum
}

func (c *Subscribe) String() string {
	fstr := "SUBSCRIBE报文标识:" + strconv.Itoa(int(c.packetId)) + "\n"
	var fs = c.filters
	for i := range fs {
		fstr += "主题过滤器" + strconv.Itoa(i) + "(qos质量" + strconv.Itoa(int(fs[i].qos)) + "):" + fs[i].filter.String() + "\n"
	}

	return fstr
}

//解包
func (c *Subscribe) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_SUBSCRIBE {
		return fmt.Errorf("SUBSCRIBE报文类型或报文控制标志错误:现为%X ,应该为%X", header, TYPE_FLAG_SUBSCRIBE)
	}
	c.remlen = len(msg)

	maxlen := c.remlen
	if maxlen == 2 {
		return fmt.Errorf("SUBSCRIBE报文有效载荷错误，订阅主题为空")
	}
	curindex := 0
	c.packetId, curindex = BytesRUint16(msg, curindex)
	var fs string
	var qos uint8
	c.filters = make([]*TopicFilter, 0, 1)
	for curindex < maxlen {
		fs, curindex = BytesRString(msg, curindex)
		qos, curindex = BytesRUint8(msg, curindex)

		if !util.VerifyZeroByMqtt([]rune(fs)) {
			return errors.New("主题过滤器包含U+0000字符关闭连接")
		}

		c.AddFilter(NewTopic(fs), QoS(qos))
	}
	return nil
}

//封包
func (c *Subscribe) Packet() []byte {
	//固定报头
	c.remlen = c.totalRemain()
	//fmt.Printf("剩余字节%d %X\n", c.remlen, c.remlen)
	remlenbyte := Remlen2Bytes(int32(c.remlen))
	//fmt.Printf("剩余字节 %X\n", remlenbyte)
	data := make([]byte, c.remlen+1+len(remlenbyte))
	curind := 0
	curind = BytesWByte(data, c.hdata, curind)
	curind = BytesWBytes(data, remlenbyte, curind)
	curind = BytesWUint16(data, c.packetId, curind)
	fs := c.filters
	for i := range fs {

		curind = BytesWBString(data, fs[i].filter.BytesValue(), curind)
		curind = BytesWUint8(data, uint8(fs[i].qos), curind)
	}
	c.totalen = len(data)
	return data
}
