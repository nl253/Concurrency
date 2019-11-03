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
