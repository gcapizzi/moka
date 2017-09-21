package moka

type Subject struct{}

func (s Subject) ToReceive(methodName string) Invocation {
	return Invocation{}
}
