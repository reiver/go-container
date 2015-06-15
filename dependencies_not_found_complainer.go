package container


import (
	"bytes"
	"fmt"
	"io"
)


const dependenciesNotFoundMessagePrefix = "Dependencies not found"


// DependenciesNotFoundComplainer is an 'error' (since it requires the same
// Error method as the 'error' interface) that represents the situation where
// the 'dependency injection container' tries to satisfy one or more dependencies
// but one or more depdencies are NOT registered with the container.
//
// How you might encouter this is with code like the following:
//
//	// Create dependency injection container.
//	container := container.New()
//	
//	// ...
//	//
//	// The "..." represents the code where the container.Register()s
//	// would happen. And the type and instantiation of the variable
//	// 'thing' would occur.
//	//
//	// ...
//	
//	// (Try to) inject dependencies into struct.
//	err := container.Inject(thing)
//	
//	// Deal with any possible errors from the attempted
//	// injection of dependencies.
//	if nil != err {
//		switch complainer := err.(type) {
//		case DependenciesNotFoundComplainer:
//			//@TODO
//		default:
//			//@TODO
//		}
//	}
//
// You can get a list of the the missing dependency names
// by caling the MissingDependencyNames method.
type DependenciesNotFoundComplainer interface {
	Error() string
	MissingDependencyNames() []string
}

// internalDependenciesNotFoundComplainer is the only underlying implementation that fits the
// DependenciesNotFoundComplainer interface, in this library.
type internalDependenciesNotFoundComplainer struct {
	missingDependencyNames map[string]struct{}
}

// newDependenciesNotFoundComplainer creates a new internalDependenciesNotFoundComplainer (struct) and
// returns it as a DependenciesNotFoundComplainer (interface).
func newDependenciesNotFoundComplainer(missingDependencies ...string) DependenciesNotFoundComplainer {
	missingDependencyNames := make(map[string]struct{})

	if 0 < len(missingDependencies) {
		for _, name := range missingDependencies {
			missingDependencyNames[name] = struct{}{}
		}
	}

	err := internalDependenciesNotFoundComplainer{
		missingDependencyNames:missingDependencyNames,
	}

	return &err
}


// Error method is necessary to satisfy the 'error' interface (and the DependenciesNotFoundComplainer
// interface).
func (err *internalDependenciesNotFoundComplainer) Error() string {
	var buffer bytes.Buffer

	io.WriteString(&buffer, dependenciesNotFoundMessagePrefix)
	i := 0
	for key,_ := range err.missingDependencyNames {
		if 0 == i {
			io.WriteString(&buffer, ": ")
		} else {
			io.WriteString(&buffer, ", ")
		}

		io.WriteString(&buffer, fmt.Sprintf("%q", key))

		i++
	}

	return buffer.String()
}

// MissingDependencyNames method is necessary to satisfy the 'DependenciesNotFoundComplainer' interface.
func (err *internalDependenciesNotFoundComplainer) MissingDependencyNames() []string {
	sliceLength := len(err.missingDependencyNames)

	slice := make([]string, sliceLength)

	i := 0
	for key := range err.missingDependencyNames {
		slice[i] = key
		i++
	}

	return slice
}


// concatenate is a helper method that adds the missing dependencies other another DependenciesNotFoundComplainer
// and adds it to its own.
func (err *internalDependenciesNotFoundComplainer) concatenate(otherComplainer DependenciesNotFoundComplainer) {
	// Get the other missing dependencies complainer's missing dependency names.
	moreMissingDependencyNames := otherComplainer.MissingDependencyNames()

	// If the other complainer is empty -- actually has no missing dependencies  -- then
	// there is nothing to do. So just return.
	if nil == moreMissingDependencyNames || 0 >= len(moreMissingDependencyNames) {
		return
	}

	for _, name := range moreMissingDependencyNames {
		err.missingDependencyNames[name] = struct{}{}
	}
}


// insert is a helper method that adds in a missing dependency.
func (err *internalDependenciesNotFoundComplainer) insert(dependencyName string) {
	err.missingDependencyNames[dependencyName] = struct{}{}
}


// len is a helper method that returns the count of the number of missing dependencies.
func (err *internalDependenciesNotFoundComplainer) len() int {
	return len(err.missingDependencyNames)
}
