#set -e
cf uninstall-plugin PCC_InDev
go build main.go helpers.go strings.go structs.go
echo Y | cf install-plugin main
