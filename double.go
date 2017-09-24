package moka

import (
	"fmt"
	"reflect"
)

type Double interface {
	AllowCall(methodName string, args []interface{}, returnValues []interface{})
	ExpectCall(methodName string, args []interface{}, returnValues []interface{})
	Call(methodName string, args ...interface{}) []interface{}
	VerifyCalls()
}

type StrictDouble struct {
	interactions []interaction
	failHandler  FailHandler
}

type FailHandler func(message string)

func NewStrictDouble() *StrictDouble {
	return &StrictDouble{interactions: []interaction{}}
}

func NewStrictDoubleWithFailHandler(failHandler FailHandler) *StrictDouble {
	return &StrictDouble{interactions: []interaction{}, failHandler: failHandler}
}

func (d *StrictDouble) Call(methodName string, args ...interface{}) []interface{} {
	for _, interaction := range d.interactions {
		interactionReturnValues, interactionMatches := interaction.Call(methodName, args)
		if interactionMatches {
			return interactionReturnValues
		}
	}

	d.failHandler(fmt.Sprintf("No stub or mock for method '%s' with arguments %v", methodName, args))
	return nil
}

func (d *StrictDouble) AllowCall(methodName string, args []interface{}, returnValues []interface{}) {
	d.interactions = append(d.interactions, stub{methodName: methodName, args: args, returnValues: returnValues})
}

func (d *StrictDouble) ExpectCall(methodName string, args []interface{}, returnValues []interface{}) {
	d.interactions = append(d.interactions, &mock{methodName: methodName, args: args, returnValues: returnValues})
}

func (d *StrictDouble) VerifyCalls() {
	for _, interaction := range d.interactions {
		err := interaction.Verify()
		if err != nil {
			d.failHandler(err.Error())
		}
	}
}

type interaction interface {
	Call(methodName string, args []interface{}) ([]interface{}, bool)
	Verify() error
}

type stub struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func (s stub) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	methodNamesAreEqual := s.methodName == methodName
	argsAreEqual := reflect.DeepEqual(s.args, args)

	if methodNamesAreEqual && argsAreEqual {
		return s.returnValues, true
	}

	return nil, false
}

func (s stub) Verify() error {
	return nil
}

type mock struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
	called       bool
}

func (m *mock) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	methodNamesAreEqual := m.methodName == methodName
	argsAreEqual := reflect.DeepEqual(m.args, args)

	if methodNamesAreEqual && argsAreEqual {
		m.called = true
		return m.returnValues, true
	}

	return nil, false
}

func (m *mock) Verify() error {
	if !m.called {
		return fmt.Errorf("Expected the method '%s' to be called with arguments %v", m.methodName, m.args)
	}

	return nil
}
