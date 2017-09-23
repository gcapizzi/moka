package moka

type Double interface {
	StubMethod(methodName string, args []interface{}, returnValue interface{})
	Call(methodName string, args ...interface{}) ReturnValues
}

type StrictDouble struct {
	stubs map[string]interface{}
}

func NewStrictDouble() StrictDouble {
	return StrictDouble{stubs: make(map[string]interface{})}
}

func (d StrictDouble) Call(methodName string, args ...interface{}) ReturnValues {
	return ReturnValues{d.stubs[methodName]}
}

func (d StrictDouble) StubMethod(methodName string, args []interface{}, returnValue interface{}) {
	d.stubs[methodName] = returnValue
}

type ReturnValues []interface{}

func (r ReturnValues) Get(index int) interface{} {
	return r[index]
}
