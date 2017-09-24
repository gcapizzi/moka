package syntax

import "github.com/gcapizzi/moka"

type AllowanceTarget struct {
	double moka.Double
}

func AllowDouble(double moka.Double) AllowanceTarget {
	return AllowanceTarget{double: double}
}

func (t AllowanceTarget) To(invocation Interaction) {
	t.double.StubMethod(invocation.methodName, invocation.args, invocation.returnValues)
}

type ExpectationTarget struct {
	double moka.Double
}

func ExpectDouble(double moka.Double) ExpectationTarget {
	return ExpectationTarget{double: double}
}

func (t ExpectationTarget) To(invocation Interaction) {
	t.double.MockMethod(invocation.methodName, invocation.args, invocation.returnValues)
}

func VerifyCalls(double moka.Double) {
	double.VerifyCalls()
}
