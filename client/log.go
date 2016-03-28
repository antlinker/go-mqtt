package client

import (
	"github.com/antlinker/alog"
)

const (
	LOG_MQTT_TAG = "MQTTClient"
)

var Mlog *alog.ALog = alog.NewALog()

func init() {
	Mlog.SetLogTag(LOG_MQTT_TAG)
	//Mlog.SetEnabled(false)
}

func MlogInit(configs string) {
	Mlog.ReloadConfig(configs)
	Mlog.SetLogTag(LOG_MQTT_TAG)
}
