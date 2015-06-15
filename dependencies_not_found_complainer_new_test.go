package container


import (
	"testing"
)


func TestNewDependenciesNotFoundComplainerWithNoParameters(t *testing.T) {

	complainer := newDependenciesNotFoundComplainer()
	if nil == complainer {
		t.Errorf("Expected the newly returned 'dependencies not found' complainer to NOT be nil, but it was %v.", complainer)
		return
	}

	internalComplainer, ok := complainer.(*internalDependenciesNotFoundComplainer)
	if !ok {
		t.Errorf("The underlying implementation for what was returned from the newDependenciesNotFoundComplainer() func was not what was expected.")
		return
	}

	if actual, expected := len(internalComplainer.missingDependencyNames), 0; expected != actual {
		t.Errorf("Expected the length of the map used to store the missing depedency names to be initialized to %d, but instead it actually was %d.", expected, actual)
		return
	}

	if actual := internalComplainer.missingDependencyNames; nil == actual {
		t.Errorf("Expected the map used to store the missing dependency names to NOT be nil, but it was actually %#v", actual)
		return
	}

	if actual, expected := len(internalComplainer.missingDependencyNames), 0; expected != actual {
		t.Errorf("Expected the length of the map used to store the missing depedency names to be initialized to %d, but instead it actually was %d.", expected, actual)
		return
	}
}


func TestNewDependenciesNotFoundComplainerWithParameters(t *testing.T) {

	tests := []struct{
		Args []string
		ExpectedLen int
	}{
		{
			Args: []string{},
			ExpectedLen: 0,
		},
		{
			Args: []string{"apple"},
			ExpectedLen: 1,
		},
		{
			Args: []string{"apple", "banana"},
			ExpectedLen: 2,
		},
		{
			Args: []string{"apple", "banana", "cherry"},
			ExpectedLen: 3,
		},

		{
			Args: []string{"apple", "apple"},
			ExpectedLen: 1,
		},
		{
			Args: []string{"banana", "banana"},
			ExpectedLen: 1,
		},
		{
			Args: []string{"cherry", "cherry"},
			ExpectedLen: 1,
		},

		{
			Args: []string{"apple", "banana", "cherry"},
			ExpectedLen: 3,
		},
		{
			Args: []string{"banana", "cherry", "apple"},
			ExpectedLen: 3,
		},
		{
			Args: []string{"cherry", "apple", "banana"},
			ExpectedLen: 3,
		},

		{
			Args: []string{"apple", "banana", "cherry", "apple", "banana", "cherry"},
			ExpectedLen: 3,
		},

		{
			Args: []string{"apple", "apple", "banana", "banana", "cherry", "cherry"},
			ExpectedLen: 3,
		},
	}

	for testNumber, test := range tests {

		complainer := newDependenciesNotFoundComplainer(test.Args...)
		if nil == complainer {
			t.Errorf("For test #%d, expected the newly returned 'dependencies not found' complainer to NOT be nil, but it was %v.\nTest: %#v", testNumber, complainer, test)
			return
		}

		internalComplainer, ok := complainer.(*internalDependenciesNotFoundComplainer)
		if !ok {
			t.Errorf("For test #%d, the underlying implementation for what was returned from the newDependenciesNotFoundComplainer() func was not what was expected.\nTest: %#v", testNumber, test)
			return
		}

		if actual := internalComplainer.missingDependencyNames; nil == actual {
			t.Errorf("For test #%d, expected the map used to store the missing dependency names to NOT be nil, but it was actually %#v\nTest: %#v", testNumber, actual, test)
			return
		}

		if actual, expected := len(internalComplainer.missingDependencyNames), test.ExpectedLen; expected != actual {
			t.Errorf("For test #%d, expected the length of the map used to store the missing depedency names to be initialized to %d, but instead it actually was %d.\nTest: %#v", testNumber, expected, actual, test)
			return
		}

		for argNumber, arg := range test.Args {
			if _,ok := internalComplainer.missingDependencyNames[arg]; !ok {
				t.Errorf("For test #%d, arg #%d was not in %#v but should have been.\nTest: %#v", testNumber, argNumber, internalComplainer.missingDependencyNames, test)
				return
			}
		}

	}
}
