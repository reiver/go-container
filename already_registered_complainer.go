package container


import (
	"fmt"
)


type internalAlreadyRegisteredComplainer struct {
	name string
}


func newAlreadyRegisteredComplainer(name string) error {
	err := internalAlreadyRegisteredComplainer{
		name:name,
	}

	return &err
}


func (err *internalAlreadyRegisteredComplainer) Error() string {
	return fmt.Sprintf("Dependency %q is already registered.", err.name)
}
