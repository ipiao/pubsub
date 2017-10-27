package pubsub

import "encoding/json"

type operation int

const (
	sub operation = iota
	subOnce
	pub
	unsub
	unsubAll
	closeTopic
	shutdown
)

// MultiResult for result
type MultiResult map[string]interface{}

// Set a resu;t
func (me *MultiResult) Set(k string, r interface{}) {
	if r == nil {
		return
	}
	(*me)[k] = r
}

// Map return itself
func (me *MultiResult) Map() map[string]interface{} {
	return *me
}

// IsNil ..
func (me *MultiResult) IsNil() bool {
	return me == nil || len(*me) == 0
}

// MultiResult for err
func (me *MultiResult) Error() string {
	var errs = make(map[string]string)
	for k, v := range *me {
		if err, ok := v.(error); ok {
			errs[k] = err.Error()
		}
	}
	bs, _ := json.Marshal(errs)
	return string(bs)
}

type handler func(interface{}) interface{}

// Handle 是一个订阅者处理一个主题事件的集
type Handle struct {
	topic   string
	asny    bool    // 是否同步，默认是异步的
	handler handler // 处理方法
}

// NewHandle create a handle
func NewHandle(topic string, asny bool, handler handler) Handle {
	return Handle{topic, asny, handler}
}

// Subber 订阅者
type Subber struct {
	name     string
	topics   []string           // 表示订阅的主题
	handlers map[string]*Handle // 每个主题的处理方法
	ch       chan Message       // 异步处理通道
}

// Message ..
type Message struct {
	Topic string
	Data  interface{}
}

func (s *Subber) run() {
	for mes := range s.ch {
		s.handlers[mes.Topic].handler(mes)
	}
}

func (s *Subber) setHandle(handle *Handle) {
	if _, ok := s.handlers[handle.topic]; !ok {
		s.topics = append(s.topics, handle.topic)
	}
	s.handlers[handle.topic] = handle
}

// PubSub 是订阅-发布集
type PubSub struct {
	topics map[string][]*Subber
	// revTopics map[*Subber][]string
	capacity int
}

// New 创建一个 pub-sub 体系
func New(capacity int) *PubSub {
	return &PubSub{
		topics:   make(map[string][]*Subber),
		capacity: capacity,
	}
}

// Sub 添加一个订阅者
func (ps *PubSub) Sub(name string, asny bool, handler handler, topics ...string) *Subber {
	var handles []Handle
	for _, topic := range topics {
		handles = append(handles, Handle{
			topic:   topic,
			asny:    asny,
			handler: handler,
		})
	}
	return ps.InitSub(name, handles...)
}

// InitSub 初始化一个订阅者
func (ps *PubSub) InitSub(name string, handles ...Handle) *Subber {
	subber := &Subber{
		name:     name,
		ch:       make(chan Message, ps.capacity),
		handlers: make(map[string]*Handle),
	}
	for i := range handles {
		handle := &handles[i]
		subber.setHandle(handle)
		ps.topics[handle.topic] = append(ps.topics[handle.topic], subber)
		//	ps.revTopics[subber] = append(ps.topics[handle.topic], subber)
	}
	go subber.run()
	return subber
}

// Pub 发布一个消息
func (ps *PubSub) Pub(topic string, msg interface{}) *MultiResult {
	return ps.send(topic, msg)
}

// NewSub 初始化一个订阅者
func (ps *PubSub) send(topic string, msg interface{}) *MultiResult {
	var res = &MultiResult{}
	for _, subber := range ps.topics[topic] {
		handle := subber.handlers[topic]
		message := Message{Topic: topic, Data: msg}
		if !handle.asny {
			res.Set(subber.name, handle.handler(message))
		} else {
			subber.ch <- message
		}
	}
	return res
}
