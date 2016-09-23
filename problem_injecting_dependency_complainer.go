package container


import (
	"fmt"
	"reflect"
)


type ProblemInjectingDependencyComplainer interface {
	error
	ProblemInjectingDependencyComplainer()
	Err() error
}


type internalProblemInjectingDependencyComplainer struct {
	dependencyName string
	reflectedValue reflect.Value
	err error
}


func newProblemInjectingDependencyComplainer(dependencyName string, reflectedValue reflect.Value, err error) error {
	complainer := internalProblemInjectingDependencyComplainer{
		dependencyName:dependencyName,
		reflectedValue:reflectedValue,
		err:err,
	}

	return &complainer
}


func (complainer *internalProblemInjectingDependencyComplainer) Error() string {
	return fmt.Sprintf("Problem injecting dependency %q with value %v (%s) (%s): %v", complainer.dependencyName, complainer.reflectedValue, complainer.reflectedValue.Kind().String(), complainer.reflectedValue.String(), complainer.err)
}


func (complainer *internalProblemInjectingDependencyComplainer) ProblemInjectingDependencyComplainer() {
	// Nothing here.
}


func (complainer *internalProblemInjectingDependencyComplainer) Err() error {
	return complainer.err
}
