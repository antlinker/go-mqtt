package mqttnet

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/antlinker/alog"

	"golang.org/x/net/websocket"

	"github.com/antlinker/go-mqtt/packet"
)

// NewWsMqttConn 创建webscoket连接
func NewWsMqttConn(conn *websocket.Conn) *MqttWsConn {
	mqttconn := &MqttWsConn{}
	mqttconn.Init(conn)
	return mqttconn
}

// MqttWsConn websocket连接
type MqttWsConn struct {
	sendwait   sync.WaitGroup
	conn       *websocket.Conn
	in         io.Reader
	writein    *io.PipeWriter
	closing    bool
	closeclock sync.Mutex

	readtimeout time.Duration

	writeclock sync.Mutex
}

func (c *MqttWsConn) Read(b []byte) (n int, err error) {

	n, err = c.in.Read(b)

	return
}
func (c *MqttWsConn) Write(b []byte) (n int, err error) {
	c.writeclock.Lock()
	defer c.writeclock.Unlock()
	err = websocket.Message.Send(c.conn, b)
	if c.readtimeout > 0 {
		c.SetReadDeadline(time.Now().Add(c.readtimeout))
	}
	n = len(b)
	return
}

// Close 关闭连接
func (c *MqttWsConn) Close() error {
	c.closeclock.Lock()
	defer c.closeclock.Unlock()
	if c.closing {
		return nil
	}
	c.closing = true
	c.WaitSendEnd()
	return c.conn.Close()

}
func (c *MqttWsConn) startCopy() {
	defer func() {
		//关闭管道
		//alog.DebugTf(LOGTAG,"startCopy，关闭管道")
		c.writein.Close()

	}()

	//	var count int
	for {
		var buffer []byte
		if c.readtimeout > 0 {
			c.SetReadDeadline(time.Now().Add(c.readtimeout))
		}

		err := websocket.Message.Receive(c.conn, &buffer)
		if err != nil {
			return
		}

		//n := len(buffer)
		//	count += n
		i, err := c.writein.Write(buffer)
		if err != nil || i < 1 {
			return
		}

	}
}

// Init 初始化连接
func (c *MqttWsConn) Init(conn net.Conn) {
	wsconn, ok := conn.(*websocket.Conn)
	if !ok {
		alog.DebugTf(LogTag, "连接类型错误")
		return
	}
	c.closing = false
	c.conn = wsconn
	c.in, c.writein = io.Pipe()

}

// LocalAddr 本地地址
func (c *MqttWsConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr 远程地址
func (c *MqttWsConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline 设置超时截止事件
func (c *MqttWsConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadTimeout 设置读取超时间隔时间
func (c *MqttWsConn) SetReadTimeout(readtimeout time.Duration) {
	c.readtimeout = readtimeout
}

// SetReadDeadline 设置读取超时截止时间
func (c *MqttWsConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时截止时间
func (c *MqttWsConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// ReadMessage 读取消息
func (c *MqttWsConn) ReadMessage() (msg packet.MessagePacket, err error) {

	return readMessage(c)

}

// SendMessage 发送消息
func (c *MqttWsConn) SendMessage(msg packet.MessagePacket) error {
	if c.closing {
		return errors.New(c.conn.RemoteAddr().String() + "连接关闭中不能继续发送")
	}

	c.sendwait.Add(1)
	defer c.sendwait.Done()
	c.writeclock.Lock()
	defer c.writeclock.Unlock()
	data := msg.Packet()

	//alog.DebugTf(LOGTAG,c.conn.RemoteAddr(), "发送报文:", msg)
	err := websocket.Message.Send(c.conn, data)
	//_, err := c.conn.Write(data)
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

// WaitSendEnd 等待发送结束
func (c *MqttWsConn) WaitSendEnd() {
	c.sendwait.Wait()
}

// GetConn 获取原始连接
func (c *MqttWsConn) GetConn() net.Conn {
	return c.conn
}
