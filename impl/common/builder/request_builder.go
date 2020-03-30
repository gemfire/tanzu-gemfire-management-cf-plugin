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

package builder

import (
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"io"
	"net/url"
	"os"
	"strings"
)

// BuildRequest implments common.RequestBuilder func type
func BuildRequest(restEndPoint domain.RestEndPoint, commandData *domain.CommandData) (requestURL string, bodyReader io.Reader, err error) {
	requestURL = commandData.ConnnectionData.LocatorAddress + "/management" + restEndPoint.URL
	var query string
	for _, param := range restEndPoint.Parameters {
		value, ok := commandData.UserCommand.Parameters["--"+param.Name]
		if ok {
			switch param.In {
			case "path":
				requestURL = strings.ReplaceAll(requestURL, "{"+param.Name+"}", url.PathEscape(value))
			case "query":
				if len(query) == 0 {
					query = "?" + param.Name + "=" + url.PathEscape(value)
				} else {
					query = query + "&" + param.Name + "=" + url.PathEscape(value)
				}
			case "body":
				bodyReader, err = getBodyReader(value)
			}
		}
	}

	requestURL = requestURL + query
	return
}

func getBodyReader(jsonFile string) (bodyReader io.Reader, err error) {
	if jsonFile[0] == '@' && len(jsonFile) > 1 {
		bodyReader, err = os.Open(jsonFile[1:])
		if err != nil {
			return
		}
	} else {
		bodyReader = strings.NewReader(jsonFile)
	}
	return
}
