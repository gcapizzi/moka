package moka

type Double interface {
	StubMethod(methodName string, args []interface{}, returnValue interface{})
	Call(methodName string, args ...interface{}) Result
}

type ConcreteDouble struct {
	stubs map[string]interface{}
}

func NewConcreteDouble() ConcreteDouble {
	return ConcreteDouble{stubs: make(map[string]interface{})}
}

func (d ConcreteDouble) Call(methodName string, args ...interface{}) Result {
	return ConcreteResult{d.stubs[methodName]}
}

func (d ConcreteDouble) StubMethod(methodName string, args []interface{}, returnValue interface{}) {
	d.stubs[methodName] = returnValue
}

type Result interface {
	Get(int) interface{}
}

type ConcreteResult []interface{}

func (r ConcreteResult) Get(index int) interface{} {
	return r[index]
}
