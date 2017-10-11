package moka

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InteractionValidator", func() {
	Describe("NullInteractionValidator", func() {
		var nullInteractionValidator NullInteractionValidator

		BeforeEach(func() {
			nullInteractionValidator = NewNullInteractionValidator()
		})

		It("never returns an error", func() {
			Expect(nullInteractionValidator.Validate(nil)).To(BeNil())
			Expect(nullInteractionValidator.Validate(NewFakeInteraction(nil, false, nil, nil))).To(BeNil())
		})
	})

	Describe("TypeInteractionValidator", func() {
		var fakeInteraction *FakeInteraction
		var typeInteractionValidator TypeInteractionValidator

		BeforeEach(func() {
			fakeInteraction = NewFakeInteraction(nil, false, nil, errors.New("CheckType failed"))
			typeInteractionValidator = NewTypeInteractionValidator(reflect.TypeOf(SomeType{}))
		})

		It("checks the interaction against the type", func() {
			err := typeInteractionValidator.Validate(fakeInteraction)

			Expect(err).To(MatchError("CheckType failed"))
			Expect(fakeInteraction.CheckTypeCalled).To(BeTrue())
			Expect(fakeInteraction.ReceivedType).To(Equal(reflect.TypeOf(SomeType{})))
		})
	})
})

type SomeType struct{}
