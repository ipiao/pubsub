package pubsub

import (
	"fmt"
	"log"
)

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

type multiError []interface{}

func (me *multiError) Add(err interface{}) {
	if err == nil {
		return
	}
	if me == nil {
		me = &multiError{}
	}
	ne := append(*me, err)
	me = &ne
	log.Println(me)
}

func (me *multiError) Error() string {
	log.Println(222, *me)
	return fmt.Sprint(*me...)
}

type handler func(interface{}) interface{}

// handle 是一个订阅者处理一个主题事件的集
type handle struct {
	topic   string
	asny    bool    // 是否同步，默认是异步的
	handler handler // 处理方法
}

// Subber 订阅者
type Subber struct {
	topics   []string           // 表示订阅的主题
	handlers map[string]*handle // 每个主题的处理方法
	ch       chan message       // 异步处理通道
}

type message struct {
	topic string
	msg   interface{}
}

func (s *Subber) run() {
	for mes := range s.ch {
		s.handlers[mes.topic].handler(mes.msg)
	}
}

func (s *Subber) setHandle(handle *handle) {
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

// // Sub 添加一个订阅者
// func (ps *PubSub) Sub(topics ...string) *Subber {

// }

// NewSub 初始化一个订阅者
func (ps *PubSub) NewSub(handles ...handle) *Subber {
	subber := &Subber{
		ch:       make(chan message, ps.capacity),
		handlers: make(map[string]*handle),
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
func (ps *PubSub) Pub(topic string, msg interface{}) error {
	return ps.send(topic, msg)
}

// NewSub 初始化一个订阅者
func (ps *PubSub) send(topic string, msg interface{}) error {
	var err *multiError
	for _, subber := range ps.topics[topic] {
		handle := subber.handlers[topic]
		if !handle.asny {
			err.Add(handle.handler(msg))
		} else {
			message := message{topic: topic, msg: msg}
			subber.ch <- message
		}
	}
	return err
}
