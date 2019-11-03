package job

import (
	"fmt"
	"sync"
)

type AsyncJob struct {
	running bool
	lk      *sync.Mutex
	f       func(...interface{}) interface{}
	val     interface{}
}

func All(js ...*AsyncJob) *AsyncJob {
	return NewClojure(func() interface{} {
		n := uint(len(js))
		xs := make([]interface{}, n)
		for i := uint(0); i < n; i++ {
			xs[i] = js[i].Await()
		}
		return xs
	})
}

func Await(js ...*AsyncJob) {
	NewClojure(func() interface{} {
		n := uint(len(js))
		for i := uint(0); i < n; i++ {
			js[i].Await()
		}
		return nil
	}).Await()
}

func Race(js ...*AsyncJob) *AsyncJob {
	return NewClojure(func() interface{} {
		var result interface{} = nil
		done := false
		n := uint(len(js))
		ws := make([]*AsyncJob, len(js), len(js))
		for idx := uint(0); idx < n; idx++ {
			ws[idx] = NewConsumer(func(args ...interface{}) {
				if !done {
					result = js[args[0].(uint)].Start().Await()
					done = true
				}
			}).Start(idx)
		}
		Await(ws...)
		return result
	})
}

func New(f func(...interface{}) interface{}) *AsyncJob {
	return &AsyncJob{
		running: false,
		lk:      &sync.Mutex{},
		f:       f,
		val:     nil,
	}
}

func NewConsumer(f func(...interface{})) *AsyncJob {
	return New(func(i ...interface{}) interface{} {
		f(i...)
		return nil
	})
}

func NewClojure(f func() interface{}) *AsyncJob {
	return New(func(i ...interface{}) interface{} {
		return f()
	})
}

func (j *AsyncJob) Start(args ...interface{}) *AsyncJob {
	j.lk.Lock()
	j.running = true
	go func() {
		j.val = j.f(args...)
		j.running = false
		j.lk.Unlock()
	}()
	return j
}

func (j *AsyncJob) Await() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	return j.val
}

func (j *AsyncJob) String() string {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.running {
		return "AsyncJob { running }"
	}
	return fmt.Sprintf("AsyncJob { result = %v (%T) }", j.val, j.val)
}
