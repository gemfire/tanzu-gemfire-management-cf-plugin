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

	"github.com/itchyny/gojq"
)

type GOJQFilter struct{}

func (filter *GOJQFilter) Filter(jsonString string, expr string) ([]json.RawMessage, error) {
	query, err := gojq.Parse(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("json query failed: %s", err.Error()))
	}

	var inputBuffer bytes.Buffer
	var returnJson []json.RawMessage

	inputBuffer.WriteString(jsonString)
	dec := json.NewDecoder(&inputBuffer)
	var input json.RawMessage
	err = dec.Decode(&input)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("problem decoding input: %s", err.Error()))
	}
	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if fmt.Sprintf("%T", v) != "json.RawMessage" {
			return nil, errors.New(fmt.Sprintf("problem decoding output: %v", v))
		}
		returnJson = append(returnJson, v.(json.RawMessage))
	}
	return returnJson, nil
}
