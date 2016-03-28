package mqttnet

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/antlinker/alog"

	"golang.org/x/net/websocket"
)

const (
	LOGTAG = "mqttSrvConn"
)

type ServerOption struct {
	MaxConnNum int
}
type Server struct {
	listeners []*mqttlistener
	option    ServerOption
	connchan  chan MQTTConner
	runing    bool
	lock      sync.Mutex
}

func Create(option *ServerOption) *Server {
	return &Server{listeners: make([]*mqttlistener, 0),
		option: *option,
		runing: false,
	}
}

func (s *Server) AddWebSocket(network string, laddr string, url string, tlsconfig *tls.Config) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.listeners = append(s.listeners, &mqttlistener{network: network, laddr: laddr, tlsconfig: tlsconfig, ws: &wsconf{url: url}})

}

func (s *Server) AddTls(network string, laddr string, tlsconfig *tls.Config) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.listeners = append(s.listeners, &mqttlistener{network: network, laddr: laddr, tlsconfig: tlsconfig})

}
func (s *Server) Add(network string, laddr string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.listeners = append(s.listeners, &mqttlistener{network: network, laddr: laddr})
}
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
			alog.ErrorTf(LOGTAG, "启动网络服务:%s %s 失败:%v", mlistener.network, mlistener.laddr, err)
			os.Exit(-1)
		}
	}
	s.runing = true
	alog.DebugTf(LOGTAG, "mqtt网络服务启动成功")
	return
}
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

func (l *mqttlistener) listen(connchan chan MQTTConner) (err error) {
	if l.ws != nil {
		return l.listenws(connchan)
	}
	if l.tlsconfig != nil {
		l.listener, err = tls.Listen(l.network, l.laddr, l.tlsconfig)
	} else {
		l.listener, err = net.Listen(l.network, l.laddr)
	}
	if err != nil {
		return err
	}
	l.closewait.Add(1)
	go func() {
		defer func() {
			alog.DebugTf(LOGTAG, l.network, ":", l.laddr, ":监听关闭退出")
			l.closewait.Done()
		}()
		for {
			//alog.DebugTf(LOGTAG,"等待客户端连入", l.laddr)
			//fmt.Println("等待客户端连入", l.laddr)
			alog.DebugTf(LOGTAG, l.network, ":监听 ", l.laddr, "成功")
			conn, e := l.listener.Accept()
			if e != nil {
				if l.closeing {
					//正在关闭中
					break
				}
				alog.Error(l.network, ":", l.laddr, ":", "错误连接:", e)
				break
			}
			alog.DebugTf(LOGTAG, l.laddr, "连入", conn.RemoteAddr(), "客户端")

			//fmt.Println(l.laddr, "连入", conn.RemoteAddr(), "客户端")
			connchan <- NewMqttConn(conn)
		}
	}()
	return nil
}
func (l *mqttlistener) listenws(connchan chan MQTTConner) (err error) {
	if l.tlsconfig != nil {
		l.listener, err = tls.Listen(l.network, l.laddr, l.tlsconfig)
	} else {
		l.listener, err = net.Listen(l.network, l.laddr)
	}
	if err != nil {
		return err
	}

	appServeMux := http.NewServeMux()

	appServeMux.Handle(l.ws.url, websocket.Handler(func(conn *websocket.Conn) {
		alog.DebugTf(LOGTAG, l.laddr, l.ws.url, "连入", conn.RemoteAddr(), "客户端", conn.IsClientConn(), conn.IsServerConn())
		// conn.Write([]byte("mqtt"))

		alog.DebugTf(LOGTAG, "连接状态：", conn.IsClientConn(), conn.IsServerConn())

		// for !conn.IsClientConn() || !conn.IsServerConn() {
		// 	alog.DebugTf(LOGTAG,"连接状态：", conn.IsClientConn(), conn.IsServerConn())
		// 	time.Sleep(10 * time.Microsecond)
		// }
		wsmqttconn := NewWsMqttConn(conn)
		connchan <- wsmqttconn
		wsmqttconn.startCopy()
		alog.DebugTf(LOGTAG, l.laddr, l.ws.url, "连入", conn.RemoteAddr(), "客户端，转发完成")
	}))
	//appServeMux.Handle(l.ws.url, wsserver)
	go http.Serve(l.listener, appServeMux)
	alog.DebugTf(LOGTAG, l.network, ":监听 ", l.laddr, "成功")
	return nil
}
func (l *mqttlistener) close() {
	l.closeing = true
	l.listener.Close()
	l.closewait.Wait()
	alog.Debugf("mqtt网络服务已经关闭%s:%s", l.network, l.laddr)
}
