package client_test

import (
	"fmt"
	"net/http"
	"testing"

	_ "net/http/pprof"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {
	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:7070", nil))
	}()
	RegisterFailHandler(Fail)
	RunSpecs(t, "进行mqtt客户端测试")
}
