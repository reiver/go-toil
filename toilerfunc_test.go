package toil


import (
	"testing"
)


func TestToilerFunc(t *testing.T) {

	fn := func() {
		// Nothing here.
	}

	var toiler Toiler = ToilerFunc(fn)

	if nil == toiler {
		t.Errorf("This should never happen.")
	}
}
