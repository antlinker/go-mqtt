package packet

import (
	//"fmt"
	//"os"
	"bufio"
	"io"
	"sync"
	"testing"
)

func BenchmarkConnectFullValue(b *testing.B) {
	var conn = NewConnect()
	for i := 0; i < b.N; i++ {
		conn.clientId = []byte("testid")
		//conn.SetClientIdByString("testid")
		conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
		conn.SetKeepAlive(80)
		conn.SetUserNameByString("中文")
		conn.SetPasswordByString("password")

	}
}
func BenchmarkConnectFullValue2(b *testing.B) {
	var conn = NewConnect()
	for i := 0; i < b.N; i++ {
		conn.SetClientIdByString("testid")
		conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
		conn.SetKeepAlive(80)
		conn.SetUserNameByString("中文")
		conn.SetPasswordByString("password")

	}
}
func BenchmarkConnectPacketLen(b *testing.B) {
	var conn = NewConnect()
	conn.SetClientIdByString("testid")
	conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
	conn.SetKeepAlive(80)
	conn.SetUserNameByString("中文")
	conn.SetPasswordByString("password")
	for i := 0; i < b.N; i++ {

		conn.totalRemain()

	}
}

func BenchmarkConnectPacket(b *testing.B) {
	var conn = NewConnect()
	conn.SetClientIdByString("testid")
	conn.SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true)
	conn.SetKeepAlive(80)
	conn.SetUserNameByString("中文")
	conn.SetPasswordByString("password")
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkConnectUnPacket(b *testing.B) {
	var conn = NewConnect()
	conn.SetClientIdByString("testid").SetWillTopicInfoByString("/topic/test", "这是遗嘱消息", QOS_1, true).SetKeepAlive(80).SetUserNameByString("中文").SetPasswordByString("password")
	var tmp = conn.Packet()

	var newconn = NewConnect()

	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {

		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewConnect(b *testing.B) {

	for i := 0; i < b.N; i++ {

		NewConnect()

	}
}

//puback 报文
func BenchmarkPubackFullValue(b *testing.B) {
	var conn = NewPuback()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(11000)
	}
}

func BenchmarkPubackPacket(b *testing.B) {
	var conn = NewPuback()
	conn.SetPacketId(11000)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPubackUnPacket(b *testing.B) {
	var conn = NewPuback()
	conn.SetPacketId(11000)
	var tmp = conn.Packet()

	var newconn = NewPuback()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {

		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPuback(b *testing.B) {

	for i := 0; i < b.N; i++ {

		NewPuback()

	}
}

//connbak　报文测试
func BenchmarkConnbakFullValue(b *testing.B) {
	var conn = NewConnbak()
	for i := 0; i < b.N; i++ {
		conn.SetSessionPresent(false)
		conn.SetReturnCode(CONNBAK_RETURN_ERROR_UNAME_PWD)
	}
}

func BenchmarkConnbakPacket(b *testing.B) {
	var conn = NewConnbak()
	conn.SetSessionPresent(false)
	conn.SetReturnCode(CONNBAK_RETURN_NO_CLIENT_ID)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkConnbakUnPacket(b *testing.B) {
	var conn = NewConnbak()
	conn.SetSessionPresent(true)
	conn.SetReturnCode(CONNBAK_RETURN_NO_SERVER)
	var tmp = conn.Packet()

	var newconn = NewConnbak()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {

		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewConnbak(b *testing.B) {

	for i := 0; i < b.N; i++ {

		NewConnbak()

	}
}

//Publish
func BenchmarkPublishFullValue(b *testing.B) {
	var conn = NewPublish()
	var data = []byte("中国人")
	for i := 0; i < b.N; i++ {
		conn.SetControlFlag(false, QOS_1, true)
		conn.SetPacketId(100)
		conn.SetTopicByString("a/b")
		conn.SetPayload(data)
	}
}
func BenchmarkPublishFullValue2(b *testing.B) {
	var conn = NewPublish()
	for i := 0; i < b.N; i++ {
		conn.SetControlFlag(false, QOS_1, true)
		conn.SetPacketId(10000)
		conn.SetTopicByString("a/b")
		conn.SetPayload([]byte("中国人"))

	}
}
func BenchmarkPublishPacketLen(b *testing.B) {
	var conn = NewPublish()
	conn.SetControlFlag(false, QOS_1, true)
	conn.SetPacketId(10000)
	conn.SetTopicByString("a/b")
	conn.SetPayload([]byte("中国人"))
	for i := 0; i < b.N; i++ {

		conn.totalRemain()

	}
}

func BenchmarkPublishPacket(b *testing.B) {
	var conn = NewPublish()
	conn.SetControlFlag(false, QOS_1, true)
	conn.SetPacketId(10000)
	conn.SetTopicByString("a/b")
	conn.SetPayload([]byte("中国人"))
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPublishUnPacket(b *testing.B) {
	var conn = NewPublish()
	conn.SetControlFlag(false, QOS_1, true)
	conn.SetPacketId(10000)
	conn.SetTopicByString("a/b")
	conn.SetPayload([]byte("中国人"))
	var tmp = conn.Packet()
	var newconn = NewPublish()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPublish(b *testing.B) {
	for i := 0; i < b.N; i++ {

		NewPublish()

	}
}

//Pubcomp
func BenchmarkPubcompFullValue(b *testing.B) {
	var conn = NewPubcomp()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(11000)
	}
}

func BenchmarkPubcompPacket(b *testing.B) {
	var conn = NewPubcomp()
	conn.SetPacketId(11000)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPubcompUnPacket(b *testing.B) {
	var conn = NewPubcomp()
	conn.SetPacketId(11000)
	var tmp = conn.Packet()

	var newconn = NewPubcomp()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {

		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPubcomp(b *testing.B) {

	for i := 0; i < b.N; i++ {

		NewPubcomp()

	}
}

//pubrec

func BenchmarkPubrecFullValue(b *testing.B) {
	var conn = NewPubrec()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(11000)
	}
}

func BenchmarkPubrecPacket(b *testing.B) {
	var conn = NewPubrec()
	conn.SetPacketId(11000)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPubrecUnPacket(b *testing.B) {
	var conn = NewPubrec()
	conn.SetPacketId(11000)
	var tmp = conn.Packet()

	var newconn = NewPubrec()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {

		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPubrec(b *testing.B) {

	for i := 0; i < b.N; i++ {

		NewPubrec()

	}
}

//pubrel

func BenchmarkPubrelFullValue(b *testing.B) {
	var conn = NewPubrel()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(11000)
	}
}

func BenchmarkPubrelPacket(b *testing.B) {
	var conn = NewPubrel()
	conn.SetPacketId(11000)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPubrelUnPacket(b *testing.B) {
	var conn = NewPubrel()
	conn.SetPacketId(11000)
	var tmp = conn.Packet()

	var newconn = NewPubrel()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {

		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPubrel(b *testing.B) {

	for i := 0; i < b.N; i++ {

		NewPubrel()

	}
}

//subscribe

func BenchmarkSubscribeFullValue(b *testing.B) {
	var conn = NewSubscribe()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(11000)
		conn.AddFilter(NewTopic("/topic/test0"), 0)
		conn.AddFilter(NewTopic("/topic/test1"), 1)
		conn.AddFilter(NewTopic("/topic/test2"), 2)
	}
}

func BenchmarkSubscribePacket(b *testing.B) {
	var conn = NewSubscribe()
	conn.SetPacketId(11000)
	conn.AddFilter(NewTopic("/topic/test0"), 0)
	conn.AddFilter(NewTopic("/topic/test1"), 1)
	conn.AddFilter(NewTopic("/topic/test2"), 2)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkSubscribeUnPacket(b *testing.B) {
	var conn = NewSubscribe()
	conn.SetPacketId(11000)
	conn.AddFilter(NewTopic("/topic/test0"), 0)
	conn.AddFilter(NewTopic("/topic/test1"), 1)
	conn.AddFilter(NewTopic("/topic/test2"), 2)
	var tmp = conn.Packet()

	var newconn = NewSubscribe()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewSubscribe(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewSubscribe()
	}
}

//unsubscribe

func BenchmarkUnSubscribeFullValue(b *testing.B) {
	var conn = NewUnSubscribe()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(11000)
		conn.AddFilter("/topic/test0")
		conn.AddFilter("/topic/test1")
		conn.AddFilter("/topic/test2")
	}
}

func BenchmarkUnSubscribePacket(b *testing.B) {
	var conn = NewUnSubscribe()
	conn.SetPacketId(11000)
	conn.AddFilter("/topic/test0")
	conn.AddFilter("/topic/test1")
	conn.AddFilter("/topic/test2")
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkUnSubscribeUnPacket(b *testing.B) {
	var conn = NewUnSubscribe()
	conn.SetPacketId(11000)
	conn.AddFilter("/topic/test0")
	conn.AddFilter("/topic/test1")
	conn.AddFilter("/topic/test2")
	var tmp = conn.Packet()

	var newconn = NewUnSubscribe()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewUnSubscribe(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewUnSubscribe()
	}
}

//suback

func BenchmarkSubackFullValue(b *testing.B) {
	var conn = NewSuback()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(10000)
		conn.AddReturnCode(SUBACK_RETURNCODE_QOS_0)
		conn.AddReturnCode(SUBACK_RETURNCODE_QOS_1)
		conn.AddReturnCode(SUBACK_RETURNCODE_QOS_2)
		conn.AddReturnCode(SUBACK_RETURNCODE_FAILURE)
	}
}

func BenchmarkSubackPacket(b *testing.B) {
	var conn = NewSuback()
	conn.SetPacketId(11000)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_0)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_1)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_2)
	conn.AddReturnCode(SUBACK_RETURNCODE_FAILURE)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkSubackUnPacket(b *testing.B) {
	var conn = NewSuback()
	conn.SetPacketId(11000)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_0)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_1)
	conn.AddReturnCode(SUBACK_RETURNCODE_QOS_2)
	conn.AddReturnCode(SUBACK_RETURNCODE_FAILURE)
	var tmp = conn.Packet()

	var newconn = NewSuback()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewSuback(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewSuback()
	}
}

//unsuback

func BenchmarkUnSubackFullValue(b *testing.B) {
	var conn = NewUnSuback()
	for i := 0; i < b.N; i++ {
		conn.SetPacketId(10000)
	}
}

func BenchmarkUnSubackPacket(b *testing.B) {
	var conn = NewUnSuback()
	conn.SetPacketId(11000)
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkUnSubackUnPacket(b *testing.B) {
	var conn = NewUnSuback()
	conn.SetPacketId(11000)
	var tmp = conn.Packet()

	var newconn = NewUnSuback()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewUnSuback(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewUnSuback()
	}
}

//pingseq

func BenchmarkPingreqPacket(b *testing.B) {
	var conn = NewPingreq()
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPingreqUnPacket(b *testing.B) {
	var conn = NewPingreq()
	var tmp = conn.Packet()

	var newconn = NewPingreq()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPingreq(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewPingreq()
	}
}

//pingsesp

func BenchmarkPingrespPacket(b *testing.B) {
	var conn = NewPingresp()
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkPingrespUnPacket(b *testing.B) {
	var conn = NewPingresp()
	var tmp = conn.Packet()

	var newconn = NewPingresp()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewPingresp(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewPingresp()
	}
}

//disconnect

func BenchmarkDisconnectPacket(b *testing.B) {
	var conn = NewDisconnect()
	for i := 0; i < b.N; i++ {

		conn.Packet()

	}
}
func BenchmarkDisconnectUnPacket(b *testing.B) {
	var conn = NewDisconnect()
	var tmp = conn.Packet()

	var newconn = NewDisconnect()
	datalen, clen := Bytes2Remlen(tmp[1:])
	newconn.remlen = int(datalen)
	for i := 0; i < b.N; i++ {
		newconn.UnPacket(tmp[0], tmp[1+clen:])
	}
}
func BenchmarkNewDisconnect(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewDisconnect()
	}
}

func BenchmarkReadMessage(b *testing.B) {
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
		for i := 0; i < b.N; i++ {
			pwriter.Write(tmp)
			//t.Logf("第%d写入", i)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			ReaderMessagePacket(reader)
			//t.Logf("第%d次读取", i)
		}
		wg.Done()
	}()
	wg.Wait()
}
