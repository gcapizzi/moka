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

// InteractionBuilder provides a fluid interface to build interactions to
// configure on a `Double`
type InteractionBuilder interface {
	build() interaction
}

// ReceiveCallTo allows to specify the method name of the interaction.
func ReceiveCallTo(methodName string) ArgsInteractionBuilder {
	return ArgsInteractionBuilder{methodName: methodName}
}

// ArgsInteractionBuilder allows to build interactions that are defined by a
// method name, a list of arguments and a list of return values
type ArgsInteractionBuilder struct {
	methodName   string
	args         []interface{}
	returnValues []interface{}
}

// With allows to specify the expected arguments of the interaction.
func (b ArgsInteractionBuilder) With(args ...interface{}) ArgsInteractionBuilder {
	return ArgsInteractionBuilder{methodName: b.methodName, args: args, returnValues: b.returnValues}
}

// AndReturn allows to specify the return value of the interaction.
func (b ArgsInteractionBuilder) AndReturn(returnValues ...interface{}) ArgsInteractionBuilder {
	return ArgsInteractionBuilder{methodName: b.methodName, args: b.args, returnValues: returnValues}
}

func (b ArgsInteractionBuilder) build() interaction {
	return newAllowedInteraction(b.methodName, b.args, b.returnValues)
}
