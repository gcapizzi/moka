package moka

type FailHandler func(message string, callerSkip ...int)

var globalFailHandler FailHandler

func RegisterDoublesFailHandler(failHandler FailHandler) {
	globalFailHandler = failHandler
}

type AllowanceTarget struct {
	double Double
}

func AllowDouble(double Double) AllowanceTarget {
	return AllowanceTarget{double: double}
}

func (t AllowanceTarget) To(interactionBuilder InteractionBuilder) {
	t.double.AddInteraction(interactionBuilder.Build())
}

type ExpectationTarget struct {
	double Double
}

func ExpectDouble(double Double) ExpectationTarget {
	return ExpectationTarget{double: double}
}

func (t ExpectationTarget) To(interactionBuilder InteractionBuilder) {
	t.double.AddInteraction(NewExpectedInteraction(interactionBuilder.Build()))
}

func VerifyCalls(double Double) {
	double.VerifyInteractions()
}

type InteractionBuilder struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

func ReceiveCallTo(methodName string) InteractionBuilder {
	return InteractionBuilder{methodName: methodName}
}

func (b InteractionBuilder) With(args ...interface{}) InteractionBuilder {
	return InteractionBuilder{methodName: b.methodName, args: args, returnValues: b.returnValues}
}

func (b InteractionBuilder) AndReturn(returnValues ...interface{}) InteractionBuilder {
	return InteractionBuilder{methodName: b.methodName, args: b.args, returnValues: returnValues}
}

func (b InteractionBuilder) Build() Interaction {
	return NewAllowedInteraction(b.methodName, b.args, b.returnValues)
}
