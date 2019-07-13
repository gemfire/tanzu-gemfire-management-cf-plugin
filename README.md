# Getting Started:

## Prerequisites:
* You have installed CF CLI
    * https://docs.cloudfoundry.org/cf-cli/install-go-cli.html
* You are logged into CF
    * `cf login…`
       * Credentials can be found in CF Ops Manager &rarr; PAS tile &rarr; UAA &rarr; Admin Credentials
* You have a PCC instance running.
    * `cf create-service p-cloudcache dev-plan myPCCInstance`
* Your PCC instance has a service key.
    * `cf create-service-key myPCCInstance myKey`

## Running the Plugin:
1. Clone the Repository 
    - `git clone git@github.com:gemfire/cloudcache-management-cf-plugin.git`
2. Build the Plugin 
    - Enter `go build` from the `cloudcache-management-cf-plugin` directory
3. Install the Plugin 
    - `cf install-plugin pcc`
4. Run Commands 

5. For Help
    - `cf pcc --help`
    

