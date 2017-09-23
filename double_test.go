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
			double.StubMethod("UltimateQuestion", []interface{}{}, []interface{}{42, nil})
		})

		It("returns the stubbed return values", func() {
			Expect(double.Call("UltimateQuestion")).To(Equal([]interface{}{42, nil}))
		})
	})
})
