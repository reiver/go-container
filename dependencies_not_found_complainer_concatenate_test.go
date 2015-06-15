package container


import (
	"testing"
)


func TestConcatenate(t *testing.T) {

	tests := []struct {
		Complainer    DependenciesNotFoundComplainer
		Complainers []DependenciesNotFoundComplainer
		Expected      DependenciesNotFoundComplainer
		ExpectedFinalLen int
	}{
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer(),
			ExpectedFinalLen: 0,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer(),
			ExpectedFinalLen: 0,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer(),
			ExpectedFinalLen: 0,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("apple")},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana")},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("apple")},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana")},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("apple"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("cherry"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("apple")},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("banana")},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("apple"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("cherry"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("apple")},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("banana")},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("apple"), newDependenciesNotFoundComplainer("apple")},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer("banana")},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("cherry"), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("apple"), newDependenciesNotFoundComplainer("apple")},
			Expected:   newDependenciesNotFoundComplainer("apple"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer("banana")},
			Expected:   newDependenciesNotFoundComplainer("banana"),
			ExpectedFinalLen: 1,
		},
		{
			Complainer: newDependenciesNotFoundComplainer("cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("cherry"), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("cherry"),
			ExpectedFinalLen: 1,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer()},
			Expected:   newDependenciesNotFoundComplainer("apple", "banana"),
			ExpectedFinalLen: 2,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer(), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("apple", "cherry"),
			ExpectedFinalLen: 2,
		},

		{
			Complainer: newDependenciesNotFoundComplainer(),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("banana", "cherry"),
			ExpectedFinalLen: 2,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana"), newDependenciesNotFoundComplainer("cherry")},
			Expected:   newDependenciesNotFoundComplainer("apple", "banana", "cherry"),
			ExpectedFinalLen: 3,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple", "banana"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("banana", "cherry"), newDependenciesNotFoundComplainer("cherry", "apple")},
			Expected:   newDependenciesNotFoundComplainer("apple", "banana", "cherry"),
			ExpectedFinalLen: 3,
		},

		{
			Complainer: newDependenciesNotFoundComplainer("apple", "banana", "cherry"),
			Complainers: []DependenciesNotFoundComplainer{newDependenciesNotFoundComplainer("one", "two"), newDependenciesNotFoundComplainer("foo", "bar")},
			Expected:   newDependenciesNotFoundComplainer("apple", "banana", "cherry", "one", "two", "foo", "bar"),
			ExpectedFinalLen: 7,
		},
	}


	for testNumber, test := range tests {

		originalComplainerCopy := newDependenciesNotFoundComplainer(test.Complainer.MissingDependencyNames()...)

		internalComplainer,_ := test.Complainer.(*internalDependenciesNotFoundComplainer)

		for complainerNumber, complainer := range test.Complainers {
			internalComplainer.concatenate(complainer)

			for _, dependencyName := range complainer.MissingDependencyNames() {
				if _, ok := internalComplainer.missingDependencyNames[dependencyName]; !ok {
					t.Errorf("For test #%d when on complainer #%d, expected dependency name %q to be in %#v, but wasn't.\nTest: %#v", testNumber, complainerNumber, dependencyName, internalComplainer.missingDependencyNames, test)
				}
			}
		}

		if expected, actual := test.ExpectedFinalLen, len(internalComplainer.missingDependencyNames); expected != actual {
			t.Errorf("For test %#d, expected the FINAL length of the map used to store the missing dependency names to be %d, but actually was %d.\nTest: %#v\nCopy of Original Complainer: %#v", testNumber, expected, actual, test, originalComplainerCopy)
			return
		}

		for _, name := range test.Expected.MissingDependencyNames() {
			if _,ok := internalComplainer.missingDependencyNames[name]; !ok {
				t.Errorf("For test #%d, after done all the concatenating, name %q was expected to be in %#v, but actually wasn't.\nTest: %#v\nCopy of Original Complainer: %#v", testNumber, name, internalComplainer.missingDependencyNames, test, originalComplainerCopy)
				return
			}
		}
	}
}
