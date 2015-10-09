package toil


// returnedNotifiableToiler is an interface that wraps the Toil and Returned methods.
//
// The purpose of the Toil method is to do work.
// The Toil method should block while it is doing work.
//
// The purpose of the ReturnedNotice method is as a means of notifying when
// the Toil returned (gracefully).
type returnedNotifiableToiler interface {
	Toiler
	ReturnedNotice()
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
