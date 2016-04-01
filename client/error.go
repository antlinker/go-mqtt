package client

import "errors"

// MqttError mqtt错误
type MqttError error

// ErrTimeout 超时错误
var ErrTimeout MqttError = errors.New("超时")
