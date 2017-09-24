package moka_test

import (
	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testFailHandlerInvoked bool
var testFailMessage string

func testFailHandler(message string) {
	testFailHandlerInvoked = true
	testFailMessage = message
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

	Context("when a stubbed method is called", func() {
		BeforeEach(func() {
			double.StubMethod(
				"UltimateQuestion",
				[]interface{}{"life", "universe", "everything"},
				[]interface{}{42, nil},
			)
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
				Expect(testFailMessage).To(Equal("No stub for method 'UltimateQuestion' with arguments [foo bar]"))
			})
		})
	})

	Context("when an unknown method is called", func() {
		It("returns nil and makes the test fail", func() {
			returnValues := double.Call("UnstubbedMethod")

			Expect(returnValues).To(BeNil())
			Expect(testFailHandlerInvoked).To(BeTrue())
			Expect(testFailMessage).To(Equal("No stub for method 'UnstubbedMethod' with arguments []"))
		})
	})
})
