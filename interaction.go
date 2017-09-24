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

type AllowedInteraction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func NewInteraction(methodName string, args []interface{}, returnValues []interface{}) Interaction {
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
	stringArgs := []string{}
	for _, arg := range i.args {
		stringArgs = append(stringArgs, fmt.Sprintf("%#v", arg))
	}

	return fmt.Sprintf("%s(%s)", i.methodName, strings.Join(stringArgs, ", "))
}

type ExpectedInteraction struct {
	Interaction
	called bool
}

func NewExpectedInteraction(interaction Interaction) Interaction {
	return &ExpectedInteraction{Interaction: interaction}
}

func (i *ExpectedInteraction) Call(methodName string, args []interface{}) ([]interface{}, bool) {
	returnValues, matches := i.Interaction.Call(methodName, args)
	i.called = matches
	return returnValues, matches
}

func (i *ExpectedInteraction) Verify() error {
	if !i.called {
		return fmt.Errorf("Expected interaction: %s", i.Interaction)
	}

	return nil
}
