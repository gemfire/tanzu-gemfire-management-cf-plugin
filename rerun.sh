cf uninstall-plugin pcc
go build main.go helpers.go strings.go structs.go formatting.go argument_parser.go commandMapping.go
echo Y | cf install-plugin main

echo "    _______________  ______________________"
echo "   / _____/ ______/ / ____  / _____/ _____/"
echo "  / /    / /___    / /___/ / /    / /    "
echo " / /____/ ____/   / ______/ /____/ /____"
echo "/______/_/       /_/     /______/______/ cf pcc 1.0.0"
