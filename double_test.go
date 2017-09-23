package moka_test

import (
	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testFailHandlerInvoked = false
var testFailMessage = ""

func testFailHandler(message string) {
	testFailHandlerInvoked = true
	testFailMessage = message
}

func resetTestFail() {
	testFailHandlerInvoked = false
	testFailMessage = ""
}

var _ = Describe("StrictDouble", func() {
	var double StrictDouble

	BeforeEach(func() {
		resetTestFail()
		double = NewStrictDoubleWithFailHandler(testFailHandler)
	})

	Context("when a stubbed method is called", func() {
		BeforeEach(func() {
			double.StubMethod("UltimateQuestion", []interface{}{}, []interface{}{42, nil})
		})

		It("returns the stubbed return values", func() {
			returnValues := double.Call("UltimateQuestion")

			Expect(returnValues).To(Equal([]interface{}{42, nil}))
			Expect(testFailHandlerInvoked).To(BeFalse())
		})
	})

	Context("when an unknown method is called", func() {
		It("fails the test", func() {
			returnValues := double.Call("UnstubbedMethod")

			Expect(returnValues).To(BeNil())
			Expect(testFailHandlerInvoked).To(BeTrue())
			Expect(testFailMessage).To(Equal("No stub for method 'UnstubbedMethod'"))
		})
	})
})
