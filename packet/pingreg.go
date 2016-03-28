/**


作者：ybq
时间:2015-09-09
**/
package packet

import (
	"fmt"
)

//mqtt　UNSUBACK类型报文
type Pingreq struct {
	//固定头
	FixedHeader
}

//构造函数
func NewPingreq() *Pingreq {
	msg := &Pingreq{}
	msg.hdata = byte(TYPE_FLAG_PINGREQ)
	return msg
}

func (c *Pingreq) String() string {
	return fmt.Sprintf("Pingreq")
}

//msg不包含固定报头
func (c *Pingreq) UnPacket(header byte, msg []byte) error {
	c.remlen = len(msg)
	if header != TYPE_FLAG_PINGREQ {
		return fmt.Errorf("PINGREQ固定头信息保留位错误:%v", header)
	}
	return nil
}
func (c *Pingreq) Packet() []byte {
	//固定报头
	data := make([]byte, 2, 2)
	data[0] = c.hdata
	data[1] = byte(0x0)
	c.totalen = 2
	c.remlen = 0
	return data
}
