package messaging

import (
	"fmt"
	"time"
)

/*
func createTestEntry(uri string, status int32) *entry {
	entry := createEntry(uri, nil)
	entry.msgs.add(CreateMessage(VirtualHost, VirtualHost, StartupEvent, status, nil))
	return entry
}
*/
var directoryTest = NewEntryDirectory()

func ExampleEntryDirectory_Add() {
	uri := "urn:test"
	uri2 := "urn:test:two"

	directoryTest.Empty()

	fmt.Printf("test: count() -> : %v\n", directoryTest.Count())
	d2 := directoryTest.Get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	directoryTest.Add(uri, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri)
	fmt.Printf("test: count() -> : %v\n", directoryTest.Count())
	d2 = directoryTest.Get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	directoryTest.Add(uri2, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri2)
	fmt.Printf("test: count() -> : %v\n", directoryTest.Count())
	d2 = directoryTest.Get(uri2)
	fmt.Printf("test: get(%v) -> : %v\n", uri2, d2)

	fmt.Printf("test: uri() -> : %v\n", directoryTest.Uri())

	//Output:
	//test: count() -> : 0
	//test: get(urn:test) -> : <nil>
	//test: add(urn:test) -> : ok
	//test: count() -> : 1
	//test: get(urn:test) -> : &{urn:test <nil>}
	//test: add(urn:test:two) -> : ok
	//test: count() -> : 2
	//test: get(urn:test:two) -> : &{urn:test:two <nil>}
	//test: uri() -> : [urn:test urn:test:two]

}

func ExampleEntryDirectory_SendError() {
	uri := "urn:test"
	directoryTest.Empty()

	fmt.Printf("test: send(%v) -> : %v\n", uri, directoryTest.Send(Message{To: uri}))

	directoryTest.Add(uri, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri)
	fmt.Printf("test: send(%v) -> : %v\n", uri, directoryTest.Send(Message{To: uri}))

	//Output:
	//test: send(urn:test) -> : entry not found: [urn:test]
	//test: add(urn:test) -> : ok
	//test: send(urn:test) -> : entry channel is nil: [urn:test]

}

func ExampleEntryDirectory_Send() {
	uri1 := "urn:test-1"
	uri2 := "urn:test-2"
	uri3 := "urn:test-3"
	c := make(chan Message, 16)
	directoryTest.Empty()

	directoryTest.Add(uri1, c)
	directoryTest.Add(uri2, c)
	directoryTest.Add(uri3, c)

	directoryTest.Send(Message{To: uri1, From: HostName, Event: StartupEvent})
	directoryTest.Send(Message{To: uri2, From: HostName, Event: StartupEvent})
	directoryTest.Send(Message{To: uri3, From: HostName, Event: StartupEvent})

	time.Sleep(time.Second * 1)
	resp1 := <-c
	resp2 := <-c
	resp3 := <-c
	fmt.Printf("test: <- c -> : [%v] [%v] [%v]\n", resp1.To, resp2.To, resp3.To)
	close(c)

	//Output:
	//test: <- c -> : [urn:test-1] [urn:test-2] [urn:test-3]

}
