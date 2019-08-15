cf uninstall-plugin pcc
go build main.go helpers.go strings.go structs.go formatting.go argument_parser.go commandMapping.go
echo Y | cf install-plugin pcc

echo "    _______  ________    ________  _______  _______"
echo "   / _____/ / ______/   / ____  / / _____/ / _____/"
echo "  / /      / /___      / /___/ / / /      / /    "
echo " / /____  / ____/     / ______/ / /____  / /____"
echo "/______/ /_/         /_/       /______/ /______/ 1.0.0"
