package client_test

import (
	"time"

	"github.com/antlinker/go-mqtt/packet"

	. "github.com/antlinker/go-mqtt/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestConnListen struct {
	stopchan chan int
}

func (l *TestConnListen) OnConnStart(event *MqttConnEvent) {
	l.stopchan <- 1
	Expect(event.GetClient().IsConnect()).To(BeTrue())
	Expect(event.GetStatus()).To(Equal(StatusConnStart))

}
func (l *TestConnListen) OnConnSuccess(event *MqttConnEvent) {
	l.stopchan <- 2
	Expect(event.GetClient().IsConnect()).To(BeTrue())
	Expect(event.GetStatus()).To(Equal(StatusConnSuccess))

}
func (l *TestConnListen) OnConnFailure(event *MqttConnEvent, returncode int, err error) {
	l.stopchan <- 3
	Expect(event.GetClient().IsConnect()).To(BeTrue())
	Expect(event.GetStatus()).To(Equal(StatusConnFailure))
}

type TestDisConnListen struct {
	stopchan chan int
}

func (l *TestDisConnListen) OnDisconning(event *MqttEvent) {
	l.stopchan <- 1
	Expect(event.GetClient().IsConnect()).To(BeFalse())
}
func (l *TestDisConnListen) OnDisconned(event *MqttEvent) {
	l.stopchan <- 2
	Expect(event.GetClient().IsConnect()).To(BeFalse())
}
func (l *TestDisConnListen) OnLostConn(event *MqttEvent, err error) {
	l.stopchan <- 3
	Expect(event.GetClient().IsConnect()).To(BeFalse())
	Expect(err).To(HaveOccurred())
}

type TestSubListen struct {
	sub      []SubFilter
	stopchan chan int
}

func (l *TestSubListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
	Expect(l.sub).To(Equal(sub))
	l.stopchan <- 1

}
func (l *TestSubListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	Expect(l.sub).To(Equal(sub))
	Expect(len(result)).To(BeNumerically("==", len(sub)))
	l.stopchan <- 2
}

type TestUnSubListen struct {
	filter   []string
	stopchan chan int
}

func (l *TestUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {
	Expect(l.filter).To(Equal(filter))
	l.stopchan <- 1

}
func (l *TestUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {
	Expect(l.filter).To(Equal(filter))
	l.stopchan <- 2
}

type TestPublishListen struct {
	stopchan chan int
}

func (l *TestPublishListen) OnPubReady(event *MqttPubEvent, mp *MqttPacket) {
	l.stopchan <- 1
	Expect(mp.Direct).To(Equal(DirectSend))

}
func (l *TestPublishListen) OnPubSuccess(event *MqttPubEvent, mp *MqttPacket) {
	Expect(mp.Direct).To(Equal(DirectSend))
	l.stopchan <- 2
}
func (l *TestPublishListen) OnPubFinal(event *MqttPubEvent, mp *MqttPacket) {
	Expect(mp.Direct).To(Equal(DirectSend))
	l.stopchan <- 3

}

type TestRecvPubListen struct {
	topic    string
	payload  []byte
	stopchan chan int
}

func (l *TestRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {
	Expect(l.topic).To(Equal(topic))
	Expect(l.payload).To(Equal(payload))
	l.stopchan <- 1

}

type TestPacketListen struct {
	stopchan chan int
}

func (l *TestPacketListen) OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64) {
	if msg.GetHeaderType() == packet.TYPE_PUBACK {
		l.stopchan <- 1
	}
}
func (l *TestPacketListen) OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error) {
	if msg.GetHeaderType() == packet.TYPE_PUBLISH {
		l.stopchan <- 2
	}
}

var _ = Describe("事件测试", func() {
	Describe("断开连接测试", func() {
		var (
			addr     = "tcp://localhost:1883"
			client   MqttClienter
			err      error
			stopchan = make(chan int)
		)
		It("断开连接事件测试", func() {
			client, err = CreateClient(MqttOption{
				Addr: addr,
			})
			Expect(err).NotTo(HaveOccurred())
			listener := &TestDisConnListen{stopchan: stopchan}
			client.AddDisConnListener(listener)
			err = client.Connect()
			Expect(err).NotTo(HaveOccurred())
			client.Disconnect()
			Expect(client.IsConnect()).To(BeFalse())
			waitChan(stopchan, 3)
		})
	})
	Describe("连接，订阅，发布，取消订阅，接收消息,接收发送报文，监听测试", func() {
		var (
			addr   = "tcp://localhost:1883"
			client MqttClienter
			err    error
		)
		BeforeEach(func() {
			client, err = CreateClient(MqttOption{
				Addr: addr,
			})
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			client.Disconnect()
			Expect(client.IsConnect()).To(BeFalse())
		})
		It("连接监听测试", func() {
			stopchan := make(chan int)
			listener := &TestConnListen{stopchan: stopchan}
			client.AddConnListener(listener)
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			waitChan(stopchan, 2)
		})
		It("订阅监听测试", func() {
			stopchan := make(chan int)
			subfilters := []SubFilter{
				CreateSubFilter("test/1", QoS1),
				CreateSubFilter("test/2", QoS2),
				CreateSubFilter("test/0", QoS0),
			}
			listener := &TestSubListen{stopchan: stopchan, sub: subfilters}
			client.AddSubListener(listener)
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			mq, err := client.Subscribes(subfilters...)
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			waitChan(stopchan, 2)

		})
		It("取消订阅监听测试", func() {
			stopchan := make(chan int)

			listener := &TestUnSubListen{stopchan: stopchan, filter: []string{"test"}}
			client.AddUnSubListener(listener)
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			mq, err := client.UnSubscribe("test")
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			waitChan(stopchan, 2)
			listener.filter = []string{"test1", "test2"}
			mq, err = client.UnSubscribes("test1", "test2")
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			waitChan(stopchan, 2)
			client.RemoveUnSubListener(listener)

		})
		It("发布消息监听测试", func() {
			stopchan := make(chan int)
			listener := &TestPublishListen{stopchan: stopchan}
			client.AddPubListener(listener)
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			mq, err := client.Publish("test0", QoS0, false, []byte("测试"))
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			waitChan(stopchan, 3)
			mq, err = client.Publish("test1", QoS1, false, []byte("测试"))
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			waitChan(stopchan, 3)
			mq, err = client.Publish("test2", QoS2, false, []byte("测试"))
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			waitChan(stopchan, 3)
		})
		It("消息接收事件测试", func() {
			stopchan := make(chan int)
			topic := "test"
			payload := []byte("测试测试")
			listener := &TestRecvPubListen{stopchan: stopchan, topic: topic, payload: payload}
			client.AddRecvPubListener(listener)
			err := client.Connect()

			Expect(err).NotTo(HaveOccurred())

			mq, err := client.Subscribe(topic, QoS2)
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			mq, err = client.Publish(topic, QoS1, false, payload)
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			waitChan(stopchan, 1)
		})
		It("接收发送报文事件监听测试", func() {
			stopchan := make(chan int)
			topic := "test"
			payload := []byte("测试测试")
			listener := &TestPacketListen{stopchan: stopchan}
			client.AddPacketListener(listener)
			err := client.Connect()
			Expect(err).NotTo(HaveOccurred())
			mq, err := client.Subscribe(topic, QoS2)
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())
			mq, err = client.Publish(topic, QoS1, false, payload)
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			waitChan(stopchan, 2)
		})
	})
})

func waitChan(c chan int, repeat int) {
	for i := 0; i < repeat; i++ {
		select {
		case <-c:
		case <-time.After(2 * time.Second):
			Expect(false).To(BeTrue())
		}
	}

}
