package example_test

import (
	"sync"
	"testing"

	"github.com/antlinker/alog"
	"github.com/antlinker/go-mqtt/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestMqttListener struct {
	client.DefaultPrintListener
	sync.WaitGroup
	subFilter   string
	subQos      client.QoS
	sendTopic   string
	sendPayload []byte
	sendQos     client.QoS
	sendRetain  bool
}

var (
	// Mlog 日志
	Mlog = alog.NewALog()
)

func (l *TestMqttListener) OnConnSuccess(event *client.MqttConnEvent) {
	Mlog.Debug("OnConnSuccess")
	event.GetClient().Subscribe(l.subFilter, l.subQos)
}

func (l *TestMqttListener) OnSubSuccess(event *client.MqttEvent, sub []client.SubFilter, result []client.QoS) {
	Mlog.Debugf("OnSubscribeSuccess:%v:%v", sub, result)
	event.GetClient().Publish(l.sendTopic, l.sendQos, l.sendRetain, l.sendPayload)
}

func (l *TestMqttListener) OnRecvPublish(event *client.MqttRecvPubEvent, topic string, payload []byte, qos client.QoS) {
	Mlog.Debugf("OnRecvPublish:%s(%d) :%s", topic, qos, string(payload))
	Expect(topic).To(Equal(l.sendTopic))
	Expect(payload).To(Equal(l.sendPayload))
	l.Done()
}

func TestCli(t *testing.T) {
	client, err := client.CreateClient(client.MqttOption{
		Addr: "tcp://192.168.175.6:1883",
	})
	listener := &TestMqttListener{
		subFilter:   "test",
		subQos:      client.QoS2,
		sendTopic:   "test",
		sendPayload: []byte("我的测试"),
		sendQos:     client.QoS1,
		sendRetain:  false,
	}
	client.AddConnListener(listener)
	client.AddPubListener(listener)
	client.AddRecvPubListener(listener)
	client.AddSubListener(listener)
	client.AddDisConnListener(listener)
	listener.Add(1)
	Expect(err).NotTo(HaveOccurred())
	By("进行连接")
	err = client.Connect()
	Expect(err).NotTo(HaveOccurred())
	listener.Wait()
	By("断开连接")

	client.Disconnect()
}
