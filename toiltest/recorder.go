package toiltest


// ToilRecorder is an implementation of toil.Toiler, as well as an implementation
// of toil.Recovereder and toil.Terminateder too; it counts the number of times
// its Toil() method has been called and has not returned (i.e., is blocking) as
// well as allows custom code to run when its Recovered(), Terminated(), or Toil()
// methods are called.
type ToilRecorder struct {
	panicCh     chan struct{value interface{}}
	terminateCh chan struct{doneCh chan struct{}}
	numToiling int

	recoveredFunc  func(interface{})
	terminatedFunc func()
	toilFunc       func()
}


// NewRecorder returns an initialized ToilRecorder.
func NewRecorder() *ToilRecorder {
	panicCh := make(chan struct{value interface{}})

	terminateCh := make(chan struct{doneCh chan struct{}})

	toilRecorder := ToilRecorder{
		panicCh:panicCh,
		terminateCh:terminateCh,
	}

	return &toilRecorder
}


// RecoveredFunc registers the "recovered function" that will be called as part of when the
// ToilRecorder's Recovered() method is called.
func (toiler *ToilRecorder) RecoveredFunc(fn func(interface{})) {
	toiler.recoveredFunc = fn
}

// TerminateFunc registers the "terminated function" that will be called as part of when the
// ToilRecorder's Terminated() method is called.
func (toiler *ToilRecorder) TerminatedFunc(fn func()) {
	toiler.terminatedFunc = fn
}

// ToilFunc registers the "toil function" that will be called as part of when the
// ToilRecorder's Toil() method is called.
func (toiler *ToilRecorder) ToilFunc(fn func()) {
	toiler.toilFunc = fn
}


// NumToiling returns the number of active calls to its Toil() method.
func (toiler *ToilRecorder) NumToiling() int {
	return toiler.numToiling
}


// Panic causes one of the still active (i.e., blocking) calls to Toil()
// on itself to panic().
//
// If there are not active (i.e., blocking) calls to Toil() on itself,
// then it will block until there is one.
//
// One use for this method is to check if its Recovered() method was
// call by the toil.Group it is in (due to the panic()).
func (toiler *ToilRecorder) Panic(value interface{}) {

	toiler.panicCh <- struct{value interface{}}{
		value:value,
	}

//@TODO: Is there a way to wait for this to complete?
}


// Terminate causes one of the still active (i.e., blocking) calls to Toil()
// on itself to return gracefully.
//
// If there are not active (i.e., blocking) calls to Toil() on itself,
// then it will block until there is one.
//
// One use for this method is to check if its Terminated() method was
// call by the toil.Group it is in (due to the gracefull return).
func (toiler *ToilRecorder) Terminate() {
	doneCh := make(chan struct{})

	toiler.terminateCh <- struct{doneCh chan struct{}}{
		doneCh:doneCh,
	}

	<-doneCh
}


// Toil is part of the toil.Toiler interface.
func (toiler *ToilRecorder) Toil() {
	toiler.numToiling++

	if nil != toiler.toilFunc {
		toiler.toilFunc()
	}

	var doneCh chan struct{}

	select {
	case panicRequest := <-toiler.panicCh:
		panic(panicRequest.value)
	case terminateRequest := <-toiler.terminateCh:
		doneCh = terminateRequest.doneCh
	}

	toiler.numToiling--

	if nil != doneCh {
		doneCh <- struct{}{}
	}
}


// Recovered is part of the toil.ToilRecovereder interface.
func (toiler *ToilRecorder) Recovered(panicValue interface{}) {
	if nil != toiler.recoveredFunc {
		toiler.recoveredFunc(panicValue)
	}
}


// Terminated is part of the toil.ToilTerminateder interface.
func (toiler *ToilRecorder) Terminated() {
	if nil != toiler.terminatedFunc {
		toiler.terminatedFunc()
	}
}




