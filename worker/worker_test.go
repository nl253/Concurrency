package worker

import (
	"math/rand"
	"testing"

	"github.com/nl253/DataStructures/stream"
	ut "github.com/nl253/Testing"
)

var fWorker = ut.Test("Worker")

const N uint = 10000

func TestWorker_Start(t *testing.T) {
	should := fWorker("Start", t)
	for i := uint(0); i < N; i++ {
		expect := rand.Int()
		should("start & send msg over output channel", expect, func() interface{} {
			worker := New(func(in *stream.Stream, out *stream.Stream) { out.PushBack(expect) })
			_, out := worker.Start()
			return out.Pull()
		})
		should("start & receive msg over input channel", expect, func() interface{} {
			worker := New(func(in *stream.Stream, out *stream.Stream) { out.PushBack(in.Pull()) })
			in, out := worker.Start()
			in.PushBack(expect)
			return out.Pull()
		})
	}
}

func TestWorker_Start2(t *testing.T) {
	should := fWorker("Start", t)
	for i := uint(0); i < N; i++ {
		should("start & receive msg over input channel", true, func() interface{} {
			worker := New(func(in *stream.Stream, out *stream.Stream) {
				n := in.Pull().(uint)
				for i := uint(0); i < n; i++ {
					out.PushBack(in.Pull().(int) + 1)
				}
			})
			in, out := worker.Start()
			in.PushBack(uint(3))
			in.PushBack(0)
			in.PushBack(1)
			in.PushBack(2)
			return out.Pull().(int) == 1 && out.Pull().(int) == 2 && out.Pull().(int) == 3
		})
	}
}
