package format_test

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/util"
	"github.com/gemfire/cloudcache-management-cf-plugin/util/format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting", func() {

	Context("Fill tests", func() {

		Context("Value is shorter than column size", func() {
			It("Fills the table with filler characters", func() {
				columnSize := 20
				value := "some string"
				filler := "-"
				response := format.Fill(columnSize, value, filler)
				expectedResponse := " some string--------"
				Expect(response).To(Equal(expectedResponse))
				Expect(len(response)).To(Equal(columnSize))
			})
		})

		Context("Value is longer than column size", func() {
			It("Truncates the value and adds Ellipsis at the end of the value", func() {
				columnSize := 20
				value := "some string that is longer than 20 characters"
				filler := "-"
				response := format.Fill(columnSize, value, filler)
				expectedResponse := " some string that" + util.Ellipsis
				Expect(response).To(Equal(expectedResponse))
				Expect(len(response)).To(Equal(columnSize))
			})
		})
	})

	Context("GetJSONFromURLResponse tests", func() {

		Context("Input string is valid JSON", func() {
			It("Returns the input as an indented string", func() {
				inputString := `{"name": "value"}`
				expectedString := `{
  "name": "value"
}`
				output, err := format.GetJSONFromURLResponse(inputString, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(Equal(expectedString))
			})
		})

		Context("Input string is not valid JSON", func() {
			It("Returns the input 'as-is'", func() {
				inputString := "foobar"
				output, err := format.GetJSONFromURLResponse(inputString, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(Equal(inputString))
			})
		})
	})
})
