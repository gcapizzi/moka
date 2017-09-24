package syntax

import "github.com/gcapizzi/moka"

type AllowanceTarget struct {
	double moka.Double
}

func AllowDouble(double moka.Double) AllowanceTarget {
	return AllowanceTarget{double: double}
}

func (t AllowanceTarget) To(invocation Interaction) {
	invocation.Apply(t.double)
}
