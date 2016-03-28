package event

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/antlinker/go-cmap"

	"github.com/antlinker/taskpool"
)

//var triggerEventChan = make(chan *TriggerEvent, 1024)

//func init() {

//	//启动事件主go程
//	go func(triggerEventChan <-chan *TriggerEvent) {
//		for {
//			select {
//			case triggerEvent := <-triggerEventChan:
//				//	Mlog.Debug( "收到触发事件:%v", triggerEvent)
//				if triggerEvent == nil {
//					Mlog.Warn( "事件调度线程退出")
//					return
//				}
//				go doTriggerEvent(triggerEvent)

//			}

//		}
//	}(triggerEventChan)
//}
// func doTriggerEvent(triggerEvent *TriggerEvent) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("======:", triggerEvent.method, err)
// 		}
// 	}()
// 	m := triggerEvent.listener.MethodByName(triggerEvent.method)
// 	in := triggerEvent.param
// 	mtype := m.Type()
// 	pcnt := mtype.NumIn()

// 	values := make([]reflect.Value, pcnt)
// 	values[0] = reflect.ValueOf(triggerEvent.event)
// 	if pcnt > 0 {
// 		for i := 1; i < pcnt; i++ {
// 			v := in[i-1]
// 			if v != nil {

// 				values[i] = reflect.ValueOf(in[i-1])
// 			} else {
// 				values[i] = reflect.Zero(mtype.In(i))

// 			}
// 		}
// 	}

// 	fmt.Println(triggerEvent.method, m.Type())
// 	fmt.Println("values:", len(values), values)
// 	fmt.Println("in:", len(in), in)
// 	m.Call(values)

// }
func createRegListen(listener Listener) *RegListen {
	return &RegListen{
		listener: listener,
		value:    reflect.ValueOf(listener),
		methods:  cmap.NewConcurrencyMap(),
	}

}

type rmethod struct {
	name     string
	mtype    reflect.Type
	otherIns []reflect.Type
	innum    int
	value    reflect.Value
}
type RegListen struct {
	listener Listener
	value    reflect.Value
	methods  cmap.ConcurrencyMap
}

func (l *RegListen) exec(event Event, method string, param []interface{}) {
	var (
		rm          *rmethod
		inum        int
		mtype       reflect.Type
		otherIns    []reflect.Type
		methodValue reflect.Value
	)

	m, _ := l.methods.Get(method)
	if m == nil {
		methodValue = l.value.MethodByName(method)
		mtype = methodValue.Type()
		inum = mtype.NumIn()
		otherIns = make([]reflect.Type, inum-1)
		if inum > 1 {
			for i := 1; i < inum; i++ {
				otherIns[i-1] = mtype.In(i)
			}
		}
		rm = &rmethod{name: method, mtype: mtype, otherIns: otherIns, innum: inum, value: methodValue}
	} else {
		rm = m.(*rmethod)
		inum = rm.innum
		mtype = rm.mtype
		otherIns = rm.otherIns
		methodValue = rm.value
	}

	params := make([]reflect.Value, inum)
	params[0] = reflect.ValueOf(event)
	if inum > 1 {
		for i := 1; i < inum; i++ {
			v := param[i-1]
			if v != nil {
				params[i] = reflect.ValueOf(param[i-1])
			} else {
				params[i] = reflect.Zero(otherIns[i-1])

			}
		}
	}

	// fmt.Println(method, mtype)
	// fmt.Println("inum:", inum)
	methodValue.Call(params)

}

type Event interface {

	//事件类型
	GetType() int
	//事件触发目标
	GetTarget() interface{}
}

type BaseEvent struct {
	etype  int
	target interface{}
}

func (be *BaseEvent) Init(etype int, target interface{}) {
	be.etype = etype
	be.target = target
}
func (be *BaseEvent) GetType() int {
	return be.etype
}
func (be *BaseEvent) GetTarget() interface{} {
	return be.target
}

type TriggerEvent struct {
	taskpool.BaseTask
	event    Event
	listener *RegListen
	method   string
	param    []interface{}
}

var triggerEventPool = sync.Pool{
	New: func() interface{} {
		return &TriggerEvent{}
	},
}

func createTriggerEvent(method string, event Event, listener *RegListen, param []interface{}) *TriggerEvent {
	e := triggerEventPool.Get().(*TriggerEvent)
	e.event = event
	e.listener = listener
	e.method = method
	e.param = param
	return e
}

type Listener interface {
}

type EventTaskExecor struct {
}

func (e *EventTaskExecor) ExecTask(task taskpool.Task) error {
	even, ok := task.(*TriggerEvent)
	if !ok {
		return fmt.Errorf("错误的任务:不是一个事件任务.", task)
	}
	even.listener.exec(even.event, even.method, even.param)
	//doTriggerEvent(even)
	triggerEventPool.Put(even)
	return nil
}

var eventTaskExecor taskpool.AsyncTaskOperater

func getEventTaskExecor() taskpool.AsyncTaskOperater {
	if eventTaskExecor == nil {
		eventTaskExecor = taskpool.CreateAsyncTaskOperater("事件任务", &EventTaskExecor{}, &taskpool.AsyncTaskOption{
			AsyncMaxWaitTaskNum: 1024,
			//最大异步任务go程数量
			MaxAsyncPoolNum: 1024,
			MinAsyncPoolNum: 32,
			//最大空闲时间
			AsyncPoolIdelTime: 10 * time.Second,
			//任务最大失败次数
			AsyncTaskMaxFaildedNum: 1,
		})
	}
	return eventTaskExecor
}

//监听组 该组是线程不安全的
type ListenerGroup struct {
	listenerList      *list.List
	listenerMap       map[Listener]*list.Element
	asyncTaskOperater taskpool.AsyncTaskOperater
}

func CreateListenerGroup() *ListenerGroup {
	group := &ListenerGroup{}
	group.listenerList = list.New()
	group.listenerMap = make(map[Listener]*list.Element, 2)
	group.asyncTaskOperater = getEventTaskExecor()
	return group
}

//向组内增加监听
func (l *ListenerGroup) AddListener(listener Listener) {
	if l.listenerList == nil {
		l.listenerList = list.New()
		l.listenerMap = make(map[Listener]*list.Element, 2)
		l.listenerMap[listener] = l.listenerList.PushFront(createRegListen(listener))
		return
	}
	_, ok := l.checkElement(listener)
	if !ok {
		l.listenerMap[listener] = l.listenerList.PushFront(createRegListen(listener))
	}

}

//从组内移除监听
func (l *ListenerGroup) RemoveListener(listener Listener) {
	elem, ok := l.checkElement(listener)
	if ok {
		l.listenerList.Remove(elem)
		delete(l.listenerMap, listener)
	}
}

//触发监听事件
func (l *ListenerGroup) FireListener(method string, event Event, param []interface{}) {
	for elem := l.listenerList.Front(); elem != nil; elem = elem.Next() {
		listener := elem.Value.(*RegListen)
		l.asyncTaskOperater.ExecAsyncTask(createTriggerEvent(method, event, listener, param))

	}

}

//检测组内是否有该监听
func (l *ListenerGroup) checkElement(listener Listener) (*list.Element, bool) {
	if l.listenerMap == nil {
		return nil, false
	}
	elem, ok := l.listenerMap[listener]
	return elem, ok

}

//事件发生器 线程安全的
type EventGenerator struct {
	sync.Mutex
	listenergroupMap map[int]*ListenerGroup
}

//增加一个监听器
func (eg *EventGenerator) AddListener(eventtype int, listener Listener) {
	eg.Lock()
	defer eg.Unlock()
	listenerGroup := eg.tryInit(eventtype)
	listenerGroup.AddListener(listener)
}

func (eg *EventGenerator) RemoveListener(eventtype int, listener Listener) {
	eg.Lock()
	defer eg.Unlock()
	if eg.listenergroupMap == nil {
		return
	}
	listenerGroup, ok := eg.listenergroupMap[eventtype]
	if !ok {
		return
	}
	listenerGroup.RemoveListener(listener)
}
func (eg *EventGenerator) FireListener(eventtype int, method string, event Event, param ...interface{}) {
	//eg.Lock()
	//defer eg.Unlock()
	lg, ok := eg.listenergroupMap[eventtype]
	if ok {
		lg.FireListener(method, event, param)
	}

}

//尝试初始化,如果已经初始化不在初始化
func (eg *EventGenerator) tryInit(eventtype int) *ListenerGroup {
	if eg.listenergroupMap == nil {

		eg.listenergroupMap = make(map[int]*ListenerGroup)
		lg := CreateListenerGroup()
		eg.listenergroupMap[eventtype] = lg
		return lg

	}
	listenerGroup, ok := eg.listenergroupMap[eventtype]
	if !ok {

		listenerGroup = CreateListenerGroup()
		eg.listenergroupMap[eventtype] = listenerGroup

	}
	return listenerGroup
}
