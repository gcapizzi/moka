package moka

type Subject struct {
	double Double
}

func (s Subject) ToReceive(methodName string) Invocation {
	return Invocation{double: s.double, methodName: methodName}
}
