package moka

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InteractionValidator", func() {
	Describe("NullInteractionValidator", func() {
		var nullInteractionValidator nullInteractionValidator

		BeforeEach(func() {
			nullInteractionValidator = newNullInteractionValidator()
		})

		It("never returns an error", func() {
			Expect(nullInteractionValidator.validate(nil)).To(BeNil())
			Expect(nullInteractionValidator.validate(newFakeInteraction(nil, false, nil, nil))).To(BeNil())
		})
	})

	Describe("TypeInteractionValidator", func() {
		var fakeInteraction *fakeInteraction
		var typeInteractionValidator typeInteractionValidator

		BeforeEach(func() {
			fakeInteraction = newFakeInteraction(nil, false, nil, errors.New("CheckType failed"))
			typeInteractionValidator = newTypeInteractionValidator(reflect.TypeOf(someType{}))
		})

		It("checks the interaction against the type", func() {
			err := typeInteractionValidator.validate(fakeInteraction)

			Expect(err).To(MatchError("CheckType failed"))
			Expect(fakeInteraction.checkTypeCalled).To(BeTrue())
			Expect(fakeInteraction.receivedType).To(Equal(reflect.TypeOf(someType{})))
		})
	})
})

type someType struct{}
