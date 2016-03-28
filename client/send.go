package client

import (
	"time"

	"github.com/antlinker/go-mqtt/packet"
)

func (c *antClient) _send(msg packet.MessagePacket) error {

	err := c.conn.SendMessage(msg)
	if err != nil {
		c.fireOnSendPacket(c, msg, err)
		return err
	}
	c.lasttime = time.Now()
	c.fireOnSendPacket(c, msg, nil)
	return nil

}

func (c *antClient) creGoSend() {
	c.disconnectWait.Add(1)
	go c.goSend()
}

func (c *antClient) goSend() {
	defer c.disconnectWait.Done()

	manger := c.packetManager
	//Mlog.Debug("启动发送go程")

	for {

		pm := manger.PopSend()
		if pm == nil {
			break

		}
		c.sendcond.L.Lock()
		for c.connected && !c.checkSend() {
			c.sendcond.Wait()
		}

		err := c.send(pm)

		c.sendcond.L.Unlock()
		if err != nil {
			pm.sending = false
			c.errClose(err)
			continue
		}
		pm.Lastime = time.Now()
		pm.sending = false

	}

}
func (c *antClient) setIssend(issend bool) {
	c.isendlock.Lock()
	if c.issend != issend {
		c.sendcond.Signal()
	}
	c.issend = issend
	c.isendlock.Unlock()
}
func (c *antClient) checkSend() bool {
	c.isendlock.RLock()
	defer c.isendlock.RUnlock()
	return c.issend
}
func (c *antClient) send(msg *MqttPacket) error {
	//Mlog.Debug("send:", msg)
	err := c._send(msg.Packet)
	if err != nil {
		return err
	}
	if msg.Packet.GetHeaderType() == packet.TYPE_PUBLISH {
		c.fireOnPubSuccess(c, msg)
		pub := msg.Packet.(*packet.Publish)
		if pub.GetQos() == packet.QOS_0 {
			msg.finalish(nil)
			c.fireOnPubFinal(c, msg)

		}
	}
	return nil
}
