package worker

import (
	"testing"
)

func BenchmarkWorker_Start(b *testing.B) {
	// a := New(func(in *stream.Stream, out *stream.Stream, err *stream.Stream) {
	// 	end := in.Pull().(int)
	// 	for i := 0; i < end; i++ {
	// 		out.PushBack(in.Pull().(int) + 1)
	// 	}
	// }).Start()
	// nMsgs := 100000
	// a.Input().PushBack(nMsgs)
	// for i := 0; i < nMsgs; i++ {
	// 	a.Input().PushBack(rand.Int())
	// }
	// for i := 0; i < nMsgs; i++ {
	// 	fmt.Printf("GOT %d\n", a.Output().Pull().(int))
	// }
}
