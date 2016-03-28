package mqttnet

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"

	"github.com/antlinker/go-mqtt/packet"
	//"fmt"
	//	"log"
	"net"
	"sync"
	"time"
)

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
}

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
			merr = createMqttError(ERR_CODE_NET, err.Error())
			return
		}
	} else {
		//net
		conn, err = net.Dial(network, address)
		if err != nil {
			merr = createMqttError(ERR_CODE_NET, err.Error())
			return
		}
	}

	conner = NewMqttConn(conn)
	if len(connectPacket) == 0 {
		return
	}
	err = conner.SendMessage(connectPacket[0])
	if err != nil {
		err = createMqttError(ERR_CODE_NET, err.Error())
		return
	}
	var msg packet.MessagePacket
	msg, err = conner.ReadMessage()
	if err != nil {
		merr = createMqttError(ERR_CODE_NET, err.Error())
		return
	}
	connbak, ok := msg.(*packet.Connbak)
	if !ok {

		merr = createMqttError(ERR_MSGPACKET, msg.String())
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

func NewMqttConn(conn net.Conn) MQTTConner {
	mqttconn := &MQTTConn{}
	mqttconn.Init(conn)
	return mqttconn
}

type MQTTConn struct {
	sendwait     sync.WaitGroup
	conn         net.Conn
	bufferReader *bufio.Reader
	//bufferWrite  *bufio.Writer
	closing    bool
	closeclock sync.Mutex

	readerChannel chan packet.MessagePacket

	readtimeout time.Duration

	writeclock sync.Mutex
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
func (c *MQTTConn) Init(conn net.Conn) {
	c.closing = false
	c.conn = conn
	//c.bufferReader = bufio.NewReaderSize(conn, 1024*1024*4)

}
func (c *MQTTConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *MQTTConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *MQTTConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
func (c *MQTTConn) SetReadTimeout(readtimeout time.Duration) {
	c.readtimeout = readtimeout
}

//设置读取超时截止时间
func (c *MQTTConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *MQTTConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

//读取消息
func (c *MQTTConn) ReadMessage() (msg packet.MessagePacket, err error) {

	return readMessage(c)

}

//发送消息
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
func (c *MQTTConn) WaitSendEnd() {
	c.sendwait.Wait()
}
func (c *MQTTConn) GetConn() net.Conn {
	return c.conn
}

func readMessage(conn io.Reader) (packet.MessagePacket, error) {
	if conn == nil {
		return nil, fmt.Errorf("连接为空")
	}

	var (
		// the message buffer
		buf []byte

		// tmp buffer to read a single byte
		b []byte = make([]byte, 1)
		// total bytes read
		l int = 0
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
