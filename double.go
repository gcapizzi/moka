package moka

import (
	"errors"
	"fmt"
	"reflect"
)

type Double interface {
	addInteraction(interaction interaction)
	Call(methodName string, args ...interface{}) ([]interface{}, error)
	verifyInteractions()
}

type StrictDouble struct {
	interactions         []interaction
	interactionValidator interactionValidator
	failHandler          failHandler
}

func NewStrictDouble() *StrictDouble {
	return &StrictDouble{
		interactions:         []interaction{},
		interactionValidator: newNullInteractionValidator(),
		failHandler:          globalFailHandler,
	}
}

func NewStrictDoubleWithTypeOf(value interface{}) *StrictDouble {
	return &StrictDouble{
		interactions:         []interaction{},
		interactionValidator: newTypeInteractionValidator(reflect.TypeOf(value)),
		failHandler:          globalFailHandler,
	}
}

func newStrictDoubleWithInteractionValidatorAndFailHandler(interactionValidator interactionValidator, failHandler failHandler) *StrictDouble {
	return &StrictDouble{
		interactions:         []interaction{},
		interactionValidator: interactionValidator,
		failHandler:          failHandler,
	}
}

func (d *StrictDouble) Call(methodName string, args ...interface{}) ([]interface{}, error) {
	for _, interaction := range d.interactions {
		interactionReturnValues, interactionMatches := interaction.call(methodName, args)
		if interactionMatches {
			return interactionReturnValues, nil
		}
	}

	errorMessage := fmt.Sprintf("Unexpected interaction: %s", formatMethodCall(methodName, args))
	d.failHandler(errorMessage)
	return nil, errors.New(errorMessage)
}

func (d *StrictDouble) addInteraction(interaction interaction) {
	validationError := d.interactionValidator.validate(interaction)

	if validationError != nil {
		d.failHandler(validationError.Error())
		return
	}

	d.interactions = append(d.interactions, interaction)
}

func (d *StrictDouble) verifyInteractions() {
	for _, interaction := range d.interactions {
		err := interaction.verify()
		if err != nil {
			d.failHandler(err.Error())
			return
		}
	}
}
