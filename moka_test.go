package moka_test

import (
	"github.com/gcapizzi/moka"
	. "github.com/gcapizzi/moka/syntax"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Moka", func() {
	var collaborator CollaboratorDouble
	var subject Subject

	BeforeEach(func() {
		collaborator = NewCollaboratorDouble()
		subject = NewSubject(collaborator)
	})

	It("allows to stub a method on a double", func() {
		AllowDouble(collaborator).To(ReceiveCallTo("Query").With("arg").AndReturn("result"))
		Expect(subject.DelegateQuery("arg")).To(Equal("result"))
	})
})

type Collaborator interface {
	Query(string) string
}

type CollaboratorDouble struct {
	moka.Double
}

func NewCollaboratorDouble() CollaboratorDouble {
	return CollaboratorDouble{Double: moka.NewStrictDouble()}
}

func (d CollaboratorDouble) Query(arg string) string {
	return d.Call("Query", arg)[0].(string)
}

type Subject struct {
	collaborator Collaborator
}

func NewSubject(collaborator Collaborator) Subject {
	return Subject{collaborator: collaborator}
}

func (s Subject) DelegateQuery(arg string) string {
	return s.collaborator.Query(arg)
}
