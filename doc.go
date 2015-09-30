/*
Package toil provides simple functionality for managing toilers (i.e., workers).

Usage

To use, create one or more types that implement the toil.Toiler interface. For example:

	type myToiler struct{}
	
	func newMyToiler() {
	
		toiler := myToiler{}
	
		return &toiler
	}
	
	func (toiler *myToiler) Toil() {
		//@TODO: Do work here.
	}

Then create a toil.Group. For example:

	var (
		myToilerGroup = toil.NewGroup()
	)

Then register one of more toilers (i.e., types that implement the toil.Toiler interface)
with the toiler group. For example:

	toiler := newMyToiler()

	myToilGroup.Register(toiler)

Then, you can call the Toil method of the toiler group in place like main(). For example:

	func main() {
	
		// ...
	
		// Calling the Toil() method on the toiler group
		// will cause it to call the Toil() method of
		// each toiler registered with it.
		//
		// Thus causing each of those toilers registered
		// with it to start doing its work (whatever that
		// happens to be).
		//
		// This will block until all the the toilers
		// registered in this toiler group's Toil() methods
		// finishes (either because it returned gracefully
		// or because it panic()ed).
		myToilGroup.Toil()
	
		// ...
	
	}

Observers

A toiler's Toil method can finish in one of two ways. Either it will return gracefully, or
it will panic().

The toiler group is OK with either.

But also, the toiler group provides the toiler with a convenient way of being notified
of each case.

If a toiler also has a Terminated() method, then the toiler group will call the toiler's
Terminated() method when the toiler's Toil() method has returned gracefully. For example:

	type myToiler struct{}
	
	func newMyToiler() {
	
		toiler := myToiler{}
	
		return &toiler
	}
	
	func (toiler *myToiler) Toil() {
		//@TODO: Do work here.
	}
	
	func (toiler *myToiler) Terminated() {
		//@TODO: Do something with this notification.
	}

If a toiler also has a Recovered() method, then the toiler group will call the toiler's
Recovered() method when the toiler's Toil() method has panic()ed. For example:

	type myToiler struct{}
	
	func newMyToiler() {
	
		toiler := myToiler{}
	
		return &toiler
	}
	
	func (toiler *myToiler) Toil() {
		//@TODO: Do work here.
	}
	
	func (toiler *myToiler) Recovered() {
		//@TODO: Do something with this notification.
	}

And of course, a toiler can take advantage of both of these notifications and have
both a Recovered() and Terminated() method. For example:

	type myToiler struct{}
	
	func newMyToiler() {
	
		toiler := myToiler{}
	
		return &toiler
	}
	
	func (toiler *myToiler) Toil() {
		//@TODO: Do work here.
	}
	
	func (toiler *myToiler) Recovered() {
		//@TODO: Do something with this notification.
	}
	
	func (toiler *myToiler) Terminated() {
		//@TODO: Do something with this notification.
	}

*/
package toil
