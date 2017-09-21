package moka

type Invocation struct {
	double     Double
	methodName string
	args       []interface{}
}

func (i Invocation) With(args ...interface{}) Invocation {
	return Invocation{double: i.double, methodName: i.methodName, args: args}
}

func (i Invocation) AndReturn(returnValue interface{}) {
	i.double.StubMethod(i.methodName, i.args, returnValue)
}
