package container


import (
	"testing"
)


func TestNew(t *testing.T) {

	container := New()

	if actual := container; nil == actual {
		t.Errorf("Received %v when trying to create new container. But should NOT have been nil.", actual)
		return
	}

	if iContainer, ok := container.(*internalContainer); !ok {
		t.Errorf("Underlying implementation for container should have been \"internalContainer\" but wasn't.")
		return
	} else if nil == iContainer {
		t.Errorf("Pointer to \"internalContainer\" is %v. But should NOT have been nil.", iContainer)
		return
	} else if actual, expected  := len(iContainer.registry), 0; expected != actual {
		t.Errorf("Expected containers registery to be empty (i.e., to have %d elements in it), but actually had %d.", expected, actual)
		return
	}
}
