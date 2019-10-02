# Getting Started:
This code is written primarily to function as a CF CLI plugin but it can also be used to execute against
any accessible Geode cluster in standalone mode. The minor differences of operation between plugin and
standalone mode will be explained below.

This code works with the V2 management API introduced in Geode 1.10. The API and this client are both
considered to be experimental code that may be subject to change.

Apart from needing the management API to be enabled the client is not bound to a particular version of
the API beyond needing a 1.10+ server. It is written to be able to adapt to the progress that is made
on the API code in a dynamic fashion.

## Prerequisites

#### In plugin mode:
* You have installed [CF CLI](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html) and you have
installed a compatible [PCC tile](https://network.pivotal.io/products/p-cloudcache/) in this environment.
(9.10+?)
* You are logged into CF
    *  `cf login --skip-ssl-validation -a https://api.sys.ENVNAME.cf-app.com -u admin -p PASSWORD`
    * `ENVNAME` corresponds to your CF environment
    * `PASSWORD` can be found in the CF Ops Manager &rarr; PAS tile &rarr; Credentials Tab &rarr; UAA
    Admin Credentials
* You have a PCC instance running
    * `cf create-service p-cloudcache dev-plan myPCCInstance`
* Your PCC instance has a service key
    * `cf create-service-key myPCCInstance myKey`

#### In standalone mode
* You have a Geode cluster running and have the co-ordinates (URI, username and password) of a `Locator`
in this cluster. This can be a development version running on `localhost`

## Running the code

#### Common instructions
 1. Have a [Go](https://golang.org/) SDK on your machine when compiling the code from the repository
 1. Clone the Repository
    - `git clone git@github.com:gemfire/cloudcache-management-cf-plugin.git`
 1. Some parameters can be replaced with environment variables to avoid having to type them in repeatedly.
 Please see the general help for details in each mode

#### As a plugin:
 1. Run the start script
    -  `./install.sh` from the `cloudcache-management-cf-plugin` directory or `./reinstall.sh` when
    replacing an existing version of the plugin
 1. For Help
    - `cf pcc --help` provides general help
    - `cf pcc <target> commands` to get a list of commands available to you. `<target>` is the PCC service instance name you are using
    - `cf pcc <target> <command> -help` to get `<command>` specific help including the format of `JSON` payload that some commands require

#### As a standalone client
 1. Run the start script
    -  `./build.sh` from the `cloudcache-management-cf-plugin` directory
 1. For Help
    - `./pcc --help` provides general help
    - `./pcc <target> commands` to get a list of commands available to you. `<target>` is the address of the `locator` you are using
    - `./pcc <target> <command> -help` to get `<command>` specific help including the format of `JSON` payload that some command require
