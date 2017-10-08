package moka

import "reflect"

type InteractionValidator interface {
	Validate(interaction Interaction) error
}

type TypeInteractionValidator struct {
	t reflect.Type
}

func NewTypeInteractionValidator(t reflect.Type) TypeInteractionValidator {
	return TypeInteractionValidator{t: t}
}

func (v TypeInteractionValidator) Validate(interaction Interaction) error {
	return interaction.CheckType(v.t)
}

type NullInteractionValidator struct{}

func NewNullInteractionValidator() NullInteractionValidator {
	return NullInteractionValidator{}
}

func (v NullInteractionValidator) Validate(interaction Interaction) error {
	return nil
}
