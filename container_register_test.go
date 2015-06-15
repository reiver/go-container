package container


import (
	"testing"

	"bytes"
	"io/ioutil"
	"log"
)


func TestRegister(t *testing.T) {

	tests := []struct {
		Name  string
		Thing interface{}
	}{
		{
			Name:"logger",
			Thing:log.New(ioutil.Discard, "we be logging: ", log.Lshortfile),
		},
		{
			Name:"com.example.logger",
			Thing:log.New(ioutil.Discard, "LOG: ", log.LstdFlags),
		},
		{
			Name:"dev-null",
			Thing:ioutil.Discard,
		},
		{
			Name:"buffer",
			Thing:new(bytes.Buffer),
		},
		{
			Name:"fruit",
			Thing:func() string {return "apple-banana-cherry"},
		},
		{
			Name:"pool-size",
			Thing:12,
		},
		{
			Name:"nothing",
			Thing:nil,
		},
		{
			Name:"TRUE",
			Thing:true,
		},
		{
			Name:"FALSE",
			Thing:false,
		},
		{
			Name:"weight",
			Thing:178.3,
		},
		{
			Name:"get-real",
			Thing:0+1i,
		},
	}


	container := New()
	if actual := container; nil == actual {
		t.Errorf("Received %v when trying to create new container. But should NOT have been nil.", actual)
		return
	}


	var iContainer *internalContainer
	var ok bool
	iContainer, ok = container.(*internalContainer)
	if !ok {
                t.Errorf("Underlying implementation for container should have been \"internalContainer\" but wasn't.")
		return
	}
        if nil == iContainer {
                t.Errorf("Pointer to \"internalContainer\" is %v. But should NOT have been nil.", iContainer)
		return
        }


	if actual, expected := len(iContainer.registry), 0; expected != actual {
		t.Errorf("Expected containers registery to be empty (i.e., to have %d elements in it), but actually had %d.", expected, actual)
		return
	}


	for testNumber,test := range tests {

		if err := container.Register(test.Name, test.Thing); nil != err {
			t.Errorf("For test #%d, Received an error when trying to register something: (%T) %v.", testNumber, err, err)
			return
		}

		if actual, expected := len(iContainer.registry), 1+testNumber; expected != actual {
			t.Errorf("For test #%d, expected containers registery to be empty (i.e., to have %d elements in it), but actually had %d.", testNumber, expected, actual)
			return
		}

	}
}
