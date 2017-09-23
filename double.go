package moka

import "fmt"

type Double interface {
	StubMethod(methodName string, args []interface{}, returnValues []interface{})
	Call(methodName string, args ...interface{}) []interface{}
}

type StrictDouble struct {
	stubs       map[string][]interface{}
	failHandler FailHandler
}

type FailHandler func(message string)

func NewStrictDouble() StrictDouble {
	return StrictDouble{stubs: make(map[string][]interface{})}
}

func NewStrictDoubleWithFailHandler(failHandler FailHandler) StrictDouble {
	return StrictDouble{stubs: make(map[string][]interface{}), failHandler: failHandler}
}

func (d StrictDouble) Call(methodName string, args ...interface{}) []interface{} {
	returnValues, methodIsStubbed := d.stubs[methodName]

	if !methodIsStubbed {
		d.failHandler(fmt.Sprintf("No stub for method '%s'", methodName))
		return nil
	}

	return returnValues
}

func (d StrictDouble) StubMethod(methodName string, args []interface{}, returnValues []interface{}) {
	d.stubs[methodName] = returnValues
}
