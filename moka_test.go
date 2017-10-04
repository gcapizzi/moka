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

	It("supports allowing a method call on a double", func() {
		AllowDouble(collaborator).To(ReceiveCallTo("Query").With("arg").AndReturn("result"))

		result := subject.DelegateQuery("arg")

		Expect(result).To(Equal("result"))
	})

	It("supports expecting a method call on a double", func() {
		ExpectDouble(collaborator).To(ReceiveCallTo("Command").With("arg").AndReturn("result"))

		result := subject.DelegateCommand("arg")

		Expect(result).To(Equal("result"))
		VerifyCalls(collaborator)
	})
})

type Collaborator interface {
	Query(string) string
	Command(string) string
}

type CollaboratorDouble struct {
	moka.Double
}

func NewCollaboratorDouble() CollaboratorDouble {
	return CollaboratorDouble{Double: moka.NewStrictDouble()}
}

func (d CollaboratorDouble) Query(arg string) string {
	returnValues, _ := d.Call("Query", arg)
	return returnValues[0].(string)
}

func (d CollaboratorDouble) Command(arg string) string {
	returnValues, _ := d.Call("Command", arg)
	return returnValues[0].(string)
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

func (s Subject) DelegateCommand(arg string) string {
	return s.collaborator.Command(arg)
}
