package client_test

import (
	"time"

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
	Expect(event.GetStatus()).To(Equal(STATUS_CONN_START))

}
func (l *TestConnListen) OnConnSuccess(event *MqttConnEvent) {
	l.stopchan <- 2
	Expect(event.GetClient().IsConnect()).To(BeTrue())
	Expect(event.GetStatus()).To(Equal(STATUS_CONN_SUCCESS))

}
func (l *TestConnListen) OnConnFailure(event *MqttConnEvent, returncode int, err error) {
	l.stopchan <- 3
	Expect(event.GetClient().IsConnect()).To(BeTrue())
	Expect(event.GetStatus()).To(Equal(STATUS_CONN_FAILURE))
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

type TestSubscribeListen struct {
	sub []SubFilter
}

func (l *TestSubscribeListen) OnSubscribeStart(event *MqttEvent, sub []SubFilter) {
	Expect(l.sub).To(Equal(sub))

}
func (l *TestSubscribeListen) OnSubscribeSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	Expect(l.sub).To(Equal(sub))
	Expect(len(result)).To(BeNumerically("==", len(sub)))
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
			select {
			case <-stopchan:
			case <-time.After(2 * time.Second):
				Expect(false).To(BeTrue())
			}
			select {
			case <-stopchan:
			case <-time.After(2 * time.Second):
				Expect(false).To(BeTrue())
			}
			select {
			case <-stopchan:
			case <-time.After(2 * time.Second):
				Expect(false).To(BeTrue())
			}
		})
	})
	Describe("连接监听测试", func() {
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
			client.Connect()
			select {
			case <-stopchan:
			case <-time.After(2 * time.Second):
				Expect(false).To(BeTrue())
			}
			select {
			case <-stopchan:
			case <-time.After(2 * time.Second):
				Expect(false).To(BeTrue())
			}
		})
		It("订阅监听测试", func() {
			subfilters := []SubFilter{
				CreateSubFilter("test/1", QoS1),
				CreateSubFilter("test/2", QoS2),
				CreateSubFilter("test/0", QoS0),
			}
			listener := &TestSubscribeListen{sub: subfilters}
			client.AddSubListener(listener)
			client.Connect()
			mq, err := client.Subscribes(subfilters...)
			Expect(err).NotTo(HaveOccurred())
			mq.Wait()
			Expect(mq.Err()).NotTo(HaveOccurred())

		})
	})
})
