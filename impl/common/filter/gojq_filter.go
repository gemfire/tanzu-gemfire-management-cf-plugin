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
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

// GOJQFilter is the placeholder struct for the Filter interface implementation
type GOJQFilter struct{}

// Filter is the implementation of the Filter interface
func (filter *GOJQFilter) Filter(jsonString string, expr string) ([]json.RawMessage, error) {
	query, err := gojq.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("json query failed: %s", err.Error())
	}

	var returnJSON []json.RawMessage

	var interfaceInput interface{}
	err = json.Unmarshal([]byte(jsonString), &interfaceInput)
	if err != nil {
		return nil, fmt.Errorf("problem decoding input: %s", err.Error())
	}
	iter := query.Run(interfaceInput)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		switch x := v.(type) {
		case error:
			return nil, fmt.Errorf("problem decoding output: %s", x.Error())
		case [2]interface{}:
			if s, ok := x[0].(string); ok {
				if s == "HALT:" {
					return returnJSON, nil
				}
				if s == "STDERR:" {
					return nil, fmt.Errorf("problem decoding output: %s", x[1])
				}
			}
		case json.RawMessage:
			returnJSON = append(returnJSON, x)
		default:
			rawJSON, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("problem decoding output: %v", v)
			}
			returnJSON = append(returnJSON, rawJSON)
		}
	}
	return returnJSON, nil
}
