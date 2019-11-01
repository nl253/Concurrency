package actor

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var RT = &sync.Map{}

type Actor struct {
	id   uuid.UUID
	done bool
	f    func(msgs <-chan interface{})
	msgs chan interface{}
}

func Async(f func(msgs <-chan interface{}), n uint) *Actor {
	id := uuid.New()
	j := &Actor{
		id:   id,
		f:    f,
		done: false,
		msgs: make(chan interface{}, n),
	}
	RT.Store(id, j)
	return j
}

func Sync(f func(msgs <-chan interface{})) *Actor {
	return Async(f, 0)
}

func (actor *Actor) Start() *Actor {
	if !actor.done {
		go func() {
			actor.f(actor.msgs)
			actor.done = true
			RT.Delete(actor.id)
		}()
	} else {
	}
	return actor
}

func (actor *Actor) Msgs() chan<- interface{} {
	return actor.msgs
}

func (actor *Actor) Id() [16]byte {
	return actor.id
}

func (actor *Actor) Done() bool {
	return actor.done
}

func (actor *Actor) String() string {
	info := "done"
	if !actor.done {
		info = "running"
	}
	return fmt.Sprintf("Actor [%s] { id = %s, #msgs = %d }", info, actor.id.String(), len(actor.msgs))
}
