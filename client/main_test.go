package client_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "进行mqtt客户端测试")
}
