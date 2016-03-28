package packet

var pingreq = NewPingreq()
var PingrespObj = NewPingresp()
var disconnect = NewDisconnect()

// //对象池中创建消息报文对象
// func creReceiveMessagePacketByPool(head uint8, data []byte) (msg MessagePacket, err error) {
// 	switch uint8(head >> 4) {
// 	case TYPE_PUBLISH:
// 		msg = PublishPool.Get().(MessagePacket)
// 		err = msg.UnPacket(head, data)
// 		//publishPool.Put(msg)
// 		return
// 	case TYPE_CONNECT:
// 		msg = ConnectPool.Get().(MessagePacket)
// 		err = msg.UnPacket(head, data)
// 		//connectPool.Put(msg)
// 		return
// 	case TYPE_PUBACK:
// 		msg = PubackPool.Get().(MessagePacket)
// 		err = msg.UnPacket(head, data)
// 		//pubackPool.Put(msg)
// 		return
// 	case TYPE_PUBCOMP:
// 		msg = PubcompPool.Get().(MessagePacket)
// 		//pubcompPool.Put(msg)
// 		err = msg.UnPacket(head, data)
// 		return
// 	case TYPE_PUBREC:
// 		msg = PubrecPool.Get().(MessagePacket)
// 		err = msg.UnPacket(head, data)
// 		//pubrecPool.Put(msg)
// 		return
// 	case TYPE_PUBREL:
// 		msg = PubrelPool.Get().(MessagePacket)
// 		err = msg.UnPacket(head, data)
// 		//pubrecPool.Put(msg)
// 		return
// 	case TYPE_SUBSCRIBE:

// 		msg = SubscribePool.Get().(MessagePacket)
// 		//subscribePool.Put(msg)
// 		err = msg.UnPacket(head, data)
// 		return
// 	case TYPE_UNSUBSCRIBE:
// 		msg = UnSubscribePool.Get().(MessagePacket)
// 		//unsubscribePool.Put(msg)
// 		err = msg.UnPacket(head, data)
// 		return
// 	case TYPE_DISCONNECT:
// 		//disconnectPool.Put(msg)
// 		if len(data) == 0 {
// 			msg = disconnect
// 		} else {
// 			return nil, fmt.Errorf("读取数据失败，Disconnect报文无剩余字节")

// 		}

// 		return
// 	case TYPE_PINGREQ:
// 		if len(data) == 0 {
// 			msg = pingreq
// 		} else {
// 			return nil, fmt.Errorf("读取数据失败，Pingreq报文无剩余字节")

// 		}
// 		//pingreqPool.Put(msg)
// 		return

// 	default:

// 		return nil, fmt.Errorf("读取数据失败，无匹配类型报文放弃该连接:%X", head)

// 	}
// }

// //对象池中创建消息报文对象
// func creReceiveMessagePacket(head uint8, data []byte) (msg MessagePacket, err error) {
// 	switch uint8(head >> 4) {
// 	case TYPE_PUBLISH:
// 		m := NewPublish()
// 		msg = &m
// 		err = msg.UnPacket(head, data)
// 		//publishPool.Put(msg)
// 		return
// 	case TYPE_CONNECT:
// 		msg = NewConnect()
// 		err = msg.UnPacket(head, data)
// 		//connectPool.Put(msg)
// 		return
// 	case TYPE_PUBACK:
// 		msg = NewPuback()
// 		err = msg.UnPacket(head, data)
// 		//pubackPool.Put(msg)
// 		return
// 	case TYPE_PUBCOMP:
// 		msg = NewPubcomp()
// 		//pubcompPool.Put(msg)
// 		err = msg.UnPacket(head, data)
// 		return
// 	case TYPE_PUBREC:
// 		msg = NewPubrec()
// 		err = msg.UnPacket(head, data)
// 		//pubrecPool.Put(msg)
// 		return
// 	case TYPE_PUBREL:
// 		msg = NewPubrel()
// 		err = msg.UnPacket(head, data)
// 		//pubrecPool.Put(msg)
// 		return
// 	case TYPE_SUBSCRIBE:

// 		msg = NewSubscribe()
// 		//subscribePool.Put(msg)
// 		err = msg.UnPacket(head, data)
// 		return
// 	case TYPE_UNSUBSCRIBE:
// 		msg = NewUnSubscribe()
// 		//unsubscribePool.Put(msg)
// 		err = msg.UnPacket(head, data)
// 		return
// 	case TYPE_DISCONNECT:
// 		//disconnectPool.Put(msg)
// 		if len(data) == 0 {
// 			msg = disconnect
// 		} else {
// 			return nil, fmt.Errorf("读取数据失败，Disconnect报文无剩余字节")

// 		}

// 		return
// 	case TYPE_PINGREQ:
// 		if len(data) == 0 {
// 			msg = pingreq
// 		} else {
// 			return nil, fmt.Errorf("读取数据失败，Pingreq报文无剩余字节")

// 		}
// 		//pingreqPool.Put(msg)
// 		return

// 	default:

// 		return nil, fmt.Errorf("读取数据失败，无匹配类型报文放弃该连接:%X", head)

// 	}
// }
