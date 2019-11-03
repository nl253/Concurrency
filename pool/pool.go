package pool

import (
	"github.com/nl253/Concurrency/job"
	"github.com/nl253/Concurrency/sema"
)

type Pool struct {
	sema *sema.Sema
	f    func(...interface{}) interface{}
}

func New(n uint, f func(...interface{}) interface{}) *Pool {
	return &Pool{
		sema: sema.New(n),
		f:    f,
	}
}

func (p *Pool) Submit(args ...interface{}) *job.Running {
	p.sema.Acquire()
	return job.NewClojure(func() interface{} {
		result := p.f(args...)
		p.sema.Release()
		return result
	}).Start()
}

func (p *Pool) String() string {
	return "Pool"
}
