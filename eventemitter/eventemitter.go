package eventemitter

import (
	"fmt"
	"strings"
	"sync"
)

type EventEmitter struct {
	lk *sync.Mutex
	fs map[string][]func(...interface{})
}

func New() *EventEmitter {
	return &EventEmitter{fs: map[string][]func(...interface{}){}, lk: &sync.Mutex{}}
}

func (e *EventEmitter) Emit(eventName string, args ...interface{}) *EventEmitter {
	e.lk.Lock()
	defer e.lk.Unlock()
	if listeners, ok := e.fs[eventName]; ok {
		for _, l := range listeners {
			l(args...)
		}
	}
	return e
}

func (e *EventEmitter) On(eventName string, fs ...func(...interface{})) *EventEmitter {
	e.lk.Lock()
	defer e.lk.Unlock()
	listeners, ok := e.fs[eventName]
	if ok {
		e.fs[eventName] = append(listeners, fs...)
	} else {
		e.fs[eventName] = fs
	}
	return e
}

func (e *EventEmitter) Off(eventName string) *EventEmitter {
	e.lk.Lock()
	defer e.lk.Unlock()
	if _, ok := e.fs[eventName]; ok {
		e.fs[eventName] = []func(...interface{}){}
	}
	return e
}

func (e *EventEmitter) Eq(x interface{}) bool {
	switch x.(type) {
	case *EventEmitter:
		return x.(*EventEmitter) == e
	default:
		return false
	}
}

func (e *EventEmitter) String() string {
	sb := strings.Builder{}
	sb.WriteString("EventEmitter {")
	sb.WriteByte(' ')
	for k, v := range e.fs {
		sb.WriteString(k)
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%d", len(v)))
		sb.WriteByte(' ')
	}
	sb.WriteRune('}')
	return sb.String()
}
