// Package moka provides a mocking framework for the Go programming language.
// Moka works very well with the Ginkgo testing framework, but can be easily
// used with any other testing framework, including the testing package from
// the standard library.
package moka

// FailHandler is the type required for Moka fail handler functions. It matches
// the type of the Ginkgo `Fail` function.
type FailHandler func(message string, callerSkip ...int)

var globalFailHandler FailHandler

// RegisterDoublesFailHandler registers a function as the global fail handler
// used by newly instantiated Moka doubles.
func RegisterDoublesFailHandler(failHandler FailHandler) {
	globalFailHandler = failHandler
}

// AllowanceTarget wraps a Double to enable the configuration of allowed
// interactions on it.
type AllowanceTarget struct {
	double Double
}

// AllowDouble wraps a Double in an `AllowanceTarget`.
func AllowDouble(double Double) AllowanceTarget {
	return AllowanceTarget{double: double}
}

// To configures the interaction built by the provided `InteractionBuilder` on
// the wrapped `Double`.
func (t AllowanceTarget) To(interactionBuilder InteractionBuilder) {
	t.double.addInteraction(interactionBuilder.build())
}

// ExpectationTarget wraps a Double to enable the configuration of expected
// interactions on it.
type ExpectationTarget struct {
	double Double
}

// ExpectDouble wraps a Double in an `ExpectationTarget`.
func ExpectDouble(double Double) ExpectationTarget {
	return ExpectationTarget{double: double}
}

// To configures the interaction built by the provided `InteractionBuilder` on
// the wrapped `Double`.
func (t ExpectationTarget) To(interactionBuilder InteractionBuilder) {
	t.double.addInteraction(newExpectedInteraction(interactionBuilder.build()))
}

// VerifyCalls verifies that all expected interactions on the wrapper `Double`
// have actually happened.
func VerifyCalls(double Double) {
	double.verifyInteractions()
}

// InteractionBuilder provides a fluid API to build interactions.
type InteractionBuilder struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

// ReceiveCallTo allows to specify the method name of the interaction.
func ReceiveCallTo(methodName string) InteractionBuilder {
	return InteractionBuilder{methodName: methodName}
}

// With allows to specify the expected arguments of the interaction.
func (b InteractionBuilder) With(args ...interface{}) InteractionBuilder {
	return InteractionBuilder{methodName: b.methodName, args: args, returnValues: b.returnValues}
}

// AndReturn allows to specify the return value of the interaction.
func (b InteractionBuilder) AndReturn(returnValues ...interface{}) InteractionBuilder {
	return InteractionBuilder{methodName: b.methodName, args: b.args, returnValues: returnValues}
}

func (b InteractionBuilder) build() interaction {
	return newAllowedInteraction(b.methodName, b.args, b.returnValues)
}
