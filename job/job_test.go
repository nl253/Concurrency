package job

import (
	"math/rand"
	"testing"

	ut "github.com/nl253/Testing"
)

var fJob = ut.Test("Worker")

const N uint = 10000

func TestActor_Start(t *testing.T) {
	should := fJob("Start", t)
	for i := uint(0); i < N; i++ {
		expect := rand.Int()
		should("not crash", expect, func() interface{} {
			return NewClojure(func() interface{} { return expect }).Start().Await()
		})
		expect2 := rand.Int()
		should("not crash when run with args", expect2+expect, func() interface{} {
			return New(func(xs ...interface{}) interface{} { return xs[0].(int) + xs[1].(int) }).Start(expect, expect2).Await()
		})
	}
}
