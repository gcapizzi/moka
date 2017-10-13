package moka

import "reflect"

type interactionValidator interface {
	validate(interaction interaction) error
}

type typeInteractionValidator struct {
	t reflect.Type
}

func newTypeInteractionValidator(t reflect.Type) typeInteractionValidator {
	return typeInteractionValidator{t: t}
}

func (v typeInteractionValidator) validate(interaction interaction) error {
	return interaction.checkType(v.t)
}

type nullInteractionValidator struct{}

func newNullInteractionValidator() nullInteractionValidator {
	return nullInteractionValidator{}
}

func (v nullInteractionValidator) validate(interaction interaction) error {
	return nil
}
