package iterator

import (
	"sync"
)

type Iterator struct {
	lk    *sync.Mutex
	f     func(interface{}) interface{}
	state interface{}
}

type finished struct{}

var EndOfIteration = &finished{}

func New(initState interface{}, f func(interface{}) interface{}) *Iterator {
	return &Iterator{
		lk:    &sync.Mutex{},
		f:     f,
		state: initState,
	}
}

func Range(initState int, step int) *Iterator {
	return New(initState, func(i interface{}) interface{} { return i.(int) + step })
}

func Ints() *Iterator {
	return Range(0, 1)
}

func (iter *Iterator) Sum() int {
	return iter.ReduceAll(func(x interface{}, y interface{}) interface{} { return x.(int) + y.(int) }).(int)
}

func (iter *Iterator) Count() int {
	return iter.Map(func(x interface{}) interface{} { return 1 }).Sum()
}

func Repeat(x interface{}) *Iterator {
	return New(x, func(_ interface{}) interface{} { return x })
}

func (iter *Iterator) Peek() interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return iter.state
}

func (iter *Iterator) Pull() interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return iter.pull()
}

func (iter *Iterator) pull() interface{} {
	result := iter.state
	iter.state = iter.f(iter.state)
	return result
}

func (iter *Iterator) PullN(n uint) []interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	result := make([]interface{}, n)
	result[0] = iter.state
	for j := uint(1); j < n; j++ {
		result[j] = iter.f(result[j-1])
	}
	iter.state = result[n-1]
	return result
}

func (iter *Iterator) PullAll() []interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	result := make([]interface{}, 1)
	result[0] = iter.state
	for focus := iter.pull(); focus != EndOfIteration; focus = iter.pull() {
		result = append(result, focus)
	}
	return result
}

func (iter *Iterator) Skip() {
	iter.SkipN(1)
}

func (iter *Iterator) SkipN(n uint) {
	iter.lk.Lock()
	for j := uint(0); j < n; j++ {
		iter.state = iter.f(iter.state)
	}
	iter.lk.Unlock()
}

func (iter *Iterator) Map(f func(interface{}) interface{}) *Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return New(f(iter.state), func(state interface{}) interface{} { return f(iter.Pull()) })
}

func (iter *Iterator) Filter(f func(interface{}) bool) *Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return New(iter.state, func(state interface{}) interface{} {
		focus := state
		for !f(focus) {
			focus = iter.Pull()
		}
		return focus
	})
}

func (iter *Iterator) ReduceN(n uint, f func(interface{}, interface{}) interface{}) interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	acc := iter.pull()
	n--
	for i := uint(0); i < n; i++ {
		acc = f(acc, iter.pull())
	}
	return acc
}

func (iter *Iterator) ReduceAll(f func(interface{}, interface{}) interface{}) interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	acc := iter.pull()
	if acc != EndOfIteration {
		for peek := iter.pull(); peek != EndOfIteration; {
			acc = f(acc, peek)
		}
	}
	return acc
}

func (iter *Iterator) Clone() *Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return New(iter.state, iter.f)
}

func (iter *Iterator) Duplicate(n uint) []*Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	newIts := make([]*Iterator, n)
	for i := uint(0); i < n; i++ {
		newIts[i] = New(iter.state, iter.f)
	}
	return newIts
}

func (iter *Iterator) Eq(x interface{}) bool {
	switch x.(type) {
	case *Iterator:
		return iter == x.(*Iterator)
	default:
		return false
	}
}

func (iter *Iterator) String() string {
	return "Iterator"
}
