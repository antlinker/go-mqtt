package mqttnet_test

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/antlinker/alog"
	"github.com/antlinker/go-mqtt/packet"

	"github.com/antlinker/go-mqtt/mqttnet"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {
	alog.RegisterAlog()
	RegisterFailHandler(Fail)
	RunSpecs(t, "进行网络连接测试")
}
func getTlsConfig(capath, certpath, keypath string) (*tls.Config, error) {
	tlsconfig := &tls.Config{}
	now := time.Now()
	tlsconfig.Time = func() time.Time { return now }
	tlsconfig.Rand = rand.Reader
	if capath != "" {
		pool := x509.NewCertPool()
		caCrt, err := ioutil.ReadFile(capath)
		if err != nil {
			return nil, fmt.Errorf("读取证书失败:%s", err)
		}
		tlsconfig.ClientCAs = pool
		tlsconfig.ClientAuth = tls.RequireAndVerifyClientCert
		pool.AppendCertsFromPEM(caCrt)
	}
	cert, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		return nil, fmt.Errorf("证书错误:%s ", err)
	}
	tlsconfig.Certificates = make([]tls.Certificate, 1)
	tlsconfig.Certificates[0] = cert

	return tlsconfig, nil
}
func createClientTlsConfig(caCertPath, certpath, keypath string) (*tls.Config, error) {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {

		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		//t.Error("Loadx509keypair err %v", err)
		return nil, err
	}

	return &tls.Config{RootCAs: pool, Certificates: []tls.Certificate{cliCrt}}, nil
}

var _ = Describe("测试mqtt服务于连接", func() {
	var (
		server     *mqttnet.Server
		caCertPath = "gokey/ca.crt"
		clientCert = "gokey/client.crt"
		clientKey  = "gokey/client.key"
		serverCert = "gokey/server.crt"
		serverKey  = "gokey/server.key"
		connchan   chan mqttnet.MQTTConner
	)
	var _ = BeforeSuite(func() {
		//mlog.Mlog.SetEnabled(true)
		server = mqttnet.Create(&mqttnet.ServerOption{MaxConnNum: 1024})
		server.Add("tcp", "0.0.0.0:10001")
		server.Add("tcp", "0.0.0.0:10002")

		tlsconfig, err := getTlsConfig(caCertPath, serverCert, serverKey)
		Ω(err).NotTo(HaveOccurred())

		server.AddTLS("tcp", "0.0.0.0:30001", tlsconfig)
		server.AddTLS("tcp", "0.0.0.0:30002", tlsconfig)

		//server.AddTLS("tcp", "0.0.0.0:8813", tlsconfig)
		connchan, _ = server.Start()
		go func(connchan chan mqttnet.MQTTConner) {
			//fmt.Println("等待连接")
			for cc := range connchan {
				if cc == nil {
					break
				}
				go func(conn mqttnet.MQTTConner) {
					defer GinkgoRecover()

					//fmt.Println("已连接连接", cc)
					for {
						msg, err := conn.ReadMessage()
						if err != nil {
							if err == io.EOF {
								break
							}
						}
						Ω(err).NotTo(HaveOccurred())
						//fmt.Println("读取消息", msg)
						err = conn.SendMessage(msg)
						//fmt.Println("回应消息", msg)
						Ω(err).NotTo(HaveOccurred())
					}

				}(cc)

			}

		}(connchan)
	})

	It("测试服务器监听", func() {

		tcphost := []string{"127.0.0.1:10001", "127.0.0.1:10002"}
		for _, host := range tcphost {
			conn, err := mqttnet.Dial("tcp", host, nil)
			Ω(err).NotTo(HaveOccurred())
			msg := packet.NewPingreq()
			conn.SendMessage(msg)
			recmsg, errs := conn.ReadMessage()
			Ω(errs).NotTo(HaveOccurred())
			Ω(msg).Should(Equal(recmsg))
			conn.Close()
		}

	})
	It("测试TLS服务器监听", func() {
		tlsconfig, err := createClientTlsConfig(caCertPath, clientCert, clientKey)
		Ω(err).NotTo(HaveOccurred())
		tcphost := []string{"localhost:30001", "localhost:30002"}
		for _, host := range tcphost {
			conn, err := mqttnet.Dial("tcp", host, tlsconfig)
			Ω(err).NotTo(HaveOccurred())
			msg := packet.NewPingreq()
			conn.SendMessage(msg)
			recmsg, errs := conn.ReadMessage()
			Ω(errs).NotTo(HaveOccurred())
			Ω(msg).Should(Equal(recmsg))
			conn.Close()
		}
	})
	It("读超时测试1", func() {

		tcphost := []string{"127.0.0.1:10001", "127.0.0.1:10002"}
		for _, host := range tcphost {
			conn, err := mqttnet.Dial("tcp", host, nil)
			Ω(err).NotTo(HaveOccurred())
			//conn.SendMessage(msg)
			conn.SetReadTimeout(1 * time.Second)
			_, errs := conn.ReadMessage()
			Ω(errs).Should(HaveOccurred())
			conn.Close()
		}

	})
	It("读超时测试2", func() {
		msg := packet.NewPingreq()
		tcphost := []string{"127.0.0.1:10001", "127.0.0.1:10002"}

		for _, host := range tcphost {
			conn, err := mqttnet.Dial("tcp", host, nil)
			Ω(err).NotTo(HaveOccurred())
			go func(conn mqttnet.MQTTConner) {
				for {
					conn.SetReadTimeout(1 * time.Second)
					_, err := conn.ReadMessage()
					if err != nil {
						if err == io.EOF {
							break
						}
						conn.Close()
						Ω(err).To(HaveOccurred())
					}

				}

			}(conn)
			time.Sleep(100 * time.Millisecond)
			err1 := conn.SendMessage(msg)
			Ω(err1).NotTo(HaveOccurred())
			time.Sleep(2 * time.Second)
			_ = conn.SendMessage(msg)
			//	Ω(err1).To(HaveOccurred())
			//conn.Close()
			break
		}

	})
	// Measure("基准测试读取发送速度", func(b Benchmarker) {

	// 	msg := packet.NewPingreq()
	// 	//fmt.Printf("有桶%d,桶中有元素数量(%d) 个 \n", tongManage.TongLen(), tongManage.Len())
	// 	rt := b.Time("runtime", func() {
	// 		tcphost := []string{"127.0.0.1:10001", "127.0.0.1:10002"}
	// 		for _, host := range tcphost {
	// 			conn, err := mqttnet.Dial("tcp", host, nil)

	// 			Ω(err).NotTo(HaveOccurred())
	// 			for i := 0; i < 1000; i++ {
	// 				start := time.Now()
	// 				conn.SendMessage(msg)
	// 				conn.ReadMessage()
	// 				//Ω(err).NotTo(HaveOccurred())
	// 				//Ω(msg).Should(Equal(recmsg))
	// 				//runtime.Gosched()
	// 				b.RecordValue("执行时间", float64(time.Now().Sub(start))/1000)
	// 			}
	// 			conn.Close()
	// 		}

	// 	})
	// 	Ω(rt.Seconds()).Should(BeNumerically("<", 0.4), "SomethingHard() shouldn't take too long.")
	// 	//b.RecordValue("执行时间:", runtime.Nanoseconds())
	// }, 10)
})
