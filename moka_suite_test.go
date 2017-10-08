package moka_test

import (
	"reflect"

	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMoka(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterDoublesFailHandler(Fail)
	RunSpecs(t, "Moka Suite")
}

type FakeInteraction struct {
	CallCalled         bool
	ReceivedMethodName string
	ReceivedArgs       []interface{}
	returnValues       []interface{}
	matches            bool
	VerifyCalled       bool
	verifyError        error
	CheckTypeCalled    bool
	ReceivedType       reflect.Type
	checkTypeError     error
}

func NewFakeInteraction(returnValues []interface{}, matches bool, verifyError error, checkTypeError error) *FakeInteraction {
	return &FakeInteraction{returnValues: returnValues, matches: matches, verifyError: verifyError, checkTypeError: checkTypeError}
}

func (i *FakeInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	i.CallCalled = true
	i.ReceivedMethodName = methodName
	i.ReceivedArgs = args
	return i.returnValues, i.matches
}

func (i *FakeInteraction) Verify() error {
	i.VerifyCalled = true
	return i.verifyError
}

func (i *FakeInteraction) CheckType(t reflect.Type) error {
	i.CheckTypeCalled = true
	i.ReceivedType = t
	return i.checkTypeError
}

func (i *FakeInteraction) String() string {
	return "<the-interaction-string-representation>"
}

type FakeInteractionValidator struct {
	validationError error
}

func NewFakeInteractionValidator(validationError error) FakeInteractionValidator {
	return FakeInteractionValidator{validationError: validationError}
}

func (v FakeInteractionValidator) Validate(interaction Interaction) error {
	return v.validationError
}
