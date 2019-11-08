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
	if listeners, ok := e.fs[eventName]; ok {
		e.fs[eventName] = append(listeners, fs...)
	} else {
		e.fs[eventName] = fs
	}
	return e
}

func (e *EventEmitter) Off(eventName string, n uint) *EventEmitter {
	e.lk.Lock()
	defer e.lk.Unlock()
	newFs := make([]func(...interface{}), 0)
	end := uint(len(e.fs))
	events := e.fs[eventName]
	for i := uint(0); i < end; i++ {
		if i != n {
			newFs = append(newFs, events[i])
		}
	}
	e.fs[eventName] = newFs
	return e
}

func (e *EventEmitter) OffAll(eventName string) *EventEmitter {
	e.lk.Lock()
	defer e.lk.Unlock()
	if _, ok := e.fs[eventName]; ok {
		e.fs[eventName] = []func(...interface{}){}
	}
	return e
}

func (e *EventEmitter) Clone() *EventEmitter {
	newM := make(map[string][]func(...interface{}), len(e.fs))
	for k, v := range e.fs {
		newM[k] = make([]func(...interface{}), len(v))
		for idx, el := range v {
			newM[k][idx] = el
		}
	}
	return &EventEmitter{
		lk: &sync.Mutex{},
		fs: newM,
	}
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
