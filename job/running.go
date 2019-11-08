package job

import (
	"fmt"
	"sync"
)

type Running struct {
	lk  *sync.Mutex
	val interface{}
}

func newRunning() *Running {
	return &Running{
		lk:  &sync.Mutex{},
		val: nil,
	}
}

func AwaitRunning(js ...*Running) []interface{} {
	n := uint(len(js))
	xs := make([]interface{}, n)
	for i := uint(0); i < n; i++ {
		xs[i] = js[i].Await()
	}
	return xs
}

func ConsumeRunning(js ...*Running) {
	n := uint(len(js))
	for i := uint(0); i < n; i++ {
		js[i].Consume()
	}
}

func RaceRunning(js ...*Running) *Running {
	return NewClojure(func() interface{} {
		var result interface{} = nil
		done := false
		n := uint(len(js))
		ws := make([]*Running, n, n)
		for idx := uint(0); idx < n; idx++ {
			ws[idx] = NewConsumer(func(args ...interface{}) {
				if !done {
					result = js[args[0].(uint)].Await()
					done = true
				}
			}).Start(idx)
		}
		ConsumeRunning(ws...)
		return result
	}).Start()
}

func (j *Running) Await() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	return j.val
}

func (j *Running) Consume() {
	j.lk.Lock()
	j.lk.Unlock()
}

func (j *Running) String() string {
	j.lk.Lock()
	defer j.lk.Unlock()
	return fmt.Sprintf("AsyncJob { result = %v (%T) }", j.val, j.val)
}
