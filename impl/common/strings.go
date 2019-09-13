package common

// Collection of common strings used by the application
const (
	NoServiceKeyMessage string = "Please create a service key for %s.\n" +
		"To create a key enter:\n\n" +
		"cf create-service-key %s <your_key_name>\n\n" +
		"use --help or -h for help"
	GenericErrorMessage       string = "Cannot retrieve credentials. Error: %s"
	InvalidServiceKeyResponse string = "The cf service-key response is invalid."

	NoJsonFileProvidedMessage string = "A JSON configuration file is required for all create/post commands.\n" +
		"Please re-enter your command appended with --body <your_json_configuration>\n\n" +
		"<your_json_configuration> is in form of @<json_file_path> OR single quoted JSON input."
	Ellipsis       string = "…"
	GeneralOptions        = "\t\t--user, -u or set 'GEODE_USERNAME' environment variable to set the username\n" +
		"\t\t--password, -p or set 'GEODE_PASSWORD' environment variable to set the password\n" +
		"\t\t--table, -t [jqFilter] to get tabular output"
)
