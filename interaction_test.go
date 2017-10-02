package moka_test

import (
	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interaction", func() {
	Describe("Default (allowed) Interaction", func() {
		var interaction Interaction

		BeforeEach(func() {
			interaction = NewInteraction(
				"UltimateQuestion",
				[]interface{}{"life", "universe", "everything"},
				[]interface{}{42, nil},
			)
		})

		Describe("Call", func() {
			var matched bool
			var returnValues []interface{}

			Context("when both the method name and the args match", func() {
				BeforeEach(func() {
					returnValues, matched = interaction.Call("UltimateQuestion", []interface{}{"life", "universe", "everything"})
				})

				It("matches and returns it return values", func() {
					Expect(returnValues).To(Equal([]interface{}{42, nil}))
					Expect(matched).To(BeTrue())
				})
			})

			Context("when the method name doesn't match", func() {
				BeforeEach(func() {
					returnValues, matched = interaction.Call("DomandaFondamentale", []interface{}{"life", "universe", "everything"})
				})

				It("doesn't match and returns nil", func() {
					Expect(returnValues).To(BeNil())
					Expect(matched).To(BeFalse())
				})
			})

			Context("when the arguments don't match", func() {
				BeforeEach(func() {
					returnValues, matched = interaction.Call("UltimateQuestion", []interface{}{"vita", "universo", "tutto quanto"})
				})

				It("doesn't match and returns nil", func() {
					Expect(returnValues).To(BeNil())
					Expect(matched).To(BeFalse())
				})
			})

			Context("when both method name and the arguments don't match", func() {
				BeforeEach(func() {
					returnValues, matched = interaction.Call("DomandaFondamentale", []interface{}{"vita", "universo", "tutto quanto"})
				})

				It("doesn't match and returns nil", func() {
					Expect(returnValues).To(BeNil())
					Expect(matched).To(BeFalse())
				})
			})
		})

		Describe("Verify", func() {
			It("does nothing and always returns nil", func() {
				Expect(interaction.Verify()).To(BeNil())
			})
		})
	})

	Describe("ExpectedInteraction", func() {
		var expectedMethodName = "UltimateQuestion"
		var expectedArgs = []interface{}{"life", "universe", "everything"}

		var returnValues []interface{}
		var matched bool

		var fakeInteraction *FakeInteraction
		var expectedInteraction Interaction

		JustBeforeEach(func() {
			expectedInteraction = NewExpectedInteraction(fakeInteraction)
			returnValues, matched = expectedInteraction.Call(expectedMethodName, expectedArgs)
		})

		Context("when called with the expected method name and args", func() {
			BeforeEach(func() {
				fakeInteraction = NewFakeInteraction([]interface{}{42, nil}, true, nil)
			})

			It("delegates to the wrapped interaction and records the call for verification", func() {
				Expect(returnValues).To(Equal([]interface{}{42, nil}))
				Expect(matched).To(Equal(true))
				Expect(expectedInteraction.Verify()).To(BeNil())
			})
		})

		Context("when called with unexpected method names or args", func() {
			BeforeEach(func() {
				fakeInteraction = NewFakeInteraction(nil, false, nil)
			})

			It("delegates to the wrapped interaction but doesn't record the call for verification", func() {
				Expect(returnValues).To(BeNil())
				Expect(matched).To(Equal(false))
				Expect(expectedInteraction.Verify()).To(MatchError("Expected interaction: <the-interaction-string-representation>"))
			})
		})
	})
})
