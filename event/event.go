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

func createRegListen(listener Listener) *regListen {
	return &regListen{
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

type regListen struct {
	listener Listener
	value    reflect.Value
	methods  cmap.ConcurrencyMap
}

func (l *regListen) exec(event Event, method string, param []interface{}) {
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

// Event 事件接口
type Event interface {

	//事件类型
	GetType() int
	//事件触发目标
	GetTarget() interface{}
}

// BaseEvent 事件基础实现
type BaseEvent struct {
	etype  int
	target interface{}
}

// Init 初始化事件
// etype 事件类型
// target 触发对象
func (be *BaseEvent) Init(etype int, target interface{}) {
	be.etype = etype
	be.target = target
}

// GetType 获取事件类型
func (be *BaseEvent) GetType() int {
	return be.etype
}

// GetTarget 获取触发对象
func (be *BaseEvent) GetTarget() interface{} {
	return be.target
}

type triggerEvent struct {
	taskpool.BaseTask
	event    Event
	listener *regListen
	method   string
	param    []interface{}
}

var triggerEventPool = sync.Pool{
	New: func() interface{} {
		return &triggerEvent{}
	},
}

func createTriggerEvent(method string, event Event, listener *regListen, param []interface{}) *triggerEvent {
	e := triggerEventPool.Get().(*triggerEvent)
	e.event = event
	e.listener = listener
	e.method = method
	e.param = param
	return e
}

// Listener 监听接口
type Listener interface {
}

type eventTaskExecor struct {
}

func (e *eventTaskExecor) ExecTask(task taskpool.Task) error {
	even, ok := task.(*triggerEvent)
	if !ok {
		return fmt.Errorf("错误的任务:%v", task)
	}
	even.listener.exec(even.event, even.method, even.param)
	//doTriggerEvent(even)
	triggerEventPool.Put(even)
	return nil
}

var eventTaskOperate taskpool.AsyncTaskOperater

func getEventTaskExecor() taskpool.AsyncTaskOperater {
	if eventTaskOperate == nil {
		eventTaskOperate = taskpool.CreateAsyncTaskOperater("事件任务", &eventTaskExecor{}, &taskpool.AsyncTaskOption{
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
	return eventTaskOperate
}

//监听组 该组是线程不安全的
type listenerGroup struct {
	listenerList      *list.List
	listenerMap       map[Listener]*list.Element
	asyncTaskOperater taskpool.AsyncTaskOperater
}

func createListenerGroup() *listenerGroup {
	group := &listenerGroup{}
	group.listenerList = list.New()
	group.listenerMap = make(map[Listener]*list.Element, 2)
	group.asyncTaskOperater = getEventTaskExecor()
	return group
}

//向组内增加监听
func (l *listenerGroup) AddListener(listener Listener) {
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
func (l *listenerGroup) RemoveListener(listener Listener) {
	elem, ok := l.checkElement(listener)
	if ok {
		l.listenerList.Remove(elem)
		delete(l.listenerMap, listener)
	}
}

//FireListener 触发监听事件
func (l *listenerGroup) FireListener(method string, event Event, param []interface{}) {
	for elem := l.listenerList.Front(); elem != nil; elem = elem.Next() {
		listener := elem.Value.(*regListen)
		l.asyncTaskOperater.ExecAsyncTask(createTriggerEvent(method, event, listener, param))

	}

}

//检测组内是否有该监听
func (l *listenerGroup) checkElement(listener Listener) (*list.Element, bool) {
	if l.listenerMap == nil {
		return nil, false
	}
	elem, ok := l.listenerMap[listener]
	return elem, ok

}

// Generator 事件发生器 线程安全的
type Generator struct {
	sync.Mutex
	listenergroupMap map[int]*listenerGroup
}

// AddListener 增加一个监听器
// eventtype 事件类型
// listener 监听器
func (eg *Generator) AddListener(eventtype int, listener Listener) {
	eg.Lock()
	defer eg.Unlock()
	listenerGroup := eg.tryInit(eventtype)
	listenerGroup.AddListener(listener)
}

// RemoveListener 移除一个监听
// eventtype 事件类型
// listener 监听器
func (eg *Generator) RemoveListener(eventtype int, listener Listener) {
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

// FireListener 触发事件
// eventtype 事件类型
// method 触发的方法
// event 事件
// param 事件自定义参数
func (eg *Generator) FireListener(eventtype int, method string, event Event, param ...interface{}) {
	//eg.Lock()
	//defer eg.Unlock()
	lg, ok := eg.listenergroupMap[eventtype]
	if ok {
		lg.FireListener(method, event, param)
	}

}

//尝试初始化,如果已经初始化不在初始化
func (eg *Generator) tryInit(eventtype int) *listenerGroup {
	if eg.listenergroupMap == nil {

		eg.listenergroupMap = make(map[int]*listenerGroup)
		lg := createListenerGroup()
		eg.listenergroupMap[eventtype] = lg
		return lg

	}
	listenerGroup, ok := eg.listenergroupMap[eventtype]
	if !ok {

		listenerGroup = createListenerGroup()
		eg.listenergroupMap[eventtype] = listenerGroup

	}
	return listenerGroup
}
