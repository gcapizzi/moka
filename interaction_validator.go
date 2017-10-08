package moka

type InteractionValidator interface {
	Validate(interaction Interaction) error
}

type NullInteractionValidator struct{}

func NewNullInteractionValidator() NullInteractionValidator {
	return NullInteractionValidator{}
}

func (v NullInteractionValidator) Validate(interaction Interaction) error {
	return nil
}
