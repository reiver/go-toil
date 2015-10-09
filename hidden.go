package toil


// toilTerminateder is an interface that wraps the Toil and Terminated methods.
//
// The purpose of the Toil method is to do work.
// The Toil method should block while it is doing work.
//
// The purpose of the Terminated method is as a means of notifying when
// the Toil returned (gracefully).
type toilTerminateder interface {
	Toiler
	Terminated()
}


// panickedNotifiableToiler is an interface that wraps the Toil and Panicked methods.
//
// The purpose of the Toil method is to do work.
// The Toil method should block while it is doing work.
//
// The purpose of the PanickedNotice method is as a means of notifying when
// the Toil method panic()ed.
type panickedNotifiableToiler interface {
	Toiler
	PanickedNotice(interface{})
}
