#Getting Started:

###Prerequisites:
* You are logged into CF
    * `cf login...`
       * Credentials can be found in Ops Manager &rarr; PAS tile &rarr; UAA &rarr; Admin Credentials
* You have a PCC instance running.
    * `cf create-service p-cloudcache dev-plan myPCCInstance`
* Your PCC instance has a service key.
    * `cf create-service-key myPCCInstance myKey`

###Running the Plugin:
1. Clone the repository &rarr; `git clone git@github.com:gemfire/cloudcache-management-cf-plugin.git`
2. Build `gf.go` &rarr; `go build gf.go`
3. Install the plugin &rarr; `cf install-plugin gf`
4. Run commands (for help &rarr; `cf gf --help`)
    

