package moka

import (
	"fmt"
	"reflect"
)

type Interaction interface {
	Call(methodName string, args []interface{}) ([]interface{}, bool)
	Verify() error
}

type AllowedInteraction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func NewAllowedInteraction(methodName string, args []interface{}, returnValues []interface{}) AllowedInteraction {
	return AllowedInteraction{methodName: methodName, args: args, returnValues: returnValues}
}

func (i AllowedInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	methodNamesAreEqual := i.methodName == methodName
	argsAreEqual := reflect.DeepEqual(i.args, args)

	if methodNamesAreEqual && argsAreEqual {
		return i.returnValues, true
	}

	return nil, false
}

func (i AllowedInteraction) Verify() error {
	return nil
}

func (i AllowedInteraction) String() string {
	return FormatMethodCall(i.methodName, i.args)
}

type ExpectedInteraction struct {
	interaction Interaction
	called      bool
}

func NewExpectedInteraction(interaction Interaction) *ExpectedInteraction {
	return &ExpectedInteraction{interaction: interaction}
}

func (i *ExpectedInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	returnValues, matches := i.interaction.Call(methodName, args)
	i.called = matches
	return returnValues, matches
}

func (i *ExpectedInteraction) Verify() error {
	if !i.called {
		return fmt.Errorf("Expected interaction: %s", i.interaction)
	}

	return nil
}
