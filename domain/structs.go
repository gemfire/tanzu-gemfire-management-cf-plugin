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

package domain

import "code.cloudfoundry.org/cli/plugin"

var VersionType = plugin.VersionType{Major: 1, Minor: 0, Build: 6}

// CommandData is all the data involved in executing plugin commands
// This data gets manipulated throughout the program
type CommandData struct {
	Target             string
	ConnnectionData    ConnectionData
	UserCommand        UserCommand
	AvailableEndpoints map[string]RestEndPoint //key is command name
}

// ConnectionData describes items required to connect to a Geode cluster
type ConnectionData struct {
	Username       string
	Password       string
	Token          string
	UseToken       bool
	LocatorAddress string
}

// ServiceKeyUsers holds the username and password for users identified in a CF service key
type ServiceKeyUsers struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// ServiceKeyUrls holds URL information for endpoints to the PCC manageability service
type ServiceKeyUrls struct {
	Management string `json:"management"`
	Gfsh       string `json:"gfsh"`
}

// ServiceKey aggregates the information returned when requesting a service key from CF
type ServiceKey struct {
	Urls  ServiceKeyUrls    `json:"urls"`
	Users []ServiceKeyUsers `json:"users"`
}

// UserCommand holds command and parameter information entered by a user
type UserCommand struct {
	Command    string
	Parameters map[string]string
}

// RestEndPoint holds endpoint information
type RestEndPoint struct {
	HTTPMethod  string
	URL         string
	CommandName string
	JQFilter    string
	Consumes    []string
	Parameters  []RestAPIParam
}

// RestAPI is used to parse the swagger json response
// first key: url | second key: method (get/post) | value: RestAPIDetail
type RestAPI struct {
	Definitions map[string]DefinitionDetail         `json:"definitions"`
	Paths       map[string]map[string]RestAPIDetail `json:"paths"`
	Info        APIInfo                             `json:"info"`
}

type APIInfo struct {
	TokenEnabled string `json:"authTokenEnabled"`
}

// RestAPIDetail provides details about an endpoint
type RestAPIDetail struct {
	CommandName string         `json:"summary"`
	JQFilter    string         `json:"jqFilter"`
	Consumes    []string       `json:"consumes"`
	Parameters  []RestAPIParam `json:"parameters"`
}

// DefinitionDetail describes the details of the type definitions
type DefinitionDetail struct {
	Properties map[string]PropertyDetail `json:"properties"`
}

// PropertyDetail describes the details of the properties of type definitions
type PropertyDetail struct {
	Type   string            `json:"type"`
	Ref    string            `json:"$ref"`
	Enum   []string          `json:"enum"`
	Format string            `json:"format"`
	Items  map[string]string `json:"items"`
}

// RestAPIParam contains the information about possible parameters for a call
type RestAPIParam struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`
	// In describes how params are submitted: "query", "body" or "path" , or "formData"
	In             string            `json:"in"`
	Schema         map[string]string `json:"schema"`
	BodyDefinition map[string]interface{}
}
