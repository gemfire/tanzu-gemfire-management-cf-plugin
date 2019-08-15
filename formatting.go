package pcc

import (
	"bytes"
	"encoding/json"
	"strings"
)

func Fill(columnSize int, value string, filler string) (response string) {
	if len(value) > columnSize-1 {
		response = " " + value[:columnSize-len([]rune(Ellipsis))-1] + Ellipsis
		return
	}
	numFillerChars := columnSize - len(value) - 1
	response = " " + value + strings.Repeat(filler, numFillerChars)
	return
}

func GetJsonFromUrlResponse(urlResponse string) (jsonOutput string, err error) {
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
