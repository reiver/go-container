package container


import (
	"bytes"
	"fmt"
	"io"
)


const wrongTypeMessagePrefix = "Wrong type for dependency "


type WrongTypeComplainer interface {
	Error() string
	DependencyName() string
}

// internalWrongTypeComplainer is the only underlying implementation that fits the
// WrongTypeComplainer interface, in this library.
type internalWrongTypeComplainer struct {
	dependencyName string
}

// newWrongTypeComplainer creates a new internalWrongTypeComplainer (struct) and
// returns it as a WrongTypeComplainer (interface).
func newWrongTypeComplainer(dependencyName string) WrongTypeComplainer {
	err := internalWrongTypeComplainer{
		dependencyName:dependencyName,
	}

	return &err
}


// Error method is necessary to satisfy the 'error' interface (and the WrongTypeComplainer
// interface).
func (err *internalWrongTypeComplainer) Error() string {
	var buffer bytes.Buffer

	io.WriteString(&buffer, wrongTypeMessagePrefix)
	io.WriteString(&buffer, fmt.Sprintf("%q", err.dependencyName))

	return buffer.String()
}

// DependencyName method is necessary to satisfy the 'WrongTypeComplainer' interface.
func (err *internalWrongTypeComplainer) DependencyName() string {
	return err.dependencyName
}
