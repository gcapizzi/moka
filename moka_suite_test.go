package moka_test

import (
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
	ReceivedMethodName string
	ReceivedArgs       []interface{}
	CallCalled         bool
	VerifyCalled       bool
	returnValues       []interface{}
	matches            bool
	verifyError        error
}

func NewFakeInteraction(returnValues []interface{}, matches bool, verifyError error) *FakeInteraction {
	return &FakeInteraction{returnValues: returnValues, matches: matches, verifyError: verifyError}
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

func (i *FakeInteraction) String() string {
	return "<the-interaction-string-representation>"
}
