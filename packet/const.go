package packet

import "errors"

const (
	//publish类型 缓存最小大小
	PUBLISH_CACHE_MINLENGTH = 1024 * 1024 * 1024
	PUBLISH_CACHE_DIR       = "/tmp/cache"
)

//报文类型
const (
	TYPE_CONNECT uint8 = 1 + iota
	TYPE_CONNACK
	TYPE_PUBLISH
	TYPE_PUBACK
	TYPE_PUBREC
	TYPE_PUBREL
	TYPE_PUBCOMP
	TYPE_SUBSCRIBE
	TYPE_SUBACK
	TYPE_UNSUBSCRIBE
	TYPE_UNSUBACK
	TYPE_PINGREQ
	TYPE_PINGRESP
	TYPE_DISCONNECT
)

func GetTypeName(btype uint8) string {
	switch btype {
	case TYPE_CONNECT:
		return "连接服务器"
	case TYPE_CONNACK:
		return "连接返回"
	case TYPE_PUBLISH:
		return "发布消息"
	case TYPE_PUBACK:
		return "发布返回"
	case TYPE_PUBREC:
		return "发布收到"
	case TYPE_PUBREL:
		return "发布释放"
	case TYPE_PUBCOMP:
		return "发布完成"
	case TYPE_SUBSCRIBE:
		return "订阅主题"
	case TYPE_SUBACK:
		return "订阅返回"
	case TYPE_UNSUBSCRIBE:
		return "取消订阅"
	case TYPE_UNSUBACK:
		return "取消订阅返回"
	case TYPE_PINGREQ:
		return "心跳请求"
	case TYPE_PINGRESP:
		return "心跳回应"
	case TYPE_DISCONNECT:
		return "连接断开"
	default:
		return "未定义"
	}
}

//报文类型＋控制报文类型的标志位
const (
	TYPE_FLAG_CONNECT     uint8 = TYPE_CONNECT << 4
	TYPE_FLAG_CONNACK     uint8 = TYPE_CONNACK << 4
	TYPE_FLAG_PUBACK      uint8 = TYPE_PUBACK << 4
	TYPE_FLAG_PUBREC      uint8 = TYPE_PUBREC << 4
	TYPE_FLAG_PUBREL      uint8 = TYPE_PUBREL<<4 | 0x2
	TYPE_FLAG_PUBCOMP     uint8 = TYPE_PUBCOMP << 4
	TYPE_FLAG_SUBSCRIBE   uint8 = TYPE_SUBSCRIBE<<4 | 0x2
	TYPE_FLAG_SUBACK      uint8 = TYPE_SUBACK << 4
	TYPE_FLAG_UNSUBSCRIBE uint8 = TYPE_UNSUBSCRIBE<<4 | 0x2
	TYPE_FLAG_UNSUBACK    uint8 = TYPE_UNSUBACK << 4
	TYPE_FLAG_PINGREQ     uint8 = TYPE_PINGREQ << 4
	TYPE_FLAG_PINGRESP    uint8 = TYPE_PINGRESP << 4
	TYPE_FLAG_DISCONNECT  uint8 = TYPE_DISCONNECT << 4
)

//服务质量等级
const (
	QOS_0 QoS = iota
	QOS_1
	QOS_2
)
const (
	RETAIN_YES uint8 = 1
	RETAIN_NO  uint8 = 0
)
const (
	PROTOCOL             = "MQTT"
	PROTOCOL_LEVEL uint8 = 0x4
)

/*
connbak 返回码的常量
*/
const (
	CONNBAK_RETURN_CODE_OK             uint8 = iota //连接已接受 连接已被服务端接受
	CONNBAK_RETURN_NO_SUPPORT_PROTOCOL              //连接已拒绝，不支持的协议版本 服务端不支持客户端请求的 MQTT 协议级别
	CONNBAK_RETURN_NO_CLIENT_ID                     //连接已拒绝，不合格的客户端标识符 客户端标识符是正确的 UTF-8 编码，但服务端不允许使用
	CONNBAK_RETURN_NO_SERVER                        //连接已拒绝，服务端不可用 网络连接已建立，但 MQTT 服务不可用
	CONNBAK_RETURN_ERROR_UNAME_PWD                  //连接已拒绝，无效的用户名或密码 用户名或密码的数据格式无效
	CONNBAK_RETURN_UNAUTHORIZED                     // 连接已拒绝，未授权 客户端未被授权连接到此服务器
	CONNBAK_RETURN_RESERVED                         // 连接已拒绝，保留码 6-255

)

//SUBACK　订阅返回码
const (
	SUBACK_RETURNCODE_QOS_0   uint8 = 0x0  //成功 最大 QoS 0
	SUBACK_RETURNCODE_QOS_1   uint8 = 0x1  //成功 – 最大 QoS 1
	SUBACK_RETURNCODE_QOS_2   uint8 = 0x2  //成功 – 最大 QoS 2
	SUBACK_RETURNCODE_FAILURE uint8 = 0x80 //Failure  失败

)

const (
	PARSE_OK              = iota //解析字节正好
	PARSE_NOT_ENOUGH_ALL         //解析字节不够，下次提供完整字节
	PARSE_NOT_ENOUGH_LACK        //解析字节不够，下次提供新增字节
	PARSE_REMAINDER              //剩余字节
)

var (
	// ErrNoTopic PUBLISH 未设置主题
	ErrNoTopic = errors.New("PUBLISH 未设置主题")
	// 主题包含U+0000字符
	ErrTopicU0000 = errors.New("主题包含U+0000字符")
)
