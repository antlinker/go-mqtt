package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/antlinker/go-mqtt/packet"

	"github.com/antlinker/go-mqtt/mqttnet"
)

//第一次连接
func (c *antClient) fisrtConnect() error {
	if err := c._connect(); err != nil {

		return err
	}
	c.creReceive()

	c.fireOnConnSuccess(c)
	return nil

}

//重新连接
func (c *antClient) reConnect() error {
	c.runGoWait.Add(1)
	go func() {
		defer c.setIssend(true)
		c.connlock.Lock()
		defer c.connlock.Unlock()
		defer c.runGoWait.Done()
		if c.reconnTimeInterval > 0 {
			time.Sleep(c.reconnTimeInterval)
		}
		for c.connected {
			if err := c._connect(); err != nil {
				if c.reconnTimeInterval > 0 {
					time.Sleep(c.reconnTimeInterval)
					continue
				}
			}
			c.connclosed = false
			c.creReceive()
			c.fireOnConnSuccess(c)

			break
		}

	}()
	return nil
}

func (c *antClient) _connect() (err error) {
	defer func() {
		if er := recover(); er != nil {
			e, ok := err.(error)
			if ok {
				err = e
			}
			Mlog.Error("连接出现异常进行恢复：", err)
			c.fireOnConnFailure(c, -1, err)
		}
	}()

	c.fireOnConnStart(c)
	u, _ := url.Parse(c.addr)
	conn, merr := mqttnet.Dial(u.Scheme, u.Host, c.tls, c.connectPacket)
	if merr != nil {

		c.fireOnConnFailure(c, merr.GetErrCode(), merr)
		return merr
	}

	c.conn = conn
	return
}
func (c *antClient) errRecover() {
	if err := recover(); err != nil {
		Mlog.Error("出现异常进行恢复：", err)
	}
}

func (c *antClient) errCloseHandle(cerr *error) {
	if err := recover(); err != nil {
		Mlog.Error("出现异常进行关闭：", err)
		e, ok := err.(error)
		if ok {
			c.errClose(e)
		} else {
			c.errClose(fmt.Errorf("出现异常进行关闭:%v", err))
		}
		return
	}
	//Mlog.Debug("cerr:", *cerr)
	if *cerr != nil {
		c.errClose(*cerr)
	}
}

func (c *antClient) errClose(err error) {
	if c.connclosed {
		return
	}
	c.connclosed = true
	c.setIssend(false)
	c.conn.Close()
	c.runGoWait.Wait()
	c.fireOnLostConn(c, err)
	if c.reconnTimeInterval > 0 && c.connected {
		c.reConnect()
	}
}

var disconnect = packet.NewDisconnect()

func (c *antClient) close() {

}
