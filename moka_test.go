package moka

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Moka", func() {
	var collaborator CollaboratorDouble
	var subject Subject

	var failHandlerCalled bool
	var failHandlerMessage string

	BeforeEach(func() {
		failHandlerCalled = false
		failHandlerMessage = ""
		RegisterDoublesFailHandler(func(message string, _ ...int) {
			failHandlerCalled = true
			failHandlerMessage = message
		})

		collaborator = NewCollaboratorDouble()
		subject = NewSubject(collaborator)
	})

	It("supports allowing a method call on a double", func() {
		AllowDouble(collaborator).To(ReceiveCallTo("Query").With("arg").AndReturn("result"))

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)

		result := subject.DelegateQuery("arg")

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)
		Expect(result).To(Equal("result"))
	})

	It("makes tests fail on unexpected interactions", func() {
		collaborator.Query("unexpected")

		Expect(failHandlerCalled).To(BeTrue())
		Expect(failHandlerMessage).To(Equal("Unexpected interaction: Query(\"unexpected\")"))
	})

	It("supports expecting a method call on a double", func() {
		ExpectDouble(collaborator).To(ReceiveCallTo("Command").With("arg").AndReturn("result", nil))

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)

		result, _ := subject.DelegateCommand("arg")

		Expect(result).To(Equal("result"))

		VerifyCalls(collaborator)

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)
	})

	It("supports allowing a method call on a double without specifying any args", func() {
		AllowDouble(collaborator).To(ReceiveCallTo("Query").AndReturn("result"))

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)

		result := subject.DelegateQuery("anything")

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)
		Expect(result).To(Equal("result"))
	})

	It("supports expecting a method call on a double without specifying any args or return values", func() {
		ExpectDouble(collaborator).To(ReceiveCallTo("CommandWithNoReturnValues"))

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)

		subject.DelegateCommandWithNoReturnValues("arg")
		VerifyCalls(collaborator)

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)
	})

	It("supports allowing a method call on a double with variadic args", func() {
		AllowDouble(collaborator).To(ReceiveCallTo("VariadicQuery").With([]string{"arg1", "arg2", "arg3"}).AndReturn("result"))

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)

		result := subject.DelegateVariadicQuery("arg1", "arg2", "arg3")

		Expect(failHandlerCalled).To(BeFalse(), failHandlerMessage)
		Expect(result).To(Equal("result"))
	})
})

type Collaborator interface {
	Query(string) string
	Command(string) (string, error)
	CommandWithNoReturnValues(string)
	VariadicQuery(...string) string
}

type CollaboratorDouble struct {
	Double
}

func NewCollaboratorDouble() CollaboratorDouble {
	return CollaboratorDouble{Double: NewStrictDoubleWithTypeOf(CollaboratorDouble{})}
}

func (d CollaboratorDouble) Query(arg string) string {
	returnValues, err := d.Call("Query", arg)
	if err != nil {
		return ""
	}

	return returnValues[0].(string)
}

func (d CollaboratorDouble) Command(arg string) (string, error) {
	returnValues, err := d.Call("Command", arg)
	if err != nil {
		return "", nil
	}

	returnedString, _ := returnValues[0].(string)
	returnedError, _ := returnValues[1].(error)

	return returnedString, returnedError
}

func (d CollaboratorDouble) CommandWithNoReturnValues(arg string) {
	d.Call("CommandWithNoReturnValues", arg)
}

func (d CollaboratorDouble) VariadicQuery(args ...string) string {
	returnValues, err := d.Call("VariadicQuery", args)
	if err != nil {
		return ""
	}

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

func (s Subject) DelegateCommand(arg string) (string, error) {
	return s.collaborator.Command(arg)
}

func (s Subject) DelegateVariadicQuery(args ...string) string {
	return s.collaborator.VariadicQuery(args...)
}

func (s Subject) DelegateCommandWithNoReturnValues(arg string) {
	s.collaborator.CommandWithNoReturnValues(arg)
}
