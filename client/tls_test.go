package client_test

import (
	"crypto/tls"

	. "github.com/antlinker/go-mqtt/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("测试客户端连接", func() {
	//alog.RegisterAlog()
	//taskpool.Tlog.SetEnabled(true)
	//Mlog.SetEnabled(false)
	var (
		addr     = "tcp://localhost:8883"
		tlsconfg = &tls.Config{InsecureSkipVerify: true}
	)

	It("异步事件测试", func() {

		client, err := CreateClient(MqttOption{
			Addr: addr,
			Tls:  tlsconfg,
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
		client.AddPublishListener(listener)
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
})
