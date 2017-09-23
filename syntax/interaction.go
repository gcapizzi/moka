package syntax

import "github.com/gcapizzi/moka"

type Interaction struct {
	methodName  string
	args        []interface{}
	returnValue interface{}
}

func ReceiveCallTo(methodName string) Interaction {
	return Interaction{methodName: methodName}
}

func (i Interaction) With(args ...interface{}) Interaction {
	return Interaction{methodName: i.methodName, args: args, returnValue: i.returnValue}
}

func (i Interaction) AndReturn(returnValue interface{}) Interaction {
	return Interaction{methodName: i.methodName, args: i.args, returnValue: returnValue}
}

func (i Interaction) Apply(double moka.Double) {
	double.StubMethod(i.methodName, i.args, i.returnValue)
}
