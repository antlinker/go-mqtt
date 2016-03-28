/**


作者：ybq
时间:2015-09-09
**/
package packet

import (
	"fmt"
)

//mqtt　UNSUBACK类型报文
type UnSuback struct {
	//固定头
	FixedHeader

	//可变头

	packetId uint16 //报文标识符

}

//构造函数
func NewUnSuback() *UnSuback {
	msg := &UnSuback{}
	msg.hdata = byte(TYPE_FLAG_UNSUBACK)
	msg.remlen = 2
	return msg
}

//设置报文标识符
func (c *UnSuback) SetPacketId(packetId uint16) *UnSuback {
	c.packetId = packetId
	return c
}

//获取报文标识符
func (c *UnSuback) GetPacketId() uint16 {
	return c.packetId
}

func (c *UnSuback) String() string {
	return fmt.Sprintf("UNSUBACK会话标识:%d ", c.packetId)
}

//msg不包含固定报头
func (c *UnSuback) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_UNSUBACK {
		return fmt.Errorf("UNSUBACK固定头信息保留位错误:%v", header)
	}
	c.remlen = 2
	c.packetId, _ = BytesRUint16(msg, 0)
	return nil
}
func (c *UnSuback) Packet() []byte {
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
