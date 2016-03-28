/**


作者：ybq
时间:2015-09-08
**/
package packet

import (
	"fmt"
)

//mqtt　PUBREC类型报文
type Pubrec struct {
	//固定头
	FixedHeader

	//可变头

	packetId uint16 //报文标识符

}

//构造函数
func NewPubrec() *Pubrec {
	msg := &Pubrec{}
	msg.hdata = byte(TYPE_FLAG_PUBREC)
	msg.remlen = 2
	return msg
}

//设置报文标识符
func (c *Pubrec) SetPacketId(packetId uint16) *Pubrec {
	c.packetId = packetId
	return c
}

//获取报文标识符
func (c *Pubrec) GetPacketId() uint16 {
	return c.packetId
}

func (c *Pubrec) String() string {
	return fmt.Sprintf("PUBREC会话标识:%d ", c.packetId)
}

//msg不包含固定报头
func (c *Pubrec) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_PUBREC {
		return fmt.Errorf("PUBREC固定头信息保留位错误:%v", header)
	}
	c.remlen = 2
	c.packetId, _ = BytesRUint16(msg, 0)
	return nil
}
func (c *Pubrec) Packet() []byte {
	//固定报头
	data := make([]byte, 4, 4)
	data[0] = c.hdata
	data[1] = byte(0x2)
	BytesWUint16(data, c.packetId, 2)
	c.totalen = 4
	c.remlen = 2
	return data
}
