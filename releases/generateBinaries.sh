#!/usr/bin/env bash

if [ ! -d "./linux64" ]; then
  mkdir ./linux64
fi
cd linux64 || exit 1
env GOOS=linux GOARCH=amd64 go build -o ./pcc ../../cmd/main
md5 ./pcc > md5.txt
echo "Built linux executable"

if [ ! -d "../osx" ]; then
  mkdir ../osx
fi
cd ../osx || exit 1
env GOOS=darwin GOARCH=amd64 go build -o ./pcc ../../cmd/main
md5 ./pcc > md5.txt
echo "Built darwin executable"

if [ ! -d "../win64" ]; then
  mkdir ../win64
fi
cd ../win64 || exit 1
env GOOS=windows GOARCH=amd64 go build -o ./pcc.exe ../../cmd/main
md5 ./pcc.exe > md5.txt
echo "Built windows exectable"
