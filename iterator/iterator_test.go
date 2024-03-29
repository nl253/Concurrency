package iterator

import (
	"testing"

	ut "github.com/nl253/Testing"
)

var fIter = ut.Test("Iterator")

const N uint = 10000

func TestIterator_Range(t *testing.T) {
	should := fIter("Range", t)
	should("make iter of ints in range [min, max)", 0, func() interface{} {
		it := Ints()
		return it.Pull()
	})
	should("make iter of ints in range [min, max)", 1, func() interface{} {
		it := Ints()
		it.Skip()
		return it.Pull()
	})
	should("make iter of ints in range [min, max)", 2, func() interface{} {
		it := Ints()
		it.SkipN(2)
		return it.Pull()
	})
}

func TestIterator_Slice(t *testing.T) {
	should := fIter("Slice", t)
	should("slice [, 0) should give EndOfIteration", EndOfIteration, func() interface{} {
		return Ints().Slice(0).Pull()
	})
	should("slice [, 1) should give 0", 0, func() interface{} {
		return Ints().Slice(1).Pull()
	})
	should("slice [, 1).Skip() should give EndOfIteration", EndOfIteration, func() interface{} {
		return Ints().Println().Slice(1).Skip().Pull()
	})
	should("slice [, 1).Skip() should give EndOfIteration", 2, func() interface{} {
		return Ints().Slice(3).SkipN(2).Pull()
	})
}

func TestIterator_Sum(t *testing.T) {
	should := fIter("Sum", t)
	should("make iter of ints in range [min, max)", 0, func() interface{} {
		return Ints().Slice(1).Sum()
	})
	should("make iter of ints in range [min, max)", 0+1+2, func() interface{} {
		return Ints().Slice(3).Sum()
	})
}

func TestIterator_Tap(t *testing.T) {
	should := fIter("Tap", t)
	should("tap to print", 2, func() interface{} {
		return Ints().Slice(2).Map(func(i interface{}) interface{} {
			return 0
		}).Sum()
	})
}

func TestIterator_Map(t *testing.T) {
	should := fIter("Map", t)
	should("id should not modify values", 0, func() interface{} {
		return Ints().Map(func(x interface{}) interface{} {
			return x
		}).Map(func(x interface{}) interface{} {
			return x
		}).Slice(1).Println().Pull()
	})
	should("id should not modify values", EndOfIteration, func() interface{} {
		return Ints().Map(func(x interface{}) interface{} {
			return x
		}).Map(func(x interface{}) interface{} {
			return x
		}).Slice(1).Skip().Println().Pull()
	})
	should("id should not modify values", 1, func() interface{} {
		return Ints().Map(func(x interface{}) interface{} {
			return x
		}).Map(func(x interface{}) interface{} {
			return x
		}).Slice(2).Skip().Println().Pull()
	})
	should("make iter of ints in range [min, max)", 0+10+1+10, func() interface{} {
		return Ints().Map(func(x interface{}) interface{} { return x.(int) + 10 }).Slice(2).Println().Sum()
	})
}
