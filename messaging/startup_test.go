package messaging

import (
	"errors"
	"fmt"
	"github.com/go-sre/core/runtime"
	"time"
)

var credFn Credentials = func() (string, string, error) {
	return "", "", nil
}

var start time.Time

func ExampleCreateToSend() {
	none := "/resource/none"
	one := "/resource/one"

	registerResourceUnchecked(none, nil)
	registerResourceUnchecked(one, nil)

	m := createToSend(nil, nil)
	msg := m[none]
	fmt.Printf("test: createToSend(nil,nil) -> [to:%v] [from:%v]\n", msg.To, msg.From)

	cm := ContentMap{one: []any{credFn}}
	m = createToSend(cm, nil)
	msg = m[one]
	fmt.Printf("test: createToSend(map,nil) -> [to:%v] [from:%v] [credentials:%v]\n", msg.To, msg.From, AccessCredentials(&msg) != nil)

	//Output:
	//test: createToSend(nil,nil) -> [to:/resource/none] [from:host]
	//test: createToSend(map,nil) -> [to:/resource/one] [from:host] [credentials:true]

}

func ExampleStartup_Success() {
	uri1 := "urn:startup:good"
	uri2 := "urn:startup:bad"
	uri3 := "urn:startup:depends"

	start = time.Now()
	directory.Empty()

	c := make(chan Message, 16)
	RegisterResource(uri1, c)
	go startupGood(c)

	c = make(chan Message, 16)
	RegisterResource(uri2, c)
	go startupBad(c)

	c = make(chan Message, 16)
	RegisterResource(uri3, c)
	go startupDepends(c, nil)

	status := Startup[runtime.DebugError, runtime.NilOutput](time.Second*2, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//test: Startup() -> [OK]

}

func ExampleStartup_Failure() {
	uri1 := "urn:startup:good"
	uri2 := "urn:startup:bad"
	uri3 := "urn:startup:depends"

	start = time.Now()
	directory.Empty()

	c := make(chan Message, 16)
	RegisterResource(uri1, c)
	go startupGood(c)

	c = make(chan Message, 16)
	RegisterResource(uri2, c)
	go startupBad(c)

	c = make(chan Message, 16)
	RegisterResource(uri3, c)
	go startupDepends(c, errors.New("startup failure error message"))

	status := Startup[runtime.DebugError, runtime.StdOutput](time.Second*2, nil)

	fmt.Printf("test: Startup() -> [%v]\n", status)

	//Output:
	//[[] github.com/go-sre/host/messaging/startup [startup failure error message]]
	//test: Startup() -> [Internal]

}

func startupGood(c chan Message) {
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

func startupBad(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second + time.Millisecond*100)
			ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}

func startupDepends(c chan Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				ReplyTo(msg, runtime.NewStatusError(startupLocation, err).SetDuration(time.Since(start)))
			} else {
				time.Sleep(time.Second + (time.Millisecond * 900))
				ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
			}

		default:
		}
	}
}
