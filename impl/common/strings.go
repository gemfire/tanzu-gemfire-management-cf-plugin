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

// Collection of common strings used by the application
const (
	NoServiceKeyMessage = "Please create a service key for %s.\n" +
		"To create a key enter:\n\n" +
		"cf create-service-key %s <your_key_name>\n\n" +
		"use --help or -h for help"
	GenericErrorMessage       = "Cannot retrieve credentials. Error: %s"
	InvalidServiceKeyResponse = "The cf service-key response is invalid."
	GeneralOptions            = "\t\t--user, -u <username> or set 'GEODE_USERNAME' environment variable to set the username\n" +
		"\t\t--password, -p <password> or set 'GEODE_PASSWORD' environment variable to set the password\n" +
		"\t\t--table, -t [<jqFilter>] to get tabular output"
)
