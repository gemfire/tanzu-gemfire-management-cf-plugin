package requests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {

	Context("GetTargetAndClusterCommand", func() {

		Context("Target and Command are provided on the command line", func() {

			It("Does nothing", func() {
				Expect(nil).To(BeNil())
			})
		})

		Context("Target is provided as an environment variable command is on command line", func() {

			It("Does nothing", func() {
				Expect(nil).To(BeNil())
			})
		})

		Context("Target or command are missing from command line", func() {

			It("Does nothing", func() {
				Expect(nil).To(BeNil())
			})
		})
	})

	Context("GetEndPoints", func() {

		It("Does nothing", func() {
			Expect(nil).To(BeNil())
		})
	})

	Context("RequestToEndPoint", func() {

		It("Does nothing", func() {
			Expect(nil).To(BeNil())
		})
	})
})
