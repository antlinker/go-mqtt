/**


作者：ybq
时间:2015-09-09
**/
package packet

import (
	"fmt"
	"strconv"
)

//mqtt　PUBREL类型报文
type Suback struct {
	//固定头
	FixedHeader

	//可变头

	packetId uint16  //报文标识符
	rcodes   []uint8 //返回码列表
}

//构造函数
func NewSuback(args ...int) *Suback {
	msg := &Suback{}
	msg.hdata = TYPE_FLAG_SUBACK
	size := 4
	if len(args) == 1 {
		size = args[0]
	}
	msg.rcodes = make([]uint8, 0, size)
	return msg
}

//设置报文标识符
func (c *Suback) SetPacketId(packetId uint16) *Suback {
	c.packetId = packetId
	return c
}

//获取报文标识符
func (c *Suback) GetPacketId() uint16 {
	return c.packetId
}

//添加返回码
//返回码常量定义
//SUBACK_RETURNCODE_QOS_0 uint8 = 0x0 //成功 最大 QoS 0
//SUBACK_RETURNCODE_QOS_1 uint8 = 0x1 //成功 – 最大 QoS 1
//	SUBACK_RETURNCODE_QOS_2 uint8 = 0x2 //成功 – 最大 QoS 2
//	SUBACK_RETURNCODE_FAILURE uint8=0x80 //Failure  失败
func (c *Suback) AddReturnCode(rcodes ...uint8) *Suback {
	for i := range rcodes {
		c.rcodes = append(c.rcodes, rcodes[i])
	}
	return c
}

//替换指定索引的返回码 自行保证索引不越界
func (c *Suback) AddReturnCodeAt(returncode uint8, index int) *Suback {
	c.rcodes[index] = returncode
	return c
}

//获取所有的返回码
func (c *Suback) GetReturnCodes() []uint8 {
	return c.rcodes
}

//获取指定位置的返回码
func (c *Suback) GetReturnCodeAt(index int) uint8 {
	return c.rcodes[index]
}
func (c *Suback) String() string {
	outs := "SUBACK会话标识:" + strconv.Itoa(int(c.packetId)) + "("
	for i := range c.rcodes {
		outs += strconv.Itoa(int(c.rcodes[i])) + " "
	}
	outs += ")"
	return outs
}

//计算剩余长度
func (c *Suback) totalRemain() int {
	return 2 + len(c.rcodes)
}

//msg不包含固定报头
func (c *Suback) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_SUBACK {
		return fmt.Errorf("SUBACK报文类型或报文控制标志错误:现为%X ,应该为%X", header, TYPE_FLAG_SUBACK)
	}
	c.remlen = len(msg)

	maxlen := c.remlen
	if maxlen == 2 {
		return fmt.Errorf("SUBACK报文有效载荷错误，返回码不能为空")
	}
	curindex := 0
	c.packetId, curindex = BytesRUint16(msg, curindex)
	var fs uint8
	for curindex < maxlen {
		fs, curindex = BytesRUint8(msg, curindex)
		c.AddReturnCode(fs)
	}
	return nil
}
func (c *Suback) Packet() []byte {
	//固定报头
	c.remlen = 2 + len(c.rcodes)
	//fmt.Printf("剩余字节%d %X\n", c.remlen, c.remlen)
	remlenbyte := Remlen2Bytes(int32(c.remlen))
	//fmt.Printf("剩余字节 %X\n", remlenbyte)
	data := make([]byte, c.remlen+1+len(remlenbyte))
	curind := 0
	curind = BytesWByte(data, c.hdata, curind)
	curind = BytesWBytes(data, remlenbyte, curind)
	curind = BytesWUint16(data, c.packetId, curind)
	fs := c.rcodes
	for i := range fs {
		curind = BytesWUint8(data, byte(fs[i]), curind)
	}
	c.totalen = len(data)
	return data
}
