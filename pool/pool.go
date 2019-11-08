package pool

import (
	"fmt"

	"github.com/nl253/Concurrency/job"
	"github.com/nl253/Concurrency/sema"
)

type Pool struct {
	sema *sema.Sema
	f    func(...interface{}) interface{}
}

func New(n uint) *Pool {
	return &Pool{sema: sema.New(n)}
}

func (p *Pool) Submit(j *job.AsyncJob) *job.Running {
	p.sema.Acquire()
	return job.NewClojure(func() interface{} {
		result := j.Start().Await()
		p.sema.Release()
		return result
	}).Start()
}

func (p *Pool) Clone() *Pool {
	return &Pool{
		sema: p.sema.Clone(),
		f:    p.f,
	}
}

func (p *Pool) String() string {
	return fmt.Sprint("Pool { sema = %u }", p.sema.NPermits())
}
