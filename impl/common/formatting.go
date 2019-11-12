/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for additional information regarding
 * copyright ownership. The ASF licenses this file to You under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License. You may obtain a
 * copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/cf/errors"
	"github.com/vito/go-interact/interact/terminal"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . JsonFilter

// JsonFilter interface provides a way to provide different json filter implementations
// or to replace the filter with a fake for testing
type JsonFilter interface {
	Filter(jsonString string, expr string) ([]json.RawMessage, error)
}

type Formatter struct {
	JsonFilter JsonFilter
}

func NewFormatter(jsonFilter JsonFilter) *Formatter {
	return &Formatter{JsonFilter: jsonFilter}
}

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
func (formatter *Formatter) FormatResponse(urlResponse string, jqFilter string, userFilter bool) (jsonOutput string, err error) {
	if jqFilter == "" {
		return indent([]byte(urlResponse))
	}
	// otherwise use the filter string to generate the list to display in table format
	filteredJSON, err := formatter.filterWithJQ(urlResponse, jqFilter)

	// if using the default jqFilter does not yield any data, then display the unfiltered result
	if filteredJSON == "[]" && !userFilter && jqFilter != "." {
		jqFilter = "."
		filteredJSON, err = formatter.filterWithJQ(urlResponse, jqFilter)
	}

	if err != nil {
		return filteredJSON, err
	}

	// convert the filtered json to tabular format
	table, err := Tabular(filteredJSON)
	if err != nil {
		return "", err
	}

	return table + "\n" + "JQFilter: " + jqFilter + "\n", nil
}

func (formatter *Formatter) filterWithJQ(jsonString string, expr string) (string, error) {
	jsonRawMessage, err := formatter.JsonFilter.Filter(jsonString, expr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to filter the response with: %s, %s", expr, err))
	}
	jsonByte, err := json.Marshal(&jsonRawMessage)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to filter the response with: %s, %s", expr, err))
	}
	return string(jsonByte), nil
}

// Tabular provides tabular string output for a given JSON string
func Tabular(jsonString string) (string, error) {
	// parse the jsonString into array of maps
	var result []map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to parse: %s, %s", jsonString, err))
	}

	if len(result) == 0 {
		return "", nil
	}

	// get all the column names
	var columnNames []string
	for _, m := range result {
		for k := range m {
			if !contains(columnNames, k) {
				columnNames = append(columnNames, k)
			}
		}
	}

	maxLengths := getMaxLength(&result)
	// this already includes the length needed by the spacers
	totalLengthNeeded := getTotalMaxLength(maxLengths)

	terminalWidth, _, _ := terminal.GetSize(int(os.Stdin.Fd()))
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
			response.WriteString(Fill(maxLengths[column], getString(m[column]), " "))
			if index < len(columnNames)-1 {
				response.WriteString("|")
			} else {
				response.WriteString("\n")
			}
		}
	}
	return response.String(), nil
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
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
func getMaxLength(result *[]map[string]interface{}) map[string]int {
	maxLengths := make(map[string]int)
	for _, m := range *result {
		for k, v := range m {
			// use the length of the column name + 2 as starting length
			if maxLengths[k] == 0 {
				maxLengths[k] = len(k) + 2
			}
			maxLength := maxLengths[k]
			// always leave space before and after
			strSize := len(getString(v)) + 2
			if strSize > maxLength {
				maxLengths[k] = strSize
			}
		}
	}
	return maxLengths
}

// DescribeEndpoint an end point with command name and required/optional parameters
func DescribeEndpoint(endPoint domain.RestEndPoint, showDetails bool) string {
	var buffer bytes.Buffer
	buffer.WriteString(endPoint.CommandName + " ")
	// show the required options first
	for _, param := range endPoint.Parameters {
		if param.Required {
			writeParam(&buffer, param)
			buffer.WriteString(" ")
		}
	}

	for _, param := range endPoint.Parameters {
		if !param.Required {
			buffer.WriteString("[")
			writeParam(&buffer, param)
			buffer.WriteString("] ")
		}

	}

	if showDetails {
		for _, param := range endPoint.Parameters {
			if len(param.BodyDefinition) > 0 {
				generateSampleBody(param, &buffer)
			}
		}
		buffer.WriteString("\n" + GeneralOptions)
	}
	return strings.Trim(buffer.String(), " ")
}

func writeParam(buffer *bytes.Buffer, param domain.RestAPIParam) {
	buffer.WriteString("--" + param.Name + " ")
	if param.In == "body" {
		buffer.WriteString("<json or @json_file_path>")
	} else {
		buffer.WriteString("<" + param.Description + ">")
	}
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

func generateSampleBody(param domain.RestAPIParam, buffer *bytes.Buffer) {
	buffer.WriteString("\n\t\t--" + param.Name + " format:\n\t\t")

	jsonBytes, err := json.MarshalIndent(param.BodyDefinition, "\t\t", "  ")
	if err != nil {
		return
	}
	buffer.Write(jsonBytes)
}
