package moka_test

import (
	. "github.com/gcapizzi/moka"

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
			Expect(nullInteractionValidator.Validate(NewFakeInteraction(nil, false, nil))).To(BeNil())
		})
	})
})
