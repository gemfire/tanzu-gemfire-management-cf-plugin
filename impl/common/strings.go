package common

// Collection of common strings used by the application
const (
	IncorrectUserInputMessage string = "The format of your request is incorrect.\nuse --help or -h for help"
	NoServiceKeyMessage       string = `Please create a service key for %s.
	To create a key enter:

		cf create-service-key %s <your_key_name>

	use --help or -h for help
	`
	GenericErrorMessage       string = `Cannot retrieve credentials. Error: %s`
	InvalidServiceKeyResponse string = `The cf service-key response is invalid.`

	NoRegionGivenMessage string = `You need to provide a region (-r flag) to execute your command.

	To see your available regions:

		cf pcc %s list regions

	use --help or -h for help
	`

	NoIDGivenMessage string = `An identifier is required for all get commands.

	Please re-enter your command appended with -id=<your_object_of_interest>
	`
	NoJsonFileProvidedMessage string = `A JSON configuration file is required for all create/post commands.

	Please re-enter your command appended with -d=<your_json_configuration_file>

	<your_json_configuration_file> is in form of @<json_file_path> OR single quoted JSON input.`

	NoEndpointFoundMessage string = `No endpoint was found for your request.

	use --help or -h for help
	`

	Ellipsis string = "…"
)
