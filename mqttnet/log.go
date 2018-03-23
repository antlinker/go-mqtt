package mqttnet

import (
	"github.com/antlinker/alog"
)

type ILogger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

var mlog ILogger

func init() {
	tmp := alog.NewALog()
	tmp.SetLogTag(LogTag)
	mlog = tmp
}
func SetLogger(l ILogger) {
	mlog = l
}
