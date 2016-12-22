package toil


// The ToilerFunc type is an adapter to allow the use of ordinary functions as toilers.
// If fn is a function with the appropriate signature, ToilerFunc(fn) is a Toiler that calls fn.
type ToilerFunc func()


// Toil calls fn().
func (fn ToilerFunc) Toil() {
	fn()
}
