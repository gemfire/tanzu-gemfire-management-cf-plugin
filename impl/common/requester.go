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
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Exchange implements the impl.RequestHelper function type
var Exchange = func(request *http.Request) (urlResponse string, statusCode int, err error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(request)
	if err != nil {
		return "", 0, err
	}

	urlResponse, err = getURLOutput(resp)
	statusCode = resp.StatusCode
	return
}

func getURLOutput(resp *http.Response) (urlResponse string, err error) {
	respInASCII, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	urlResponse = fmt.Sprintf("%s", respInASCII)

	return
}
