package moka

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMoka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Moka Suite")
}

type fakeInteraction struct {
	callCalled         bool
	receivedMethodName string
	receivedArgs       []interface{}
	returnValues       []interface{}
	matches            bool
	verifyCalled       bool
	verifyError        error
	checkTypeCalled    bool
	receivedType       reflect.Type
	checkTypeError     error
}

func newFakeInteraction(returnValues []interface{}, matches bool, verifyError error, checkTypeError error) *fakeInteraction {
	return &fakeInteraction{returnValues: returnValues, matches: matches, verifyError: verifyError, checkTypeError: checkTypeError}
}

func (i *fakeInteraction) call(methodName string, args []interface{}) ([]interface{}, bool) {
	i.callCalled = true
	i.receivedMethodName = methodName
	i.receivedArgs = args
	return i.returnValues, i.matches
}

func (i *fakeInteraction) verify() error {
	i.verifyCalled = true
	return i.verifyError
}

func (i *fakeInteraction) checkType(t reflect.Type) error {
	i.checkTypeCalled = true
	i.receivedType = t
	return i.checkTypeError
}

func (i *fakeInteraction) String() string {
	return "<the-interaction-string-representation>"
}

type fakeInteractionValidator struct {
	validationError error
}

func newFakeInteractionValidator(validationError error) fakeInteractionValidator {
	return fakeInteractionValidator{validationError: validationError}
}

func (v fakeInteractionValidator) validate(interaction interaction) error {
	return v.validationError
}
