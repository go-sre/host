package messaging

import (
	"errors"
	"fmt"
	"github.com/go-sre/core/runtime"
	"time"
)

type messageMap map[string]Message

var startupLocation = pkgPath + "/startup"

var directory = NewEntryDirectory()

// RegisterResource - function to register a resource uri
func RegisterResource(uri string, c chan Message) error {
	if uri == "" {
		return errors.New("invalid argument: uri is empty")
	}
	if c == nil {
		return errors.New(fmt.Sprintf("invalid argument: channel is nil for [%v]", uri))
	}
	registerResourceUnchecked(uri, c)
	return nil
}

func registerResourceUnchecked(uri string, c chan Message) error {
	directory.Add(uri, c)
	return nil
}

// Shutdown - resource shutdown
func Shutdown() {
	directory.Shutdown()
}

// Startup - templated function to start all registered resources.
func Startup[E runtime.ErrorHandler, O runtime.OutputHandler](duration time.Duration, content ContentMap) (status *runtime.Status) {
	var e E
	var failures []string
	var count = directory.Count()

	if count == 0 {
		return runtime.NewStatusOK()
	}
	cache := NewMessageCache()
	toSend := createToSend(content, NewMessageCacheHandler(cache))
	sendMessages(toSend)
	for wait := time.Duration(float64(duration) * 0.25); duration >= 0; duration -= wait {
		time.Sleep(wait)
		// Check for completion
		if cache.Count() < count {
			continue
		}
		// Check for failed resources
		failures = cache.Exclude(StartupEvent, runtime.StatusOK)
		if len(failures) == 0 {
			handleStatus[O](cache)
			return runtime.NewStatusOK()
		}
		break
	}
	Shutdown()
	if len(failures) > 0 {
		handleErrors[E](failures, cache)
		return runtime.NewStatusCode(runtime.StatusInternal)
	}
	return e.Handle(startupLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), directory.Count()))).SetCode(runtime.StatusDeadlineExceeded)
}

func createToSend(cm ContentMap, fn MessageHandler) messageMap {
	m := make(messageMap)
	for _, k := range directory.Uri() {
		msg := Message{To: k, From: HostName, Event: StartupEvent, Status: nil, ReplyTo: fn}
		if cm != nil {
			if content, ok := cm[k]; ok {
				msg.Content = append(msg.Content, content...)
			}
		}
		m[k] = msg
	}
	return m
}

func sendMessages(msgs messageMap) {
	for k := range msgs {
		directory.Send(msgs[k])
	}
}

func handleErrors[E runtime.ErrorHandler](failures []string, cache *MessageCache) {
	var e E
	for _, uri := range failures {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			e.HandleStatus(msg.Status)
		}
	}
}

func handleStatus[O runtime.OutputHandler](cache *MessageCache) {
	var o O
	for _, uri := range cache.Uri() {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			o.Write(fmt.Sprintf("startup successful for resource [%v] : %s", uri, msg.Status.Duration()))
		}
	}
}
