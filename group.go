package toil


// Group is an interface that wraps the Len, Register and Toil methods.
type Group interface {

	// Len returns the number of toilers registered with this Group.
	Len() int

	// Register registers a toiler with this Group.
	Register(Toiler)

	// Toil makes all the toilers registered with this Group toil (i.e., do work),
	// by calling each of the registered toilers' Toil methods.
	Toil()
}


type internalGroup struct {
	daemon *internalGroupDaemon
	panicCh chan interface{}
}


// NewGroup returns an initialized Group.
func NewGroup() Group {
	panicCh := make(chan interface{})

	groupDaemon := newGroupDaemon(panicCh)

	group := internalGroup{
		daemon:groupDaemon,
		panicCh:panicCh,
	}

	return &group
}


func (group *internalGroup) Len() int {
	lengthReturnCh := make(chan int)
	defer close(lengthReturnCh)

	group.daemon.LengthCh() <- struct{returnCh chan int}{
		returnCh:lengthReturnCh,
	}

	length := <-lengthReturnCh

	return length
}


func (group *internalGroup) Register(toiler Toiler) {
	doneCh := make(chan struct{})

	group.daemon.RegisterCh() <- struct{doneCh chan struct{}; toiler Toiler}{
		doneCh:doneCh,
		toiler:toiler,
	}

	<-doneCh // NOTE that we are waiting on this before we call
	         // the Wait() method below to avoid a race condition.
}


func (group *internalGroup) Toil() {
	// By sending on this channel, we make all the toilers
	// registered in this group toil.
	doneCh := make(chan struct{})

	group.daemon.ToilCh() <- struct{doneCh chan struct{}}{
		doneCh:doneCh,
	}

	<-doneCh // NOTE that we are waiting on this before we call
	         // the Wait() method below to avoid a race condition.


	// Block while any toiler in this group is still toiling and
	// none of them have panic()ed.
	//
	// If any panic() then this panic()s.
	waitForThem := func() (<-chan struct{}) {
		ch := make(chan struct{})
		go func() {
			group.daemon.Waiter().Wait()
			ch <- struct{}{}
		}()
		return ch
	}

	select {
	case panicValue := <-group.panicCh:
		panic(panicValue)
	case <-waitForThem():
	}
}
