package moka_test

import (
	. "github.com/gcapizzi/moka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StrictDouble", func() {
	var double StrictDouble

	BeforeEach(func() {
		double = NewStrictDouble()
	})

	Context("when a method is stubbed", func() {
		BeforeEach(func() {
			double.StubMethod("UltimateQuestion", []interface{}{}, 42)
		})

		It("returns the stubbed return values", func() {
			Expect(double.Call("UltimateQuestion").Get(0)).To(Equal(42))
		})
	})
})
