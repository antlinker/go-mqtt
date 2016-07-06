package mqttnet

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/antlinker/alog"
	"github.com/kavu/go_reuseport"

	"golang.org/x/net/websocket"
)

const (
	// LogTag 日志标签
	LogTag = "mqttSrvConn"
)

// ServerOption 服务
type ServerOption struct {
	MaxConnNum int
}

// Server 服务实现
type Server struct {
	listeners []*mqttlistener
	option    ServerOption
	connchan  chan MQTTConner
	runing    bool
	lock      sync.Mutex
}

// Create 创建服务
func Create(option *ServerOption) *Server {
	return &Server{listeners: make([]*mqttlistener, 0),
		option: *option,
		runing: false,
	}
}

// AddWebSocket 增加webscoket监听
func (s *Server) AddWebSocket(network string, laddr string, url string, tlsconfig *tls.Config) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.listeners = append(s.listeners, &mqttlistener{network: network, laddr: laddr, tlsconfig: tlsconfig, ws: &wsconf{url: url}})

}

// AddTLS 增加tcp/tls监听
func (s *Server) AddTLS(network string, laddr string, tlsconfig *tls.Config) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.listeners = append(s.listeners, &mqttlistener{network: network, laddr: laddr, tlsconfig: tlsconfig})

}

// Add 增加tcp监听
func (s *Server) Add(network string, laddr string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.listeners = append(s.listeners, &mqttlistener{network: network, laddr: laddr})
}

// Start 开始服务
func (s *Server) Start() (connchan chan MQTTConner, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.runing {
		return
	}
	connchan = make(chan MQTTConner, s.option.MaxConnNum)
	s.connchan = connchan
	//connchan = make(chan MQTTConner, s.option.MaxConnNum)
	for _, mlistener := range s.listeners {
		err = mlistener.listen(s.connchan)
		if err != nil {
			alog.ErrorTf(LogTag, "启动网络服务:%s %s 失败:%v", mlistener.network, mlistener.laddr, err)
			os.Exit(-1)
		}
	}
	s.runing = true
	alog.DebugTf(LogTag, "mqtt网络服务启动成功")
	return
}

// Stop 停止服务
func (s *Server) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.runing {
		return
	}
	for _, mlistener := range s.listeners {
		mlistener.close()
	}
	close(s.connchan)
	return
}

type wsconf struct {
	url string
}
type mqttlistener struct {
	listener  net.Listener
	network   string
	laddr     string
	tlsconfig *tls.Config
	closeing  bool
	closewait sync.WaitGroup
	ws        *wsconf
}

func (l *mqttlistener) _listen() (err error) {
	alog.DebugTf(LogTag, "监听%s:%s", l.network, l.laddr)
	if l.network == "tcp4" || l.network == "tcp6" {
		l.listener, err = reuseport.NewReusablePortListener(l.network, l.laddr)
		if err != nil {
			alog.DebugTf(LogTag, "端口复用失败：%v", err)
			l.listener, err = net.Listen(l.network, l.laddr)
			if err != nil {
				alog.DebugTf(LogTag, "端口绑定失败：%v", err)
				return err
			}
		}
	} else {
		l.listener, err = net.Listen(l.network, l.laddr)
		if err != nil {
			alog.DebugTf(LogTag, "端口绑定失败：%v", err)
			return err
		}
	}

	if l.tlsconfig != nil {
		l.listener = tls.NewListener(l.listener, l.tlsconfig)
	}
	return nil
}
func (l *mqttlistener) listen(connchan chan MQTTConner) (err error) {
	err = l._listen()
	if err != nil {
		return nil
	}
	if l.ws != nil {
		return l.listenws(connchan)
	}
	var tlstr = ""
	if l.tlsconfig != nil {
		tlstr = "(TLS)"
	}

	l.closewait.Add(1)
	go func() {
		defer func() {
			alog.DebugTf(LogTag, l.network, ":", l.laddr, ":监听关闭退出")
			l.closewait.Done()
		}()
		for {
			//alog.DebugTf(LogTag,"等待客户端连入", l.laddr)
			//fmt.Println("等待客户端连入", l.laddr)
			alog.DebugTf(LogTag, "监听%s:%s 成功", tlstr, l.laddr)
			conn, e := l.listener.Accept()
			if e != nil {
				if l.closeing {
					//正在关闭中
					break
				}
				alog.Error(l.network, ":", l.laddr, ":", "错误连接:", e)
				break
			}
			alog.DebugTf(LogTag, "%s连入%s客户端", l.laddr, conn.RemoteAddr())

			//fmt.Println(l.laddr, "连入", conn.RemoteAddr(), "客户端")
			connchan <- NewMqttConn(conn)
		}
	}()
	return nil
}
func (l *mqttlistener) listenws(connchan chan MQTTConner) (err error) {
	var tlstr = "(ws)"
	if l.tlsconfig != nil {
		tlstr = "(ws|TLS)"
	}

	appServeMux := http.NewServeMux()

	appServeMux.Handle(l.ws.url, websocket.Handler(func(conn *websocket.Conn) {
		alog.DebugT(LogTag, l.laddr, l.ws.url, "连入", conn.RemoteAddr(), "客户端", conn.IsClientConn(), conn.IsServerConn())
		// conn.Write([]byte("mqtt"))

		alog.DebugT(LogTag, "连接状态：", conn.IsClientConn(), conn.IsServerConn())

		// for !conn.IsClientConn() || !conn.IsServerConn() {
		// 	alog.DebugTf(LogTag,"连接状态：", conn.IsClientConn(), conn.IsServerConn())
		// 	time.Sleep(10 * time.Microsecond)
		// }
		wsmqttconn := NewWsMqttConn(conn)
		connchan <- wsmqttconn
		wsmqttconn.startCopy()
		alog.DebugT(LogTag, l.laddr, l.ws.url, "连入", conn.RemoteAddr(), "客户端，转发完成")
	}))
	//appServeMux.Handle(l.ws.url, wsserver)
	go http.Serve(l.listener, appServeMux)
	alog.DebugTf(LogTag, "监听%s:%s 成功", tlstr, l.laddr)
	return nil
}
func (l *mqttlistener) close() {
	l.closeing = true
	l.listener.Close()
	l.closewait.Wait()
	alog.Debugf("mqtt网络服务已经关闭%s:%s", l.network, l.laddr)
}
