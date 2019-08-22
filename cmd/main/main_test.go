package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context("There are no tests", func() {
		It("Does nothing", func() {
			Expect("nothing").To(Equal("nothing"))
		})
	})
})
