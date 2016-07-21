package packet

/**


作者：ybq
时间:2015-09-08
**/
import (
	"fmt"
)

// Pubrel mqtt　PUBREL类型报文
type Pubrel struct {
	//固定头
	FixedHeader

	//可变头
	packetId uint16 //报文标识符

}

// NewPubrel 构造函数
func NewPubrel() *Pubrel {
	msg := &Pubrel{}
	msg.hdata = byte(TYPE_FLAG_PUBREL)
	msg.remlen = 2
	return msg
}

// SetPacketId 设置报文标识符
func (c *Pubrel) SetPacketId(packetId uint16) {
	c.packetId = packetId
}

// GetPacketId 获取报文标识符
func (c *Pubrel) GetPacketId() uint16 {
	return c.packetId
}

func (c *Pubrel) String() string {
	return fmt.Sprintf("PUBREL会话标识:%d ", c.packetId)
}

// UnPacket msg不包含固定报头
func (c *Pubrel) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_PUBREL {
		return fmt.Errorf("PUBREL固定头信息保留位错误:%v", header)
	}
	c.remlen = 2
	c.packetId, _ = BytesRUint16(msg, 0)
	return nil
}
func (c *Pubrel) Packet() []byte {
	//固定报头
	data := make([]byte, 4, 4)
	data[0] = c.hdata
	data[1] = byte(0x2)
	BytesWUint16(data, c.packetId, 2)
	c.totalen = 4
	c.remlen = 2
	return data
}
