package toil


import (
	"testing"

	"github.com/reiver/go-toil/toiltest"

	"math/rand"
	"sync"
	"time"
)


func TestNewGroupDaemon(t *testing.T) {

	daemon := newGroupDaemon()
	if nil == daemon {
		t.Errorf("After creating a new daemon, expected it to not be nil, but it was: %v", daemon)
		return
	}

	if nil == daemon.pingCh {
		t.Errorf("After creating a new daemon, expected daemon.pingCh to not be nil, but it was: %v", daemon.pingCh)
		return
	}
	if nil == daemon.toilCh {
		t.Errorf("After creating a new daemon, expected daemon.toilCh to not be nil, but it was: %v", daemon.toilCh)
		return
	}
	if nil == daemon.registerCh {
		t.Errorf("After creating a new daemon, expected daemon.registerCh to not be nil, but it was: %v", daemon.registerCh)
		return
	}

}


func TestPingCh(t *testing.T) {

	daemon := newGroupDaemon()

	const NUM_PING_TESTS = 20
	doneCh := make(chan struct{})
	for i:=0; i<NUM_PING_TESTS; i++ {
		daemon.PingCh() <- struct{doneCh chan struct{}}{
			doneCh:doneCh,
		}
		<-doneCh // THE TEST WE ARE DOING IS MAKING SURE THIS DOES NOT RESULT IN A DEADLOCK.
		         // THUS WE ARE NOT CALLING ANYTHING LIKE t.Errorf() IN THIS CASE.
	}
}


func TestLengthCh(t *testing.T) {

	daemon := newGroupDaemon()

	lengthReturnCh := make(chan int)
	daemon.LengthCh() <- struct{returnCh chan int}{
		returnCh:lengthReturnCh,
	}
	length := <-lengthReturnCh

	if expected, actual := 0, length; expected != actual {
		t.Errorf("After creating a new daemon, expected the number of registered toilers to be %d, but actually was %d.", expected, actual)
		return
	}
}


func TestRegisterCh(t *testing.T) {

	toiler := toiltest.NewRecorder()

	daemon := newGroupDaemon()

	const NUM_REGISTER_TESTS = 20
	doneCh := make(chan struct{})
	for testNumber:=0; testNumber<NUM_REGISTER_TESTS; testNumber++ {
		daemon.RegisterCh() <- struct{doneCh chan struct{}; toiler Toiler}{
			doneCh:doneCh,
			toiler:toiler,
		}
		<-doneCh // THE TEST WE ARE DOING IS MAKING SURE THIS DOES NOT RESULT IN A DEADLOCK.
		         // THUS WE ARE NOT CALLING ANYTHING LIKE t.Errorf() IN THIS CASE.


		lengthReturnCh := make(chan int)
		daemon.LengthCh() <- struct{returnCh chan int}{
			returnCh:lengthReturnCh,
		}
		length := <-lengthReturnCh

		if expected, actual := 1+testNumber, length; expected != actual {
			t.Errorf("For test #%d, after registering toilers with new daemon, expected the number of registered toilers to be %d, but actually was %d.", testNumber, expected, actual)
			continue
		}
	}
}


func TestToilCh(t *testing.T) {

	// Initialize.
	randomness := rand.New( rand.NewSource( time.Now().UTC().UnixNano() ) )


	// Do tests.
	const NUM_TOIL_TESTS = 20
	doneCh := make(chan struct{})
	for testNumber:=0; testNumber<NUM_TOIL_TESTS; testNumber++ {

		numberOfTimesToToil := randomness.Intn(44)

		var waitGroup sync.WaitGroup
		waitGroup.Add(numberOfTimesToToil)

		toiler := toiltest.NewRecorder()
		toiler.ToilFunc(func(){
			waitGroup.Done()
		})


		daemon := newGroupDaemon()


		for i:=0; i<numberOfTimesToToil; i++ {

			daemon.RegisterCh() <- struct{doneCh chan struct{}; toiler Toiler}{
				doneCh:doneCh,
				toiler:toiler,
			}
			<-doneCh // THE TEST WE ARE DOING IS MAKING SURE THIS DOES NOT RESULT IN A DEADLOCK.
			         // THUS WE ARE NOT CALLING ANYTHING LIKE t.Errorf() IN THIS CASE.
		}


		daemon.ToilCh() <- struct{doneCh chan struct{}}{
			doneCh:doneCh,
		}
		<-doneCh // THE TEST WE ARE DOING IS MAKING SURE THIS DOES NOT RESULT IN A DEADLOCK.
		         // THUS WE ARE NOT CALLING ANYTHING LIKE t.Errorf() IN THIS CASE.


		waitGroup.Wait() // Make sure all the calls on the Toil() method are done before continuing.


		if expected, actual := numberOfTimesToToil, toiler.NumToiling(); expected != actual {
			t.Errorf("For test #%d with, expected the number of toiling toilers to be %d, but actually was %d.", testNumber, expected, actual)
			continue
		}
	}
}
