package moka

import (
	"errors"
	"fmt"
)

type Double interface {
	AddInteraction(interaction Interaction)
	Call(methodName string, args ...interface{}) ([]interface{}, error)
	VerifyInteractions()
}

type StrictDouble struct {
	interactions         []Interaction
	interactionValidator InteractionValidator
	failHandler          FailHandler
}

func NewStrictDouble() *StrictDouble {
	return &StrictDouble{interactions: []Interaction{}, interactionValidator: NewNullInteractionValidator(), failHandler: globalFailHandler}
}

func NewStrictDoubleWithInteractionValidatorAndFailHandler(interactionValidator InteractionValidator, failHandler FailHandler) *StrictDouble {
	return &StrictDouble{interactions: []Interaction{}, interactionValidator: interactionValidator, failHandler: failHandler}
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
	validationError := d.interactionValidator.Validate(interaction)

	if validationError != nil {
		d.failHandler(validationError.Error())
		return
	}

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
