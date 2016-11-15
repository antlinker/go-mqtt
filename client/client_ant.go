package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/antlinker/conshash"

	"github.com/antlinker/go-mqtt/packet"

	"github.com/antlinker/go-mqtt/mqttnet"
)

const (
	// clientIDPre 随机客户端标志前缀
	clientIDPre = "antMqtt_"
)

var _zfc = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var _zfclen = len(_zfc)
var _zfcrand = rand.New(rand.NewSource(time.Now().Unix()))

func createRandomString(size int) string {

	out := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		out = append(out, _zfc[_zfcrand.Int()%_zfclen])
	}
	return string(out)
}
func createRandomClientid() string {

	return clientIDPre + createRandomString(10)
}
func parseAddr(addrstr string) ([]string, error) {
	addrs := strings.Split(addrstr, ",")
	for _, addr := range addrs {
		_, err := url.Parse(addr)
		if err != nil {
			return nil, errors.New("未设置mqtt连接格式错误")
		}

	}

	return addrs, nil
}

// CreateClient 创建mqtt客户端
func CreateClient(option MqttOption) (MqttClienter, error) {
	if option.Addr == "" {
		return nil, errors.New("未设置mqtt连接")
	}
	addrs, err := parseAddr(option.Addr)
	if err != nil {
		return nil, errors.New("设置mqtt连接格式错误")
	}
	client := &antClient{addrs: addrs, tls: option.TLS}

	if option.ReconnTimeInterval > 0 {

		client.reconnTimeInterval = time.Duration(option.ReconnTimeInterval) * time.Second
	}
	err = client.init()
	if err != nil {
		return nil, err
	}
	if option.KeepAlive > 0 {
		client.keepAlive = time.Duration(option.KeepAlive) * time.Second
	}
	if option.HeartbeatCheckInterval > 0 {
		client.heartbeatCheckInterval = time.Duration(option.HeartbeatCheckInterval) * time.Second
	} else {
		client.heartbeatCheckInterval = 5 * time.Second
	}
	client.option = option
	return client, nil
}
func (c *antClient) createConnectPacket() *packet.Connect {
	option := c.option
	if c.connectPacket == nil {
		connectPacket := packet.NewConnect()

		connectPacket.SetCleanSession(option.CleanSession)
		if option.Clientid != "" {
			connectPacket.SetClientIdByString(option.Clientid)
		} else {
			connectPacket.SetClientIdByString(createRandomClientid())
			connectPacket.SetCleanSession(true)
		}
		if option.UserName != "" {

			connectPacket.SetUserNameByString(option.UserName)
		}
		if option.WillTopic != "" {
			connectPacket.SetWillTopicInfo([]byte(option.WillTopic), option.WillPayload, packet.QoS(option.WillQos), option.WillRetain)
		}
		connectPacket.SetKeepAlive(option.KeepAlive)
		if option.Password != "" {
			connectPacket.SetPasswordByString(option.Password)
		}
		c.connectPacket = connectPacket
	}

	if option.Password == "" && option.PasswordHandler != nil {
		c.connectPacket.SetPasswordByString(option.PasswordHandler())
	}
	return c.connectPacket
}

// MqttOption mqtt连接配置
type MqttOption struct {
	Addr                   string        //服务器ip端口
	TLS                    *tls.Config   //TLS配置
	ReconnTimeInterval     int           //设置重连时间重连时间<=0不重连　单位秒
	Clientid               string        //客户端标志
	UserName               string        //用户名
	Password               string        //密码
	PasswordHandler        func() string //密码生成函数
	WillTopic              string        //保留消息主题
	WillPayload            []byte        //保留消息有效载荷
	CleanSession           bool          //清理会话标志
	WillQos                QoS           //保留消息服务质量
	WillRetain             bool          //保留消息保留标志
	KeepAlive              uint16        //保持会话时间
	HeartbeatCheckInterval int           //心跳间隔间隔　，发出心跳后，检测心跳间隔时间单位秒
}

type baseClientStatus struct {
	//当前重连次数统计
	curRecnt int64
	//总重连次数统计
	totalRecnt int64
	//接收发布消息报文统计
	recvPubCnt int64
	//发布消息次数统计
	pubcnt map[CntType]*pubCnt

	//接收报文数量
	recvPacketCnt int64
	//发送报文数量
	sendPacketCnt int64
}

func (c *baseClientStatus) init() {
	c.pubcnt = make(map[CntType]*pubCnt)
	c.pubcnt[PubCntQoS0] = &pubCnt{}
	c.pubcnt[PubCntQoS1] = &pubCnt{}
	c.pubcnt[PubCntQoS2] = &pubCnt{}
	c.pubcnt[PubCntTOTAL] = &pubCnt{}

}

type antClient struct {
	mqttListen
	addr                   string   //连接网址
	addrs                  []string //tls配置
	tls                    *tls.Config
	reconnTimeInterval     time.Duration             //重连间隔
	conn                   mqttnet.MQTTConner        //服务器连接
	connectPacket          *packet.Connect           //连接报文
	runGoWait              sync.WaitGroup            //运行中需要等待
	disconnectWait         sync.WaitGroup            //关闭时需要等待
	recvChan               chan packet.MessagePacket //接收通道
	sendChan               chan packet.MessagePacket
	closeChan              chan struct{}
	keepAlive              time.Duration //保持时间
	lasttime               time.Time
	heartbeatCheckInterval time.Duration //心跳间隔间隔　，发出心跳后，检测心跳间隔时间
	connlock               sync.Mutex
	packetManager          PacketManager
	connected              bool
	issend                 bool
	connclosed             bool
	isendlock              sync.RWMutex
	sendcond               *sync.Cond
	hasher                 conshash.ConsistentHashinger
	option                 MqttOption
}

func (c *antClient) SetPacketManager(manager PacketManager) {

	c.packetManager = manager
}

func (c *antClient) addAddr(addrs []string) {
	for _, addr := range addrs {
		c.hasher.Put(addr, addr)
	}
}
func (c *antClient) selectAddr(clientid string) string {
	addr, _ := c.hasher.Get(clientid)
	return addr
}

func (c *antClient) init() error {
	c.baseClientStatus.init()
	c.hasher = conshash.CreateConsistentHashinger(10)
	c.addAddr(c.addrs)
	return nil
}
func (c *antClient) IsConnect() bool {
	return c.connected
}

//开始连接
func (c *antClient) Connect() error {
	c.connlock.Lock()
	defer c.connlock.Unlock()
	if c.connected {
		return nil
	}

	c.connected = true
	if c.packetManager == nil {
		c.packetManager = NewMemPacketManager(c)
		c.packetManager.Start()

		c.sendcond = sync.NewCond(&c.connlock)

	}
	c.recvChan = make(chan packet.MessagePacket)
	clientid := c.option.Clientid
	c.addr, _ = c.hasher.Get(clientid)

	err := c.fisrtConnect()
	if err != nil {
		c.connected = false
		return err
	}
	c.creGoSend()
	c.creDoReceive()
	c.setIssend(true)
	return nil
}

//断开连接
func (c *antClient) Disconnect() {
	c.connlock.Lock()
	if !c.connected {
		return
	}
	c.fireOnDisconning(c)
	c.connected = false
	c.connlock.Unlock()
	c.setIssend(false)
	c.conn.SendMessage(disconnect)
	c.conn.Close()
	c.runGoWait.Wait()
	c.packetManager.Stop()
	c.sendcond.L.Lock()
	c.sendcond.Signal()
	c.sendcond.L.Unlock()

	close(c.recvChan)

	c.disconnectWait.Wait()
	c.fireOnDisconned(c)
}

//发布消息
func (c *antClient) Publish(topic string, qos QoS, retain bool, payload interface{}) (*MqttPacket, error) {

	var mpayload []byte
	switch payload.(type) {
	case []byte:
		mpayload = payload.([]byte)
	case string:
		mpayload = []byte(payload.(string))
	default:
		return nil, fmt.Errorf("有效载荷类型错误:%v", payload)
	}
	msg := packet.NewPublishAll(topic, mpayload, packet.QoS(qos), retain)
	mp, err := c.addPacket(msg)
	if err == nil {

		c.fireOnPubReady(c, mp)
	}
	return mp, err
}

//订阅
func (c *antClient) Subscribe(filter string, qos QoS) (*MqttPacket, error) {
	sub := packet.NewSubscribe(1)
	sub.AddFilter(filter, packet.QoS(qos))
	c.fireOnSubStart(c, []SubFilter{{filter, qos}})
	return c.addPacket(sub)
}

//批量订阅
func (c *antClient) Subscribes(filters ...SubFilter) (*MqttPacket, error) {
	if !c.connected {
		return nil, fmt.Errorf("已经停止和服务器的连接")
	}
	sub := packet.NewSubscribe(1)

	for _, sf := range filters {
		sub.AddFilter(sf.filter, packet.QoS(sf.qos))
	}
	c.fireOnSubStart(c, filters)
	return c.addPacket(sub)
}

//取消订阅
func (c *antClient) UnSubscribe(filter string) (*MqttPacket, error) {
	if !c.connected {
		return nil, fmt.Errorf("已经停止和服务器的连接")
	}
	unsub := packet.NewUnSubscribe(1)
	unsub.AddFilter(filter)
	defer c.fireOnUnSubStart(c, []string{filter})
	return c.addPacket(unsub)
}

//批量取消订阅
func (c *antClient) UnSubscribes(filters ...string) (*MqttPacket, error) {
	if !c.connected {
		return nil, fmt.Errorf("已经停止和服务器的连接")
	}
	unsub := packet.NewUnSubscribe(1)
	for _, f := range filters {
		unsub.AddFilter(f)
	}
	defer c.fireOnUnSubStart(c, filters)
	return c.addPacket(unsub)
}
func (c *antClient) addPacket(msg packet.PacketIdMessage) (*MqttPacket, MqttError) {
	mp := CreateMqttPacket(DirectSend, msg)

	err := c.packetManager.AddSendPacket(mp)
	return mp, err
}
