package moka_test

import (
	"reflect"

	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interaction", func() {
	Describe("AllowedInteraction", func() {
		var interaction Interaction

		BeforeEach(func() {
			interaction = NewAllowedInteraction(
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

		Describe("CheckType", func() {
			var t = reflect.TypeOf((*DeepThought)(nil)).Elem()

			var checkTypeError error

			JustBeforeEach(func() {
				checkTypeError = interaction.CheckType(t)
			})

			Context("when the method exists and all types match", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"UltimateQuestion",
						[]interface{}{"life", "universe", "everything"},
						[]interface{}{42, nil},
					)
				})

				It("succeeds", func() {
					Expect(checkTypeError).NotTo(HaveOccurred())
				})
			})

			Context("when the method is not defined", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"WorstQuestion",
						[]interface{}{"life", "universe", "everything"},
						[]interface{}{42, nil},
					)
				})

				It("fails", func() {
					Expect(checkTypeError).To(MatchError("Invalid interaction: type 'DeepThought' has no method 'WorstQuestion'"))
				})
			})

			Context("when the number of arguments doesn't match", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"UltimateQuestion",
						[]interface{}{"life", "universe"},
						[]interface{}{42, nil},
					)
				})

				It("fails", func() {
					Expect(checkTypeError).To(MatchError("Invalid interaction: method 'DeepThought.UltimateQuestion' takes 3 arguments, 2 specified"))
				})
			})

			Context("when the type of arguments don't match", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"UltimateQuestion",
						[]interface{}{"life", "universe", 0},
						[]interface{}{42, nil},
					)
				})

				It("fails", func() {
					Expect(checkTypeError).To(MatchError("Invalid interaction: type of argument 3 of method 'DeepThought.UltimateQuestion' is 'string', 'int' given"))
				})
			})

			Context("when a nil is specified for a non-nillable type", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"UltimateQuestion",
						[]interface{}{"life", "universe", nil},
						[]interface{}{42, nil},
					)
				})

				It("fails", func() {
					Expect(checkTypeError).To(MatchError("Invalid interaction: type of argument 3 of method 'DeepThought.UltimateQuestion' is 'string', 'nil' given"))
				})
			})

			Context("when the number of return values doesn't match", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"UltimateQuestion",
						[]interface{}{"life", "universe", "everything"},
						[]interface{}{42},
					)
				})

				It("fails", func() {
					Expect(checkTypeError).To(MatchError("Invalid interaction: method 'DeepThought.UltimateQuestion' returns 2 values, 1 specified"))
				})
			})

			Context("when the type of return values don't match", func() {
				BeforeEach(func() {
					interaction = NewAllowedInteraction(
						"UltimateQuestion",
						[]interface{}{"life", "universe", "everything"},
						[]interface{}{"forty-two", nil},
					)
				})

				It("fails", func() {
					Expect(checkTypeError).To(MatchError("Invalid interaction: type of return value 1 of method 'DeepThought.UltimateQuestion' is 'int', 'string' given"))
				})
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
				fakeInteraction = NewFakeInteraction([]interface{}{42, nil}, true, nil, nil)
			})

			It("delegates to the wrapped interaction and records the call for verification", func() {
				Expect(returnValues).To(Equal([]interface{}{42, nil}))
				Expect(matched).To(Equal(true))
				Expect(expectedInteraction.Verify()).To(BeNil())
			})
		})

		Context("when called with unexpected method names or args", func() {
			BeforeEach(func() {
				fakeInteraction = NewFakeInteraction(nil, false, nil, nil)
			})

			It("delegates to the wrapped interaction but doesn't record the call for verification", func() {
				Expect(returnValues).To(BeNil())
				Expect(matched).To(Equal(false))
				Expect(expectedInteraction.Verify()).To(MatchError("Expected interaction: <the-interaction-string-representation>"))
			})
		})
	})
})

type DeepThought interface {
	UltimateQuestion(topicOne, topicTwo, topicThree string) (int, error)
}
