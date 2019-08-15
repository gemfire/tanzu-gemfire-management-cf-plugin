cd linux64
env GOOS=linux GOARCH=amd64 go build -o ./pcc ../../main.go ../../helpers.go ../../strings.go ../../structs.go ../../formatting.go ../../argument_parser.go ../../commandMapping.go
md5 ./pcc > md5.txt
cd ../osx
env GOOS=darwin GOARCH=amd64 go build -o ./pcc ../../main.go ../../helpers.go ../../strings.go ../../structs.go ../../formatting.go ../../argument_parser.go ../../commandMapping.go
md5 ./pcc > md5.txt
cd ../win64
env GOOS=windows GOARCH=amd64 go build -o ./pcc ../../main.go ../../helpers.go ../../strings.go ../../structs.go ../../formatting.go ../../argument_parser.go ../../commandMapping.go
md5 ./pcc > md5.txt
