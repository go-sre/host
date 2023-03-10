package messaging

import (
	"errors"
	"fmt"
	"github.com/gotemplates/core/runtime"
	"time"
)

func ExamplePing() {
	uri1 := "urn:ping:good"
	uri2 := "urn:ping:bad"
	uri3 := "urn:ping:error"
	uri4 := "urn:ping:delay"

	start = time.Now()
	directory.Empty()

	c := make(chan Message, 16)
	RegisterResource(uri1, c)
	go pingGood(c)
	status := Ping[runtime.DebugError](nil, uri1)
	duration := status.Duration()
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan Message, 16)
	RegisterResource(uri2, c)
	go pingBad(c)
	status = Ping[runtime.DebugError](nil, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan Message, 16)
	RegisterResource(uri3, c)
	go pingError(c, errors.New("ping depends error message"))
	status = Ping[runtime.DebugError](nil, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan Message, 16)
	RegisterResource(uri4, c)
	go pingDelay(c)
	status = Ping[runtime.DebugError](nil, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status, duration)

	//Output:
	//test: Ping(good) -> [OK] [duration:0s]
	//[[] github.com/idiomatic-go/motif/messaging/ping [ping response time out: [urn:ping:bad]]]
	//test: Ping(bad) -> [DeadlineExceeded] [duration:0s]
	//test: Ping(error) -> [Internal [ping depends error message]] [duration:0s]
	//test: Ping(delay) -> [OK] [duration:0s]

}

func pingGood(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}

func pingBad(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second * 4)
			ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}

func pingError(c chan Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				ReplyTo(msg, runtime.NewStatusError(pingLocation, err).SetDuration(time.Since(start)))
			}
		default:
		}
	}
}

func pingDelay(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second)
			ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}
