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
	"bytes"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/domain"
	"github.com/gemfire/tanzu-gemfire-management-cf-plugin/impl/common"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// BuildRequest implements common.RequestBuilder func type
func BuildRequest(restEndPoint domain.RestEndPoint, commandData *domain.CommandData) (request *http.Request, err error) {
	connectionData := commandData.ConnnectionData
	requestURL := connectionData.LocatorAddress + "/management" + restEndPoint.URL
	var query string
	var multiPartForm bool

	// for content body reader
	var bodyReader io.Reader

	// for multipart form body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// form data are handled differently
	if common.Contains(restEndPoint.Consumes, "multipart/form-data") {
		multiPartForm = true
	}

	for _, param := range restEndPoint.Parameters {
		value, ok := commandData.UserCommand.Parameters["--"+param.Name]
		if ok {
			switch param.In {
			case "path":
				requestURL = strings.ReplaceAll(requestURL, "{"+param.Name+"}", url.PathEscape(value))
			case "query":
				if multiPartForm {
					err = writer.WriteField(param.Name, value)
					if err != nil {
						return nil, err
					}
				}
				if len(query) == 0 {
					query = "?" + param.Name + "=" + url.PathEscape(value)
				} else {
					query = query + "&" + param.Name + "=" + url.PathEscape(value)
				}
			case "body":
				bodyReader, err = getBodyReader(value)
				if err != nil {
					return nil, err
				}
			case "formData":
				if param.Type == "string" {
					err = writer.WriteField(param.Name, value)
				} else if param.Type == "file" {
					file, err := os.Open(value)
					if err != nil {
						return nil, err
					}
					defer file.Close()
					part, err := writer.CreateFormFile(param.Name, filepath.Base(value))
					if err != nil {
						return nil, err
					}
					_, err = io.Copy(part, file)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	err = writer.Close()
	requestURL = requestURL + query
	httpAction := strings.ToUpper(restEndPoint.HTTPMethod)
	if multiPartForm {
		request, err = http.NewRequest(httpAction, requestURL, body)
	} else {
		request, err = http.NewRequest(httpAction, requestURL, bodyReader)
	}

	if connectionData.UseToken {
		var bearer = "Bearer " + connectionData.Token
		request.Header.Add("Authorization", bearer)
	} else {
		request.SetBasicAuth(connectionData.Username, connectionData.Password)
	}

	if multiPartForm {
		request.Header.Set("content-type", writer.FormDataContentType())
	} else {
		request.Header.Set("content-type", "application/json")
	}
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
