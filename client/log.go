package client

import (
	"github.com/antlinker/alog"
)

const (
	// LogMqttTag 日志标签
	LogMqttTag = "MQTTClient"
)

var (
	// Mlog 日志
	Mlog = alog.NewALog()
)

func init() {
	Mlog.SetLogTag(LogMqttTag)
	//Mlog.SetEnabled(false)
}

// MlogInit 初始化日志配置
func MlogInit(configs string) {
	Mlog.ReloadConfig(configs)
	Mlog.SetLogTag(LogMqttTag)
}
