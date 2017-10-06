package moka

import (
	"errors"
	"fmt"
	"reflect"
)

type Double interface {
	AddInteraction(interaction Interaction)
	Call(methodName string, args ...interface{}) ([]interface{}, error)
	VerifyInteractions()
}

type StrictDouble struct {
	interactions []Interaction
	failHandler  FailHandler
}

func NewStrictDouble() *StrictDouble {
	return &StrictDouble{interactions: []Interaction{}, failHandler: globalFailHandler}
}

func NewStrictDoubleWithFailHandler(failHandler FailHandler) *StrictDouble {
	return &StrictDouble{interactions: []Interaction{}, failHandler: failHandler}
}

func (d *StrictDouble) Call(methodName string, args ...interface{}) ([]interface{}, error) {
	for _, interaction := range d.interactions {
		interactionReturnValues, interactionMatches := interaction.Call(methodName, args)
		if interactionMatches {
			return interactionReturnValues, nil
		}
	}

	errorMessage := fmt.Sprintf("Unexpected interaction: %s", FormatMethodCall(methodName, args))
	d.failHandler(errorMessage)
	return nil, errors.New(errorMessage)
}

func (d *StrictDouble) AddInteraction(interaction Interaction) {
	d.interactions = append(d.interactions, interaction)
}

func (d *StrictDouble) VerifyInteractions() {
	for _, interaction := range d.interactions {
		err := interaction.Verify()
		if err != nil {
			d.failHandler(err.Error())
			return
		}
	}
}
