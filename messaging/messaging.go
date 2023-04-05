package messaging

import (
	"github.com/go-sre/core/runtime"
)

const (
	StartupEvent  = "event:startup"
	ShutdownEvent = "event:shutdown"
	PingEvent     = "event:ping"
	StatusEvent   = "event:status"
	HostName      = "host"
)

// MessageHandler - function type to process a Message
type MessageHandler func(msg Message)

// Message - message access data
type Message struct {
	To      string
	From    string
	Event   string
	Status  *runtime.Status
	Content []any
	ReplyTo MessageHandler
}

// ReplyTo - function used by message recipient to reply with a runtime.Status
func ReplyTo(msg Message, status *runtime.Status) {
	if msg.ReplyTo == nil {
		return
	}
	msg.ReplyTo(Message{
		To:      msg.From,
		From:    msg.To,
		Event:   msg.Event,
		Status:  status,
		Content: nil,
		ReplyTo: nil,
	})
}

// NewMessageCacheHandler - handler to receive messages into a cache.
func NewMessageCacheHandler(cache *MessageCache) MessageHandler {
	return func(msg Message) {
		cache.Add(msg)
	}
}
