package job

import (
	"testing"

	ut "github.com/nl253/Testing"
)

var fJob = ut.Test("Worker")

func TestActor_Start(t *testing.T) {
	should := fJob("Start", t)
}

func TestActor_Done(t *testing.T) {
	should := fJob("Running", t)
}

func TestActor_String(t *testing.T) {
	should := fJob("String", t)
	should("stringify converts to string", true, func() interface{} {
		return ""
	})
}
