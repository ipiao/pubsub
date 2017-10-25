package pubsub

import (
	"log"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ps := New(5)
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
			return "error1"
		}},
		handle{topic: "4", asny: false, handler: func(interface{}) interface{} {
			log.Println("topic 4 from handles1")
			return "error2"
		}},
	}
	ps.NewSub(handles1...)

	t.Log(ps.Pub("1", ""))
	t.Log(ps.Pub("2", ""))
	t.Log(ps.Pub("3", ""))
	t.Log(ps.Pub("4", ""))

	time.Sleep(time.Second * 1)
}
