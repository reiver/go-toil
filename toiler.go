package toil


// Toiler is an interface that wraps the Toil method.
//
// The purpose of the Toil method is to do work.
// The Toil method should block while it is doing work.
type Toiler interface {
	Toil()
}
