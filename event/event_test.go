package event

import

//"fmt"
"testing"

//"time"

type TestEvent struct {
	BaseEvent
}

func CreateTestEvent(to *TestObject) *TestEvent {
	event := &TestEvent{}
	event.etype = 1
	event.target = to
	return event
}

type TestListener interface {
	OnNo(event *TestEvent)
	OnOk(event *TestEvent, a string)
}

type TestObject struct {
	EventGenerator
}

func (to *TestObject) AddTestListener(tlistener TestListener) {
	to.AddListener(1, tlistener)
}
func (to *TestObject) FireOk() {
	//Mlog.Debug( "调用了FireOk方法")
	to.FireListener(1, "OnOk", CreateTestEvent(to), "测试ok")
}
func (to *TestObject) FireNo() {
	//Mlog.Debug( "调用了FireNo方法")
	to.FireListener(1, "OnNo", CreateTestEvent(to))
}

type DefaultTestListener struct {
}

func (d *DefaultTestListener) OnNo(event *TestEvent) {
	//Mlog.Info( "调用了OnNo方法\n")
}
func (d *DefaultTestListener) OnOk(event *TestEvent, a string) {
	//Mlog.Info( "调用了OnOk方法:", a)
}

//func TestListenerOp(t *testing.T) {
//	t.Log("开始测试")
//	to := &TestObject{}
//	dt := &DefaultTestListener{}
//	to.AddTestListener(dt)
//	to.FireOk()
//	to.FireNo()
//	to.RemoveListener(1, dt)
//	to.FireOk()
//	to.FireNo()
//}

func BenchmarkOp(b *testing.B) {
	to := &TestObject{}
	dt := &DefaultTestListener{}
	to.AddTestListener(dt)
	for i := 0; i < b.N; i++ {

		to.FireNo()

	}
}
