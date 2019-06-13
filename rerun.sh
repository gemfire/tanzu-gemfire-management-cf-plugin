cf uninstall-plugin CLI_InDev
go build cli.go
echo Y | cf install-plugin cli
#cf pcc-inspector -e -u jaccc -p pwww cert43975493875cert list-regions jjack
cf cli list-members jjack