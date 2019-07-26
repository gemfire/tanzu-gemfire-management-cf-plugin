#set -e
cf uninstall-plugin PCC_InDev
go build main.go helpers.go strings.go structs.go formatting.go argument_parser.go
echo Y | cf install-plugin main
