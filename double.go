package moka

type Double interface {
	StubMethod(methodName string, args []interface{}, returnValues []interface{})
	Call(methodName string, args ...interface{}) []interface{}
}

type StrictDouble struct {
	stubs map[string][]interface{}
}

func NewStrictDouble() StrictDouble {
	return StrictDouble{stubs: make(map[string][]interface{})}
}

func (d StrictDouble) Call(methodName string, args ...interface{}) []interface{} {
	return d.stubs[methodName]
}

func (d StrictDouble) StubMethod(methodName string, args []interface{}, returnValues []interface{}) {
	d.stubs[methodName] = returnValues
}
