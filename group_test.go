package toil


import (
	"testing"

	"github.com/reiver/go-toil/toiltest"

	"fmt"
	"math/rand"
	"sync"
	"time"
)


func TestNewGroup(t *testing.T) {

	group := NewGroup()
	if nil == group {
		t.Errorf("After creating a new group, expected it to not be nil, but it was: %v", group)
		return
	}
}


func TestLen(t *testing.T) {

	group := NewGroup()

	length := group.Len()

	if expected, actual := 0, length; expected != actual {
		t.Errorf("After creating new groupp, expected the number of registered toilers to be %d, but actually was %d.", expected, actual)
		return
	}
}


func TestRegister(t *testing.T) {

	toiler := toiltest.NewRecorder()

	group := NewGroup()

	const NUM_REGISTER_TESTS = 20
	for testNumber:=0; testNumber<NUM_REGISTER_TESTS; testNumber++ {

		group.Register(toiler)

		length := group.Len()

		if expected, actual := 1+testNumber, length; expected != actual {
			t.Errorf("For test #%d, after registering toilers with new group, expected the number of registered toilers to be %d, but actually was %d.", testNumber, expected, actual)
			continue
		}
	}
}


func TestToil(t *testing.T) {

	// Initialize.
	randomness := rand.New( rand.NewSource( time.Now().UTC().UnixNano() ) )


	// Do tests.
	const NUM_TOIL_TESTS = 20
	for testNumber:=0; testNumber<NUM_TOIL_TESTS; testNumber++ {

		numberOfTimesToToil := randomness.Intn(44)

		var waitGroup sync.WaitGroup
		waitGroup.Add(numberOfTimesToToil)

		toiler := toiltest.NewRecorder()
		toiler.ToilFunc(func(){
			waitGroup.Done()
		})


		group := NewGroup()


		for i:=0; i<numberOfTimesToToil; i++ {
			group.Register(toiler)
		}


		go group.Toil()


		waitGroup.Wait() // Make sure all the calls on the Toil() method are done before continuing.


		if expected, actual := numberOfTimesToToil, toiler.NumToiling(); expected != actual {
			t.Errorf("For test #%d with, expected the number of toiling toilers to be %d, but actually was %d.", testNumber, expected, actual)
			continue
		}
	}
}


func TestToilRecovereder(t *testing.T) {

	// Initialize.
	randomness := rand.New( rand.NewSource( time.Now().UTC().UnixNano() ) )


	// Do test.
	toiler := toiltest.NewRecorder()

	var toilWaitGroup sync.WaitGroup
	toilWaitGroup.Add(1)
	toiler.ToilFunc(func(){
		toilWaitGroup.Done()
	})

	var receivedPanicValue interface{}
	var recoveredWaitGroup sync.WaitGroup
	recoveredWaitGroup.Add(1)
	toiler.RecoveredFunc(func(panicValue interface{}){
		receivedPanicValue = panicValue
		recoveredWaitGroup.Done()
	})


	group := NewGroup()


	group.Register(toiler)


	go group.Toil()
	toilWaitGroup.Wait() // Make sure all the calls on the Toil() method are done before continuing.


	panicValue := fmt.Sprintf("Panic Value with some random stuff: %d", randomness.Intn(999999999))
	toiler.Panic(panicValue)
	recoveredWaitGroup.Wait() // Make sure all the calls on the Recovered() method are done before continuing.


	if expected, actual := panicValue, receivedPanicValue; expected != actual {
		t.Errorf("Expected recovered panic value to be %v, but actually was %v.", expected, actual)
	}
}


func TestToilTerminateder(t *testing.T) {

	toiler := toiltest.NewRecorder()

	var toilWaitGroup sync.WaitGroup
	toilWaitGroup.Add(1)
	toiler.ToilFunc(func(){
		toilWaitGroup.Done()
	})

	var terminatedWaitGroup sync.WaitGroup
	terminatedWaitGroup.Add(1)
	toiler.TerminatedFunc(func(){
		terminatedWaitGroup.Done()
	})


	group := NewGroup()


	group.Register(toiler)


	go group.Toil()
	toilWaitGroup.Wait() // Make sure all the calls on the Toil() method are done before continuing.


	toiler.Terminate()
	terminatedWaitGroup.Wait() // Make sure all the calls on the Terminated() method are done before continuing.
}
