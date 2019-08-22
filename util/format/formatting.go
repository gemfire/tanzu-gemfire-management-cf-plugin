package format

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/util"
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
func GetJSONFromURLResponse(urlResponse string) (jsonOutput string, err error) {
	var raw map[string]interface{}
	err = json.Unmarshal([]byte(urlResponse), &raw)
	if err != nil {
		return urlResponse, nil
	}
	out, err := json.Marshal(raw)
	if err != nil {
		return urlResponse, nil
	}
	var buf bytes.Buffer
	json.Indent(&buf, out, "", "  ")
	jsonOutput = string(buf.Bytes())
	return
}
