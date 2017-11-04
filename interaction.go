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

type argsInteraction struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func newArgsInteraction(methodName string, args []interface{}, returnValues []interface{}) argsInteraction {
	return argsInteraction{methodName: methodName, args: args, returnValues: returnValues}
}

func (i argsInteraction) call(methodName string, args []interface{}) ([]interface{}, bool) {
	methodNamesAreEqual := i.methodName == methodName
	argsAreEqual := i.args == nil || reflect.DeepEqual(i.args, args)

	if methodNamesAreEqual && argsAreEqual {
		return i.returnValues, true
	}

	return nil, false
}

func (i argsInteraction) verify() error {
	return nil
}

func (i argsInteraction) String() string {
	return formatMethodCall(i.methodName, i.args)
}

func (i argsInteraction) checkType(t reflect.Type) error {
	method, methodExists := t.MethodByName(i.methodName)

	if !methodExists {
		return fmt.Errorf("Invalid interaction: type '%s' has no method '%s'", t.Name(), i.methodName)
	}

	if i.args != nil {
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
	}

	expectedNumberOfReturnValues := method.Type.NumOut()
	numberOfReturnValues := len(i.returnValues)
	if numberOfReturnValues != expectedNumberOfReturnValues {
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

type bodyInteraction struct {
	methodName string
	body       interface{}
}

func newBodyInteraction(methodName string, body interface{}) bodyInteraction {
	return bodyInteraction{methodName: methodName, body: body}
}

func (i bodyInteraction) call(methodName string, args []interface{}) ([]interface{}, bool) {
	if methodName == i.methodName {
		bodyAsValue := reflect.ValueOf(i.body)
		argsAsValues := interfacesToValues(args)
		returnValuesAsValues := bodyAsValue.Call(argsAsValues)
		returnValuesAsInterfaces := valuesToInterfaces(returnValuesAsValues)
		return returnValuesAsInterfaces, true
	}

	return nil, false
}

func interfacesToValues(interfaces []interface{}) []reflect.Value {
	values := []reflect.Value{}
	for _, i := range interfaces {
		values = append(values, reflect.ValueOf(i))
	}
	return values
}

func valuesToInterfaces(values []reflect.Value) []interface{} {
	interfaces := []interface{}{}
	for _, v := range values {
		interfaces = append(interfaces, v.Interface())
	}
	return interfaces
}

func (i bodyInteraction) verify() error {
	return nil
}

func (i bodyInteraction) checkType(t reflect.Type) error {
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
