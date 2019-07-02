#Getting Started:

###Prerequisites:
* You have installed CF CLI
    * https://docs.cloudfoundry.org/cf-cli/install-go-cli.html
* You are logged into CF
    * `cf login…`
       * Credentials can be found in CF Ops Manager &rarr; PAS tile &rarr; UAA &rarr; Admin Credentials
* You have a PCC instance running.
    * `cf create-service p-cloudcache dev-plan myPCCInstance`
* Your PCC instance has a service key.
    * `cf create-service-key myPCCInstance myKey`

###Running the Plugin:
1. Clone the repository 
  - `git clone git@github.com:gemfire/cloudcache-management-cf-plugin.git`
2. Build `gf.go` 
  - `go build gf.go`
3. Install the plugin 
  - `cf install-plugin gf`
4. Run commands 

5. help
  - `cf gf --help`
    

