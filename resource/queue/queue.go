package queue

import "github.com/alileza/gebet/resource/queue/rabbitmq"

const Name = "queue"

type Client interface {
	Listen(target string) error
	Count(target string, count int) error
	Publish(target string, payload []byte) error
	Message(target string) []byte
}

func Cast(i interface{}) Client {
	return i.(Client)
}

func New(params map[string]string) Client {
	driver, ok := params["driver"]
	if !ok {
		panic("queue: driver is required")
	}
	switch driver {
	case "rabbitmq":
		return rabbitmq.New(params)
	}
	panic("queue: invalid driver > " + driver)
}
