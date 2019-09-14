package common

import (
	"bytes"
	"encoding/json"
	"github.com/vito/go-interact/interact/terminal"
	"os"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	jq "github.com/threatgrid/jqpipe-go"
)

// Fill ensures that a column is filled with desired filler characters to desired size
func Fill(columnSize int, value string, filler string) (response string) {
	if len(filler) != 1 {
		// if invalid filler, use default filler
		filler = " "
	}
	// always leave one space before and after the value
	if len(value)+2 <= columnSize {
		numFillerChars := columnSize - (len(value) + 2)
		return filler + value + filler + strings.Repeat(filler, numFillerChars)
	}
	// when space is limited
	if columnSize < 5 {
		columnSize = 5
	}
	return filler + value[:columnSize-5] + "..." + filler
}

// FormatResponse extracts JSON from a response
func FormatResponse(urlResponse string, jqFilter string) (jsonOutput string, err error) {
	if jqFilter == "" {
		return indent([]byte(urlResponse))
	}
	// otherwise use the filter string to generate the list for display.
	filteredJson, err := filterWithJQ(urlResponse, jqFilter)
	if err != nil {
		return filteredJson, err
	}

	// convert the filtered json to tabular format
	return Tabular(filteredJson)
}

func filterWithJQ(jsonString string, jqFilter string) (string, error) {
	jsonRawMessage, err := jq.Eval(jsonString, jqFilter)
	if err != nil {
		return "unable to filter the response with: " + jqFilter, err
	}
	jsonByte, err := json.Marshal(&jsonRawMessage)
	if err != nil {
		return "unable to filter the response with: " + jqFilter, err
	}
	return string(jsonByte), nil
}

func Tabular(jsonString string) (string, error) {
	// parse the jsonString into array of maps
	var result []map[string]string
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return "unable to parse: " + jsonString, err
	}

	// get all the column names
	var columnNames []string
	for _, m := range result {
		for k, _ := range m {
			if !contains(columnNames, k) {
				columnNames = append(columnNames, k)
			}
		}
	}

	maxLengths := getMaxLength(result)
	// this already includes the length needed by the spacers
	totalLengthNeeded := getTotalMaxLength(maxLengths)

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	// if unable to get the width, assume the longest needed
	if terminalWidth <= 0 {
		terminalWidth = totalLengthNeeded
	}

	deficit := totalLengthNeeded - terminalWidth
	// shrink the longer columns to fit the screen
	if deficit > 0 {
		averageColumnWidth := int(float64(terminalWidth) / float64(len(columnNames)))
		// find the number of columns that are longer than the average, and the deficit
		// should be shared equally by these culprit columns
		longColumnCount := 0
		for _, length := range maxLengths {
			if length > averageColumnWidth {
				longColumnCount++
			}
		}
		shrinkage := int(float64(deficit) / float64(longColumnCount))
		for column, length := range maxLengths {
			if length > averageColumnWidth {
				maxLengths[column] = length - shrinkage
			}
		}
		totalLengthNeeded = getTotalMaxLength(maxLengths)
	}

	// print out the column names
	var response strings.Builder
	response.Grow(20)
	for index, column := range columnNames {
		response.WriteString(Fill(maxLengths[column], column, " "))
		if index < len(columnNames)-1 {
			response.WriteString("|")
		} else {
			response.WriteString("\n")
		}
	}

	// print out the divider
	response.WriteString(Fill(totalLengthNeeded, "", "-") + "\n")

	// print out the values
	for _, m := range result {
		for index, column := range columnNames {
			response.WriteString(Fill(maxLengths[column], m[column], " "))
			if index < len(columnNames)-1 {
				response.WriteString("|")
			} else {
				response.WriteString("\n")
			}
		}
	}
	return response.String(), nil
}

func getTotalMaxLength(maxLengths map[string]int) int {
	total := 0
	for _, length := range maxLengths {
		// need the spacer between each column
		total += length + 1
	}
	return total - 1
}

// get the max length of each column, and the total max length
func getMaxLength(result []map[string]string) map[string]int {
	maxLengths := make(map[string]int)
	for _, m := range result {
		for k, v := range m {
			maxLength := maxLengths[k]
			// always leave space before and after
			if (len(v) + 2) > maxLength {
				maxLengths[k] = len(v) + 2
			}
		}
	}
	return maxLengths
}

// Describe an end point with command name and required/optional parameters
func Describe(endPoint domain.RestEndPoint) string {
	var buffer bytes.Buffer
	buffer.WriteString(endPoint.CommandName + " ")
	// show the required options first
	for _, param := range endPoint.Parameters {
		if param.Required {
			buffer.WriteString(getOption(param))
		}
	}

	for _, param := range endPoint.Parameters {
		if !param.Required {
			buffer.WriteString("[" + strings.Trim(getOption(param), " ") + "] ")
		}
	}
	return buffer.String()
}

func getOption(param domain.RestAPIParam) string {
	if param.In == "body" {
		return "--body  "
	}
	return "--" + param.Name + " "
}

func indent(rawJSON []byte) (indented string, err error) {
	dst := &bytes.Buffer{}
	err = json.Indent(dst, rawJSON, "", "  ")
	if err != nil {
		return string(rawJSON), nil
	}
	return dst.String(), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
