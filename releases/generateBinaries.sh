cd linux64
env GOOS=linux GOARCH=amd64 go build -o ./pcc ../../cmd/main
md5 ./pcc > md5.txt
cd ../osx
env GOOS=darwin GOARCH=amd64 go build -o ./pcc ../../cmd/main
md5 ./pcc > md5.txt
cd ../win64
env GOOS=windows GOARCH=amd64 go build -o ./pcc.exe ../../cmd/main
md5 ./pcc.exe > md5.txt
