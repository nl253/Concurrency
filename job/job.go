package job

import (
	"fmt"
	"sync"
)

type AsyncJob struct {
	f func(...interface{}) interface{}
}

type Running struct {
	lk  *sync.Mutex
	val interface{}
}

func Await(js ...*AsyncJob) *AsyncJob {
	return NewClojure(func() interface{} {
		return AwaitRunning(Dispatch(js...)...)
	})
}

func AwaitRunning(js ...*Running) []interface{} {
	n := uint(len(js))
	xs := make([]interface{}, n)
	for i := uint(0); i < n; i++ {
		xs[i] = js[i].Await()
	}
	return xs
}

func Consume(js ...*AsyncJob) {
	ConsumeRunning(Dispatch(js...)...)
}

func ConsumeRunning(js ...*Running) {
	n := uint(len(js))
	for i := uint(0); i < n; i++ {
		js[i].Await()
	}
}

func Dispatch(js ...*AsyncJob) []*Running {
	n := uint(len(js))
	ws := make([]*Running, n)
	for i := uint(0); i < n; i++ {
		ws[i] = js[i].Start()
	}
	return ws
}

func Race(js ...*AsyncJob) *AsyncJob {
	return NewClojure(func() interface{} {
		return RaceRunning(Dispatch(js...)...)
	})
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

func New(f func(...interface{}) interface{}) *AsyncJob {
	return &AsyncJob{f}
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

func (j *AsyncJob) Start(args ...interface{}) *Running {
	p := &Running{
		lk:  &sync.Mutex{},
		val: nil,
	}
	p.lk.Lock()
	go func() {
		p.val = j.f(args...)
		p.lk.Unlock()
	}()
	return p
}

func (j *Running) Await() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	return j.val
}

func (j *Running) String() string {
	j.lk.Lock()
	defer j.lk.Unlock()
	return fmt.Sprintf("AsyncJob { result = %v (%T) }", j.val, j.val)
}

func (j *AsyncJob) String() string {
	return "AsyncJob"
}
