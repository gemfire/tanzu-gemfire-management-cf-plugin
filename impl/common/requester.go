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
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"io"
	"io/ioutil"
	"net/http"
)

// Requester is the receiver for the RequestHelper implementation
type Requester struct{}

// Exchange implements the RequestHelper interface
func (requester *Requester) Exchange(url string, method string, bodyReader io.Reader, connectionData *domain.ConnectionData) (urlResponse string, err error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return "", err
	}

	if connectionData != nil {
		if connectionData.UseToken {
			var bearer = "Bearer " + connectionData.Token
			req.Header.Add("Authorization", bearer)
		} else {
			req.SetBasicAuth(connectionData.Username, connectionData.Password)
		}
	}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	return getURLOutput(resp)
}

func getURLOutput(resp *http.Response) (urlResponse string, err error) {
	respInASCII, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	urlResponse = fmt.Sprintf("%s", respInASCII)

	return
}
