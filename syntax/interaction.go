package syntax

import "github.com/gcapizzi/moka"

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

func (b InteractionBuilder) Build() moka.Interaction {
	return moka.NewInteraction(b.methodName, b.args, b.returnValues)
}
