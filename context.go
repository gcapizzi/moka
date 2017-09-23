package moka

import "testing"

type Context struct{}

func New(testingT *testing.T) Context {
	return Context{}
}

func (Context) Allow(double Double) Subject {
	return Subject{double: double}
}
