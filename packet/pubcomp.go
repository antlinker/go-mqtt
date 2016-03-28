/**


作者：ybq
时间:2015-09-08
**/
package packet

import (
	"fmt"
)

//mqtt　PUBCOMP类型报文
type Pubcomp struct {
	//固定头
	FixedHeader

	//可变头

	packetId uint16 //报文标识符

}

//构造函数
func NewPubcomp() *Pubcomp {
	msg := &Pubcomp{}
	msg.hdata = byte(TYPE_FLAG_PUBCOMP)
	msg.remlen = 2
	return msg
}

//设置报文标识符
func (c *Pubcomp) SetPacketId(packetId uint16) *Pubcomp {
	c.packetId = packetId
	return c
}

//获取报文标识符
func (c *Pubcomp) GetPacketId() uint16 {
	return c.packetId
}

func (c *Pubcomp) String() string {
	return fmt.Sprintf("PUBCOMP会话标识:%d ", c.packetId)
}

//msg不包含固定报头
func (c *Pubcomp) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_PUBCOMP {
		return fmt.Errorf("PUBCOMP固定头信息保留位错误:%X", header)
	}
	c.remlen = 2
	c.packetId, _ = BytesRUint16(msg, 0)
	return nil
}
func (c *Pubcomp) Packet() []byte {
	//固定报头
	data := make([]byte, 4, 4)
	data[0] = c.hdata
	data[1] = byte(0x2)
	BytesWUint16(data, c.packetId, 2)
	c.totalen = 4
	c.remlen = 2
	return data
}
