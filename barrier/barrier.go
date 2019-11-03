package barrier

import (
    "sync"

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

func (s *Barrier) Wait() {
    s.l.Lock()
    if s.n > 0  {
        l := &sync.Mutex{}
        l.Lock()
        s.wsLks.Append(l)
        s.l.Unlock()
        l.Lock()
        l.Unlock()
        s.l.Lock()
    }
    s.l.Unlock()
}


func (s *Barrier) Done() {
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
