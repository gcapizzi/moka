package moka

type Double interface {
	StubMethod(methodName string, args []interface{}, returnValue interface{})
	Call(methodName string, args ...interface{}) Result
}

type StrictDouble struct {
	stubs map[string]interface{}
}

func NewStrictDouble() StrictDouble {
	return StrictDouble{stubs: make(map[string]interface{})}
}

func (d StrictDouble) Call(methodName string, args ...interface{}) Result {
	return ConcreteResult{d.stubs[methodName]}
}

func (d StrictDouble) StubMethod(methodName string, args []interface{}, returnValue interface{}) {
	d.stubs[methodName] = returnValue
}

type Result interface {
	Get(int) interface{}
}

type ConcreteResult []interface{}

func (r ConcreteResult) Get(index int) interface{} {
	return r[index]
}
