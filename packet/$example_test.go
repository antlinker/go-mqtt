package packet

import (
	"fmt"
)

//测试封包，解包
func ExampleConnect() {

	var conn = NewConnect()
	conn.SetClientIdByString("testid")
	conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
	conn.SetKeepAlive(80)
	conn.SetUserNameByString("中文")
	conn.SetPasswordByString("password")
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_CONNECT {
		var newconn = NewConnect()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Print("CONNECT解析后结果:\n")
		fmt.Printf("保持时间:%d\n", newconn.GetKeepAlive())
		fmt.Printf("客户端ID:%s\n", newconn.GetClientIdByString())
		fmt.Printf("遗嘱标志:%t qos:%d 保留:%t\n", newconn.GetWillFlag(), newconn.GetWillQoS(), newconn.GetWillRetain())
		fmt.Printf("遗嘱主题:%s\n", newconn.GetWillTopicByString())
		fmt.Printf("遗嘱消息:%s\n", newconn.GetWillMessageByString())
		fmt.Printf("用户标志:%t %s\n", newconn.GetUnameFlag(), newconn.GetUserName())
		fmt.Printf("密码标志:%t %s\n", newconn.GetPassFlag(), newconn.GetPasswordByString())

	}
	//  Output:CONNECT解析后结果:
	//保持时间:80
	//客户端ID:testid
	//遗嘱标志:true qos:1 保留:true
	//遗嘱主题:/topic/test
	//遗嘱消息:这是遗嘱消息
	//用户标志:true 中文
	//密码标志:true password
}

//测试封包，解包
func ExamplePublish() {

	var conn = NewPublish()
	conn.SetControlFlag(false, QOS_1, true)
	conn.SetPacketId(100)
	conn.SetTopicByString("a/b")
	conn.SetPayload([]byte("中国人"))
	var tmp = conn.Packet()
	if uint8(tmp[0]>>4) == TYPE_PUBLISH {
		var newconn = NewPublish()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Print("PUBLISH解析后结果:\n")
		fmt.Printf("重发标志:%t\n", newconn.GetDupFlag())
		fmt.Printf("服务质量等级(Qos):%d\n", newconn.GetQos())
		fmt.Printf("保留标志:%t\n", newconn.GetRetain())
		fmt.Printf("主题名:%s\n", newconn.GetTopicByString())
		fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		fmt.Printf("有效载荷:%s", string(newconn.GetPayload()))

	}
	//  Output:PUBLISH解析后结果:
	//重发标志:false
	//服务质量等级(Qos):1
	//保留标志:true
	//主题名:a/b
	//报文标识:100
	//有效载荷:中国人
}

func ExamplePuback() {

	var conn = NewPuback()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	if tmp[0] == TYPE_FLAG_PUBACK {
		var newconn = NewPuback()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Printf("PUBBACK解析失败: %s\n", err)
		} else {
			fmt.Print("PUBBACK解析后结果:\n")
			fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		}
	}
	//  Output:PUBBACK解析后结果:
	//报文标识:10000

}

func ExamplePubrec() {

	var conn = NewPubrec()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	if tmp[0] == TYPE_FLAG_PUBREC {
		var newconn = NewPubrec()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Printf("PUBREC解析失败: %s\n", err)
		} else {
			fmt.Print("PUBREC解析后结果:\n")
			fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		}
	}
	//  Output:PUBREC解析后结果:
	//报文标识:10000
}

func ExamplePubrel() {
	var conn = NewPubrel()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	if tmp[0] == TYPE_FLAG_PUBREL {
		var newconn = NewPubrel()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Printf("PUBREL解析失败: %s\n", err)
		} else {
			fmt.Print("PUBREL解析后结果:\n")
			fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		}
	}
	//  Output:PUBREL解析后结果:
	//报文标识:10000
}

func ExamplePubcomp() {
	var conn = NewPubcomp()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	if tmp[0] == TYPE_FLAG_PUBCOMP {
		var newconn = NewPubcomp()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Printf("PUBCOMP解析失败: %s\n", err)
		} else {
			fmt.Print("PUBCOMP解析后结果:\n")
			fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		}
	}
	//  Output:PUBCOMP解析后结果:
	//报文标识:10000
}

//测试订阅封包，解包
func ExampleSubscribe() {

	var conn = NewSubscribe(3)
	conn.SetPacketId(10000)
	conn.AddFilter(NewTopic("/topic/test0"), 0)
	conn.AddFilter(NewTopic("/topic/test1"), 1)
	conn.AddFilter(NewTopic("/topic/test2"), 2)
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_SUBSCRIBE {
		var newconn = NewSubscribe()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Print("SUBSCRIBE解析后结果:\n")
		fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		fs := newconn.GetFilters()
		for i := range fs {
			fmt.Printf("主题过滤器:%s Qos: %d\n", fs[i].filter, fs[i].qos)
		}

	}
	//  Output:SUBSCRIBE解析后结果:
	//报文标识:10000
	//主题过滤器:/topic/test0 Qos: 0
	//主题过滤器:/topic/test1 Qos: 1
	//主题过滤器:/topic/test2 Qos: 2
}

//测试订阅响应封包，解包
func ExampleSuback() {

	var conn = NewSuback(3)
	conn.SetPacketId(10000)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_0)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_1)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_2)
	conn.AddReturnCode(SUBACK_RETURNCODE_FAILURE)
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_SUBACK {
		var newconn = NewSuback()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Print("SUBACK解析后结果:\n")
		fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		fs := newconn.GetReturnCodes()
		for i := range fs {
			fmt.Printf("返回码:  %d\n", fs[i])
		}

	}
	//  Output:SUBACK解析后结果:
	//报文标识:10000
	//返回码:  0
	//返回码:  1
	//返回码:  2
	//返回码:  128

}

//测试取消订阅封包，解包
func ExampleUnSubscribe() {

	var conn = NewUnSubscribe(3)
	conn.SetPacketId(10000)
	conn.AddFilter("/topic/test0")
	conn.AddFilter("/topic/test1")
	conn.AddFilter("/topic/test2")
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_UNSUBSCRIBE {
		var newconn = NewUnSubscribe()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Print("UNSUBSCRIBE解析后结果:\n")
		fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
		fs := newconn.GetFilters()
		for i := range fs {
			fmt.Printf("主题过滤器:%s\n", fs[i])
		}

	}
	//  Output:UNSUBSCRIBE解析后结果:
	//报文标识:10000
	//主题过滤器:/topic/test0
	//主题过滤器:/topic/test1
	//主题过滤器:/topic/test2
}

//测试取消订阅响应封包，解包
func ExampleUnSuback() {

	var conn = NewUnSuback()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_UNSUBACK {
		var newconn = NewUnSuback()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Print("UNSUBACK解析后结果:\n")
		fmt.Printf("报文标识:%d\n", newconn.GetPacketId())
	}
	//  Output:UNSUBACK解析后结果:
	//报文标识:10000
}

//测试心跳封包，解包
func ExamplePingreq() {

	var conn = NewPingreq()
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_PINGREQ {
		var newconn = NewPingreq()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Print("PINGREQ解析失败.\n")
		} else {
			fmt.Print("PINGREQ解析成功.\n")
		}
	}
	//  Output:PINGREQ解析成功.
}

//测试心跳响应封包，解包
func ExamplePingresp() {

	var conn = NewPingresp()
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_PINGRESP {
		var newconn = NewPingresp()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Print("PINGRESP解析失败.\n")
		} else {
			fmt.Print("PINGRESP解析成功.\n")
		}
	}
	//  Output:PINGRESP解析成功.
}

//测试断开连接封包，解包
func ExampleDisconnect() {

	var conn = NewDisconnect()
	var tmp = conn.Packet()
	if uint8(tmp[0]) == TYPE_FLAG_DISCONNECT {
		var newconn = NewDisconnect()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		err := newconn.UnPacket(tmp[0], tmp[1+clen:])
		if err != nil {
			fmt.Print("DISCONNECT解析失败.\n")
		} else {
			fmt.Print("DISCONNECT解析成功.\n")
		}
	}
	//  Output:DISCONNECT解析成功.
}
