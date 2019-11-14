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

package filter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	jq "github.com/jmelchio/gojq/cli"
)

type GOJQFilter struct{}

func (filter *GOJQFilter) Filter(jsonString string, expr string) ([]json.RawMessage, error) {
	var inputBuffer bytes.Buffer
	var outputBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	inputBuffer.WriteString(jsonString)

	args := []string{expr, "-M"}
	returnCode := jq.RunLib(&inputBuffer, &outputBuffer, &errBuffer, args)

	if returnCode != 0 {
		return nil, errors.New(fmt.Sprintf("json query failed: %s, %d", errBuffer.String(), returnCode))
	}
	if outputBuffer.Len() == 0 {
		return []json.RawMessage{}, nil
	}
	dec := json.NewDecoder(&outputBuffer)
	var output json.RawMessage
	err := dec.Decode(&output)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("problem decoding output: %s", err.Error()))
	}
	return []json.RawMessage{output}, nil
}
