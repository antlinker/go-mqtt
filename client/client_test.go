package client_test

import (
	//"strconv"
	"strconv"
	"sync"

	. "github.com/antlinker/go-mqtt/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestMqttListener struct {
	DefaultPrintListener
	sync.WaitGroup
	subFilter   string
	subQos      QoS
	sendTopic   string
	sendPayload []byte
	sendQos     QoS
	sendRetain  bool
}

func (l *TestMqttListener) OnConnSuccess(event *MqttConnEvent) {
	Mlog.Debug("OnConnSuccess")
	event.GetClient().Subscribe(l.subFilter, l.subQos)
}

func (l *TestMqttListener) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	Mlog.Debugf("OnSubscribeSuccess:%v:%v", sub, result)
	event.GetClient().Publish(l.sendTopic, l.sendQos, l.sendRetain, l.sendPayload)
}

func (l *TestMqttListener) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {
	Mlog.Debugf("OnRecvPublish:%s(%d) :%s", topic, qos, string(payload))
	Expect(topic).To(Equal(l.sendTopic))
	Expect(payload).To(Equal(l.sendPayload))
	l.Done()
}

type TestMqttListener2 struct {
	DefaultConnListen
	DefaultPubListen
	sync.WaitGroup
}

func (l *TestMqttListener2) OnConnSuccess(event *MqttConnEvent) {
	Mlog.Debug("OnConnSuccess")
	l.Done()
}

type TestMqttListener3 struct {
	DefaultConnListen
	sync.WaitGroup
}

func (l *TestMqttListener3) OnConnSuccess(event *MqttConnEvent) {
	Mlog.Debug("OnConnSuccess")
	l.Done()
}

type TestMqttListener4 struct {
	DefaultPubListen
	sync.WaitGroup
}

func (*TestMqttListener4) OnPubReady(event *MqttPubEvent, mp *MqttPacket) {
	//Mlog.Debugf("OnPubReady:%v", event.GetSendCnt(PubCnt_TOTAL))
}
func (*TestMqttListener4) OnPubSuccess(event *MqttPubEvent, mp *MqttPacket) {
	//Mlog.Debugf("OnPubSuccess:%v", mp.Packet)

}
func (l *TestMqttListener4) OnPubFinal(event *MqttPubEvent, mp *MqttPacket) {
	//Mlog.Debugf("OnPubFinal:%v", mp.Packet)
	l.Done()
}

var _ = Describe("测试客户端连接", func() {
	//alog.RegisterAlog()
	//taskpool.Tlog.SetEnabled(true)
	Mlog.SetEnabled(false)
	var (
		addr = "tcp://192.168.175.6:1883"
	)

	It("异步事件测试", func() {

		client, err := CreateClient(MqttOption{
			Addr: addr,
		})
		listener := &TestMqttListener{
			subFilter:   "test",
			subQos:      QoS2,
			sendTopic:   "test",
			sendPayload: []byte("我的测试"),
			sendQos:     QoS1,
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
	})
	It("同步测试QoS0", func() {

		client, err := CreateClient(MqttOption{
			Addr: addr,
		})
		listener := &TestMqttListener2{}
		client.AddConnListener(listener)
		//client.AddPubListener(listener)
		//client.AddRecvPubListener(listener)
		//client.AddSubListener(listener)
		//client.AddDisConnListener(listener)
		listener.Add(1)
		Expect(err).NotTo(HaveOccurred())
		By("进行连接")
		err = client.Connect()
		Expect(err).NotTo(HaveOccurred())
		listener.Wait()
		for i := 0; i < 10; i++ {
			mp, err := client.Publish("topic", QoS0, false, []byte("继续我的测试"+strconv.Itoa(i)))
			Expect(err).NotTo(HaveOccurred())
			mp.Wait()
		}

		By("断开连接")

		client.Disconnect()
	})
	It("同步测试QoS1", func() {

		client, err := CreateClient(MqttOption{
			Addr: addr,
		})
		listener := &TestMqttListener2{}
		client.AddConnListener(listener)
		//client.AddPubListener(listener)
		//client.AddRecvPubListener(listener)

		//client.AddDisConnListener(listener)
		listener.Add(1)
		Expect(err).NotTo(HaveOccurred())
		By("进行连接")
		err = client.Connect()
		Expect(err).NotTo(HaveOccurred())
		listener.Wait()

		for i := 0; i < 10; i++ {
			mp, err := client.Publish("topic", QoS1, false, []byte("继续我的测试"+strconv.Itoa(i)))
			Expect(err).NotTo(HaveOccurred())
			mp.Wait()
		}
		By("断开连接")

		client.Disconnect()
	})
	It("同步测试QoS2", func() {

		client, err := CreateClient(MqttOption{
			Addr: addr,
		})
		listener := &TestMqttListener2{}
		client.AddConnListener(listener)
		//client.AddPubListener(listener)
		//client.AddRecvPubListener(listener)
		//client.AddSubListener(listener)
		//client.AddDisConnListener(listener)
		listener.Add(1)
		Expect(err).NotTo(HaveOccurred())
		By("进行连接")
		err = client.Connect()
		Expect(err).NotTo(HaveOccurred())

		listener.Wait()
		for i := 0; i < 10; i++ {
			mp, err := client.Publish("topic", QoS2, false, []byte("继续我的测试"+strconv.Itoa(i)))
			Expect(err).NotTo(HaveOccurred())
			mp.Wait()
		}
		By("断开连接")

		client.Disconnect()
	})
	It("同步测试订阅取消订阅", func() {

		client, err := CreateClient(MqttOption{
			Addr: addr,
		})
		listener := &TestMqttListener2{}
		client.AddConnListener(listener)

		sublistener := &DefaultSubscribeListen{}
		unsublistener := &DefaultUnSubListen{}
		client.AddSubListener(sublistener)
		client.AddUnSubListener(unsublistener)
		listener.Add(1)
		Expect(err).NotTo(HaveOccurred())
		By("进行连接")
		err = client.Connect()
		Expect(err).NotTo(HaveOccurred())
		listener.Wait()
		submp, err := client.Subscribe("test", QoS1)
		Expect(err).NotTo(HaveOccurred())
		submp.Wait()
		unsubmp, err := client.UnSubscribe("test")
		Expect(err).NotTo(HaveOccurred())
		unsubmp.Wait()
		for i := 0; i < 10; i++ {
			mp, err := client.Publish("topic", QoS2, false, []byte("继续我的测试"+strconv.Itoa(i)))
			Expect(err).NotTo(HaveOccurred())
			mp.Wait()
		}
		By("断开连接")

		client.Disconnect()
	})

	It("重连测试", func() {

		client1, err := CreateClient(MqttOption{
			Addr:               addr,
			ReconnTimeInterval: 1,
			Clientid:           "aaa",
		})

		Expect(err).NotTo(HaveOccurred())
		client2, err := CreateClient(MqttOption{
			Addr:               addr,
			ReconnTimeInterval: 2,
			Clientid:           "aaa",
		})

		Expect(err).NotTo(HaveOccurred())
		listener := &TestMqttListener3{}
		client1.AddConnListener(listener)
		client2.AddConnListener(listener)
		client1.AddDisConnListener(&DefaultDisConnListen{})
		client2.AddDisConnListener(&DefaultDisConnListen{})
		listener.Add(3)
		By("进行连接client1")
		err = client1.Connect()
		Expect(err).NotTo(HaveOccurred())
		err = client2.Connect()
		Expect(err).NotTo(HaveOccurred())
		listener.Wait()

		By("断开连接")

		client1.Disconnect()
		client2.Disconnect()
	})

	Measure("进行发布压测QoS1", func(b Benchmarker) {
		var client MqttClienter
		var err error
		client, err = CreateClient(MqttOption{
			Addr: addr,
		})
		lis := &TestMqttListener4{}
		client.AddPubListener(lis)
		Expect(err).NotTo(HaveOccurred())
		err = client.Connect()

		Expect(err).NotTo(HaveOccurred())
		runtime := b.Time("runtime", func() {

			for i := int64(0); i < 1000; i++ {

				_, err := client.Publish("topic", QoS1, false, []byte("继续我的测试"))
				Expect(err).NotTo(HaveOccurred())
				lis.Add(1)
			}
			lis.Wait()

		})
		Ω(runtime.Seconds()).Should(BeNumerically("<", 4), "SomethingHard() shouldn't take too long.")
		b.RecordValue("发送时间:", float64(runtime.Nanoseconds())/float64(1000))
	}, 1)
	Measure("进行发布压测QoS2", func(b Benchmarker) {
		var client MqttClienter
		var err error
		client, err = CreateClient(MqttOption{
			Addr: addr,
		})
		lis := &TestMqttListener4{}
		client.AddPubListener(lis)
		Expect(err).NotTo(HaveOccurred())
		err = client.Connect()

		Expect(err).NotTo(HaveOccurred())
		runtime := b.Time("runtime", func() {

			for i := int64(0); i < 1000; i++ {

				_, err := client.Publish("topic", QoS1, false, []byte("继续我的测试"))
				Expect(err).NotTo(HaveOccurred())
				lis.Add(1)
			}
			lis.Wait()

		})
		Ω(runtime.Seconds()).Should(BeNumerically("<", 4), "SomethingHard() shouldn't take too long.")
		b.RecordValue("发送时间:", float64(runtime.Nanoseconds())/float64(1000))
	}, 1)
	Measure("进行发布压测QoS0", func(b Benchmarker) {
		var client MqttClienter
		var err error
		client, err = CreateClient(MqttOption{
			Addr: addr,
		})
		lis := &TestMqttListener4{}
		client.AddPubListener(lis)
		Expect(err).NotTo(HaveOccurred())
		err = client.Connect()

		Expect(err).NotTo(HaveOccurred())
		runtime := b.Time("runtime", func() {

			for i := int64(0); i < 1000; i++ {

				_, err := client.Publish("topic", QoS1, false, []byte("继续我的测试"))
				Expect(err).NotTo(HaveOccurred())
				lis.Add(1)
			}
			lis.Wait()

		})
		Ω(runtime.Seconds()).Should(BeNumerically("<", 4), "SomethingHard() shouldn't take too long.")
		b.RecordValue("发送时间:", float64(runtime.Nanoseconds())/float64(1000))
	}, 1)
})
