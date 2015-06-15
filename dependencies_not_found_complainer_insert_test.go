package container


import (
	"testing"
)


func TestInsert(t *testing.T) {

	tests := []struct {
		Names []string
		ExpectedLens []int
		ExpectedFinalLen int
	}{
		{
			Names:     []string{},
			ExpectedLens: []int{},
			ExpectedFinalLen: 0,
		},
		{
			Names:     []string{"apple"},
			ExpectedLens: []int{1},
			ExpectedFinalLen: 1,
		},
		{
			Names:     []string{"apple", "banana"},
			ExpectedLens: []int{1,       2},
			ExpectedFinalLen: 2,
		},
		{
			Names:     []string{"apple", "banana", "cherry"},
			ExpectedLens: []int{1,       2,        3},
			ExpectedFinalLen: 3,
		},
		{
			Names:     []string{"apple", "banana", "cherry", "one-two-three-four"},
			ExpectedLens: []int{1,       2,        3,        4},
			ExpectedFinalLen: 4,
		},
		{
			Names:     []string{"apple-banana-cherry"},
			ExpectedLens: []int{1},
			ExpectedFinalLen: 1,
		},
		{
			Names:     []string{"cherry", "cherry"},
			ExpectedLens: []int{1,         1},
			ExpectedFinalLen: 1,
		},
		{
			Names:     []string{"cherry", "cherry", "cherry"},
			ExpectedLens: []int{1,         1,       1},
			ExpectedFinalLen: 1,
		},
		{
			Names:     []string{"apple", "banana", "cherry", "apple", "banana", "cherry"},
			ExpectedLens: []int{1,       2,        3,        3,       3,        3},
			ExpectedFinalLen: 3,
		},
		{
			Names:     []string{"apple", "apple", "banana", "banana", "cherry", "cherry"},
			ExpectedLens: []int{1,       1,       2,         2,       3,       3},
			ExpectedFinalLen: 3,
		},
	}


	for testNumber, test := range tests {

		complainer := newDependenciesNotFoundComplainer()

		internalComplainer,_ := complainer.(*internalDependenciesNotFoundComplainer)

		for nameNumber, name := range test.Names {
			internalComplainer.insert(name)

			if actual, expected := len(internalComplainer.missingDependencyNames), test.ExpectedLens[nameNumber]; expected != actual {
				t.Errorf("For test #%d and name number #%d, expected length of the map used to store the missing depedency names to be %d, but actually was %d.\nName: %q.\nTest: %#v", testNumber, nameNumber, expected, actual, name, test)
				return
			}

			if _,ok := internalComplainer.missingDependencyNames[name]; !ok {
				t.Errorf("For test #%d and name number #%d, name %q was expected to be in %#v, but actually wasn't.\nTest: %#v", testNumber, nameNumber, name, internalComplainer.missingDependencyNames, test)
				return
			}
		}

		if expected, actual := test.ExpectedFinalLen, len(internalComplainer.missingDependencyNames); expected != actual {
			t.Errorf("For test %#d, expected the FINAL length of the map used to store the missing dependency names to be %d, but actually was %d.\nTest: %#v", testNumber, expected, actual, test)
			return
		}

		for _, name := range test.Names {
			if _,ok := internalComplainer.missingDependencyNames[name]; !ok {
				t.Errorf("For test #%d, after done all the inserting, name %q was expected to be in %#v, but actually wasn't.\nTest: %#v", testNumber, name, internalComplainer.missingDependencyNames, test)
				return
			}
		}
	}
}
