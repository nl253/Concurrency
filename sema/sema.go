package sema

import (
	"sync"

	"github.com/nl253/DataStructures/list"
)

type Sema struct {
	n     uint
	max   uint
	l     *sync.Mutex
	wsLks *list.ConcurrentList
}

func New(n uint) *Sema {
	return &Sema{
		n:     n,
		max:   n,
		l:     &sync.Mutex{},
		wsLks: list.New(),
	}
}

func (s *Sema) Acquire() {
	for {
		s.l.Lock()
		if s.n > 1 {
			s.n--
			break
		}
		l := &sync.Mutex{}
		l.Lock()
		s.wsLks.Append(l)
		s.l.Unlock()
		l.Lock()
	}
	s.l.Unlock()
}

func (s *Sema) AcquireN(n uint) {
	for i := uint(0); i < n; i++ {
		s.Acquire()
	}
}

func (s *Sema) AcquireAll() {
	s.AcquireN(s.max)
}

func (s *Sema) Release() {
	s.l.Lock()
	s.n++
	if !s.wsLks.Empty() {
		l := s.wsLks.PopFront().(*sync.Mutex)
		l.Unlock()
	}
	s.l.Unlock()
}

func (s *Sema) ReleaseN(n uint) {
	for i := uint(0); i < n; i++ {
		s.Release()
	}
}

func (s *Sema) ReleaseAll() {
	s.ReleaseN(s.max)
}
