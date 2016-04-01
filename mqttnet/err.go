package mqttnet

import (
	"strconv"
)

// MqttError 错误
type MqttError struct {
	err  string
	code int
}

// GetError 获取错误
func (e *MqttError) GetError() string {
	return e.err
}

// GetErrCode 获取错误代码
func (e *MqttError) GetErrCode() int {
	return e.code
}

//  Error 实现错误接口
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
	// ErrCodeNet 网络错误
	ErrCodeNet = 99
	// ErrMsgPacket 收到了不符合协议的报文
	ErrMsgPacket = 100
)
