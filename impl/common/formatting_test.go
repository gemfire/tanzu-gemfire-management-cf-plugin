package common_test

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting", func() {

	Describe("Fill tests", func() {
		It("Replaces invalid filler with 'space' characters", func() {
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
		It("Returns the input as table format", func() {
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

	Context("DescribeEndpoint", func() {
		var (
			endPoint domain.RestEndPoint
		)

		It("Shows the expected command name in the output", func() {
			endPoint = domain.RestEndPoint{}
			endPoint.CommandName = "testcommand"

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("testcommand"))
		})

		It("Shows all expected parameters in the output", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramTwo := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramTwo.Name = "paramTwo"
			endPoint.Parameters = []domain.RestAPIParam{paramOne, paramTwo}

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("--paramOne"))
			Expect(result).To(ContainSubstring("--paramTwo"))
		})

		It("Shows if parameters are optional in the output", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramOne.Description = "first parameter"
			paramOne.Required = false
			endPoint.Parameters = []domain.RestAPIParam{paramOne}

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("[--paramOne <first parameter>]"))
		})

		It("Shows if parameters are required in the output", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramOne.Required = true
			endPoint.Parameters = []domain.RestAPIParam{paramOne}

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).NotTo(ContainSubstring("["))
			Expect(result).NotTo(ContainSubstring("]"))
			Expect(result).To(ContainSubstring("--paramOne"))
		})

		It("Shows GeneralOptions when showDetails flag set to true", func() {
			endPoint = domain.RestEndPoint{}

			result := common.DescribeEndpoint(endPoint, true)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring(common.GeneralOptions))
		})

		It("Hides GeneralOptions when showDetails flag set to false", func() {
			endPoint = domain.RestEndPoint{}

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(ContainSubstring(common.GeneralOptions))
		})

		It("Correctly display expected body format if 'body' parameter present and showDetails flag set to true", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramOne.Required = true
			paramOne.In = "body"
			bodyDefinition := make(map[string]interface{})

			bodyDefinition["propDetail1"] = "string-value"
			bodyDefinition["propDetail2"] = 42
			bodyDefinition["propDetail3"] = []int{21, 22}
			bodyDefinition["propDetail4"] = true
			bodyDefinition["propDetail5"] = map[string]interface{}{"1someString": "stringValue", "2someNumber": 23, "3someBool": false}

			paramOne.BodyDefinition = bodyDefinition
			endPoint.Parameters = []domain.RestAPIParam{paramOne}
			expectedOutput := `--paramOne format:
		{
		  "propDetail1": "string-value",
		  "propDetail2": 42,
		  "propDetail3": [
		    21,
		    22
		  ],
		  "propDetail4": true,
		  "propDetail5": {
		    "1someString": "stringValue",
		    "2someNumber": 23,
		    "3someBool": false
		  }
		}`

			result := common.DescribeEndpoint(endPoint, true)
			Expect(result).To(ContainSubstring(expectedOutput))
			Expect(result).To(ContainSubstring(common.GeneralOptions))
		})

		It("describe the rest end point without body param", func() {
			var endPoint domain.RestEndPoint
			endPoint.CommandName = "test"
			endPoint.Parameters = make([]domain.RestAPIParam, 2)

			var param1, param2 domain.RestAPIParam
			param1.In = "query"
			param1.Name = "id"
			param1.Description = "id"
			param1.Required = true

			param2.In = "query"
			param2.Name = "group"
			param2.Description = "group"
			param2.Required = false
			endPoint.Parameters[0] = param2
			endPoint.Parameters[1] = param1

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).To(Equal("test --id <id> [--group <group>]"))
		})

		It("describe the rest end point with body param", func() {
			var endPoint domain.RestEndPoint
			endPoint.CommandName = "test"
			endPoint.Parameters = make([]domain.RestAPIParam, 2)

			var param1, param2 domain.RestAPIParam
			param1.In = "body"
			param1.Name = "config"
			param1.Description = "config"
			param1.Required = true

			param2.In = "query"
			param2.Name = "group"
			param2.Description = "group"
			param2.Required = false
			endPoint.Parameters[0] = param2
			endPoint.Parameters[1] = param1

			result := common.DescribeEndpoint(endPoint, false)
			Expect(result).To(Equal("test --config <json or @json_file_path> [--group <group>]"))
		})
	})
})
