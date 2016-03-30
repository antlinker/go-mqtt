# MQTT　golang 客户端 
[![GoDoc](https://godoc.org/github.com/antlinker/go-mqtt/client?status.svg)](https://godoc.org/github.com/antlinker/go-mqtt/client)
## 支持

* 支持mqtt 3.1.1协议
* 支持同步模式
* 支持事件模式
* 支持自动重连
* 支持报文数量统计，支持tcp tls连接
* 目前仅支持内存数据存储，可以扩展其他存储

## 获取

``` bash
$ go get github.com/antlinker/go-mqtt/client
```

## 例子

* 参考例子项目[mqttgroup](https://github.com/antlinker/mqttgroup/tree/go-mqtt)的go-mqtt分支
* 参考例子项目[mqttpersonal](https://github.com/antlinker/mqttpersonal/tree/go-mqt)的go-mqtt分支

## API使用
	
### 客户端创建
```go

	import (
		fmt
		 "github.com/antlinker/go-mqtt/client"
	)
	.
	.
	.
	//开始日志　false则关闭日志显示
	client.Mlog.SetEnabled(true)
	client, err := client.CreateClient(client.MqttOption{
			Addr:  "tcp://localhost:8883",
			//断开连接１秒后自动连接，０不自动重连
			ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	.
	.
	.
	//断开连接
	client.Disconnect()
```
### 连接事件


```go

	//连接事件接口
	type MqttDisConnListener interface {
		event.Listener
		OnDisconning(event *MqttEvent)
		OnDisconned(event *MqttEvent)
		OnLostConn(event *MqttEvent, err error)
	}
	//默认连接事件实现，可以打印事件信息
	type DefaultPrintDisConnListen struct {
	}
	
	func (l *DefaultPrintDisConnListen) OnLostConn(event *MqttEvent, err error) {
		Mlog.Debugf("OnLostConn :%v", err)
	}
    
	func (l *DefaultPrintDisConnListen) OnDisconning(event *MqttEvent) {
		Mlog.Debugf("OnDisconning")
	}
	func (l *DefaultPrintDisConnListen) OnDisconned(event *MqttEvent) {
	 	Mlog.Debugf("OnDisconned")
	}
	//默认事件实现，不做任何事情
	type DefaultDisConnListen struct {
	}

	func (l *DefaultDisConnListen) OnDisconning(event *MqttEvent) {

	}
	func (l *DefaultDisConnListen) OnDisconned(event *MqttEvent) {

	}
	func (l *DefaultDisConnListen) OnLostConn(event *MqttEvent, err error) {

	}
    
```
#### 例子
```go

    .
    .
    .
	 client, err := client.CreateClient(client.MqttOption{
		Addr:  "tcp://localhost:1883",
		//断开连接１秒后自动连接，０不自动重连
		ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&DefaultPrintDisConnListen{}
	//注册事件
	client.AddConnListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	//移除事件
	client.RemoveConnListener(listener)   
	.
	.
	.
```
### 断开连接事件


```go

	//事件接口
	type MqttDisConnListener interface {
		event.Listener
		OnDisconning(event *MqttEvent)
		OnDisconned(event *MqttEvent)
		OnLostConn(event *MqttEvent, err error)
	}
	//默认事件实现，可以打印事件信息
	type DefaultPrintDisConnListen struct {
	}
	func (l *DefaultPrintDisConnListen) OnLostConn(event *MqttEvent, err error) {
		Mlog.Debugf("OnLostConn :%v", err)
	}
	func (l *DefaultPrintDisConnListen) OnDisconning(event *MqttEvent) {
		Mlog.Debugf("OnDisconning")
	}
	func (l *DefaultPrintDisConnListen) OnDisconned(event *MqttEvent) {
		Mlog.Debugf("OnDisconned")
	}
	//默认事件实现，不做任何事情
	type DefaultDisConnListen struct {
	}
	func (l *DefaultDisConnListen) OnDisconning(event *MqttEvent) {
	}
	func (l *DefaultDisConnListen) OnDisconned(event *MqttEvent) {
	}
	func (l *DefaultDisConnListen) OnLostConn(event *MqttEvent, err error) {
	}
    
```
#### 例子
```go

	.
	.
	.
	 client, err := client.CreateClient(client.MqttOption{
		Addr:  "tcp://localhost:1883",
		//断开连接１秒后自动连接，０不自动重连
		ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&DefaultPrintDisConnListen{}
	//注册断开连接事件
	client.AddDisConnListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	client.Disconnect()
	//移除事件
	client.RemoveDisConnListener(listener)   
	.
	.
	.
```
### 订阅事件


```go

	//事件接口
	type MqttSubListener interface {
		event.Listener
		OnSubStart(event *MqttEvent, sub []SubFilter)
		OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS)
	}
    	//默认事件实现，可以打印事件信息
	type DefaultPrintSubscribeListen struct {
	}
	func (*DefaultPrintSubscribeListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
		Mlog.Debugf("OnSubStart:%v", sub)
	}
	func (*DefaultPrintSubscribeListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
		Mlog.Debugf("OnSubSuccess:%v:%v", sub, result)
	}
	//默认事件实现，不做任何事情
	type DefaultSubscribeListen struct {
	}
	func (*DefaultSubscribeListen) OnSubStart(event *MqttEvent, sub []SubFilter) {
	}
	func (*DefaultSubscribeListen) OnSubSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	}
    
```
#### 例子
```go

	.
	.
	.
	 client, err := client.CreateClient(client.MqttOption{
			Addr:  "tcp://localhost:1883",
			//断开连接１秒后自动连接，０不自动重连
			ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&DefaultPrintSubscribeListen{}
	//注册订阅事件
	client.AddSubListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	mq,err:= client.Subscribe("Test/1", client.QoS1)
	if err!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",err))
	}
	
	//等待订阅成功
	mq.Wait()
	if mq.Err()!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",mqt.Err()))
	}
	mq,err= client.Subscribes(client.CreateSubFilter("test/1", QoS1),
		client.CreateSubFilter("test/2", QoS2),
		client.CreateSubFilter("test/0", QoS0))
	if err!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",err))
	}
	//等待订阅成功
	mq.Wait()
	if mq.Err()!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",mqt.Err()))
	}
	.
	.
	.
	//移除事件
	client.RemoveSubListener(listener)   
	.
	.
	.
```
### 发布消息事件


```go

	//事件接口
	//发布消息事件监听
	type MqttPublishListener interface {
		event.Listener
		//准备发送消息
		OnPubReady(event *MqttPublishEvent, packet *MqttPacket)
		//消息发布成功
		OnPubSuccess(event *MqttPublishEvent, packet *MqttPacket)
		//消息发布完成
		OnPubFinal(event *MqttPublishEvent, packet *MqttPacket)
	}
    	//默认事件实现，可以打印事件信息
	type DefaultPrintPublishListen struct {
	}
	func (*DefaultPrintPublishListen) OnPubReady(event *MqttPublishEvent, mp *MqttPacket) {
		Mlog.Debugf("OnPubReady:%v", event.GetSendCnt(PubCnt_TOTAL))
	}
	func (*DefaultPrintPublishListen) OnPubSuccess(event *MqttPublishEvent, mp *MqttPacket) {
		Mlog.Debugf("OnPubSuccess:%v", mp.Packet)

	}
	func (*DefaultPrintPublishListen) OnPubFinal(event *MqttPublishEvent, mp *MqttPacket) {
		Mlog.Debugf("OnPubFinal:%v", mp.Packet)
	}

	//默认事件实现，不做任何事情
	type DefaultPublishListen struct {
	}

	func (*DefaultPublishListen) OnPubReady(event *MqttPublishEvent, mp *MqttPacket) {

	}
	func (*DefaultPublishListen) OnPubSuccess(event *MqttPublishEvent, mp *MqttPacket) {

	}
	func (*DefaultPublishListen) OnPubFinal(event *MqttPublishEvent, mp *MqttPacket) {

	}
    
```
#### 例子
```go

	.
	.
	.
	 client, err := client.CreateClient(client.MqttOption{
			Addr:  "tcp://localhost:1883",
			//断开连接１秒后自动连接，０不自动重连
			ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&DefaultPrintPublishListen{}
	//注册订阅事件
	client.AddPubListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	mq, err := client.Publish("test0", QoS0, false, []byte("测试"))
	if err!=nil {
		//发布消息失败
		panic(fmt.Errorf("发布消息失败:%v",err))
	}
	
	//等待发布消息完成
	mq.Wait()
	if mq.Err()!=nil {
		//发布消息失败
		panic(fmt.Errorf("发布消息失败:%v",mqt.Err()))
	}
	
	.
	.
	.
	//移除事件
	client.RemoveSubListener(listener)   
	.
	.
	.
```
### 接收消息事件


```go

	//事件接口
	//接收消息事件监听
	type MqttRecvPubListener interface {
		event.Listener
		OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS)
	}
    	//默认事件实现，可以打印事件信息
	type DefaultPrintRecvPubListen struct {
	}
	func (l *DefaultPrintRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {
		Mlog.Debugf("OnRecvPublish:%s(%d) :%s", topic, qos, string(payload))
	}
	//默认事件实现，不做任何事情
	type DefaultRecvPubListen struct {
	}

	func (l *DefaultRecvPubListen) OnRecvPublish(event *MqttRecvPubEvent, topic string, payload []byte, qos QoS) {

	}
    
```
#### 例子
```go

	.
	.
	.
	 client, err := client.CreateClient(client.MqttOption{
			Addr:  "tcp://localhost:1883",
			//断开连接１秒后自动连接，０不自动重连
			ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&MqttRecvPubListener{}
	//注册订阅事件
	client.AddRecvPubListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	mq,err= client.Subscribe("test",QoS2)
	if err!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",err))
	}
	//等待订阅成功
	mq.Wait()
	if mq.Err()!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",mqt.Err()))
	}
	mq, err := client.Publish("test", QoS1, false, []byte("测试"))
	if err!=nil {
		//发布消息失败
		panic(fmt.Errorf("发布消息失败:%v",err))
	}
	
	//等待发布消息完成
	mq.Wait()
	if mq.Err()!=nil {
		//发布消息失败
		panic(fmt.Errorf("发布消息失败:%v",mqt.Err()))
	}
	
	.
	.
	.
	//移除事件
	client.RemoveRecvPubListener(listener)   
	.
	.
	.
```
### 取消订阅事件


```go

	//事件接口
	type MqttUnSubListener interface {
		event.Listener
		OnUnSubStart(event *MqttEvent, filter []string)
		OnUnSubSuccess(event *MqttEvent, filter []string)
	}
    	//默认事件实现，可以打印事件信息
	type DefaultPrintUnSubListen struct {
	}

	func (*DefaultPrintUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {
		Mlog.Debugf("OnUnSubStart:%v", filter)
	}
	func (*DefaultPrintUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {
		Mlog.Debugf("OnUnSubSuccess:%v", filter)
	}

	//默认事件实现，不做任何事情
	type DefaultUnSubListen struct {
	}
	func (*DefaultUnSubListen) OnUnSubStart(event *MqttEvent, filter []string) {
	}
	func (*DefaultUnSubListen) OnUnSubSuccess(event *MqttEvent, filter []string) {
	}
    
```
#### 例子
```go

	.
	.
	.
	 client, err := client.CreateClient(client.MqttOption{
			Addr:  "tcp://localhost:1883",
			//断开连接１秒后自动连接，０不自动重连
			ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&DefaultPrintUnSubListen{}
	//注册订阅事件
	client.AddUnSubListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	mq,err:= client.UnSubscribe("Test/1")
	if err!=nil {
		//取消订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",err))
	}
	
	//等待取消订阅成功
	mq.Wait()
	if mq.Err()!=nil {
		//取消订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",mqt.Err()))
	}
	mq,err= client.UnSubscribes("test1","test2")
	if err!=nil {
		//订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",err))
	}
	//等待取消订阅成功
	mq.Wait()
	if mq.Err()!=nil {
		//订阅失败
		panic(fmt.Errorf("订阅失败:%v",mqt.Err()))
	}
	.
	.
	.
	//移除事件
	client.RemoveUnSubListener(listener)   
	.
	.
	.
```

### 报文发送接收事件


```go

	//发送接收报文接口
	type MqttPacketListener interface {
		event.Listener
		//接收到报文
		OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64)
		//发送报文
		OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error)
	}
    	//默认事件实现，可以打印事件信息
	type DefaultPrintPacketListen struct {
	}
	func (l *DefaultPrintPacketListen) OnRecvPacket(event *MqttEvent, msg packet.MessagePacket, recvPacketCnt int64) {
		Mlog.Debugf("OnRecvPacket:(%d) :%v", recvPacketCnt, msg)
	}
	func (l *DefaultPrintPacketListen) OnSendPacket(event *MqttEvent, msg packet.MessagePacket, sendPacketCnt int64, err error) {
		Mlog.Debugf("OnSendPacket:(%d) :%v\n %v", sendPacketCnt, msg, err)
	}

	//默认事件实现，不做任何事情
	type DefaultPacketListen struct {
	}
	func (l *DefaultPacketListen) OnRecvPacket(event *MqttEvent, packet packet.MessagePacket, recvPacketCnt int64) {
	}
	func (l *DefaultPacketListen) OnSendPacket(event *MqttEvent, packet packet.MessagePacket, sendPacketCnt int64) {
	}
    
```
#### 例子
```go

	.
	.
	.
	 client, err := client.CreateClient(client.MqttOption{
			Addr:  "tcp://localhost:1883",
			//断开连接１秒后自动连接，０不自动重连
			ReconnTimeInterval: 1,
	})
	if err!=nil {
		//配置文件解析失败
		panic("配置文件解析失败")
	}
	listener:=&DefaultPrintPacketListen{}
	//注册订阅事件
	client.AddPacketListener(listener)
	//建立连接
	err:=client.Connect()
	if err!=nil {
		//连接失败，不会进入自动重连状态
		panic(fmt.Errorf("连接失败:%v",err))
	}
	mq,err:= client.Subscribe("Test/1",client.QoS1)
	if err!=nil {
		//取消订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",err))
	}
	
	//等待取消订阅成功
	mq.Wait()
	if mq.Err()!=nil {
		//取消订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",mqt.Err()))
	}
	mq,err= client.Publish("Test/1", QoS1, false, []byte("测试"))
	if err!=nil {
		//订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",err))
	}
	mq,err= client.Publish("Test/2", QoS2, false, []byte("测试"))
	if err!=nil {
		//订阅失败
		panic(fmt.Errorf("取消订阅失败:%v",err))
	}
	.
	.
	.
	//移除事件
	client.RemovePacketListener(listener)   
	.
	.
	.
```
