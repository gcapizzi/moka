package moka

import (
	"errors"
	"fmt"
	"reflect"
)

// Double is the interface implemented by all Moka double types.
type Double interface {
	addInteraction(interaction interaction)
	Call(methodName string, args ...interface{}) ([]interface{}, error)
	verifyInteractions()
}

// StrictDouble is a strict implementation of the Double interface.
// Any invocation of the `Call` method that won't match any of the configured
// interactions will trigger a test failure and return an error.
type StrictDouble struct {
	interactions         []interaction
	interactionValidator interactionValidator
	failHandler          FailHandler
}

// NewStrictDouble instantiates a new `StrictDouble`, using the global fail
// handler and no validation on the configured interactions.
func NewStrictDouble() *StrictDouble {
	return newStrictDoubleWithInteractionValidatorAndFailHandler(
		newNullInteractionValidator(),
		globalFailHandler,
	)
}

// NewStrictDoubleWithTypeOf instantiates a new `StrictDouble`, using the
// global fail handler and validating that any configured interaction matches
// the specified type.
func NewStrictDoubleWithTypeOf(value interface{}) *StrictDouble {
	return newStrictDoubleWithInteractionValidatorAndFailHandler(
		newTypeInteractionValidator(reflect.TypeOf(value)),
		globalFailHandler,
	)
}

func newStrictDoubleWithInteractionValidatorAndFailHandler(interactionValidator interactionValidator, failHandler FailHandler) *StrictDouble {
	if failHandler == nil {
		panic("You are trying to instantiate a double, but Moka's fail handler is nil.\n" +
			"If you're using Ginkgo, make sure you instantiate your doubles in a BeforeEach(), JustBeforeEach() or It() block.\n" +
			"Alternatively, you may have forgotten to register a fail handler with RegisterDoublesFailHandler().")
	}

	return &StrictDouble{
		interactions:         []interaction{},
		interactionValidator: interactionValidator,
		failHandler:          failHandler,
	}
}

// Call performs a method call on the double. If a matching interaction is
// found, its return values will be returned. If no configured interaction
// matches, an error will be returned.
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
