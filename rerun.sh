#set -e
cf uninstall-plugin PCC_InDev
go build pcc.go
echo Y | cf install-plugin pcc
