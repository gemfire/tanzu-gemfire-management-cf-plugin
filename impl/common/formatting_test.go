package common_test

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting", func() {

	Context("Fill tests", func() {
		It("invalid filler", func() {
			columnSize := 5
			value := "id"
			filler := "test"
			response := common.Fill(columnSize, value, filler)
			expectedResponse := " id  "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})
		It("small column size", func() {
			columnSize := 4
			value := "id"
			filler := " "
			response := common.Fill(columnSize, value, filler)
			expectedResponse := " id "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})
		It("small column size", func() {
			columnSize := 3
			value := "i"
			filler := " "
			response := common.Fill(columnSize, value, filler)
			expectedResponse := " i "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})
		It("small column size", func() {
			columnSize := 2
			value := "idle"
			filler := " "
			response := common.Fill(columnSize, value, filler)
			expectedResponse := " ... "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(5))
		})
		It("Fills the table with filler characters", func() {
			columnSize := 20
			value := "some string"
			filler := "-"
			response := common.Fill(columnSize, value, filler)
			expectedResponse := "-some string--------"
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})

		It("Truncates the value and adds Ellipsis at the end of the value", func() {
			columnSize := 20
			value := "some strings that is longer than 20 characters"
			filler := "-"
			response := common.Fill(columnSize, value, filler)
			expectedResponse := "-some strings th...-"
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})
	})

	Context("FormatResponse tests", func() {
		It("Returns the input as an indented string", func() {
			inputString := `{"name": "value"}`
			expectedString := "{\n  \"name\": \"value\"\n}"
			output, err := common.FormatResponse(inputString, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(expectedString))
		})
		It("Returns the input 'as-is'", func() {
			inputString := "foobar"
			output, err := common.FormatResponse(inputString, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(inputString))
		})
	})

	Context("tabular output", func() {
		It("Returns the input as an indented string", func() {
			json := `[{
				"id": "server",
				"status": "online"},
			  {"id": "locator",
				"status": "online"}]`
			output, _ := common.Tabular(json)
			expected := " id      | status \n" +
				"------------------\n" +
				" server  | online \n" +
				" locator | online \n"
			Expect(output).To(Equal(expected))
		})
		It("different attributes", func() {
			json := `[{"id": "server"},{"status": "online"}]`
			output, _ := common.Tabular(json)
			expected := " id     | status \n" +
				"-----------------\n" +
				" server |        \n" +
				"        | online \n"
			Expect(output).To(Equal(expected))
		})
	})
})
