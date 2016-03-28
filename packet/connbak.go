/**


作者：ybq
时间:2015-09-08
**/
package packet

import (
	"fmt"
)

//mqtt　CONNBAK类型报文
type Connbak struct {
	//固定头
	FixedHeader

	//可变头
	//当前会话标志
	//位置：连接确认标志的第 0 位。
	//		如果服务端收到清理会话（CleanSession）标志为 1 的连接，
	//除了将 CONNACK 报文中的返回码设置为 0之外，
	//还必须将 CONNACK 报文中的当前会话设置（Session Present）标志为 0  [MQTT-3.2.2-1]。
	//
	// 	 如果服务端收到一个 CleanSession 为 0 的连接，
	//当前会话标志的值取决于服务端是否已经保存了 ClientId对应客户端的会话状态。
	//	如果服务端已经保存了会话状态，它必须将 CONNACK 报文中的当前会话标志设置为 1  [MQTT-3.2.2-2]。
	//	如果服务端没有已保存的会话状态，它必须将 CONNACK 报文中的当前会话设置为 0。
	//还需要将 CONNACK 报文中的返回码设置为 0  [MQTT-3.2.2-3]。
	//
	// 		当前会话标志使服务端和客户端在是否有已存储的会话状态上保持一致。
	// 		一旦完成了会话的初始化设置，已经保存会话状态的客户端将期望服务端维持它存储的会话状态。
	//如果客户端从服务端收到的当前的值与预期的不同，客户端可以选择继续这个会话或者断开连接。
	//客户端可以丢弃客户端和服务端之间的会话状态，方法是，断开连接，将清理会话标志设置为 1，再次连接，然后再次断开连接。
	//  如果服务端发送了一个包含非零返回码的 CONNACK 报文，它必须将当前会话标志设置为 0  [MQTT-3.2.2-4]
	sessionPresent bool //当前会话标志

	//返回码
	//连接返回码字段使用一个字节的无符号值，在 表格 3.1 –连接返回码的值 中列出。
	//如果服务端收到一个合法的 CONNECT 报文，但出于某些原因无法处理它，
	//服务端应该尝试发送一个包含非零返回码（表格中的某一个）的 CONNACK 报文。
	//如果服务端发送了一个包含非零返回码的 CONNACK 报文，那么它必须关闭网络连接
	returnCode uint8 //返回码

}

//构造函数
func NewConnbak() *Connbak {
	msg := &Connbak{}
	msg.hdata = byte(TYPE_FLAG_CONNACK)
	msg.remlen = 2
	return msg
}

//设置当前会话标志
func (c *Connbak) SetSessionPresent(sessionPresent bool) *Connbak {
	c.sessionPresent = sessionPresent
	return c
}

//获取当前会话标志
func (c *Connbak) GetSessionPresent() bool {
	return c.sessionPresent
}

//设置返回码
//返回码常量：
//		CONNBAK_RETURN_CODE_OK             = iota //连接已接受 连接已被服务端接受
//		CONNBAK_RETURN_NO_SUPPORT_PROTOCOL        //连接已拒绝，不支持的协议版本 服务端不支持客户端请求的 MQTT 协议级别
//		CONNBAK_RETURN_NO_CLIENT_ID               //连接已拒绝，不合格的客户端标识符 客户端标识符是正确的 UTF-8 编码，但服务端不允许使用
//		CONNBAK_RETURN_NO_SERVER                  //连接已拒绝，服务端不可用 网络连接已建立，但 MQTT 服务不可用
//		CONNBAK_RETURN_ERROR_UNAME_PWD            //连接已拒绝，无效的用户名或密码 用户名或密码的数据格式无效
//		CONNBAK_RETURN_UNAUTHORIZED 				  // 连接已拒绝，未授权 客户端未被授权连接到此服务器
func (c *Connbak) SetReturnCode(returnCode uint8) *Connbak {
	c.returnCode = returnCode
	return c
}

//获取返回码
//返回码常量：
//		CONNBAK_RETURN_CODE_OK             = iota //连接已接受 连接已被服务端接受
//		CONNBAK_RETURN_NO_SUPPORT_PROTOCOL        //连接已拒绝，不支持的协议版本 服务端不支持客户端请求的 MQTT 协议级别
//		CONNBAK_RETURN_NO_CLIENT_ID               //连接已拒绝，不合格的客户端标识符 客户端标识符是正确的 UTF-8 编码，但服务端不允许使用
//		CONNBAK_RETURN_NO_SERVER                  //连接已拒绝，服务端不可用 网络连接已建立，但 MQTT 服务不可用
//		CONNBAK_RETURN_ERROR_UNAME_PWD            //连接已拒绝，无效的用户名或密码 用户名或密码的数据格式无效
//		CONNBAK_RETURN_UNAUTHORIZED 				  // 连接已拒绝，未授权 客户端未被授权连接到此服务器
func (c *Connbak) GetReturnCode() uint8 {
	return c.returnCode
}

func (c *Connbak) String() string {
	return fmt.Sprintf("CONNACK当前会话标志:%t 返回码:%d ", c.sessionPresent, c.returnCode)
}

//msg不包含固定报头
func (c *Connbak) UnPacket(header byte, msg []byte) error {

	if header != TYPE_FLAG_CONNACK {
		return fmt.Errorf("CONNACK固定头信息保留位错误:%v", header)
	}

	c.remlen = 2
	if (msg[0] | 1) != 1 {
		return fmt.Errorf("CONNACK连接确认标志错误:收到%d 应收到0或1", msg[0])
	}
	c.sessionPresent = msg[0] == 1
	c.returnCode = uint8(msg[1])
	return nil
}
func (c *Connbak) Packet() []byte {
	//固定报头
	data := make([]byte, 4, 4)
	data[0] = c.hdata
	data[1] = byte(0x2)
	if c.sessionPresent {
		data[2] = 0x1
	}
	data[3] = c.returnCode
	c.totalen = 4
	c.remlen = 2
	return data
}
