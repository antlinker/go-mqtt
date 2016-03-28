package packet

import (
	"fmt"
	//"os"
	"bufio"
	//"bytes"
	"io"
	"sync"
	"testing"
)

//测试封包，解包
func TestConnectPacketAndUnPacket(t *testing.T) {

	var conn = NewConnect()
	conn.SetClientIdByString("testid")
	conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
	conn.SetKeepAlive(80)
	conn.SetUserNameByString("中文")
	conn.SetPasswordByString("password")
	var tmp = conn.Packet()
	fmt.Print(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if (tmp[0] >> 4) == TYPE_CONNECT {
		if tmp[0]&0xf == 0 {
			var newconn = NewConnect()
			datalen, clen := Bytes2Remlen(tmp[1:])
			newconn.remlen = int(datalen)
			fmt.Printf("%X\n", tmp[1+clen:])
			newconn.UnPacket(tmp[0], tmp[1+clen:])
			fmt.Printf("解析后结果:\n%s", newconn)
		}
	}

}

//测试封包，解包
func TestConnbakPacketAndUnPacket(t *testing.T) {

	var conn = NewConnbak()
	conn.SetSessionPresent(true)
	conn.SetReturnCode(CONNBAK_RETURN_CODE_OK)
	var tmp = conn.Packet()
	fmt.Print(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if (tmp[0] >> 4) == TYPE_CONNACK {
		if tmp[0]&0xf == 0 {
			var newconn = NewConnbak()
			datalen, clen := Bytes2Remlen(tmp[1:])
			newconn.remlen = int(datalen)
			fmt.Printf("%X\n", tmp[1+clen:])
			newconn.UnPacket(tmp[0], tmp[1+clen:])
			fmt.Printf("解析后结果:\n%s", newconn)
		}
	}

}

//测试封包，解包
func TestPublishPacketAndUnPacket(t *testing.T) {

	var conn = NewPublish()
	conn.SetControlFlag(false, QOS_1, true)
	conn.SetPacketId(100)
	conn.SetTopicByString("a/b")
	conn.SetPayload([]byte("中国人"))
	var tmp = conn.Packet()
	fmt.Print(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if uint8(tmp[0]>>4) == TYPE_PUBLISH {
		var newconn = NewPublish()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		fmt.Printf("%X\n", tmp[1+clen:])
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Printf("解析后结果:\n%s", newconn)

	}

}

//测试封包，解包
func TestPubackPacketAndUnPacket(t *testing.T) {

	var conn = NewPuback()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	fmt.Print(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if tmp[0] == TYPE_FLAG_PUBACK {
		var newconn = NewPuback()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		fmt.Printf("%X\n", tmp[1+clen:])
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Printf("解析后结果:\n%s", newconn)

	}

}

//测试封包，解包
func TestPubcompPacketAndUnPacket(t *testing.T) {

	var conn = NewPubcomp()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	fmt.Print(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if tmp[0] == TYPE_FLAG_PUBCOMP {
		var newconn = NewPubcomp()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		fmt.Printf("%X\n", tmp[1+clen:])
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Printf("解析后结果:\n%s", newconn)
	}

}

//测试封包，解包
func TestPubrecPacketAndUnPacket(t *testing.T) {

	var conn = NewPubrec()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	fmt.Print(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if tmp[0] == TYPE_FLAG_PUBREC {
		var newconn = NewPubrec()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		fmt.Printf("%X\n", tmp[1+clen:])
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Printf("解析后结果:\n%s", newconn)

	}

}

//测试封包，解包
func TestPubrelPacketAndUnPacket(t *testing.T) {
	var conn = NewPubrel()
	conn.SetPacketId(10000)
	var tmp = conn.Packet()
	fmt.Println(conn)
	for i := range tmp {
		fmt.Printf("%X ", tmp[i])
	}
	fmt.Printf("%X\n", tmp)
	if tmp[0] == TYPE_FLAG_PUBREL {
		var newconn = NewPubrel()
		datalen, clen := Bytes2Remlen(tmp[1:])
		newconn.remlen = int(datalen)
		fmt.Printf("%X\n", tmp[1+clen:])
		newconn.UnPacket(tmp[0], tmp[1+clen:])
		fmt.Printf("解析后结果:\n%s", newconn)
	}

}

func TestReadMessage(t *testing.T) {
	var conn = NewConnect()
	conn.SetClientIdByString("testid")
	conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
	conn.SetKeepAlive(80)
	conn.SetUserNameByString("中文")
	conn.SetPasswordByString("password")
	var tmp = conn.Packet()
	preader, pwriter := io.Pipe()
	reader := bufio.NewReader(preader)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			pwriter.Write(tmp)
			t.Logf("第%d写入", i)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			t.Log(ReaderMessagePacket(reader))
			t.Logf("第%d次读取", i)
		}
		wg.Done()
	}()
	wg.Wait()
}
func TestSpilt(t *testing.T) {
	tmp := splitTopic("$aa/bb/c/d/e")
	t.Logf("%s %d", tmp, len(tmp))
}
func TestFilter(t *testing.T) {
	a := NewTopicFilter(NewTopic("#a"), QOS_1)
	if a.IsValidTopicFilter() {
		t.Error("#a 这是一个无效规则")
	}
	a = NewTopicFilter(NewTopic("a#"), QOS_1)
	if a.IsValidTopicFilter() {
		t.Error("a# 这是一个无效规则")
	}
	a = NewTopicFilter(NewTopic("aa/#/aa"), QOS_1)
	if a.IsValidTopicFilter() {
		t.Error("a/#/a 这是一个无效规则")
	}
	a = NewTopicFilter(NewTopic("$#"), QOS_1)
	if !a.IsValidTopicFilter() {
		t.Error("$# 这是一个有效规则")
	}
	a = NewTopicFilter(NewTopic("aa/#"), QOS_1)
	if !a.IsValidTopicFilter() {
		t.Error("aa/# 这是一个有效规则")
	}

	a = NewTopicFilter(NewTopic("aa+"), QOS_1)
	if a.IsValidTopicFilter() {
		t.Error("aa+ 这是一个无效规则")
	}
	a = NewTopicFilter(NewTopic("+aa"), QOS_1)
	if a.IsValidTopicFilter() {
		t.Error("+aa 这是一个无效规则")
	}
	a = NewTopicFilter(NewTopic("+/aa"), QOS_1)
	if !a.IsValidTopicFilter() {
		t.Error("+/aa 这是一个有效规则")
	}
}
func TestCompart(t *testing.T) {
	a := NewTopic("aa/bb/c")
	f := NewTopic("aa/bb/c")
	_compart(a, f, t, true)
	f = NewTopic("+/bb/c")
	_compart(a, f, t, true)
	f = NewTopic("+/+/c")
	_compart(a, f, t, true)
	f = NewTopic("+/+/cc")
	_compart(a, f, t, false)
	f = NewTopic("+/+/+")
	_compart(a, f, t, true)
	f = NewTopic("aa/+/+")
	_compart(a, f, t, true)
	f = NewTopic("bb/+/+")
	_compart(a, f, t, false)
	f = NewTopic("aa/bb/+")
	_compart(a, f, t, true)
	f = NewTopic("aa/d/+")
	_compart(a, f, t, false)
	f = NewTopic("aa/bb/c")
	_compart(a, f, t, true)
	f = NewTopic("aa/bb/d")
	_compart(a, f, t, false)
	f = NewTopic("+/+/+/+")
	_compart(a, f, t, false)

	f = NewTopic("+/+/+/")
	_compart(a, f, t, false)
	f = NewTopic("#/+/+")
	_compart(a, f, t, false)
	f = NewTopic("#/+/+/")
	_compart(a, f, t, false)
	f = NewTopic("#/aa/bb")
	_compart(a, f, t, false)
	f = NewTopic("#/aa/bb/c")
	_compart(a, f, t, false)
	f = NewTopic("#")
	_compart(a, f, t, true)
	f = NewTopic("#/#")
	_compart(a, f, t, false)
	f = NewTopic("#/bb/c")
	_compart(a, f, t, false)
	f = NewTopic("#/c")
	_compart(a, f, t, false)
	a = NewTopic("a/b/c/d/e")
	f = NewTopic("#/c/#")
	_compart(a, f, t, false)
	f = NewTopic("#/c/#/e")
	_compart(a, f, t, false)
	f = NewTopic("#/c/#/d/e")
	_compart(a, f, t, false)
	f = NewTopic("#/f/#/d/e")
	_compart(a, f, t, false)
	f = NewTopic("#/c/d/#")
	_compart(a, f, t, false)
	a = NewTopic("/aa/bb/c/d/e")
	f = NewTopic("#")
	_compart(a, f, t, true)
	f = NewTopic("/#")
	_compart(a, f, t, true)
	f = NewTopic("/aa/#")
	_compart(a, f, t, true)
	f = NewTopic("+/bb/c/d/e")
	_compart(a, f, t, true)
	f = NewTopic("/+/bb/c/d/e")
	_compart(a, f, t, true)
	a = NewTopic("$aa/bb/c/d/e")
	f = NewTopic("#")
	_compart(a, f, t, false)
	f = NewTopic("+/bb/c/d/e")
	_compart(a, f, t, false)
	f = NewTopic("$+/bb/c/d/e")
	_compart(a, f, t, true)
	f = NewTopic("$#")
	_compart(a, f, t, false)
}

func _compart(topic *Topic, filter *Topic, t *testing.T, result bool) {
	if topic.CompareFilter(filter) != result {
		t.Errorf("%s <> %s ====%v  |预期结果%v", topic, filter, topic.CompareFilter(filter), result)
	}

}
