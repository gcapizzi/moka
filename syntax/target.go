package syntax

import "github.com/gcapizzi/moka"

type Target struct {
	double moka.Double
}

func Allow(double moka.Double) Target {
	return Target{double: double}
}

func (t Target) To(invocation Interaction) {
	invocation.Apply(t.double)
}
