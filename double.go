package moka

import (
	"fmt"
)

type Double interface {
	AddInteraction(interaction Interaction)
	Call(methodName string, args ...interface{}) []interface{}
	VerifyInteractions()
}

type StrictDouble struct {
	interactions []Interaction
	failHandler  FailHandler
}

type FailHandler func(message string)

func NewStrictDouble() *StrictDouble {
	return &StrictDouble{interactions: []Interaction{}}
}

func NewStrictDoubleWithFailHandler(failHandler FailHandler) *StrictDouble {
	return &StrictDouble{interactions: []Interaction{}, failHandler: failHandler}
}

func (d *StrictDouble) Call(methodName string, args ...interface{}) []interface{} {
	for _, interaction := range d.interactions {
		interactionReturnValues, interactionMatches := interaction.Call(methodName, args)
		if interactionMatches {
			return interactionReturnValues
		}
	}

	d.failHandler(fmt.Sprintf("Unexpected interaction: %s", FormatMethodCall(methodName, args)))
	return nil
}

func (d *StrictDouble) AddInteraction(interaction Interaction) {
	d.interactions = append(d.interactions, interaction)
}

func (d *StrictDouble) VerifyInteractions() {
	for _, interaction := range d.interactions {
		err := interaction.Verify()
		if err != nil {
			d.failHandler(err.Error())
		}
	}
}
