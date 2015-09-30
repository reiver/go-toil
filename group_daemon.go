package toil


import (
	"sync"
)


type internalGroupDaemon struct {
	waitGroup       sync.WaitGroup
	lengthCh   chan struct{returnCh chan int}
	pingCh     chan struct{doneCh   chan struct{}}
	registerCh chan struct{doneCh   chan struct{}; toiler Toiler}
	toilCh     chan struct{doneCh   chan struct{}}
}


func newGroupDaemon() *internalGroupDaemon {

	lengthCh   := make(chan struct{returnCh chan int})
	pingCh     := make(chan struct{doneCh   chan struct{}})
	registerCh := make(chan struct{doneCh   chan struct{}; toiler Toiler})
	toilCh     := make(chan struct{doneCh   chan struct{}})

	daemon := internalGroupDaemon{
		lengthCh:lengthCh,
		pingCh:pingCh,
		toilCh:toilCh,
		registerCh:registerCh,
	}

	go daemon.animate()

	return &daemon
}


func (daemon *internalGroupDaemon) Waiter() waiter {
	return &daemon.waitGroup
}



func (daemon *internalGroupDaemon) PingCh() chan<- struct{doneCh chan struct{}} {
	return daemon.pingCh
}

func (daemon *internalGroupDaemon) LengthCh() chan<- struct{returnCh chan int} {
	return daemon.lengthCh
}

func (daemon *internalGroupDaemon) RegisterCh() chan<- struct{doneCh chan struct{}; toiler Toiler} {
	return daemon.registerCh
}

func (daemon *internalGroupDaemon) ToilCh() chan<- struct{doneCh chan struct{}} {
	return daemon.toilCh
}



func (daemon *internalGroupDaemon) animate() {

	toilers := make([]Toiler, 0, 8)

	toiling := false

	for {
		select {
		case lengthRequest := <-daemon.lengthCh:
			lengthRequest.returnCh <- len(toilers)
		case pingRequest := <-daemon.pingCh:
			pingRequest.doneCh <- struct{}{}
		case registrationRequest := <-daemon.registerCh:
			toiler := registrationRequest.toiler

			toilers = append(toilers, toiler)
			if toiling {
				daemon.spawn(toiler)
			}

			registrationRequest.doneCh <- struct{}{}
		case toilRequest := <-daemon.toilCh:
			if !toiling {
				toiling = true
				for _,toiler := range toilers {
					daemon.spawn(toiler)
				}
				toilRequest.doneCh <- struct{}{}
			}
		}
	}
}


// spawn does the hard work of making a toiler toil.
func (daemon *internalGroupDaemon) spawn(toiler Toiler) {

	// We increment the wait group for each goroutine we spawn.
	//
	// This wait group is used by the "Group" type in its Toil()
	// method to make it so Toil() blocks (and does not return)
	// while there are toilers still toiling.
	//
	// Of course, the "Group" type's Toil() method does NOT have
	// direct access to this wait group, but instead gets indirect
	// access to it via this daemon's Waiter() method.
	daemon.waitGroup.Add(1)


	// Spawn a goroutine, and make the toiler toil within the spawned goroutine.
	go func(toiler Toiler){

		// We decrement the wait group each time a goroutine (of this type)
		// exits, by either panic()ing or the toiler.Toil() method returning.
		//
		// This wait group is used by the "Group" type in its Toil()
		// method to make it so Toil() blocks (and does not return)
		// while there are toilers still toiling.
		//
		// Of course, the "Group" type's Toil() method does NOT have
		// direct access to this wait group, but instead gets indirect
		// access to it via this daemon's Waiter() method.
		defer daemon.waitGroup.Done()


		// We do this so that we can capture a panic() that could happen from the
		// toiler's Toil() method.
		defer func() {
			if r := recover(); nil != r {

				// If we got to this point in the code, then the toiler's Toil()
				// method has panic()ed (rather than returning gracefully).
				//
				// At this point we see if the toiler supports us telling it that its
				// Toil() method panic()ed.
				//
				// We do this by trying to cast it to another type of interface.
				// Specifically, the toilRecovereder interface.
				//
				// This can be useful for adding in logging, tracking, etc.
				//
				// We do the actual call to the toilRecovereder's Recovered() method
				// in a goroutine, since we don't want it to block or panic() here!
				//
				// NOTE THAT THIS IS A POTENTIAL SOURCE OF A RESOURCE LEAK!!!!!!
				if recoveredableToiler, ok := toiler.(toilRecovereder); ok {
					go func(recoveredableToiler toilRecovereder){
						recoveredableToiler.Recovered(r)
					}(recoveredableToiler)
				}
			}
		}()

		// Make the toiler toil. (I.e., do work.)
		//
		// This method call is expected to be blocking!
		toiler.Toil()


		// If we got to this point in the code, then the toiler's Toil()
		// method has gracefully returned (rather than panic()ing).
		//
		// At this point we see if the toiler supports us telling it that its
		// Toil() method return (gracefully).
		//
		// We do this by trying to cast it to another type of interface.
		// Specifically, the toilTerminateder interface.
		//
		// This can be useful for adding in logging, tracking, etc.
		//
		// We do the actual call to the toilTerminateder's Terminated() method
		// in a goroutine, since we don't want it to block or panic() here!
		//
		// NOTE THAT THIS IS A POTENTIAL SOURCE OF A RESOURCE LEAK!!!!!!
		if terminatedableToiler, ok := toiler.(toilTerminateder); ok {
			go func(terminatedableToiler toilTerminateder){
				terminatedableToiler.Terminated()
			}(terminatedableToiler)
		}

	}(toiler)
}
