package barrier

import (
	"sync"

	"github.com/nl253/Concurrency/job"
	"github.com/nl253/DataStructures/list"
)

type Barrier struct {
	n     uint
	l     *sync.Mutex
	wsLks *list.ConcurrentList
}

func New(n uint) *Barrier {
	return &Barrier{
		n:     n,
		l:     &sync.Mutex{},
		wsLks: list.New(),
	}
}

func NewFromJobs(js ...*job.AsyncJob) *Barrier {
	end := uint(len(js))
	b := New(end)
	for i := uint(0); i < end; i++ {
		b.Submit(js[i])
	}
	return b
}

func (s *Barrier) Wait() {
	for {
		s.l.Lock()
		if s.n == 0 {
			break
		}
		l := &sync.Mutex{}
		l.Lock()
		s.wsLks.Append(l)
		s.l.Unlock()
		l.Lock()
		l.Unlock()
	}
	s.l.Unlock()
}

func (s *Barrier) Submit(j *job.AsyncJob) *job.Running {
	s.l.Lock()
	defer s.l.Unlock()
	s.n++
	return job.NewClojure(func() interface{} {
		defer s.done()
		return j.Start().Await()
	}).Start()
}

func (s *Barrier) done() {
	s.l.Lock()
	s.n--
	if s.n == 0 {
		for !s.wsLks.Empty() {
			l := s.wsLks.PopFront().(*sync.Mutex)
			l.Unlock()
		}
	}
	s.l.Unlock()
}
