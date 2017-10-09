package moka

import (
	"fmt"
	"reflect"
)

type Interaction interface {
	Call(methodName string, args []interface{}) ([]interface{}, bool)
	Verify() error
	CheckType(t reflect.Type) error
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

func (i AllowedInteraction) CheckType(t reflect.Type) error {
	method, methodExists := t.MethodByName(i.methodName)

	if !methodExists {
		return fmt.Errorf("Invalid interaction: type '%s' has no method '%s'", t.Name(), i.methodName)
	}

	methodNumberOfArgs := method.Type.NumIn()
	numberOfArgs := len(i.args)
	if methodNumberOfArgs != numberOfArgs {
		return fmt.Errorf(
			"Invalid interaction: method '%s.%s' takes %d arguments, %d specified",
			t.Name(),
			method.Name,
			methodNumberOfArgs,
			numberOfArgs,
		)
	}

	for i, arg := range i.args {
		argType := reflect.TypeOf(arg)
		expectedType := method.Type.In(i)
		if !assignable(argType, expectedType) {
			return fmt.Errorf(
				"Invalid interaction: type of argument %d of method '%s.%s' is '%s', '%s' given",
				i+1,
				t.Name(),
				method.Name,
				typeString(expectedType),
				typeString(argType),
			)
		}
	}

	methodNumberOfReturnValues := method.Type.NumOut()
	numberOfReturnValues := len(i.returnValues)
	if method.Type.NumOut() != len(i.returnValues) {
		return fmt.Errorf(
			"Invalid interaction: method '%s.%s' returns %d values, %d specified",
			t.Name(),
			method.Name,
			methodNumberOfReturnValues,
			numberOfReturnValues,
		)
	}

	for i, returnValue := range i.returnValues {
		returnValueType := reflect.TypeOf(returnValue)
		expectedType := method.Type.Out(i)
		if !assignable(returnValueType, expectedType) {
			return fmt.Errorf(
				"Invalid interaction: type of return value %d of method '%s.%s' is '%s', '%s' given",
				i+1,
				t.Name(),
				method.Name,
				typeString(expectedType),
				typeString(returnValueType),
			)
		}
	}

	return nil
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

func (i *ExpectedInteraction) CheckType(t reflect.Type) error {
	return nil
}

func assignable(leftType, rightType reflect.Type) bool {
	if leftType == nil {
		return isNillable(rightType)
	}

	return leftType.AssignableTo(rightType)
}

func isNillable(t reflect.Type) bool {
	return reflect.Zero(t).Interface() == nil
}

func typeString(t reflect.Type) string {
	if t == nil {
		return "nil"
	}

	return t.String()
}
