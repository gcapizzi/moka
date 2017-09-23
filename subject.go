package moka

type Target struct {
	double Double
}

func (t Target) ToReceive(methodName string) Invocation {
	return Invocation{double: t.double, methodName: methodName}
}
