/**


作者：ybq
时间:2015-09-09
**/
package packet

import (
	"fmt"
)

//mqtt　DISCONNECT类型报文
type Disconnect struct {
	//固定头
	FixedHeader
}

//构造函数
func NewDisconnect() *Disconnect {
	msg := &Disconnect{}
	msg.hdata = byte(TYPE_FLAG_DISCONNECT)
	return msg
}

func (c *Disconnect) String() string {
	return fmt.Sprintf("DISCONNECT")
}

//msg不包含固定报头
func (c *Disconnect) UnPacket(header byte, msg []byte) error {
	c.remlen = len(msg)
	if header != TYPE_FLAG_DISCONNECT {
		return fmt.Errorf("DISCONNECT固定头信息保留位错误:%v", header)
	}
	return nil
}
func (c *Disconnect) Packet() []byte {
	//固定报头
	data := make([]byte, 2, 2)
	data[0] = c.hdata
	data[1] = byte(0x0)
	c.totalen = 2
	c.remlen = 0
	return data
}
