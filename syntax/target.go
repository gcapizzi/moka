package syntax

import "github.com/gcapizzi/moka"

type AllowanceTarget struct {
	double moka.Double
}

func AllowDouble(double moka.Double) AllowanceTarget {
	return AllowanceTarget{double: double}
}

func (t AllowanceTarget) To(interactionBuilder InteractionBuilder) {
	t.double.AddInteraction(interactionBuilder.Build())
}

type ExpectationTarget struct {
	double moka.Double
}

func ExpectDouble(double moka.Double) ExpectationTarget {
	return ExpectationTarget{double: double}
}

func (t ExpectationTarget) To(interactionBuilder InteractionBuilder) {
	t.double.AddInteraction(moka.NewExpectedInteraction(interactionBuilder.Build()))
}

func VerifyCalls(double moka.Double) {
	double.VerifyInteractions()
}
