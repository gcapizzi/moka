package moka_test

import (
	"testing"

	"github.com/gcapizzi/moka"
)

func TestQueryDelegation(t *testing.T) {
	m := moka.New(t)

	collaborator := NewCollaboratorDouble()
	m.Allow(collaborator).ToReceive("Query").With("arg").AndReturn("result")

	subject := NewSubject(collaborator)

	queryResult := subject.DelegateQuery("arg")

	if queryResult != "result" {
		t.Errorf("Query result: '%s'", queryResult)
	}
}

type Collaborator interface {
	Query(string) string
}

type CollaboratorDouble struct{}

func NewCollaboratorDouble() CollaboratorDouble {
	return CollaboratorDouble{}
}

func (c CollaboratorDouble) Query(arg string) string {
	return ""
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
