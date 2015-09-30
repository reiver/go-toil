package toil


// Toiler is an interface that wraps the Toil method.
//
// The purpose of the Toil method is to do work.
// The Toil method should block while it is doing work.
type Toiler interface {
	Toil()
}


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


// toilRecovereder is an interface that wraps the Toil and Recovered methods.
//
// The purpose of the Toil method is to do work.
// The Toil method should block while it is doing work.
//
// The purpose of the Recovered method is as a means of notifying when
// the Toil method panic()ed.
type toilRecovereder interface {
	Toiler
	Recovered(interface{})
}
