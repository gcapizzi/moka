package moka_test

import (
	"errors"

	. "github.com/gcapizzi/moka"

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
	var firstInteraction *FakeInteraction
	var secondInteraction *FakeInteraction
	var thirdInteraction *FakeInteraction
	var double *StrictDouble

	BeforeEach(func() {
		resetTestFail()
		double = NewStrictDoubleWithFailHandler(testFailHandler)
	})

	JustBeforeEach(func() {
		double.AddInteraction(firstInteraction)
		double.AddInteraction(secondInteraction)
		double.AddInteraction(thirdInteraction)
	})

	Describe("Call", func() {
		var returnValues []interface{}
		var err error

		JustBeforeEach(func() {
			returnValues, err = double.Call("UltimateQuestion", "life", "universe", "everything")
		})

		Context("when some interactions match", func() {
			BeforeEach(func() {
				firstInteraction = NewFakeInteraction(nil, false, nil)
				secondInteraction = NewFakeInteraction([]interface{}{42, nil}, true, nil)
				thirdInteraction = NewFakeInteraction([]interface{}{43, nil}, true, nil)
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
				firstInteraction = NewFakeInteraction(nil, false, nil)
				secondInteraction = NewFakeInteraction(nil, false, nil)
				thirdInteraction = NewFakeInteraction(nil, false, nil)
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
		JustBeforeEach(func() {
			double.VerifyInteractions()
		})

		Context("when all interactions are verified", func() {
			BeforeEach(func() {
				firstInteraction = NewFakeInteraction(nil, false, nil)
				secondInteraction = NewFakeInteraction(nil, false, nil)
				thirdInteraction = NewFakeInteraction(nil, false, nil)
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
				firstInteraction = NewFakeInteraction(nil, false, nil)
				secondInteraction = NewFakeInteraction(nil, false, errors.New("nope"))
				thirdInteraction = NewFakeInteraction(nil, false, nil)
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
