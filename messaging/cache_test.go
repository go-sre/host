package messaging

import (
	"fmt"
	"github.com/gotemplates/core/runtime"
	"net/http"
)

func ExampleMessageCache_Add() {
	resp := NewMessageCache()

	resp.Add(Message{To: "to-uri", From: "from-uri-0", Event: StartupEvent, Status: runtime.NewStatusCode(runtime.StatusNotProvided)})
	resp.Add(Message{To: "to-uri", From: "from-uri-1", Event: StartupEvent, Status: runtime.NewStatusOK()})
	resp.Add(Message{To: "to-uri", From: "from-uri-2", Event: PingEvent, Status: runtime.NewStatusCode(runtime.StatusNotProvided)})
	resp.Add(Message{To: "to-uri", From: "from-uri-3", Event: PingEvent, Status: runtime.NewStatusCode(runtime.StatusNotProvided)})
	resp.Add(Message{To: "to-uri", From: "from-uri-4", Event: PingEvent, Status: runtime.NewHttpStatusCode(http.StatusOK)})

	fmt.Printf("test: count() -> : %v\n", resp.Count())

	m, err := resp.Get("invalid")
	fmt.Printf("test: Get(%v) -> : [error:%v] [msg:%v]\n", "invalid", err, m)

	m, err = resp.Get("from-uri-3")
	fmt.Printf("test: Get(%v) -> : [error:%v] [msg:%v]\n", "from-uri-3", err, m)

	fmt.Printf("test: include(%v,%v) -> : %v\n", ShutdownEvent, runtime.StatusNotProvided, resp.Include(ShutdownEvent, runtime.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", ShutdownEvent, runtime.StatusNotProvided, resp.Exclude(ShutdownEvent, runtime.StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", StartupEvent, runtime.StatusNotProvided, resp.Include(StartupEvent, runtime.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", StartupEvent, runtime.StatusNotProvided, resp.Exclude(StartupEvent, runtime.StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", PingEvent, runtime.StatusNotProvided, resp.Include(PingEvent, runtime.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", PingEvent, runtime.StatusNotProvided, resp.Exclude(PingEvent, runtime.StatusNotProvided))

	//Output:
	//test: count() -> : 5
	//test: Get(invalid) -> : [error:invalid argument: uri not found [invalid]] [msg:{   <nil> [] <nil>}]
	//test: Get(from-uri-3) -> : [error:<nil>] [msg:{to-uri from-uri-3 event:ping Not provided [] <nil>}]
	//test: include(event:shutdown,Code(93)) -> : []
	//test: exclude(event:shutdown,Code(93)) -> : [from-uri-0 from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:startup,Code(93)) -> : [from-uri-0]
	//test: exclude(event:startup,Code(93)) -> : [from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:ping,Code(93)) -> : [from-uri-2 from-uri-3]
	//test: exclude(event:ping,Code(93)) -> : [from-uri-0 from-uri-1 from-uri-4]

}

func ExampleMessageCache_Uri() {
	resp := NewMessageCache()

	resp.Add(Message{To: "to-uri", From: "from-uri-0", Event: StartupEvent, Status: runtime.NewStatusCode(runtime.StatusNotProvided)})
	resp.Add(Message{To: "to-uri", From: "from-uri-1", Event: StartupEvent, Status: runtime.NewStatusOK()})
	resp.Add(Message{To: "to-uri", From: "from-uri-2", Event: PingEvent, Status: runtime.NewStatusCode(runtime.StatusNotProvided)})
	resp.Add(Message{To: "to-uri", From: "from-uri-3", Event: PingEvent, Status: runtime.NewStatusCode(runtime.StatusNotProvided)})
	resp.Add(Message{To: "to-uri", From: "from-uri-4", Event: PingEvent, Status: runtime.NewHttpStatusCode(http.StatusOK)})

	fmt.Printf("test: count() -> : %v\n", resp.Count())

	m, err := resp.Get("invalid")
	fmt.Printf("test: Get(%v) -> : [error:%v] [msg:%v]\n", "invalid", err, m)

	m, err = resp.Get("from-uri-3")
	fmt.Printf("test: Get(%v) -> : [error:%v] [msg:%v]\n", "from-uri-3", err, m)

	fmt.Printf("test: include(%v,%v) -> : %v\n", ShutdownEvent, runtime.StatusNotProvided, resp.Include(ShutdownEvent, runtime.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", ShutdownEvent, runtime.StatusNotProvided, resp.Exclude(ShutdownEvent, runtime.StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", StartupEvent, runtime.StatusNotProvided, resp.Include(StartupEvent, runtime.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", StartupEvent, runtime.StatusNotProvided, resp.Exclude(StartupEvent, runtime.StatusNotProvided))

	fmt.Printf("test: include(%v,%v) -> : %v\n", PingEvent, runtime.StatusNotProvided, resp.Include(PingEvent, runtime.StatusNotProvided))
	fmt.Printf("test: exclude(%v,%v) -> : %v\n", PingEvent, runtime.StatusNotProvided, resp.Exclude(PingEvent, runtime.StatusNotProvided))

	//Output:
	//test: count() -> : 5
	//test: Get(invalid) -> : [error:invalid argument: uri not found [invalid]] [msg:{   <nil> [] <nil>}]
	//test: Get(from-uri-3) -> : [error:<nil>] [msg:{to-uri from-uri-3 event:ping Not provided [] <nil>}]
	//test: include(event:shutdown,Code(93)) -> : []
	//test: exclude(event:shutdown,Code(93)) -> : [from-uri-0 from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:startup,Code(93)) -> : [from-uri-0]
	//test: exclude(event:startup,Code(93)) -> : [from-uri-1 from-uri-2 from-uri-3 from-uri-4]
	//test: include(event:ping,Code(93)) -> : [from-uri-2 from-uri-3]
	//test: exclude(event:ping,Code(93)) -> : [from-uri-0 from-uri-1 from-uri-4]

}
