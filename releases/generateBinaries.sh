#!/usr/bin/env bash

#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if [ ! -d "./linux64" ]; then
  mkdir ./linux64
fi
cd linux64 || exit 1
env GOOS=linux GOARCH=amd64 go build -o ./pcc ../../cmd/main
md5 ./pcc >md5.txt
echo "Built linux executable"

if [ ! -d "../osx" ]; then
  mkdir ../osx
fi
cd ../osx || exit 1
env GOOS=darwin GOARCH=amd64 go build -o ./pcc ../../cmd/main
md5 ./pcc >md5.txt
echo "Built darwin executable"

if [ ! -d "../win64" ]; then
  mkdir ../win64
fi
cd ../win64 || exit 1
env GOOS=windows GOARCH=amd64 go build -o ./pcc.exe ../../cmd/main
md5 ./pcc.exe >md5.txt
echo "Built windows executable"
