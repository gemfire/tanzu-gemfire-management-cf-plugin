package common

// Collection of common strings used by the application
const (
	NoServiceKeyMessage = "Please create a service key for %s.\n" +
		"To create a key enter:\n\n" +
		"cf create-service-key %s <your_key_name>\n\n" +
		"use --help or -h for help"
	GenericErrorMessage       = "Cannot retrieve credentials. Error: %s"
	InvalidServiceKeyResponse = "The cf service-key response is invalid."
	GeneralOptions            = "\t\t--user, -u <username> or set 'GEODE_USERNAME' environment variable to set the username\n" +
		"\t\t--password, -p <password> or set 'GEODE_PASSWORD' environment variable to set the password\n" +
		"\t\t--table, -t [<jqFilter>] to get tabular output"
)
