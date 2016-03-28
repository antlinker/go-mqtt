package client

import "errors"

type MqttError error

var TimeoutError MqttError = errors.New("超时")
