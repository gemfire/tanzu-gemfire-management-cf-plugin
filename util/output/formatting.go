package output

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/util"
	jq "github.com/threatgrid/jqpipe-go"
)

// Fill ensures that a column is filled with desired filler characters to desired size
func Fill(columnSize int, value string, filler string) (response string) {
	if len(value) > columnSize-1 {
		response = " " + value[:columnSize-len(util.Ellipsis)-1] + util.Ellipsis
		return
	}
	numFillerChars := columnSize - len(value) - 1
	response = " " + value + strings.Repeat(filler, numFillerChars)
	return
}

// GetJSONFromURLResponse extracts JSON from a response
func GetJSONFromURLResponse(urlResponse string, jqFilter string) (jsonOutput string, err error) {
	if jqFilter == "" {
		return indent([]byte(urlResponse))
	}
	// otherwise use the filter string to generate the list for display.
	//".result[] | .runtimeInfo[] | {id:.id,status:.status}"
	jsonRawMessage, err := jq.Eval(urlResponse, jqFilter)
	jsonByte, err := json.Marshal(&jsonRawMessage)

	return indent(jsonByte)
}

func indent(rawJson []byte) (indented string, err error) {
	dst := &bytes.Buffer{}
	err = json.Indent(dst, rawJson, "", "  ")
	if err != nil {
		return string(rawJson), nil
	}
	return dst.String(), nil
}
