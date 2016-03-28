/**
作者：ybq
时间:2015-09-09
**/
package packet

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/antlinker/go-mqtt/util"
)

//取消订阅主题过滤器列表
type UnsTopicFilterList struct {
	filters []string
}

//增加主题过滤器
func (t *UnsTopicFilterList) AddFilter(filters ...string) *UnsTopicFilterList {
	for i := range filters {
		t.filters = append(t.filters, filters[i])
	}
	return t
}
func (t *UnsTopicFilterList) GetFilters() []string {
	return t.filters
}
func (t *UnsTopicFilterList) GetFilter(ind int) (string, error) {
	if ind >= len(t.filters) {
		return "", fmt.Errorf("参数越界，最大为%d 当前为%d", len(t.filters), ind)
	}
	return t.filters[ind], nil
}

/**
mqtt　UNSUBSCRIBE类型报文

作者：ybq
时间:2015-09-09
**/
type UnSubscribe struct {
	//固定头
	FixedHeader
	//可变头
	packetId uint16 //报文标识符
	UnsTopicFilterList
}

//构造函数
func NewUnSubscribe(args ...int) *UnSubscribe {

	msg := &UnSubscribe{}
	msg.hdata = TYPE_FLAG_UNSUBSCRIBE
	size := 1
	if len(args) == 1 {
		size = args[0]
	}
	msg.filters = make([]string, 0, size)

	return msg
}

//设置报文标识符
func (c *UnSubscribe) SetPacketId(packetId uint16) {
	c.packetId = packetId
}

//获取报文标识符
func (c *UnSubscribe) GetPacketId() uint16 {
	return c.packetId
}

//计算剩余长度
func (c *UnSubscribe) totalRemain() int {
	sum := 2
	var fs = c.filters
	for i := range fs {
		sum += len([]byte(fs[i])) + 2
	}
	return sum
}

func (c *UnSubscribe) String() string {
	fstr := "UNSUBSCRIBE报文标识:" + strconv.Itoa(int(c.packetId)) + "\n"
	var fs = c.filters
	for i := range fs {
		fstr += "主题过滤器" + strconv.Itoa(i) + ":" + fs[i] + "\n"
	}

	return fstr
}

//解包
func (c *UnSubscribe) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_UNSUBSCRIBE {
		return fmt.Errorf("UNSUBSCRIBE报文类型或报文控制标志错误:现为%X ,应该为%X", header, TYPE_FLAG_UNSUBSCRIBE)
	}
	c.remlen = len(msg)
	maxlen := c.remlen
	if maxlen == 2 {
		return fmt.Errorf("UNSUBSCRIBE报文有效载荷错误，取消订阅主题为空")
	}
	curindex := 0
	c.packetId, curindex = BytesRUint16(msg, curindex)
	c.filters = make([]string, 0, 1)
	var fs string
	for curindex < maxlen {
		fs, curindex = BytesRString(msg, curindex)
		if !util.VerifyZeroByMqtt([]rune(fs)) {
			return errors.New("主题过滤器包含U+0000字符关闭连接")
		}
		c.AddFilter(fs)
	}
	return nil
}

//封包
func (c *UnSubscribe) Packet() []byte {
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
		curind = BytesWBString(data, []byte(fs[i]), curind)
	}
	c.totalen = len(data)
	return data
}
