package job

import (
	"fmt"
	"sync"
)

type AsyncJob struct {
	done bool
	lk   *sync.Mutex
	val  interface{}
}

func Function(f func(...interface{}) interface{}) *AsyncJob {
	return &AsyncJob{
		done: false,
		lk:   &sync.Mutex{},
		val:  f,
	}
}

func Consumer(f func(...interface{})) *AsyncJob {
	return Function(func(i ...interface{}) interface{} {
		f(i...)
		return nil
	})
}

func Clojure(f func() interface{}) *AsyncJob {
	return Function(func(i ...interface{}) interface{} {
		return f()
	})
}

func (j *AsyncJob) Start(args ...interface{}) *AsyncJob {
	j.lk.Lock()
	if !j.done {
		go func() {
			j.val = j.val.(func(...interface{}) interface{})(args...)
			j.done = true
			j.lk.Unlock()
		}()
	} else {
		j.lk.Unlock()
	}
	return j
}

func (j *AsyncJob) Await() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	return j.val
}

func (j *AsyncJob) Done() bool {
	return j.done
}

func (j *AsyncJob) String() string {
	return fmt.Sprintf("AsyncJob { done = %v, result = %v }", j.done, j.val)
}
