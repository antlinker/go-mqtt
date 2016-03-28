package mqttnet

import (
	"strconv"
)

type MqttError struct {
	err  string
	code int
}

func (e *MqttError) GetError() string {
	return e.err
}
func (e *MqttError) GetErrCode() int {
	return e.code
}
func (e *MqttError) Error() string {
	return "[" + strconv.Itoa(e.code) + "]" + e.err
}
func createMqttError(code int, msg string) *MqttError {
	return &MqttError{
		err:  msg,
		code: code,
	}
}

const (
	ERR_CODE_NET = 99
	//收到了不符合协议的报文
	ERR_MSGPACKET = 100
)
