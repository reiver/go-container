package container


import (
	"testing"

	"fmt"
	"log"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)


// TestInjectShouldNotInject makes sure that if no struct tags are used, then nothing gets injected.
func TestInjectShouldNotInject(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	type Thing struct {
		Length   int
		Logger  *log.Logger
		PoolSize int
	}

	thing := new(Thing)

	var beforeLength   int        = randomness.Int()
	var beforeLogger  *log.Logger = nil
	var beforePoolSize int        = randomness.Int()

	thing.Length   = beforeLength
	thing.Logger   = beforeLogger
	thing.PoolSize = beforePoolSize

	// Create the dependency injection container.
	container := New()

	// Register some stuff into the dependency injection container.
	container.Register("logger", log.New(ioutil.Discard, "we be logging: ", log.Lshortfile))
	container.Register("pool-size", 20)
	container.Register("easter-egg", "apple-banana-cherry")
	container.Register("out", ioutil.Discard)

	// Inject.
	if err := container.Inject(thing); nil != err {
		t.Errorf("Received error when injecting. Error: %q", err)
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			missingDependencyNames := complainer.MissingDependencyNames()
			t.Errorf("Had %d missing dependencies: %#v", len(missingDependencyNames), missingDependencyNames)
		default:
			t.Errorf("Unknown complaint from injecting.")
		}
		return
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v, but didn't. Changed to %v.", expected, actual)
		return
	}

	if expected, actual := beforeLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to NOT have changed from %v, but didn't. Changed to %v.", expected, actual)
		return
	}

	if expected, actual := beforePoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to NOT have changed from %v, but didn't. Changed to %v.", expected, actual)
		return
	}
}


// TestInjectShouldInject makes sure that if (some) struct tags are used, then the appropriate gets injected.
func TestInjectShouldInject(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	type Thing struct {
		Length   int
		Logger  *log.Logger `inject:"logger"`
		PoolSize int        `inject:"pool-size"`
	}

	thing := new(Thing)

	var beforeLength   int        = randomness.Int()
	var beforeLogger  *log.Logger = nil
	var beforePoolSize int        = randomness.Int()

	thing.Length   = beforeLength
	thing.Logger   = beforeLogger
	thing.PoolSize = beforePoolSize

	// Create the dependency injection container.
	container := New()

	// Register some stuff into the dependency injection container.
	expectedLogger   := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize := 20

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", "apple-banana-cherry")
	container.Register("out", ioutil.Discard)

	// Inject.
	if err := container.Inject(thing); nil != err {
		t.Errorf("Received error when injecting. Error: %q", err)
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			missingDependencyNames := complainer.MissingDependencyNames()
			t.Errorf("Had %d missing dependencies: %#v", len(missingDependencyNames), missingDependencyNames)
		default:
			t.Errorf("Unknown complaint from injecting.")
		}
		return
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeLogger, thing.Logger; notExpected == actual {
		t.Errorf("Expected thing.Logger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforePoolSize, thing.PoolSize; notExpected == actual {
		t.Errorf("Expected thing.PoolSize to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	// Confirm that the things that changed, due to injection, changed to what we expected them to change to.
	if expected, actual := expectedLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedPoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to be %d, but actually was %d.", expected, actual)
		return
	}
}



type thingDependencies_TestInjectShouldAlsoInjectSingularHiddenDependencies struct {
	ShouldNotBeInjected  int
	HiddenLogger        *log.Logger `inject:"logger"`
	EasterEgg            string     `inject:"easter-egg"`
}

type Thing_TestInjectShouldAlsoInjectSingularHiddenDependencies struct {
	dependencies thingDependencies_TestInjectShouldAlsoInjectSingularHiddenDependencies
	Length   int
	Logger  *log.Logger `inject:"logger"`
	PoolSize int        `inject:"pool-size"`
}

func (thing *Thing_TestInjectShouldAlsoInjectSingularHiddenDependencies) Dependencies() interface{} {
	// Notice that it is returning a pointer
	// to the dependenies. That's important!
	return &thing.dependencies
}

// TestInjectShouldAlsoInjectSingularHiddenDependencies makes sure that a singular hidden dependency, then it is injected as well.
func TestInjectShouldAlsoInjectSingularHiddenDependencies(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	thing := new(Thing_TestInjectShouldAlsoInjectSingularHiddenDependencies)

	var beforeDependenciesShouldNotBeInjected int        = randomness.Int()
	var beforeDependenciesHiddenLogger       *log.Logger = nil
	var beforeDependenciesEasterEgg           string     = fmt.Sprintf("%d", randomness.Int())

	var beforeLength   int        = randomness.Int()
	var beforeLogger  *log.Logger = nil
	var beforePoolSize int        = randomness.Int()

	thing.dependencies.ShouldNotBeInjected = beforeDependenciesShouldNotBeInjected
	thing.dependencies.HiddenLogger        = beforeDependenciesHiddenLogger
	thing.dependencies.EasterEgg           = beforeDependenciesEasterEgg

	thing.Length   = beforeLength
	thing.Logger   = beforeLogger
	thing.PoolSize = beforePoolSize

	// Create the dependency injection container.
	container := New()

	// Register some stuff into the dependency injection container.
	expectedLogger    := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize  := 20
	expectedEasterEgg := "apple-banana-cherry"

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", expectedEasterEgg)
	container.Register("out", ioutil.Discard)

	// Inject.
	if err := container.Inject(thing); nil != err {
		t.Errorf("Received error when injecting. Error: %q", err)
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			missingDependencyNames := complainer.MissingDependencyNames()
			t.Errorf("Had %d missing dependencies: %#v", len(missingDependencyNames), missingDependencyNames)
		default:
			t.Errorf("Unknown complaint from injecting.")
		}
		return
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeDependenciesShouldNotBeInjected, thing.dependencies.ShouldNotBeInjected; expected != actual {
		t.Errorf("Expected thing.dependencies.ShouldNotBeInjected to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeDependenciesHiddenLogger, thing.dependencies.HiddenLogger; notExpected == actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforeDependenciesEasterEgg, thing.dependencies.EasterEgg; notExpected == actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}


	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeLogger, thing.Logger; notExpected == actual {
		t.Errorf("Expected thing.Logger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforePoolSize, thing.PoolSize; notExpected == actual {
		t.Errorf("Expected thing.PoolSize to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	// Confirm that the things that changed, due to injection, changed to what we expected them to change to.
	if expected, actual := expectedLogger, thing.dependencies.HiddenLogger; expected != actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedEasterEgg, thing.dependencies.EasterEgg; expected != actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to be %q, but actually was %q.", expected, actual)
		return
	}


	if expected, actual := expectedLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedPoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to be %d, but actually was %d.", expected, actual)
		return
	}
}



type thingDependencies1_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies struct {
	ShouldNotBeInjected  int
	HiddenLogger        *log.Logger `inject:"logger"`
	EasterEgg            string     `inject:"easter-egg"`
}

type thingDependencies2_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies struct {
	Out io.Writer `inject:"out"`
}

type Thing_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies struct {
	dependencies1 thingDependencies1_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies
	dependencies2 thingDependencies2_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies
	Length   int
	Logger  *log.Logger `inject:"logger"`
	PoolSize int        `inject:"pool-size"`
}

func (thing *Thing_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies) Dependencies() interface{} {
	// Notice that it is returning a pointer
	// to the dependenies. That's important!
	return []interface{}{ &thing.dependencies1 , &thing.dependencies2 }
}

// TestInjectShouldAlsoInjectPluralHiddenSliceDependencies makes sure that a plural hidden dependency, then it is injected as well.
func TestInjectShouldAlsoInjectPluralHiddenSliceDependencies(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	thing := new(Thing_TestInjectShouldAlsoInjectPluralHiddenSliceDependencies)

	var beforeDependencies1ShouldNotBeInjected int        = randomness.Int()
	var beforeDependencies1HiddenLogger       *log.Logger = nil
	var beforeDependencies1EasterEgg           string     = fmt.Sprintf("%d", randomness.Int())

	var beforeDependencies2Out io.Writer = nil

	var beforeLength   int        = randomness.Int()
	var beforeLogger  *log.Logger = nil
	var beforePoolSize int        = randomness.Int()

	thing.dependencies1.ShouldNotBeInjected = beforeDependencies1ShouldNotBeInjected
	thing.dependencies1.HiddenLogger        = beforeDependencies1HiddenLogger
	thing.dependencies1.EasterEgg           = beforeDependencies1EasterEgg

	thing.dependencies2.Out = beforeDependencies2Out

	thing.Length   = beforeLength
	thing.Logger   = beforeLogger
	thing.PoolSize = beforePoolSize

	// Create the dependency injection container.
	container := New()

	// Register some stuff into the dependency injection container.
	expectedLogger    := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize  := 20
	expectedEasterEgg := "apple-banana-cherry"

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", expectedEasterEgg)
	container.Register("out", ioutil.Discard)

	// Inject.
	if err := container.Inject(thing); nil != err {
		t.Errorf("Received error when injecting. Error: %q", err)
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			missingDependencyNames := complainer.MissingDependencyNames()
			t.Errorf("Had %d missing dependencies: %#v", len(missingDependencyNames), missingDependencyNames)
		default:
			t.Errorf("Unknown complaint from injecting.")
		}
		return
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeDependencies1ShouldNotBeInjected, thing.dependencies1.ShouldNotBeInjected; expected != actual {
		t.Errorf("Expected thing.dependencies.ShouldNotBeInjected to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1HiddenLogger, thing.dependencies1.HiddenLogger; notExpected == actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1EasterEgg, thing.dependencies1.EasterEgg; notExpected == actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}


	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeLogger, thing.Logger; notExpected == actual {
		t.Errorf("Expected thing.Logger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforePoolSize, thing.PoolSize; notExpected == actual {
		t.Errorf("Expected thing.PoolSize to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	// Confirm that the things that changed, due to injection, changed to what we expected them to change to.
	if expected, actual := expectedLogger, thing.dependencies1.HiddenLogger; expected != actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedEasterEgg, thing.dependencies1.EasterEgg; expected != actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to be %q, but actually was %q.", expected, actual)
		return
	}


	if expected, actual := expectedLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedPoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to be %d, but actually was %d.", expected, actual)
		return
	}
}



type thingDependencies1_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies struct {
	ShouldNotBeInjected  int
	HiddenLogger        *log.Logger `inject:"logger"`
	EasterEgg            string     `inject:"easter-egg"`
}

type thingDependencies2_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies struct {
	Out io.Writer `inject:"out"`
}

type Thing_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies struct {
	dependencies1 thingDependencies1_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies
	dependencies2 thingDependencies2_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies
	Length   int
	Logger  *log.Logger `inject:"logger"`
	PoolSize int        `inject:"pool-size"`
}

func (thing *Thing_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies) Dependencies() interface{} {
	// Notice that it is returning a pointer
	// to the dependenies. That's important!
	return [2]interface{}{ &thing.dependencies1 , &thing.dependencies2 }
}

// TestInjectShouldAlsoInjectPluralHiddenArrayDependencies makes sure that a plural hidden dependency, then it is injected as well.
func TestInjectShouldAlsoInjectPluralHiddenArrayDependencies(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	thing := new(Thing_TestInjectShouldAlsoInjectPluralHiddenArrayDependencies)

	var beforeDependencies1ShouldNotBeInjected int        = randomness.Int()
	var beforeDependencies1HiddenLogger       *log.Logger = nil
	var beforeDependencies1EasterEgg           string     = fmt.Sprintf("%d", randomness.Int())

	var beforeDependencies2Out io.Writer = nil

	var beforeLength   int        = randomness.Int()
	var beforeLogger  *log.Logger = nil
	var beforePoolSize int        = randomness.Int()

	thing.dependencies1.ShouldNotBeInjected = beforeDependencies1ShouldNotBeInjected
	thing.dependencies1.HiddenLogger        = beforeDependencies1HiddenLogger
	thing.dependencies1.EasterEgg           = beforeDependencies1EasterEgg

	thing.dependencies2.Out = beforeDependencies2Out

	thing.Length   = beforeLength
	thing.Logger   = beforeLogger
	thing.PoolSize = beforePoolSize

	// Create the dependency injection container.
	container := New()

	// Register some stuff into the dependency injection container.
	expectedLogger    := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize  := 20
	expectedEasterEgg := "apple-banana-cherry"

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", expectedEasterEgg)
	container.Register("out", ioutil.Discard)

	// Inject.
	if err := container.Inject(thing); nil != err {
		t.Errorf("Received error when injecting. Error: %q", err)
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			missingDependencyNames := complainer.MissingDependencyNames()
			t.Errorf("Had %d missing dependencies: %#v", len(missingDependencyNames), missingDependencyNames)
		default:
			t.Errorf("Unknown complaint from injecting.")
		}
		return
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeDependencies1ShouldNotBeInjected, thing.dependencies1.ShouldNotBeInjected; expected != actual {
		t.Errorf("Expected thing.dependencies.ShouldNotBeInjected to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1HiddenLogger, thing.dependencies1.HiddenLogger; notExpected == actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1EasterEgg, thing.dependencies1.EasterEgg; notExpected == actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}


	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeLogger, thing.Logger; notExpected == actual {
		t.Errorf("Expected thing.Logger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforePoolSize, thing.PoolSize; notExpected == actual {
		t.Errorf("Expected thing.PoolSize to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	// Confirm that the things that changed, due to injection, changed to what we expected them to change to.
	if expected, actual := expectedLogger, thing.dependencies1.HiddenLogger; expected != actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedEasterEgg, thing.dependencies1.EasterEgg; expected != actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to be %q, but actually was %q.", expected, actual)
		return
	}


	if expected, actual := expectedLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedPoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to be %d, but actually was %d.", expected, actual)
		return
	}
}



type thingDependencies1_TestInjectShouldAlsoInjectPluralHiddenMapDependencies struct {
	ShouldNotBeInjected  int
	HiddenLogger        *log.Logger `inject:"logger"`
	EasterEgg            string     `inject:"easter-egg"`
}

type thingDependencies2_TestInjectShouldAlsoInjectPluralHiddenMapDependencies struct {
	Out io.Writer `inject:"out"`
}

type Thing_TestInjectShouldAlsoInjectPluralHiddenMapDependencies struct {
	dependencies1 thingDependencies1_TestInjectShouldAlsoInjectPluralHiddenMapDependencies
	dependencies2 thingDependencies2_TestInjectShouldAlsoInjectPluralHiddenMapDependencies
	Length   int
	Logger  *log.Logger `inject:"logger"`
	PoolSize int        `inject:"pool-size"`
}

func (thing *Thing_TestInjectShouldAlsoInjectPluralHiddenMapDependencies) Dependencies() interface{} {
	// Notice that it is returning a pointer
	// to the dependenies. That's important!
	return map[string]interface{}{ "apple": &thing.dependencies1 , "banana": &thing.dependencies2 }
}

// TestInjectShouldAlsoInjectPluralHiddenMapDependencies makes sure that a plural hidden dependency, then it is injected as well.
func TestInjectShouldAlsoInjectPluralHiddenMapDependencies(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	thing := new(Thing_TestInjectShouldAlsoInjectPluralHiddenMapDependencies)

	var beforeDependencies1ShouldNotBeInjected int        = randomness.Int()
	var beforeDependencies1HiddenLogger       *log.Logger = nil
	var beforeDependencies1EasterEgg           string     = fmt.Sprintf("%d", randomness.Int())

	var beforeDependencies2Out io.Writer = nil

	var beforeLength   int        = randomness.Int()
	var beforeLogger  *log.Logger = nil
	var beforePoolSize int        = randomness.Int()

	thing.dependencies1.ShouldNotBeInjected = beforeDependencies1ShouldNotBeInjected
	thing.dependencies1.HiddenLogger        = beforeDependencies1HiddenLogger
	thing.dependencies1.EasterEgg           = beforeDependencies1EasterEgg

	thing.dependencies2.Out = beforeDependencies2Out

	thing.Length   = beforeLength
	thing.Logger   = beforeLogger
	thing.PoolSize = beforePoolSize

	// Create the dependency injection container.
	container := New()

	// Register some stuff into the dependency injection container.
	expectedLogger    := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize  := 20
	expectedEasterEgg := "apple-banana-cherry"

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", expectedEasterEgg)
	container.Register("out", ioutil.Discard)

	// Inject.
	if err := container.Inject(thing); nil != err {
		t.Errorf("Received error when injecting. Error: %q", err)
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			missingDependencyNames := complainer.MissingDependencyNames()
			t.Errorf("Had %d missing dependencies: %#v", len(missingDependencyNames), missingDependencyNames)
		default:
			t.Errorf("Unknown complaint from injecting.")
		}
		return
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeDependencies1ShouldNotBeInjected, thing.dependencies1.ShouldNotBeInjected; expected != actual {
		t.Errorf("Expected thing.dependencies.ShouldNotBeInjected to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1HiddenLogger, thing.dependencies1.HiddenLogger; notExpected == actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1EasterEgg, thing.dependencies1.EasterEgg; notExpected == actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}


	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeLogger, thing.Logger; notExpected == actual {
		t.Errorf("Expected thing.Logger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforePoolSize, thing.PoolSize; notExpected == actual {
		t.Errorf("Expected thing.PoolSize to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	// Confirm that the things that changed, due to injection, changed to what we expected them to change to.
	if expected, actual := expectedLogger, thing.dependencies1.HiddenLogger; expected != actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedEasterEgg, thing.dependencies1.EasterEgg; expected != actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to be %q, but actually was %q.", expected, actual)
		return
	}


	if expected, actual := expectedLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedPoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to be %d, but actually was %d.", expected, actual)
		return
	}
}



type thingDependencies1_TestInjectWithMissingDependencies struct {
	ShouldNotBeInjected   int
	HiddenLogger         *log.Logger `inject:"logger"`
	EasterEgg             string     `inject:"easter-egg"`
	LoggerShoudBeMissing *log.Logger `inject:"not-there"`
}

type thingDependencies2_TestInjectWithMissingDependencies struct {
	Out io.Writer `inject:"out"`
}

type Thing_TestInjectWithMissingDependencies struct {
	dependencies1 thingDependencies1_TestInjectWithMissingDependencies
	dependencies2 thingDependencies2_TestInjectWithMissingDependencies
	Length             int
	ShoudAlsoBeMissing int        `inject:"also-not-there"`
	Logger            *log.Logger `inject:"logger"`
	PoolSize           int        `inject:"pool-size"`
}

func (thing *Thing_TestInjectWithMissingDependencies) Dependencies() interface{} {
	// Notice that it is returning a pointer
	// to the dependenies. That's important!
	return map[string]interface{}{ "apple": &thing.dependencies1 , "banana": &thing.dependencies2 }
}

func TestInjectWithMissingDependencies(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	thing := new(Thing_TestInjectWithMissingDependencies)

	var beforeDependencies1ShouldNotBeInjected   int        = randomness.Int()
	var beforeDependencies1HiddenLogger         *log.Logger = nil
	var beforeDependencies1EasterEgg             string     = fmt.Sprintf("%d", randomness.Int())
	var beforeDependencies1LoggerShoudBeMissing *log.Logger = nil

	var beforeDependencies2Out io.Writer = nil

	var beforeLength             int        = randomness.Int()
	var beforeShoudAlsoBeMissing int        = randomness.Int()
	var beforeLogger            *log.Logger = nil
	var beforePoolSize           int        = randomness.Int()

	thing.dependencies1.ShouldNotBeInjected  = beforeDependencies1ShouldNotBeInjected
	thing.dependencies1.HiddenLogger         = beforeDependencies1HiddenLogger
	thing.dependencies1.EasterEgg            = beforeDependencies1EasterEgg
	thing.dependencies1.LoggerShoudBeMissing = beforeDependencies1LoggerShoudBeMissing

	thing.dependencies2.Out = beforeDependencies2Out

	thing.Length             = beforeLength
	thing.ShoudAlsoBeMissing = beforeShoudAlsoBeMissing
	thing.Logger             = beforeLogger
	thing.PoolSize           = beforePoolSize

	// Create the dependency injection container.
	container := New()

//container2 := New()
//container2.Register("logger", log.New(os.Stdout, "DI CONTAINER TEST> ", log.Lshortfile))
//container2.Inject(container)

	// Register some stuff into the dependency injection container.
	expectedLogger    := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize  := 20
	expectedEasterEgg := "apple-banana-cherry"

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", expectedEasterEgg)
	container.Register("out", ioutil.Discard)

	// Inject.
	//
	// We expect it to return an error here!
	if err := container.Inject(thing); nil == err {
		t.Errorf("We expected Inject to return an error, due to the missing depedencies, but didn't get one.")
		return
	} else if complainer, ok := err.(DependenciesNotFoundComplainer); !ok {
		t.Errorf("We expected the error returned by Inject to fit the MissingDependenciesComplainer interface, but it didn't.")
		return
	} else {

		missingDependencyNames := complainer.MissingDependencyNames()

		if actual, expected := len(missingDependencyNames), 2; expected != actual {
			t.Errorf("We expected to get %d missing dependency names, but actually got %d.", expected, actual)
			return
		}

		inStringSlice := func(x string, slice []string) bool {
			for _, y := range slice {
				if x == y {
					return true
				}
			}

			return false
		}

		if expected, actualCollection := "not-there", missingDependencyNames; !inStringSlice(expected, actualCollection) {
			t.Errorf("We expected %q to be in the missing dependency names %#v, but it wasn't.", expected, actualCollection)
		}
		if expected, actualCollection := "also-not-there", missingDependencyNames; !inStringSlice(expected, actualCollection) {
			t.Errorf("We expected %q to be in the missing dependency names %#v, but it wasn't.", expected, actualCollection)
		}
	}

	// Confirm what was expected to change changed, and what was expected to stay the same stayed the same.
	if expected, actual := beforeDependencies1ShouldNotBeInjected, thing.dependencies1.ShouldNotBeInjected; expected != actual {
		t.Errorf("Expected thing.dependencies.ShouldNotBeInjected to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1HiddenLogger, thing.dependencies1.HiddenLogger; notExpected == actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforeDependencies1EasterEgg, thing.dependencies1.EasterEgg; notExpected == actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if expected, actual := beforeDependencies1LoggerShoudBeMissing, thing.dependencies1.LoggerShoudBeMissing; expected != actual {
		t.Errorf("Expected thing.dependencies.LoggerShoudBeMissing to NOT have changed from %v (since the value for its an \"inject\" struct tag does not have a registration), but didn't. Changed to %v.", expected, actual)
		return
	}


	if expected, actual := beforeLength, thing.Length; expected != actual {
		t.Errorf("Expected thing.Length to NOT have changed from %v (since it doesn't have an \"inject\" struct tag), but didn't. Changed to %v.", expected, actual)
		return
	}

	if expected, actual := beforeShoudAlsoBeMissing, thing.ShoudAlsoBeMissing; expected != actual {
		t.Errorf("Expected thing.ShoudAlsoBeMissing to NOT have changed from %v (since the value for its an \"inject\" struct tag does not have a registration), but didn't. Changed to %v.", expected, actual)
		return
	}

	if notExpected, actual := beforeLogger, thing.Logger; notExpected == actual {
		t.Errorf("Expected thing.Logger to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	if notExpected, actual := beforePoolSize, thing.PoolSize; notExpected == actual {
		t.Errorf("Expected thing.PoolSize to have changed from %v, but didn't. Stayed as %v.", notExpected, actual)
		return
	}

	// Confirm that the things that changed, due to injection, changed to what we expected them to change to.
	if expected, actual := expectedLogger, thing.dependencies1.HiddenLogger; expected != actual {
		t.Errorf("Expected thing.dependencies.HiddenLogger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedEasterEgg, thing.dependencies1.EasterEgg; expected != actual {
		t.Errorf("Expected thing.dependencies.EasterEgg to be %q, but actually was %q.", expected, actual)
		return
	}


	if expected, actual := expectedLogger, thing.Logger; expected != actual {
		t.Errorf("Expected thing.Logger to point to %p, but actually pointed to %p.", expected, actual)
		return
	}

	if expected, actual := expectedPoolSize, thing.PoolSize; expected != actual {
		t.Errorf("Expected thing.PoolSize to be %d, but actually was %d.", expected, actual)
		return
	}

}



func TestInjectWithWrongType(t *testing.T) {

	// Create a random number generator.
	//
	// (There will be some randomness to our testing here.)
	randomness := rand.New(rand.NewSource( time.Now().UTC().UnixNano() ))

	// Create a struct, that we can use as a "target" for our dependency injection.
	type Thing struct {
		WrongType        int         `inject:"logger"`
		AnotherWrongType *log.Logger `inject:"pool-size"`
	}

	thing := new(Thing)

	var beforeWrongType         int        = randomness.Int()
	var beforeAnotherWrongType *log.Logger = nil

	thing.WrongType        = beforeWrongType
	thing.AnotherWrongType = beforeAnotherWrongType

	// Create the dependency injection container.
	container := New()

//container2 := New()
//container2.Register("logger", log.New(os.Stdout, "DI CONTAINER TEST> ", log.Lshortfile))
//container2.Inject(container)

	// Register some stuff into the dependency injection container.
	expectedLogger    := log.New(ioutil.Discard, "we be logging: ", log.Lshortfile)
	expectedPoolSize  := 20
	expectedEasterEgg := "apple-banana-cherry"

	container.Register("logger", expectedLogger)
	container.Register("pool-size", expectedPoolSize)
	container.Register("easter-egg", expectedEasterEgg)
	container.Register("out", ioutil.Discard)


	// Inject.
	//
	// We expect it to return an error here!
	if err := container.Inject(thing); nil == err {
		t.Errorf("We expected Inject to return an error, due to a wrong type on a dependenct, but didn't get one.")
		return
	} else if _, ok := err.(WrongTypeComplainer); !ok {
		t.Errorf("We expected the error returned by Inject to fit the WrongType Complainer interface, but it didn't.")
		return
	} else {
		// We don't have a (good) way of inspecting what it failed on (at least for now)
		// so we won't try checking.
	}
}
