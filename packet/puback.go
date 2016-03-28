/**


作者：ybq
时间:2015-09-08
**/
package packet

import (
	"fmt"
)

//mqtt　PUBACK类型报文
type Puback struct {
	//固定头
	FixedHeader

	//可变头

	packetId uint16 //报文标识符

}

//构造函数
func NewPuback() *Puback {
	msg := &Puback{}
	msg.hdata = byte(TYPE_FLAG_PUBACK)
	msg.remlen = 2
	return msg
}

//设置报文标识符
func (c *Puback) SetPacketId(packetId uint16) *Puback {
	c.packetId = packetId
	return c
}

//获取报文标识符
func (c *Puback) GetPacketId() uint16 {
	return c.packetId
}

func (c *Puback) String() string {
	return fmt.Sprintf("PUBACK会话标识:%d ", c.packetId)
}

//msg不包含固定报头
func (c *Puback) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_PUBACK {
		return fmt.Errorf("PUBACK固定头信息保留位错误%d 需要%d", header, TYPE_FLAG_PUBACK)
	}
	c.remlen = 2
	c.packetId = (uint16(msg[0]) << 8) | uint16(msg[1])
	return nil
}
func (c *Puback) Packet() []byte {
	//固定报头
	data := make([]byte, 4, 4)
	data[0] = c.hdata
	data[1] = byte(0x2)
	data[2] = byte(c.packetId >> 8)
	data[3] = byte(c.packetId)
	c.totalen = 4
	c.remlen = 2
	return data
}
