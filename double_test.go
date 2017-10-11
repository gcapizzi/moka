package moka

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testFailHandlerInvoked bool
var testFailMessage string

func testFailHandler(message string, callerSkip ...int) {
	if !testFailHandlerInvoked {
		testFailHandlerInvoked = true
		testFailMessage = message
	}
}

func resetTestFail() {
	testFailHandlerInvoked = false
	testFailMessage = ""
}

var _ = Describe("StrictDouble", func() {
	var interactionValidator FakeInteractionValidator
	var double *StrictDouble

	BeforeEach(func() {
		resetTestFail()
	})

	JustBeforeEach(func() {
		double = NewStrictDoubleWithInteractionValidatorAndFailHandler(interactionValidator, testFailHandler)
	})

	Describe("AddInteraction", func() {
		JustBeforeEach(func() {
			double.AddInteraction(NewFakeInteraction([]interface{}{"result"}, true, nil, nil))
		})

		Context("when the interaction is valid", func() {
			BeforeEach(func() {
				interactionValidator = NewFakeInteractionValidator(nil)
			})

			It("succeeds", func() {
				By("not making the test fail", func() {
					Expect(testFailHandlerInvoked).To(BeFalse())
				})

				By("adding the interaction to the double", func() {
					result, err := double.Call("", []interface{}{})

					Expect(result).To(Equal([]interface{}{"result"}))
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the interaction is not valid", func() {
			BeforeEach(func() {
				interactionValidator = NewFakeInteractionValidator(errors.New("invalid interaction!"))
			})

			It("fails", func() {
				By("making the test fail", func() {
					Expect(testFailHandlerInvoked).To(BeTrue())
					Expect(testFailMessage).To(Equal("invalid interaction!"))
				})

				By("not adding the interaction to the double", func() {
					result, err := double.Call("", []interface{}{})

					Expect(result).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("Call", func() {
		var firstInteraction *FakeInteraction
		var secondInteraction *FakeInteraction
		var thirdInteraction *FakeInteraction

		var returnValues []interface{}
		var err error

		BeforeEach(func() {
			interactionValidator = NewFakeInteractionValidator(nil)
		})

		JustBeforeEach(func() {
			double.AddInteraction(firstInteraction)
			double.AddInteraction(secondInteraction)
			double.AddInteraction(thirdInteraction)

			returnValues, err = double.Call("UltimateQuestion", "life", "universe", "everything")
		})

		Context("when some interactions match", func() {
			BeforeEach(func() {
				firstInteraction = NewFakeInteraction(nil, false, nil, nil)
				secondInteraction = NewFakeInteraction([]interface{}{42, nil}, true, nil, nil)
				thirdInteraction = NewFakeInteraction([]interface{}{43, nil}, true, nil, nil)
			})

			It("returns the configured return values", func() {
				By("stopping at the first matching interaction", func() {
					Expect(firstInteraction.CallCalled).To(BeTrue())
					Expect(secondInteraction.CallCalled).To(BeTrue())
					Expect(thirdInteraction.CallCalled).To(BeFalse())
				})

				By("returning its return values", func() {
					Expect(returnValues).To(Equal([]interface{}{42, nil}))
				})

				By("not returning an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				By("not invoking the fail handler", func() {
					Expect(testFailHandlerInvoked).To(BeFalse())
				})
			})
		})

		Context("when no interaction matches", func() {
			BeforeEach(func() {
				firstInteraction = NewFakeInteraction(nil, false, nil, nil)
				secondInteraction = NewFakeInteraction(nil, false, nil, nil)
				thirdInteraction = NewFakeInteraction(nil, false, nil, nil)
			})

			It("makes the test fail", func() {
				By("calling all interactions", func() {
					Expect(firstInteraction.CallCalled).To(BeTrue())
					Expect(secondInteraction.CallCalled).To(BeTrue())
					Expect(thirdInteraction.CallCalled).To(BeTrue())
				})

				By("returning nil", func() {
					Expect(returnValues).To(BeNil())
				})

				By("calling the fail handler", func() {
					Expect(testFailHandlerInvoked).To(BeTrue())
					Expect(testFailMessage).To(Equal("Unexpected interaction: UltimateQuestion(\"life\", \"universe\", \"everything\")"))
				})

				By("returning an error", func() {
					Expect(err).To(MatchError("Unexpected interaction: UltimateQuestion(\"life\", \"universe\", \"everything\")"))
				})
			})
		})
	})

	Describe("VerifyInteractions", func() {
		var firstInteraction *FakeInteraction
		var secondInteraction *FakeInteraction
		var thirdInteraction *FakeInteraction

		JustBeforeEach(func() {
			double.AddInteraction(firstInteraction)
			double.AddInteraction(secondInteraction)
			double.AddInteraction(thirdInteraction)

			double.VerifyInteractions()
		})

		Context("when all interactions are verified", func() {
			BeforeEach(func() {
				firstInteraction = NewFakeInteraction(nil, false, nil, nil)
				secondInteraction = NewFakeInteraction(nil, false, nil, nil)
				thirdInteraction = NewFakeInteraction(nil, false, nil, nil)
			})

			It("lets the test pass", func() {
				By("verifying all interactions", func() {
					Expect(firstInteraction.VerifyCalled).To(BeTrue())
					Expect(secondInteraction.VerifyCalled).To(BeTrue())
					Expect(thirdInteraction.VerifyCalled).To(BeTrue())
				})

				By("not invoking the fail handler", func() {
					Expect(testFailHandlerInvoked).To(BeFalse())
				})
			})
		})

		Context("when some interactions are not verified", func() {
			BeforeEach(func() {
				firstInteraction = NewFakeInteraction(nil, false, nil, nil)
				secondInteraction = NewFakeInteraction(nil, false, errors.New("nope"), nil)
				thirdInteraction = NewFakeInteraction(nil, false, nil, nil)
			})

			It("makes the test fail", func() {
				By("stopping at the first unverified interaction", func() {
					Expect(firstInteraction.VerifyCalled).To(BeTrue())
					Expect(secondInteraction.VerifyCalled).To(BeTrue())
					Expect(thirdInteraction.VerifyCalled).To(BeFalse())
				})

				By("invoking the fail handler", func() {
					Expect(testFailHandlerInvoked).To(BeTrue())
					Expect(testFailMessage).To(Equal("nope"))
				})
			})
		})
	})
})
