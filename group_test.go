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


func TestToilPanicked(t *testing.T) {

	// Initialize.
	randomness := rand.New( rand.NewSource( time.Now().UTC().UnixNano() ) )


	// Do test.
	toiler := toiltest.NewRecorder()

	numToiled := 0
	var toiledWaitGroup sync.WaitGroup
	toiledWaitGroup.Add(1)
	toiler.ToilFunc(func(){
		numToiled++
		toiledWaitGroup.Done()
	})


	numPanicked := 0
	var panickedWaitGroup sync.WaitGroup
	toiler.PanickedNoticeFunc(func(panicValue interface{}){
		numPanicked++
		panickedWaitGroup.Done()
	})


	group := NewGroup()


	group.Register(toiler)


	var expectedPanicValue interface{} = nil
	go func() {
		defer func() {
			if expected, actual := 1, numToiled; expected != actual {
				t.Errorf("Expected number of times toiled to be %d, but actually was %d.", expected, actual)
			}

			if panicValue := recover(); nil != panicValue {
				if expected, actual := panicValue, expectedPanicValue; expected != actual {
					t.Errorf("Expected caught panic value to be [%v], but actually was [%v].", expected, actual)
				}
			} else {
				t.Errorf("This should NOT get to this part of the code either!!")
			}
		}()

		group.Toil()
	}()
	toiledWaitGroup.Wait() // Make sure all the calls on the Toil() method are done before continuing.



	if expected, actual := 1, numToiled; expected != actual {
		t.Errorf("Expected number of times toiled to be %d, but actually was %d.", expected, actual)
	}

	if expected, actual := 0, numPanicked; expected != actual {
		t.Errorf("Expected number of times panicked to be %d, but actually was %d.", expected, actual)
	}



	panickedWaitGroup.Add(1)
	panicValue := fmt.Sprintf("Panic Value with some random stuff: %d", randomness.Intn(999999999))
	expectedPanicValue = panicValue // <----------------- NOTE we set the expectedPanicValue
	toiler.Panic(panicValue)
	panickedWaitGroup.Wait()



	//                     V---------- NOTE that sayed as 1.
	if expected, actual := 1, numToiled; expected != actual {
		t.Errorf("Expected number of times toiled to be %d, but actually was %d.", expected, actual)
	}

	if expected, actual := 1, numPanicked; expected != actual {
		t.Errorf("Expected number of times panicked to be %d, but actually was %d.", expected, actual)
	}
}
