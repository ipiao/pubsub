package pubsub

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ps := New(1)
	handles1 := []handle{
		handle{topic: "1", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 1 from handles1")
			return nil
		}},
		handle{topic: "2", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 2 from handles1")
			return nil
		}},
		handle{topic: "3", asny: true, handler: func(interface{}) interface{} {
			log.Println("topic 3 from handles1")
			return errors.New("hello1")
		}},
		handle{topic: "4", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 4 from handles1")
			return errors.New("hello2")
		}},
	}

	handles2 := []handle{
		handle{topic: "1", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 1 from handles2")
			return nil
		}},
		handle{topic: "2", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 2 from handles2")
			return nil
		}},
		handle{topic: "3", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 3 from handles2")
			return errors.New("hello1")
		}},
		handle{topic: "4", asny: false, handler: func(msg interface{}) interface{} {
			mes := msg.(Message)
			log.Println(mes)
			log.Println("topic 4 from handles2")
			return "hello2"
		}},
	}

	ps.InitSub("sub1", handles1...)
	ps.InitSub("sub2", handles2...)

	t.Log(ps.Pub("1", "").Map())
	t.Log(ps.Pub("2", "").Map())
	t.Log(ps.Pub("3", "").Map())
	t.Log(ps.Pub("4", "this is data").Map())

	time.Sleep(time.Second * 1)
}
