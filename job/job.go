package job

type AsyncJob struct {
	f func(...interface{}) interface{}
}

func Await(js ...*AsyncJob) *AsyncJob {
	return NewClojure(func() interface{} {
		return AwaitRunning(Dispatch(js...)...)
	})
}

func Consume(js ...*AsyncJob) {
	ConsumeRunning(Dispatch(js...)...)
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
	return New(func(i ...interface{}) interface{} { return f() })
}

func (j *AsyncJob) Start(args ...interface{}) *Running {
	p := newRunning()
	p.lk.Lock()
	go func() {
		p.val = j.f(args...)
		p.lk.Unlock()
	}()
	return p
}

func (j *AsyncJob) String() string {
	return "AsyncJob"
}
