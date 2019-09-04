package domain

var ()

// CommandData is all the data involved in executing plugin commands
// This data gets manipulated throughout the program
type CommandData struct {
	Username           string
	Password           string
	LocatorAddress     string
	Target             string
	ServiceKey         string
	Region             string
	JSONFile           string
	Group              string
	ID                 string
	HasGroup           bool
	IsJSONOutput       bool
	ExplicitTarget     bool
	ConnnectionData    ConnectionData
	UserCommand        UserCommand
	FirstResponse      SwaggerInfo
	AvailableEndpoints []IndividualEndpoint
	Endpoint           IndividualEndpoint
}

// ConnectionData describes items required to connect to a Geode cluster
type ConnectionData struct {
	Username       string
	Password       string
	Token          string
	LocatorAddress string
}

// ServiceKeyUsers holds the username and password for users identified in a CF service key
type ServiceKeyUsers struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// ServiceKeyUrls holds URL information for endpoints to the PCC manageability service
type ServiceKeyUrls struct {
	Management string `json:"management"`
	Gfsh       string `json:"gfsh"`
}

// ServiceKey aggregates the information returned when requesting a service key from CF
type ServiceKey struct {
	Urls  ServiceKeyUrls    `json:"urls"`
	Users []ServiceKeyUsers `json:"users"`
}

// UserCommand holds command and parameter information entered by a user
type UserCommand struct {
	Command    string
	Parameters map[string]interface{}
}

// IndividualEndpoint holds endpoint information
type IndividualEndpoint struct {
	HTTPMethod  string `json:"httpMethod"`
	URL         string `json:"url"`
	CommandCall string `json:"summary"`
}

// SwaggerInfo holds information returned by calls to the Swagger endpoint for the PCC manageability service
type SwaggerInfo struct {
	Paths map[string]map[string]FurtherEndpointDetails `json:"paths"`
}

// FurtherEndpointDetails provides details about an endpoint
type FurtherEndpointDetails struct {
	Summary string `json:"summary"`
}
