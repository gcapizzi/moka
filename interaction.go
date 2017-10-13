package moka

import (
	"fmt"
	"reflect"
)

type interaction interface {
	call(methodName string, args []interface{}) ([]interface{}, bool)
	verify() error
	checkType(t reflect.Type) error
}

type allowedInteraction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func newAllowedInteraction(methodName string, args []interface{}, returnValues []interface{}) allowedInteraction {
	return allowedInteraction{methodName: methodName, args: args, returnValues: returnValues}
}

func (i allowedInteraction) call(methodName string, args []interface{}) ([]interface{}, bool) {
	methodNamesAreEqual := i.methodName == methodName
	argsAreEqual := reflect.DeepEqual(i.args, args)

	if methodNamesAreEqual && argsAreEqual {
		return i.returnValues, true
	}

	return nil, false
}

func (i allowedInteraction) verify() error {
	return nil
}

func (i allowedInteraction) String() string {
	return formatMethodCall(i.methodName, i.args)
}

func (i allowedInteraction) checkType(t reflect.Type) error {
	method, methodExists := t.MethodByName(i.methodName)

	if !methodExists {
		return fmt.Errorf("Invalid interaction: type '%s' has no method '%s'", t.Name(), i.methodName)
	}

	expectedArgTypes := methodArgTypes(t, method)

	expectedNumberOfArgs := len(expectedArgTypes)
	numberOfArgs := len(i.args)
	if expectedNumberOfArgs != numberOfArgs {
		return fmt.Errorf(
			"Invalid interaction: method '%s.%s' takes %d arguments, %d specified",
			t.Name(),
			method.Name,
			expectedNumberOfArgs,
			numberOfArgs,
		)
	}

	for i, arg := range i.args {
		argType := reflect.TypeOf(arg)
		expectedType := expectedArgTypes[i]
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

	expectedNumberOfReturnValues := method.Type.NumOut()
	numberOfReturnValues := len(i.returnValues)
	if method.Type.NumOut() != len(i.returnValues) {
		return fmt.Errorf(
			"Invalid interaction: method '%s.%s' returns %d values, %d specified",
			t.Name(),
			method.Name,
			expectedNumberOfReturnValues,
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

type expectedInteraction struct {
	interaction interaction
	called      bool
}

func newExpectedInteraction(interaction interaction) *expectedInteraction {
	return &expectedInteraction{interaction: interaction}
}

func (i *expectedInteraction) call(methodName string, args []interface{}) ([]interface{}, bool) {
	returnValues, matches := i.interaction.call(methodName, args)
	i.called = matches
	return returnValues, matches
}

func (i *expectedInteraction) verify() error {
	if !i.called {
		return fmt.Errorf("Expected interaction: %s", i.interaction)
	}

	return nil
}

func (i *expectedInteraction) checkType(t reflect.Type) error {
	return nil
}

func assignable(leftType, rightType reflect.Type) bool {
	if leftType == nil {
		return isNillable(rightType)
	}

	return leftType.AssignableTo(rightType)
}

func isNillable(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr ||
		t.Kind() == reflect.Func ||
		t.Kind() == reflect.Interface ||
		t.Kind() == reflect.Slice ||
		t.Kind() == reflect.Chan ||
		t.Kind() == reflect.Map
}

func typeString(t reflect.Type) string {
	if t == nil {
		return "nil"
	}

	return t.String()
}

func methodArgTypes(t reflect.Type, method reflect.Method) []reflect.Type {
	argTypes := []reflect.Type{}
	fromIndex := 0
	if t.Kind() != reflect.Interface {
		fromIndex = 1
	}

	for i := fromIndex; i < method.Type.NumIn(); i++ {
		argTypes = append(argTypes, method.Type.In(i))
	}

	return argTypes
}
