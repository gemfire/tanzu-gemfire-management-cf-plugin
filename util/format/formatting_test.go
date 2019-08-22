package format_test

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/util/format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting", func() {

	Context("Safekeeping tests", func() {
		It("Validate table filling", func() {
			columnSize := 20
			value := "some string"
			filler := "-"
			response := format.Fill(columnSize, value, filler)
			expectedResponse := " some string--------"
			Expect(response).To(Equal(expectedResponse))
		})
	})
})
