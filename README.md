# Getting Started:

## Prerequisites:
* You have installed CF CLI
    * https://docs.cloudfoundry.org/cf-cli/install-go-cli.html
* You are logged into CF
    *  `cf login --skip-ssl-validation -a https://api.sys.ENVNAME.cf-app.com -u admin -p PASSWORD`
    * `ENVNAME` corresponds to your CF environment
    * `PASSWORD` can be found in the CF Ops Manager &rarr; PAS tile &rarr; Credentials Tab &rarr; UAA Admin Credentials

* You have a PCC instance running.
    * `cf create-service p-cloudcache dev-plan myPCCInstance`
* Your PCC instance has a service key.
    * `cf create-service-key myPCCInstance myKey`

## Running the Plugin:
1. Clone the Repository 
    - `git clone git@github.com:gemfire/cloudcache-management-cf-plugin.git`
2. Start the Plugin 
    - Run `./start.sh` from the `cloudcache-management-cf-plugin` directory
3. Run Commands 

4. For Help
    - `cf pcc --help`
    
## command format
 1. "cf pcc cluster1 get members"
     - there are 3 parts in this command
        1. "cf pcc" - commmand prefix
        2. "cluster1" - service name
        3. "get members" - actual command
 1. "cf pcc cluster1 get member -id=member1"
     - there are 4 parts in this command
        1. "cf pcc" - commmand prefix
        2. "cluster1" - service name
        3. "get member" - actual command
        1. "-id=member1" - param for this command,  also you can input json as param in ""
  1. "cf pcc cluster1 get member -id=member1 -j"
     - there are 4 parts in this command
        1. "cf pcc" - commmand prefix
        2. "cluster1" - service name
        3. "get member" - actual command
        1. "-id=member1" - param for this command,  also you can input json as param in ""
        1. "-j" - output option. 
           - "-j" means that all the output will be shown in Json format
           - "-t" means that all the output will be shown in table format. and also, you can specify which columns will be shown, by similar like this: "-t "name, address.zipcode, mobileNo""
