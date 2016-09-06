package mqttnet

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/antlinker/go-mqtt/packet"
	//"fmt"
	//	"log"
	"net"
	"sync"
	"time"
)

var (
	// ErrPacketSize 超过最大报文长度限制
	ErrPacketSize = errors.New("超过最大报文长度限制")
)

// MQTTConner mqtt连接接口
type MQTTConner interface {
	net.Conn

	//设置读取超时时间(当前时间之后多久超时)
	//如果有写入完成,正在等待读取,则重置等待时间
	SetReadTimeout(time.Duration)
	//Init(conn net.Conn)
	//同步读取消息报文
	ReadMessage() (packet.MessagePacket, error)
	//发送消息报文
	SendMessage(msg packet.MessagePacket) error
	//设置最大报文大小，默认０不限制
	SetMaxPacketSize(size int32)
	// 获取最大报文大小
	GetMaxPacketSize() int32
}

// Dial mqtt客户端连接
//network 网络协议
// address 网址
// tlsconfig TLS配置
// connectPacket 连接报文,只能设置一个或者不设置
func Dial(network, address string, tlsconfig *tls.Config, connectPacket ...*packet.Connect) (conner MQTTConner, merr *MqttError) {
	var conn net.Conn
	var err error

	if len(connectPacket) > 1 {
		panic("参数错误，只允许传递一个connectPacket")

	}
	if tlsconfig != nil {
		//tls
		conn, err = tls.Dial(network, address, tlsconfig)
		if err != nil {
			merr = createMqttError(ErrCodeNet, err.Error())
			return
		}
	} else {
		//net
		conn, err = net.Dial(network, address)
		if err != nil {
			merr = createMqttError(ErrCodeNet, err.Error())
			return
		}
	}

	conner = NewMqttConn(conn)
	if len(connectPacket) == 0 {
		return
	}
	err = conner.SendMessage(connectPacket[0])
	if err != nil {
		err = createMqttError(ErrCodeNet, err.Error())
		return
	}
	var msg packet.MessagePacket
	msg, err = conner.ReadMessage()
	if err != nil {
		merr = createMqttError(ErrCodeNet, err.Error())
		return
	}
	connbak, ok := msg.(*packet.Connbak)
	if !ok {

		merr = createMqttError(ErrMsgPacket, msg.String())
		conn.Close()
		return
	}
	returncode := connbak.GetReturnCode()
	if returncode == packet.CONNBAK_RETURN_CODE_OK {
		return
	}
	merr = createMqttError(int(returncode), "")
	conn.Close()

	return

}

// NewMqttConn 新建mqtt连接
func NewMqttConn(conn net.Conn) MQTTConner {
	mqttconn := &MQTTConn{}
	mqttconn.Init(conn)
	return mqttconn
}

// MQTTConn mqtt连接实现
type MQTTConn struct {
	sendwait     sync.WaitGroup
	conn         net.Conn
	bufferReader *bufio.Reader
	//bufferWrite  *bufio.Writer
	closing       bool
	maxPacketSize int32
	closeclock    sync.Mutex
	readerChannel chan packet.MessagePacket
	readtimeout   time.Duration
	writeclock    sync.Mutex
}

// SetMaxPacketSize 设置最大报文大小，默认０不限制
func (c *MQTTConn) SetMaxPacketSize(size int32) {
	c.maxPacketSize = size
}

// GetMaxPacketSize 设置最大报文大小，默认０不限制
func (c *MQTTConn) GetMaxPacketSize() int32 {
	return c.maxPacketSize
}
func (c *MQTTConn) Read(b []byte) (n int, err error) {
	if c.readtimeout > 0 {
		c.SetReadDeadline(time.Now().Add(c.readtimeout))
	}
	return c.conn.Read(b)
}
func (c *MQTTConn) Write(b []byte) (n int, err error) {
	c.writeclock.Lock()
	defer c.writeclock.Unlock()
	n, err = c.conn.Write(b)
	if c.readtimeout > 0 {
		c.SetReadDeadline(time.Now().Add(c.readtimeout))
	}
	return
}

// Close 关闭连接
func (c *MQTTConn) Close() error {
	c.closeclock.Lock()
	defer c.closeclock.Unlock()
	if c.closing {
		return nil
	}
	c.closing = true
	c.WaitSendEnd()
	return c.conn.Close()

}

// Init 初始化连接
func (c *MQTTConn) Init(conn net.Conn) {
	c.closing = false
	c.conn = conn
	//c.bufferReader = bufio.NewReaderSize(conn, 1024*1024*4)

}

// LocalAddr 获取本地地址
func (c *MQTTConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr 获取远程地址
func (c *MQTTConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline 设置超时截止时间
func (c *MQTTConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadTimeout 设置超时到期时间
func (c *MQTTConn) SetReadTimeout(readtimeout time.Duration) {
	c.readtimeout = readtimeout
}

// SetReadDeadline 设置读取超时截止时间
func (c *MQTTConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时截止时间
func (c *MQTTConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// ReadMessage 读取消息
func (c *MQTTConn) ReadMessage() (msg packet.MessagePacket, err error) {

	return readMessage(c)

}

// SendMessage 发送消息
func (c *MQTTConn) SendMessage(msg packet.MessagePacket) error {
	if c.closing {
		return errors.New(c.conn.RemoteAddr().String() + "连接关闭中不能继续发送")
	}

	c.sendwait.Add(1)
	defer c.sendwait.Done()
	c.writeclock.Lock()
	defer c.writeclock.Unlock()
	data := msg.Packet()

	//alog.Debug(c.conn.RemoteAddr(), "发送报文:", msg)
	_, err := c.conn.Write(data)
	//time.Sleep(20 * time.Microsecond)
	if err != nil {
		return err
	}
	if c.readtimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readtimeout))
	}
	//c.bufferWrite.Flush()
	return nil
}

// WaitSendEnd 等待发送完成
func (c *MQTTConn) WaitSendEnd() {
	c.sendwait.Wait()
}

// GetConn 获取原始的连接
func (c *MQTTConn) GetConn() net.Conn {
	return c.conn
}

func readMessage(conn MQTTConner) (packet.MessagePacket, error) {
	if conn == nil {
		return nil, fmt.Errorf("连接为空")
	}

	var (
		// the message buffer
		buf []byte

		// tmp buffer to read a single byte
		b = make([]byte, 1)
		// total bytes read
		l = 0
	)
	var recivetime time.Time
	//alog.Debug("++++++++++++Read message: %v", reflect.TypeOf(conn))
	//conn.SetReadDeadline(time.Now().Add(180 * time.Second))
	// Let's read enough bytes to get the message header (msg type, remaining length)
	for {
		// If we have read 5 bytes and still not done, then there's a problem.
		if l > 5 {
			err := fmt.Errorf("connect/getMessage: 4th byte of remaining length has continuation bit set")

			return nil, err
		}
		n, err := conn.Read(b[0:])
		//alog.Debugf("Read message: %d [%X] %v", n, b, err)
		//conn.SetReadDeadline(time.Now().Add(conn.))
		if err != nil {
			//alog.Debug("Read error: %v,%s", err)
			return nil, err
		}
		recivetime = time.Now()
		// Technically i don't think we will ever get here
		if n == 0 {
			continue
		}

		buf = append(buf, b...)
		l += n

		// Check the remlen byte (1+) to see if the continuation bit is set. If so,
		// increment cnt and continue reading. Otherwise break.
		if l > 1 && b[0] < 0x80 {
			break
		}
	}

	// Get the remaining length of the message

	remlen, _ := packet.Bytes2Remlen(buf[1:])
	msize := conn.GetMaxPacketSize()
	if msize > 0 && msize < remlen {
		return nil, ErrPacketSize
	}
	//alog.Debugf("Read message remlen: %d [%X] ", remlen, b)

	rembuf := make([]byte, remlen)
	l = 0
	for l < int(remlen) {
		n, err := conn.Read(rembuf[l:])

		if err != nil {
			return nil, err
		}
		l += n
	}
	//alog.Debugf("+++Read message end: %d==>%d %d [%X] ", l, buf[0], remlen, rembuf)
	pack, err := packet.NewMessageUnPacket(buf[0], int(remlen), rembuf).Parse()
	if err != nil {
		return pack, err
	}
	pub, ok := pack.(*packet.Publish)
	if ok {
		pub.ReciveStartTime = recivetime
	}
	return pack, nil
}
