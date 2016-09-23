package container


import (
	"fmt"
)


type ProblemInjectingDependencyComplainer interface {
	error
	ProblemInjectingDependencyComplainer()
	Err() error
}


type internalProblemInjectingDependencyComplainer struct {
	dependencyName string
	err error
}


func newProblemInjectingDependencyComplainer(dependencyName string, err error) error {
	complainer := internalProblemInjectingDependencyComplainer{
		dependencyName:dependencyName,
		err:err,
	}

	return &complainer
}


func (complainer *internalProblemInjectingDependencyComplainer) Error() string {
	return fmt.Sprintf("Problem injecting dependency %q: %v", complainer.dependencyName, complainer.err)
}


func (complainer *internalProblemInjectingDependencyComplainer) ProblemInjectingDependencyComplainer() {
	// Nothing here.
}


func (complainer *internalProblemInjectingDependencyComplainer) Err() error {
	return complainer.err
}
