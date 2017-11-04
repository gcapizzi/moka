package moka

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("interaction", func() {
	Describe("argsInteraction", func() {
		var interaction interaction

		Describe("call", func() {
			var matched bool
			var returnValues []interface{}

			BeforeEach(func() {
				interaction = newArgsInteraction(
					"UltimateQuestion",
					[]interface{}{"life", "universe", "everything"},
					[]interface{}{42, nil},
				)
			})

			Context("when both the method name and the args match", func() {
				JustBeforeEach(func() {
					returnValues, matched = interaction.call("UltimateQuestion", []interface{}{"life", "universe", "everything"})
				})

				It("matches and returns its return values", func() {
					Expect(returnValues).To(Equal([]interface{}{42, nil}))
					Expect(matched).To(BeTrue())
				})
			})

			Context("when the method name doesn't match", func() {
				JustBeforeEach(func() {
					returnValues, matched = interaction.call("DomandaFondamentale", []interface{}{"life", "universe", "everything"})
				})

				It("doesn't match and returns nil", func() {
					Expect(returnValues).To(BeNil())
					Expect(matched).To(BeFalse())
				})
			})

			Context("when the arguments don't match", func() {
				JustBeforeEach(func() {
					returnValues, matched = interaction.call("UltimateQuestion", []interface{}{"vita", "universo", "tutto quanto"})
				})

				It("doesn't match and returns nil", func() {
					Expect(returnValues).To(BeNil())
					Expect(matched).To(BeFalse())
				})
			})

			Context("when both method name and the arguments don't match", func() {
				JustBeforeEach(func() {
					returnValues, matched = interaction.call("DomandaFondamentale", []interface{}{"vita", "universo", "tutto quanto"})
				})

				It("doesn't match and returns nil", func() {
					Expect(returnValues).To(BeNil())
					Expect(matched).To(BeFalse())
				})
			})
		})

		Describe("verify", func() {
			BeforeEach(func() {
				interaction = newArgsInteraction("", nil, nil)
			})

			It("does nothing and always returns nil", func() {
				Expect(interaction.verify()).To(BeNil())
			})
		})

		Describe("checkType", func() {
			var TestCheckType = func(t reflect.Type) {
				var checkTypeError error

				JustBeforeEach(func() {
					checkTypeError = interaction.checkType(t)
				})

				Context("when the method is defined and all types match", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
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
						interaction = newArgsInteraction(
							"WorstQuestion",
							[]interface{}{"life", "universe", "everything"},
							[]interface{}{42, nil},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: type '%s' has no method 'WorstQuestion'", t.Name())))
					})
				})

				Context("when the number of arguments doesn't match", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestion",
							[]interface{}{"life", "universe"},
							[]interface{}{42, nil},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: method '%s.UltimateQuestion' takes 3 arguments, 2 specified", t.Name())))
					})
				})

				Context("when the type of some arguments doesn't match", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestion",
							[]interface{}{"life", "universe", 0},
							[]interface{}{42, nil},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: type of argument 3 of method '%s.UltimateQuestion' is 'string', 'int' given", t.Name())))
					})
				})

				Context("when nil is specified for a non-nillable type argument", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestion",
							[]interface{}{"life", "universe", nil},
							[]interface{}{42, nil},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: type of argument 3 of method '%s.UltimateQuestion' is 'string', 'nil' given", t.Name())))
					})
				})

				Context("when nil is specified for a nillable type argument", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestionWithSlice",
							[]interface{}{nil},
							[]interface{}{42, nil},
						)
					})

					It("succeeds", func() {
						Expect(checkTypeError).NotTo(HaveOccurred())
					})
				})

				Context("when the number of return values doesn't match", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestion",
							[]interface{}{"life", "universe", "everything"},
							[]interface{}{42},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: method '%s.UltimateQuestion' returns 2 values, 1 specified", t.Name())))
					})
				})

				Context("when the type of return values don't match", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestion",
							[]interface{}{"life", "universe", "everything"},
							[]interface{}{"forty-two", nil},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: type of return value 1 of method '%s.UltimateQuestion' is 'int', 'string' given", t.Name())))
					})
				})

				Context("when nil is specified for a non-nillable type return value", func() {
					BeforeEach(func() {
						interaction = newArgsInteraction(
							"UltimateQuestion",
							[]interface{}{"life", "universe", "everything"},
							[]interface{}{nil, nil},
						)
					})

					It("fails", func() {
						Expect(checkTypeError).To(MatchError(fmt.Sprintf("Invalid interaction: type of return value 1 of method '%s.UltimateQuestion' is 'int', 'nil' given", t.Name())))
					})
				})
			}

			TestCheckType(reflect.TypeOf((*deepThought)(nil)).Elem())
			TestCheckType(reflect.TypeOf(myDeepThought{}))
		})

		Context("when no arguments are specified", func() {
			var matched bool
			var returnValues []interface{}
			var checkTypeError error

			BeforeEach(func() {
				interaction = newArgsInteraction(
					"UltimateQuestion",
					nil,
					[]interface{}{42, nil},
				)
			})

			Describe("call", func() {
				Context("when the method name matches", func() {
					JustBeforeEach(func() {
						returnValues, matched = interaction.call("UltimateQuestion", []interface{}{"anything"})
					})

					It("matches and returns its return values", func() {
						Expect(returnValues).To(Equal([]interface{}{42, nil}))
						Expect(matched).To(BeTrue())
					})
				})

				Context("when the method name doesn't match", func() {
					JustBeforeEach(func() {
						returnValues, matched = interaction.call("DomandaFondamentale", []interface{}{"anything"})
					})

					It("doesn't match and returns nil", func() {
						Expect(returnValues).To(BeNil())
						Expect(matched).To(BeFalse())
					})
				})
			})

			Describe("checkType", func() {
				JustBeforeEach(func() {
					checkTypeError = interaction.checkType(reflect.TypeOf(myDeepThought{}))
				})

				Context("when the method is defined", func() {
					It("succeeds", func() {
						Expect(checkTypeError).NotTo(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("expectedInteraction", func() {
		var expectedMethodName = "UltimateQuestion"
		var expectedArgs = []interface{}{"life", "universe", "everything"}

		var returnValues []interface{}
		var matched bool

		var fakeInteraction *fakeInteraction
		var expectedInteraction interaction

		JustBeforeEach(func() {
			expectedInteraction = newExpectedInteraction(fakeInteraction)
			returnValues, matched = expectedInteraction.call(expectedMethodName, expectedArgs)
		})

		Context("when called with the expected method name and args", func() {
			BeforeEach(func() {
				fakeInteraction = newFakeInteraction([]interface{}{42, nil}, true, nil, nil)
			})

			It("delegates to the wrapped interaction and records the call for verification", func() {
				Expect(returnValues).To(Equal([]interface{}{42, nil}))
				Expect(matched).To(Equal(true))
				Expect(expectedInteraction.verify()).To(BeNil())
			})
		})

		Context("when called with unexpected method names or args", func() {
			BeforeEach(func() {
				fakeInteraction = newFakeInteraction(nil, false, nil, nil)
			})

			It("delegates to the wrapped interaction but doesn't record the call for verification", func() {
				Expect(returnValues).To(BeNil())
				Expect(matched).To(Equal(false))
				Expect(expectedInteraction.verify()).To(MatchError("Expected interaction: <the-interaction-string-representation>"))
			})
		})
	})
})

type deepThought interface {
	UltimateQuestion(topicOne, topicTwo, topicThree string) (int, error)
	UltimateQuestionWithSlice(things []string) (int, error)
}

type myDeepThought struct{}

func (dt myDeepThought) UltimateQuestion(topicOne, topicTwo, topicThree string) (int, error) {
	return 42, nil
}

func (dt myDeepThought) UltimateQuestionWithSlice(things []string) (int, error) {
	return 42, nil
}
