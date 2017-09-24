package syntax

import "github.com/gcapizzi/moka"

type Target struct {
	double moka.Double
}

func AllowDouble(double moka.Double) AllowanceTarget {
	return Target{double: double}
}

func (t Target) To(invocation Interaction) {
	invocation.Apply(t.double)
}
