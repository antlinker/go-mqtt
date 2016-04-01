package client

import (
	"fmt"
	"time"

	"github.com/antlinker/go-mqtt/packet"
)

func (c *antClient) creReceive() {
	c.runGoWait.Add(1)
	go func() {
		var err error
		var msg packet.MessagePacket
		defer c.errCloseHandle(&err)
		defer c.runGoWait.Done()
		for {
			msg, err = c.conn.ReadMessage()
			if err != nil {
				//出现错误关闭重连
				//Mlog.Debug("连接关闭:", err)
				break
			}
			c.lasttime = time.Now()
			//Mlog.Debug("receive", msg)
			c.fireOnRecvPacket(c, msg)
			c.recvChan <- msg
		}
	}()
}

var pingreq = packet.NewPingreq()

func (c *antClient) pingreq() {
	c._send(pingreq)
	c.runGoWait.Add(1)
	go func() {
		defer c.runGoWait.Done()
		<-time.After(c.heartbeatCheckInterval)
		if c.lasttime.Add(c.keepAlive).Add(c.heartbeatCheckInterval).Sub(time.Now()) < 0 {
			//服务器已经超时准备断开连接
			c.errClose(ErrTimeout)
		}
	}()
}
func (c *antClient) creDoReceive() {
	c.disconnectWait.Add(1)
	go c.doReceive()
}
func (c *antClient) doReceive() {
	defer c.disconnectWait.Done()
	if c.keepAlive > 0 {
		//Mlog.Debug("doReceive开始接收数据")
		for {
			select {
			case msg := <-c.recvChan:
				if msg == nil {
					break
				}
				c.doMsgPacket(msg)
			case <-time.After(c.keepAlive):
				c.pingreq()
			}
		}
	} else {
		//Mlog.Debug("doReceive开始接收数据")
		for msg := range c.recvChan {
			c.doMsgPacket(msg)
		}
	}
}
func (c *antClient) doMsgPacket(msg packet.MessagePacket) {

	switch msg.GetHeaderType() {
	case packet.TYPE_PUBLISH:
		c.doPublish(msg.(*packet.Publish))
	case packet.TYPE_PUBACK:
		c.doPuback(msg.(*packet.Puback))
	case packet.TYPE_PUBREC:
		c.doPubrec(msg.(*packet.Pubrec))
	case packet.TYPE_PUBREL:
		c.doPubrel(msg.(*packet.Pubrel))
	case packet.TYPE_PUBCOMP:
		c.doPubcomp(msg.(*packet.Pubcomp))
	case packet.TYPE_SUBACK:
		c.doSuback(msg.(*packet.Suback))
	case packet.TYPE_UNSUBACK:

		c.doUnSuback(msg.(*packet.UnSuback))
	case packet.TYPE_PINGRESP:
	default:
		//收到异常报文，丢弃关闭连接
		Mlog.Warnf("收到异常报文％v", msg)
		c.errClose(fmt.Errorf("收到异常报文:%s", msg))

	}
}
func (c *antClient) doPublish(publish *packet.Publish) {

	packetid := publish.GetPacketId()
	switch publish.GetQos() {
	case packet.QOS_0:
		c.fireOnRecvPublish(c, publish.GetTopicByString(), publish.GetPayload(), QoS(publish.GetQos()))
		//释放对象到对象池
		//packet.PublishPool.Put(msg)
		return
	case packet.QOS_1:
		pubak := packet.NewPuback()
		pubak.SetPacketId(packetid)
		c._send(pubak)
		c.fireOnRecvPublish(c, publish.GetTopicByString(), publish.GetPayload(), QoS(publish.GetQos()))
		return
	case packet.QOS_2:
		pubrec := packet.NewPubrec()
		pubrec.SetPacketId(packetid)
		c.packetManager.AddReceivePacket(CreateMqttPacket(DirectRecive, publish))
		c._send(pubrec)

		return
	}
}
func (c *antClient) doPubrel(pubrel *packet.Pubrel) {
	pid := pubrel.GetPacketId()
	mp := c.packetManager.RemoveReceivePacket(pid)
	pubcomp := packet.NewPubcomp()
	pubcomp.SetPacketId(pid)
	c._send(pubcomp)
	if mp != nil {
		publish := mp.Packet.(*packet.Publish)
		c.fireOnRecvPublish(c, publish.GetTopicByString(), publish.GetPayload(), QoS(publish.GetQos()))
	}

}

//收到发布发布结束
func (c *antClient) doPubcomp(pubcomp *packet.Pubcomp) {
	pid := pubcomp.GetPacketId()
	pkt := c.packetManager.RemoveSendQos2Pakcet(pid)
	//Mlog.Debug("doPubcomp:", pkt)
	if pkt != nil {

		pkt.finalish(nil)
		c.fireOnPubFinal(c, pkt)
	}

}

func (c *antClient) doPubrec(pubrec *packet.Pubrec) {
	pid := pubrec.GetPacketId()
	pkt := c.packetManager.RemoveSendUnfinal(pid)
	if pkt != nil {
		pkt.rectime = time.Now()
		c.packetManager.AddSendQos2Pakcet(pkt)
		//发送pubrel
		pubrel := packet.NewPubrel()
		pubrel.SetPacketId(pid)
		c._send(pubrel)
	}

}

func (c *antClient) doPuback(pubback *packet.Puback) {
	pid := pubback.GetPacketId()
	pkt := c.packetManager.RemoveSendUnfinal(pid)
	if pkt != nil {
		pkt.finalish(nil)
		c.fireOnPubFinal(c, pkt)
	}
}

func (c *antClient) doUnSuback(unsuback *packet.UnSuback) {
	pid := unsuback.GetPacketId()
	pkt := c.packetManager.RemoveSendUnfinal(pid)
	if pkt != nil {
		unsub := pkt.Packet.(*packet.UnSubscribe)
		pkt.finalish(nil)
		c.fireOnUnSubSuccess(c, unsub.GetFilters())
	}
}

func (c *antClient) doSuback(suback *packet.Suback) {
	pid := suback.GetPacketId()
	pkt := c.packetManager.RemoveSendUnfinal(pid)
	//Mlog.Debug("doSuback:", pkt)

	if pkt != nil {
		sub := pkt.Packet.(*packet.Subscribe)
		subfs := make([]SubFilter, 0, len(sub.GetFilters()))
		for _, filter := range sub.GetFilters() {
			subfs = append(subfs, SubFilter{
				filter: filter.GetFilter().String(),
				qos:    QoS(filter.GetQos()),
			})

		}
		result := make([]QoS, 0, len(sub.GetFilters()))
		for _, qos := range suback.GetReturnCodes() {
			result = append(result, QoS(qos))
		}
		pkt.finalish(nil)
		c.fireOnSubSuccess(c, subfs, result)
	} else {
		Mlog.Warn("doSuback get faild:", pkt)
	}
}
