package messaging

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"sort"
	"sync"
)

// MessageCache - message cache by resource uri
type MessageCache struct {
	m  map[string]Message
	mu sync.RWMutex
}

// NewMessageCache - create a message cache
func NewMessageCache() *MessageCache {
	return &MessageCache{m: make(map[string]Message)}
}

func (r *MessageCache) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, _ = range r.m {
		count++
	}
	return count
}

func (r *MessageCache) Filter(event string, code codes.Code, include bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var uri []string
	for u, resp := range r.m {
		if include {
			if resp.Status != nil && resp.Status.Code() == code && resp.Event == event {
				uri = append(uri, u)
			}
		} else {
			if resp.Status != nil && resp.Status.Code() != code || resp.Event != event {
				uri = append(uri, u)
			}
		}
	}
	sort.Strings(uri)
	return uri
}

func (r *MessageCache) Include(event string, status codes.Code) []string {
	return r.Filter(event, status, true)
}

func (r *MessageCache) Exclude(event string, status codes.Code) []string {
	return r.Filter(event, status, false)
}

func (r *MessageCache) Add(msg Message) error {
	if msg.From == "" {
		return errors.New("invalid argument: message from is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[msg.From]; !ok {
		r.m[msg.From] = msg
	}
	return nil
}

func (r *MessageCache) Get(uri string) (Message, error) {
	if uri == "" {
		return Message{}, errors.New("invalid argument: uri is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[uri]; ok {
		return r.m[uri], nil
	}
	return Message{}, errors.New(fmt.Sprintf("invalid argument: uri not found [%v]", uri))
}

func (r *MessageCache) Uri() []string {
	var uri []string
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, _ := range r.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}
