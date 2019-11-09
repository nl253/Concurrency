package iterator

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/nl253/DataStructures/list"
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
		lk: &sync.Mutex{},
		f: func(state interface{}) interface{} {
			if state == EndOfIteration {
				return state
			} else {
				return f(state)
			}
		},
		state: initState,
	}
}

func Range(initState int, step int) *Iterator {
	return New(initState, func(i interface{}) interface{} { return i.(int) + step })
}

func Ints() *Iterator {
	return Range(0, 1)
}

func Floats() *Iterator {
	return New(rand.Float64(), func(_ interface{}) interface{} { return rand.Float64() })
}

func Bytes() *Iterator {
	return New(rand.Uint32(), func(_ interface{}) interface{} { return byte(rand.Uint32()) })
}

func (iter *Iterator) Sum() int {
	return iter.ReduceAll(func(x interface{}, y interface{}) interface{} { return x.(int) + y.(int) }).(int)
}

func (iter *Iterator) Count() int {
	return iter.Map(func(x interface{}) interface{} { return 1 }).Sum()
}

func Repeat(x interface{}) *Iterator {
	return New(x, func(y interface{}) interface{} { return y })
}

func (iter *Iterator) Slice(n uint) *Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	if n == 0 {
		return Repeat(EndOfIteration)
	}
	return New(iter.pull(), func(state interface{}) interface{} {
		if n <= 1 {
			return EndOfIteration
		} else {
			n--
			return iter.pull()
		}
	})
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

func (iter *Iterator) tryAdvance() bool {
	ok := iter.state != EndOfIteration
	if ok {
		iter.state = iter.f(iter.state)
	}
	return ok
}

func (iter *Iterator) pull() interface{} {
	defer iter.tryAdvance()
	return iter.state
}

func (iter *Iterator) PullN(n uint) []interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	result := make([]interface{}, n)
	for i := uint(0); i < n; i++ {
		result[i] = iter.pull()
	}
	return result
}

func (iter *Iterator) PullAll() *list.ConcurrentList {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	result := list.New()
	for iter.state != EndOfIteration {
		result.Append(iter.pull())
	}
	return result
}

func (iter *Iterator) Skip() *Iterator {
	return iter.SkipN(1)
}

func (iter *Iterator) skip() *Iterator {
	return iter.skipN(1)
}

func (iter *Iterator) SkipN(n uint) *Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return iter.skipN(n)
}

func (iter *Iterator) skipN(n uint) *Iterator {
	for i := uint(0); i < n; i++ {
		if !iter.tryAdvance() {
			return iter
		}
	}
	return iter
}

func (iter *Iterator) Map(f func(interface{}) interface{}) *Iterator {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	return New(f(iter.pull()), func(state interface{}) interface{} { return f(iter.Pull()) })
}

func (iter *Iterator) Printf(format string) *Iterator {
	return iter.Tap(func(x interface{}) { fmt.Printf(format, x) })
}

func (iter *Iterator) Println() *Iterator {
	return iter.Printf("%v\n")
}

func (iter *Iterator) Tap(f func(interface{})) *Iterator {
	return iter.Map(func(x interface{}) interface{} {
		if x != EndOfIteration {
			f(x)
		}
		return x
	})
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
	for i := uint(1); i < n; i++ {
		acc = f(acc, iter.pull())
	}
	return acc
}

func (iter *Iterator) ReduceAll(f func(interface{}, interface{}) interface{}) interface{} {
	iter.lk.Lock()
	defer iter.lk.Unlock()
	acc := iter.pull()
	for {
		if x := iter.pull(); x != EndOfIteration {
			acc = f(acc, x)
		} else {
			return acc
		}
	}
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
