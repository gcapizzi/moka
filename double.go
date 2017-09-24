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

	d.failHandler(fmt.Sprintf("Unexpected call to method '%s' with arguments %v", methodName, args))
	return nil
}

func (d *StrictDouble) AllowCall(methodName string, args []interface{}, returnValues []interface{}) {
	d.interactions = append(d.interactions, allowedInteraction{methodName: methodName, args: args, returnValues: returnValues})
}

func (d *StrictDouble) ExpectCall(methodName string, args []interface{}, returnValues []interface{}) {
	d.interactions = append(d.interactions, &expectedInteraction{allowedInteraction: allowedInteraction{methodName: methodName, args: args, returnValues: returnValues}})
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

type allowedInteraction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func (i allowedInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	methodNamesAreEqual := i.methodName == methodName
	argsAreEqual := reflect.DeepEqual(i.args, args)

	if methodNamesAreEqual && argsAreEqual {
		return i.returnValues, true
	}

	return nil, false
}

func (i allowedInteraction) Verify() error {
	return nil
}

type expectedInteraction struct {
	allowedInteraction
	called bool
}

func (i *expectedInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	returnValues, matches := i.allowedInteraction.Call(methodName, args)
	i.called = matches
	return returnValues, matches
}

func (i *expectedInteraction) Verify() error {
	if !i.called {
		return fmt.Errorf("Expected the method '%s' to be called with arguments %v", i.methodName, i.args)
	}

	return nil
}
