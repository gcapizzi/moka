package moka

import (
	"fmt"
	"reflect"
	"strings"
)

type Interaction interface {
	Call(methodName string, args []interface{}) ([]interface{}, bool)
	Verify() error
}

type allowedInteraction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func NewInteraction(methodName string, args []interface{}, returnValues []interface{}) Interaction {
	return allowedInteraction{methodName: methodName, args: args, returnValues: returnValues}
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

func (i allowedInteraction) String() string {
	stringArgs := []string{}
	for _, arg := range i.args {
		stringArgs = append(stringArgs, fmt.Sprintf("%#v", arg))
	}

	return fmt.Sprintf("%s(%s)", i.methodName, strings.Join(stringArgs, ", "))
}

type expectedInteraction struct {
	interaction Interaction
	called      bool
}

func NewExpectedInteraction(interaction Interaction) Interaction {
	return &expectedInteraction{interaction: interaction}
}

func (i *expectedInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	returnValues, matches := i.interaction.Call(methodName, args)
	i.called = matches
	return returnValues, matches
}

func (i *expectedInteraction) Verify() error {
	if !i.called {
		return fmt.Errorf("Expected interaction: %s", i.interaction)
	}

	return nil
}
