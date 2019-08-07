go build main.go helpers.go strings.go structs.go formatting.go argument_parser.go commandMapping.go
echo Y | cf install-plugin main
cf pcc --help
