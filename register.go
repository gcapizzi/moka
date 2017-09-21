package moka

import "testing"

type Register struct{}

func New(testingT *testing.T) Register {
	return Register{}
}

func (r Register) Allow(double Double) Subject {
	return Subject{double: double}
}

func (r Register) Expect(double Double) Subject {
	return Subject{}
}

func (r Register) Verify() {}
