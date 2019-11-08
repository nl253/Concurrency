package worker

import (
	"github.com/nl253/DataStructures/stream"
)

type Worker struct {
	f func(in *stream.Stream, out *stream.Stream)
}

func New(f func(in *stream.Stream, out *stream.Stream)) *Worker { return &Worker{f} }

func (w *Worker) Start() (*stream.Stream, *stream.Stream) {
	in := stream.New()
	out := stream.New()
	go func() {
		w.f(in, out)
		in.Close()
		out.Close()
	}()
	return in, out
}

func (w *Worker) String() string {
	return "Worker"
}

func (w *Worker) Eq(x interface{}) bool {
	switch x.(type) {
	case Worker:
		return w == x.(*Worker)
	default:
		return false
	}
}
