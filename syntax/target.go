package syntax

import "github.com/gcapizzi/moka"

type AllowanceTarget struct {
	double moka.Double
}

func AllowDouble(double moka.Double) AllowanceTarget {
	return AllowanceTarget{double: double}
}

func (t AllowanceTarget) To(interaction Interaction) {
	t.double.AllowCall(interaction.methodName, interaction.args, interaction.returnValues)
}

type ExpectationTarget struct {
	double moka.Double
}

func ExpectDouble(double moka.Double) ExpectationTarget {
	return ExpectationTarget{double: double}
}

func (t ExpectationTarget) To(interaction Interaction) {
	t.double.ExpectCall(interaction.methodName, interaction.args, interaction.returnValues)
}

func VerifyCalls(double moka.Double) {
	double.VerifyCalls()
}
