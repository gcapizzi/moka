package moka

import (
	"fmt"
	"reflect"
)

type Double interface {
	StubMethod(methodName string, args []interface{}, returnValues []interface{})
	Call(methodName string, args ...interface{}) []interface{}
}

type StrictDouble struct {
	stubs       []Stub
	failHandler FailHandler
}

type FailHandler func(message string)

func NewStrictDouble() *StrictDouble {
	return &StrictDouble{stubs: []Stub{}}
}

func NewStrictDoubleWithFailHandler(failHandler FailHandler) *StrictDouble {
	return &StrictDouble{stubs: []Stub{}, failHandler: failHandler}
}

func (d *StrictDouble) Call(methodName string, args ...interface{}) []interface{} {
	for _, stub := range d.stubs {
		if stub.Matches(methodName, args) {
			return stub.returnValues
		}
	}

	d.failHandler(fmt.Sprintf("No stub for method '%s' with arguments %v", methodName, args))
	return nil
}

func (d *StrictDouble) StubMethod(methodName string, args []interface{}, returnValues []interface{}) {
	d.stubs = append(d.stubs, Stub{methodName: methodName, args: args, returnValues: returnValues})
}

type Stub struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func (s Stub) Matches(methodName string, args []interface{}) bool {
	methodNamesAreEqual := s.methodName == methodName
	argsAreEqual := reflect.DeepEqual(s.args, args)

	return methodNamesAreEqual && argsAreEqual
}
