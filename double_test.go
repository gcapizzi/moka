package moka_test

import (
	"fmt"

	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testFailHandlerInvoked bool
var testFailMessage string

func testFailHandler(message string) {
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
	var double *StrictDouble
	var returnValues []interface{}

	BeforeEach(func() {
		resetTestFail()
		double = NewStrictDoubleWithFailHandler(testFailHandler)
	})

	Context("when an allowed method call is performed", func() {
		BeforeEach(func() {
			double.AddInteraction(NewInteraction(
				"UltimateQuestion",
				[]interface{}{"life", "universe", "everything"},
				[]interface{}{42, nil},
			))
		})

		Context("with the right arguments", func() {
			BeforeEach(func() {
				returnValues = double.Call("UltimateQuestion", "life", "universe", "everything")
			})

			It("returns the stubbed return values", func() {
				Expect(returnValues).To(Equal([]interface{}{42, nil}))
				Expect(testFailHandlerInvoked).To(BeFalse())
			})
		})

		Context("with the wrong arguments", func() {
			BeforeEach(func() {
				returnValues = double.Call("UltimateQuestion", "foo", "bar")
			})

			It("returns nil and makes the test fail", func() {
				Expect(returnValues).To(BeNil())
				Expect(testFailHandlerInvoked).To(BeTrue())
				Expect(testFailMessage).To(Equal("Unexpected interaction: UltimateQuestion(\"foo\", \"bar\")"))
			})
		})
	})

	Context("when a method call is expected", func() {
		BeforeEach(func() {
			double.AddInteraction(NewExpectedInteraction(NewInteraction(
				"MakeMeASandwich",
				[]interface{}{"bacon", "lettuce", "tomatoes"},
				[]interface{}{fmt.Errorf("ain't got no bacon mate")},
			)))
		})

		Context("and it is performed", func() {
			Context("with the right arguments", func() {
				BeforeEach(func() {
					returnValues = double.Call("MakeMeASandwich", "bacon", "lettuce", "tomatoes")
					double.VerifyInteractions()
				})

				It("returns the mocked return values and records the call", func() {
					Expect(returnValues).To(Equal([]interface{}{fmt.Errorf("ain't got no bacon mate")}))
					Expect(testFailHandlerInvoked).To(BeFalse())
				})
			})

			Context("with the wrong arguments", func() {
				BeforeEach(func() {
					returnValues = double.Call("MakeMeASandwich", "peanut butter", "jelly")
					double.VerifyInteractions()
				})

				It("returns nil and makes the test fail", func() {
					Expect(returnValues).To(BeNil())
					Expect(testFailHandlerInvoked).To(BeTrue())
					Expect(testFailMessage).To(Equal("Unexpected interaction: MakeMeASandwich(\"peanut butter\", \"jelly\")"))
				})
			})
		})

		Context("but it is not performed", func() {
			BeforeEach(func() {
				double.VerifyInteractions()
			})

			It("makes the test fail", func() {
				Expect(testFailHandlerInvoked).To(BeTrue())
				Expect(testFailMessage).To(Equal("Expected interaction: MakeMeASandwich(\"bacon\", \"lettuce\", \"tomatoes\")"))
			})
		})
	})

	Context("when an unknown method is called", func() {
		It("returns nil and makes the test fail", func() {
			returnValues := double.Call("UnstubbedMethod")

			Expect(returnValues).To(BeNil())
			Expect(testFailHandlerInvoked).To(BeTrue())
			Expect(testFailMessage).To(Equal("Unexpected interaction: UnstubbedMethod()"))
		})
	})
})
