cf uninstall-plugin GF_InDev
go build gf.go
echo Y | cf install-plugin gf
cf gf list-regions jjack -g=group1
