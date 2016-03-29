# MQTT　golang 客户端 
## 支持
支持mqtt 3.1.1协议,可以支持同步异步模式，支持事件模式，支持自动重连
## 获取

``` bash
$ go get github.com/antlinker/go-mqtt/client
```
## API使用

### 客户端创建
```

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


```

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
```

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


```

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
```

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


```

	//事件接口
	type MqttSubListener interface {
		event.Listener
		OnSubscribeStart(event *MqttEvent, sub []SubFilter)
		OnSubscribeSuccess(event *MqttEvent, sub []SubFilter, result []QoS)
	}
    	//默认事件实现，可以打印事件信息
	type DefaultPrintSubscribeListen struct {
	}
	func (*DefaultPrintSubscribeListen) OnSubscribeStart(event *MqttEvent, sub []SubFilter) {
		Mlog.Debugf("OnSubscribeStart:%v", sub)
	}
	func (*DefaultPrintSubscribeListen) OnSubscribeSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
		Mlog.Debugf("OnSubscribeSuccess:%v:%v", sub, result)
	}
	//默认事件实现，不做任何事情
	type DefaultSubscribeListen struct {
	}
	func (*DefaultSubscribeListen) OnSubscribeStart(event *MqttEvent, sub []SubFilter) {
	}
	func (*DefaultSubscribeListen) OnSubscribeSuccess(event *MqttEvent, sub []SubFilter, result []QoS) {
	}
    
```
#### 例子
```

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