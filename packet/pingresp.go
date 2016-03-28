/**


作者：ybq
时间:2015-09-09
**/
package packet

import (
	"fmt"
)

//mqtt　PINGRESP类型报文
type Pingresp struct {
	//固定头
	FixedHeader
}

//构造函数
func NewPingresp() *Pingresp {
	msg := &Pingresp{}
	msg.hdata = byte(TYPE_FLAG_PINGRESP)
	return msg
}

func (c *Pingresp) String() string {
	return fmt.Sprintf("Pingresp")
}

//msg不包含固定报头
func (c *Pingresp) UnPacket(header byte, msg []byte) error {
	c.remlen = len(msg)
	if header != TYPE_FLAG_PINGRESP {
		return fmt.Errorf("PINGRESP固定头信息保留位错误:%v", header)
	}

	return nil
}
func (c *Pingresp) Packet() []byte {
	//固定报头
	data := make([]byte, 2, 2)
	data[0] = c.hdata
	data[1] = byte(0x0)
	c.totalen = 2
	c.remlen = 0
	return data
}
