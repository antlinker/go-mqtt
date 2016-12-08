/**
作者：ybq
时间:2015-09-08
**/
package packet

import (
	"errors"
	"fmt"

	"github.com/antlinker/go-mqtt/util"
)

/**
mqtt　Connect类型报文

作者：ybq
时间:2015-09-08
**/
type Connect struct {
	//固定头
	FixedHeader

	//可变头

	protocol     []byte //协议
	level        uint8  //协议级别
	cleanSession bool   //清理会话
	willFlag     bool   //遗嘱标志
	willQoS      QoS    //遗嘱服务质量等级
	willRetain   bool   //遗嘱保留标志
	passFlag     bool   //密码标志
	unameFlag    bool   //用户名
	keepAlive    uint16 //保持时间单位秒

	//有效载荷
	clientId    []byte //客户端标识符
	willTopic   []byte //遗嘱主题
	willMessage []byte //遗嘱消息
	userName    []byte //用户名
	password    []byte //密码
}

//将常量转换为字节数组
var protocol = []byte(PROTOCOL)

//构造函数
func NewConnect() *Connect {
	msg := new(Connect)
	msg.hdata = byte(TYPE_FLAG_CONNECT)
	msg.protocol = protocol
	msg.level = PROTOCOL_LEVEL
	return msg
}

//获取协议
func (c *Connect) GetProtocol() []byte {
	return c.protocol
}

//获取协议级别
func (c *Connect) GetProtocolLevel() uint8 {
	return c.level
}

//设置清理会话标志
//true清理
//false不清理默认false
func (c *Connect) SetCleanSession(flag bool) *Connect {
	c.cleanSession = flag
	return c
}

//获取清理标志
func (c *Connect) GetCleanSession() bool {
	return c.cleanSession
}

////设置遗嘱标志
//func (c *Connect) SetWillFlag(flag bool) *Connect {
//	c.willFlag = flag
//	return c
//}
//获取遗嘱标志
func (c *Connect) GetWillFlag() bool {
	return c.willFlag
}

////设置遗嘱qos
//func (c *Connect) SetWillQoS(flag uint8) *Connect {
//	c.willQoS = flag
//	return c
//}
//获取遗嘱qos
func (c *Connect) GetWillQoS() QoS {
	return c.willQoS
}

////设置遗嘱保留标志
//func (c *Connect) SetWillRetain(flag bool) *Connect {
//	c.willRetain = flag
//	return c
//}
//获取遗嘱保留标志
func (c *Connect) GetWillRetain() bool {
	return c.willRetain
}

////设置密码标志
//func (c *Connect) SetPassFlag(flag bool) *Connect {
//	c.passFlag = flag
//	return c
//}
//获取密码标志
func (c *Connect) GetPassFlag() bool {
	return c.passFlag
}

////设置用户名标志
//func (c *Connect) SetUnameFlag(flag bool) *Connect {
//	c.unameFlag = flag
//	return c
//}
//获取户名标志
func (c *Connect) GetUnameFlag() bool {
	return c.unameFlag
}

//设置连接保持单位秒
func (c *Connect) SetKeepAlive(flag uint16) *Connect {
	c.keepAlive = flag
	return c
}

//获取连接保持单位秒
func (c *Connect) GetKeepAlive() uint16 {
	return c.keepAlive
}

//设置客户端标识符
func (c *Connect) SetClientId(flag []byte) *Connect {
	c.clientId = flag
	return c
}

//设置客户端标识符
func (c *Connect) SetClientIdByString(flag string) *Connect {
	c.clientId = []byte(flag)
	return c
}

//获取客户端标识符
func (c *Connect) GetClientId() []byte {
	return c.clientId
}

//获取客户端标识符
func (c *Connect) GetClientIdByString() string {
	return byte2str(c.clientId)
}

//设置遗嘱主题修改遗嘱标志为true
func (c *Connect) SetWillTopicInfo(willTopic []byte, willMessage []byte, willQoS QoS, willRetain bool) *Connect {
	c.willTopic = willTopic
	c.willMessage = willMessage
	c.willFlag = true
	c.willQoS = willQoS
	c.willRetain = willRetain
	return c
}
func (c *Connect) SetWillTopicInfoByString(willTopic string, willMessage string, willQoS QoS, willRetain bool) *Connect {
	return c.SetWillTopicInfo([]byte(willTopic), []byte(willMessage), willQoS, willRetain)
}

//设置遗嘱主题 不更改遗嘱标志及遗嘱qos遗嘱保留标志
func (c *Connect) SetWillTopic(willTopic []byte) *Connect {
	c.willTopic = willTopic
	return c
}

//设置遗嘱主题 不更改遗嘱标志及遗嘱qos遗嘱保留标志
func (c *Connect) SetWillTopicByString(willTopic string) *Connect {
	c.willTopic = []byte(willTopic)
	return c
}

//设置遗嘱消息 不更改遗嘱标志及遗嘱qos遗嘱保留标志
func (c *Connect) SetWillMessage(willMessage []byte) *Connect {
	c.willMessage = willMessage
	return c
}

//设置遗嘱消息 不更改遗嘱标志及遗嘱qos遗嘱保留标志
func (c *Connect) SetWillMessageByString(willMessage string) *Connect {
	c.willMessage = []byte(willMessage)
	return c
}

//获取遗嘱主题
func (c *Connect) GetWillTopic() []byte {
	return c.willTopic
}

//获取遗嘱主题
func (c *Connect) GetWillTopicByString() string {
	return byte2str(c.willTopic)
}

//获取遗嘱消息
func (c *Connect) GetWillMessage() []byte {
	return c.willMessage
}

//获取遗嘱消息
func (c *Connect) GetWillMessageByString() string {
	return byte2str(c.willMessage)
}

//设置用户名
//同时设置用户标志为true
func (c *Connect) SetUserName(userName []byte) *Connect {
	c.userName = userName
	c.unameFlag = true
	return c
}

//设置用户名
//同时设置用户标志为true
func (c *Connect) SetUserNameByString(userName string) *Connect {
	return c.SetUserName([]byte(userName))
}

//获取用户名
func (c *Connect) GetUserName() []byte {
	return c.userName
}

//获取用户名
func (c *Connect) GetUserNameByString() string {
	return byte2str(c.userName)
}

//设置密码
//同时设置密码为true
func (c *Connect) SetPassword(password []byte) *Connect {
	c.password = password
	c.passFlag = true
	return c
}

//设置密码
//同时设置密码为true
func (c *Connect) SetPasswordByString(password string) *Connect {
	return c.SetPassword([]byte(password))
}

//获取密码
func (c *Connect) GetPassword() []byte {
	return c.password
}

//获取密码
func (c *Connect) GetPasswordByString() string {
	return byte2str(c.password)
}
func (c *Connect) totalRemain() int {
	sum := 0
	sum += len(c.protocol)
	sum += len(c.clientId) + 8
	if c.willFlag {
		sum += len(c.willTopic)
		sum += len(c.willMessage) + 4
	}
	if c.passFlag {
		sum += len(c.password) + 2
	}
	if c.unameFlag {
		sum += len(c.userName) + 2
	}
	return sum
}

func (c *Connect) String() string {
	return fmt.Sprintf("协议:%s 协议级别:%d 清理会话:%t 遗嘱标志:%t 遗嘱Qos:%d 遗嘱保留标志:%t 密码标志:%t 用户名标志: %t 保持连接时间:%d 客户端id:%s 遗嘱主题:%s 遗嘱消息:%s 用户名:%s 密码:%s ", byte2str(c.protocol), c.level, c.cleanSession, c.willFlag, c.willQoS, c.willRetain, c.passFlag, c.unameFlag, c.keepAlive, c.clientId, c.willTopic, c.willMessage, c.userName, c.password)
}

func (c *Connect) UnPacket(header byte, msg []byte) error {
	if header != TYPE_FLAG_CONNECT {
		return fmt.Errorf("CONNECT固定头信息保留位错误:%v", header)
	}
	c.remlen = len(msg)
	curindex := 0
	c.protocol, curindex = BytesRBString(msg, curindex)
	//	if string(c.protocol) != PROTOCOL {
	//		//return fmt.E
	//		return fmt.Errorf("CONNECT协议错误:收到%s 应收到%s", c.protocol, PROTOCOL)
	//	}
	c.level, curindex = BytesRUint8(msg, curindex)
	//	if c.level != PROTOCOL_LEVEL {
	//		return fmt.Errorf("CONNECT协议级别错误错误:收到%d 应收到%d", c.level, PROTOCOL_LEVEL)
	//	}
	var connectflag uint8 = 0
	connectflag, curindex = BytesRUint8(msg, curindex)
	c.cleanSession = ByteRBool(connectflag, 1)
	c.willFlag = ByteRBool(connectflag, 2)
	c.willQoS = QoS((connectflag >> 3) & 0x3)
	c.willRetain = ByteRBool(connectflag, 5)
	c.passFlag = ByteRBool(connectflag, 6)
	c.unameFlag = ByteRBool(connectflag, 7)
	c.keepAlive, curindex = BytesRUint16(msg, curindex)
	//fmt.Printf("+++=%d/%d,%X %X\n", curindex, len(msg), msg[curindex], msg[curindex+1])
	c.clientId, curindex = BytesRBString(msg, curindex)
	if c.willFlag {
		c.willTopic, curindex = BytesRBString(msg, curindex)
		c.willMessage, curindex = BytesRBString(msg, curindex)
		if len(c.willTopic) == 0 {
			return errors.New("遗嘱主题最少一个字节")
		}
	} else {
		c.willTopic = nil
		c.willRetain = false
		c.willMessage = nil
	}
	if c.unameFlag {
		c.userName, curindex = BytesRBString(msg, curindex)
		if c.userName == nil || len(c.userName) == 0 {
			return errors.New("用户名不能为空")
		}
		if !util.VerifyUnicodeByMqtt([]rune(byte2str(c.userName))) {
			return errors.New("用户名中有非法字符")
		}
	} else {
		c.userName = nil
	}
	if c.passFlag {
		c.password, curindex = BytesRBString(msg, curindex)
		if c.password == nil || len(c.password) == 0 {
			return errors.New("密码不能为空")
		}
		if !util.VerifyUnicodeByMqtt([]rune(byte2str(c.password))) {
			return errors.New("密码中有非法字符")
		}
	} else {
		c.password = nil
	}
	if c.remlen != curindex {
		return errors.New("协议错误有剩余字节")
	}
	return nil
}
func (c *Connect) Packet() []byte {
	//固定报头
	c.remlen = c.totalRemain()
	//fmt.Printf("剩余字节%d %X\n", c.remlen, c.remlen)
	remlenbyte := Remlen2Bytes(int32(c.remlen))
	//fmt.Printf("剩余字节 %X\n", remlenbyte)
	data := make([]byte, c.remlen+1+len(remlenbyte))
	curind := 0
	curind = BytesWByte(data, c.hdata, curind)
	curind = BytesWBytes(data, remlenbyte, curind)
	curind = BytesWBString(data, c.protocol, curind)
	curind = BytesWUint8(data, c.level, curind)

	//标志位
	var flag uint8 = 0
	//	cleanSession bool   //清理会话
	if c.cleanSession {
		flag |= 0x2
	}
	//	willFlag     bool   //遗嘱标志
	if c.willFlag {
		flag |= 0x4
		//	willQoS      uint8  //遗嘱服务质量等级
		flag |= uint8(c.willQoS) << 3
		//willRetain   bool   //遗嘱保留标志
		if c.willRetain {
			flag |= 0x20
		}
	}
	//passFlag     bool   //密码标志
	if c.passFlag {
		flag |= 0x40
	}
	//unameFlag    bool   //用户名
	if c.unameFlag {
		flag |= 0x80
	}
	//连接标志
	curind = BytesWUint8(data, flag, curind)
	//keepAlive    uint16 //保持时间单位秒
	curind = BytesWUint16(data, c.keepAlive, curind)
	//有效载荷
	//客户端ｉｄ
	if len(c.clientId) == 0 {
		curind = BytesWUint16(data, 0, curind)
	} else {
		curind = BytesWBString(data, c.clientId, curind)
	}
	//遗嘱
	if c.willFlag {
		curind = BytesWBString(data, c.willTopic, curind)
		curind = BytesWBString(data, c.willMessage, curind)

	}
	//userName string//用户名
	if c.unameFlag {
		curind = BytesWBString(data, c.userName, curind)

	}
	//password string //密码
	if c.passFlag {
		curind = BytesWBString(data, c.password, curind)
	}
	c.totalen = len(data)
	return data
}
