package container


import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)


// Container is an abstration that represents a 'dependency injection container'.
//
// Use the New func to get a new (dependenc injection) container.
type Container interface {
	Register(string, interface{}) error

	Get(string) (interface{}, error)

	Inject(interface{}) error
}

type internalContainerDependencies struct {
	Logger *log.Logger `inject:"logger"`
}

type internalContainer struct {
	registry map[string]interface{}
	dependencies internalContainerDependencies
}


func (container *internalContainer) Dependencies() interface{} {
	return &container.dependencies
}


// New returns a new 'dependency injection container'.
func New() Container {
	logger := log.New(ioutil.Discard, "dependency injection container> ", log.Lshortfile)

	registry  := make(map[string]interface{})

	container := internalContainer{
		registry:registry,
		dependencies:internalContainerDependencies{
			Logger:logger,
		},
	}

	return &container
}


func (container *internalContainer) Register(dependencyName string, dependency interface{}) error {

	logger := container.dependencies.Logger

	logger.Printf("[BEGIN] Register(%q, <dependency> %T)", dependencyName, dependency)

	if _,ok := container.registry[dependencyName]; ok {
		err := newAlreadyRegisteredComplainer(dependencyName)

		logger.Printf("[END]   Register(%q, <dependency> %T) with ERROR: %q", dependencyName, dependency, err)
		return err
	}

	container.registry[dependencyName] = dependency

	logger.Printf("[END]   Register(%q, <dependency> %T)", dependencyName, dependency)

	return nil
}


func (container *internalContainer) Get(dependencyName string) (interface{}, error) {

	logger := container.dependencies.Logger

	logger.Printf("[BEGIN] Get(%q)", dependencyName)

	dependency,ok := container.registry[dependencyName]
	if !ok {
		err := newDependenciesNotFoundComplainer(dependencyName)

		logger.Printf("[END]   Get(%q) with ERROR: %q", dependencyName, err)
		return nil, err
	}

	logger.Printf("[END]   Get(%q)", dependencyName)

	return dependency, nil
}


func (container *internalContainer) Inject(thing interface{}) (errr error) {

	logger := container.dependencies.Logger

	logger.Printf("[BEGIN] Inject(??? %T)", thing)

	// Initialize.
	dependenciesNotFoundComplainer := newDependenciesNotFoundComplainer()
	internalDependenciesNotFoundComplainer, _ := dependenciesNotFoundComplainer.(*internalDependenciesNotFoundComplainer)

	// If the 'thing' passed to this Inject method fits a Depender (interface)
	// (and thus has a Dependencies method) then we "inject" what is returned
	// by the Dependencies method in addition to the 'thing' passed as an argument
	// to this method.
	//
	// The point of this is so that a struct can hide where it stores its
	// dependencies.
	if depender,ok := thing.(Depender); ok {
		if otherThing := depender.Dependencies(); nil != otherThing {
			if err := container.inject(otherThing); nil != err {
				switch complainer := err.(type) {
				case DependenciesNotFoundComplainer:
					logger.Printf("[INSIDE] Inject(??? %T) Intermediate error: %q", thing, complainer)
					internalDependenciesNotFoundComplainer.concatenate(complainer)
					logger.Printf("[INSIDE] Inject(??? %T) Accumulative error: %q", thing, dependenciesNotFoundComplainer)
				default:
					logger.Printf("[END]   Inject(??? %T) with ERROR: %q", thing, err)
					return err
				}
			}
		}
	}

	// Inject the actual thing.
	if err := container.inject(thing); nil != err {
		switch complainer := err.(type) {
		case DependenciesNotFoundComplainer:
			logger.Printf("[INSIDE] Inject(??? %T) Intermediate error: %q", thing, complainer)
			internalDependenciesNotFoundComplainer.concatenate(complainer)
			logger.Printf("[INSIDE] Inject(??? %T) Accumulative error: %q", thing, dependenciesNotFoundComplainer)
		default:
			logger.Printf("[END]   Inject(??? %T) with ERROR: %q", thing, err)
			return err
		}
	}

	// If we had any missing dependencies, then return an error.
	if 0 < internalDependenciesNotFoundComplainer.len() {
		err := dependenciesNotFoundComplainer

		logger.Printf("[END]   Inject(??? %T) with ERROR: %q", thing, err)
		return err
	}

	logger.Printf("[END]   Inject(??? %T)", thing)

	// Return (no errors).
	return nil
}

func (container *internalContainer) inject(thing interface{}) error {

	value := reflect.ValueOf(thing)

	switch value.Kind() {
		case reflect.Array:
			return container.injectArrayOrSlice(thing)

		case reflect.Map:
			return container.injectMap(thing)

		case reflect.Slice:
			return container.injectArrayOrSlice(thing)

		case reflect.Ptr:
			return container.injectPtr(thing)

		default:
			// Nothing here.
	}

	// Return (no errors).
	return nil
}

func (container *internalContainer) injectPtr(thing interface{}) error {

	// Initialize.
	dependenciesNotFoundComplainer := newDependenciesNotFoundComplainer().(*internalDependenciesNotFoundComplainer)

	// Reflection!
	value := reflect.ValueOf(thing)
	x := value.Elem()
	typeOfX := x.Type()

	// Go through each field of the struct, and if it has a dependency-tag
	// indicating that a dependency should be injected, then try to do so.
	numFields := x.NumField()
	for i:=0; i<numFields; i++ {
		field := typeOfX.Field(i)

		fieldTag  := field.Tag

		dependencyName := fieldTag.Get("inject")

		// See if the dependency is registered. If it is, then
		// inject it. Else, make a note of that error.
		//
		// Note that reflection will alwys return a value for the
		// 'struct tag' we ask for. Regardless of whether the
		// programmer added our specific 'struct tag' or not.
		//
		// Even if the programmer doesn't include the specific
		// 'struct tag' we are interested in, reflection will
		// still return a value of "" (i.e., empty string) when
		// we ask it for it.
		//
		// Note that this can be potentially confusing because
		// the programmer could include the specific struct
		// tag we are interested in, but give it "" (i.e., the
		// empty string) as a value, and we couldn't tell the
		// difference between that and when the programmer does
		// not include the 'struct tag'.
		//
		// So, because of this, in the else-clause when we are
		// checking for errors, we ignore the case where the
		// 'dependency name' is "" (i.e., the empty string),
		// and do not consider it an error.
		if dependency,ok := container.registry[dependencyName]; ok {
			err := func(value reflect.Value, dependencyName string) (err error) {

				defer func() {

					if r := recover(); nil != r {
						// See if we received a message of the form:
						//
						//	reflect.Set: value of type ??? is not assignable to type ???
						//
						// If we did then we interpret this as the programmer using the
						// dependency container trying to inject a dependency into a struct
						// field of the wrong type.
						//
						// We return a special error for that.
						if s, ok := r.(string); ok {
							needle := "reflect.Set: value of type "

							if strings.HasPrefix(s, needle) {
								needle = " is not assignable to type "

								if strings.Contains(s, needle) {
									err = newWrongTypeComplainer(dependencyName)
									return
								}
							}
						}

						err = fmt.Errorf("%T %v", r, r)
						return
					}

				}()

				value.Set( reflect.ValueOf(dependency) )

				return nil
			}(x.Field(i), dependencyName)

			if nil != err {
				return err
			}
		} else if "" != dependencyName {
			dependenciesNotFoundComplainer.insert(dependencyName)
		}
	}

	// If we had any missing dependencies, then return an error.
	if 0 < dependenciesNotFoundComplainer.len() {
		return dependenciesNotFoundComplainer
	}

	// Return (no errors).
	return nil
}


func (container *internalContainer) injectArrayOrSlice(thing interface{}) error {

	// Initialize.
	dependenciesNotFoundComplainer := newDependenciesNotFoundComplainer().(*internalDependenciesNotFoundComplainer)

	// Reflection!
	value := reflect.ValueOf(thing)

	// Go through each item in the array or slice, and inject each of them.
	length := value.Len()
	for i:=0; i<length; i++ {
		element := value.Index(i)

		if err := container.Inject(element.Interface()); nil != err {
			switch complainer := err.(type) {
			case DependenciesNotFoundComplainer:
				dependenciesNotFoundComplainer.concatenate(complainer)
			default:
				return err
			}
		}
	}

	// If we had any missing dependencies, then return an error.
	if 0 < dependenciesNotFoundComplainer.len() {
		return dependenciesNotFoundComplainer
	}

	// Return (no errors).
	return nil
}


func (container *internalContainer) injectMap(thing interface{}) error {

	// Initialize.
	dependenciesNotFoundComplainer := newDependenciesNotFoundComplainer().(*internalDependenciesNotFoundComplainer)

	// Reflection!
	value := reflect.ValueOf(thing)

	// Go through each item in the map, and inject each of them.
	keys := value.MapKeys()
	for _,key := range keys {
		element := value.MapIndex(key)

		if err := container.Inject(element.Interface()); nil != err {
			switch complainer := err.(type) {
			case DependenciesNotFoundComplainer:
				dependenciesNotFoundComplainer.concatenate(complainer)
			default:
				return err
			}
		}
	}

	// If we had any missing dependencies, then return an error.
	if 0 < dependenciesNotFoundComplainer.len() {
		return dependenciesNotFoundComplainer
	}

	// Return (no errors).
	return nil
}
