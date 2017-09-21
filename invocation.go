package moka

type Invocation struct{}

func (i Invocation) With(args ...interface{}) Invocation {
	return i
}

func (i Invocation) AndReturn(returnValue interface{}) Invocation {
	return i
}
