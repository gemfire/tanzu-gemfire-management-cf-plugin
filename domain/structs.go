package domain

// CommandData is all the data involved in executing plugin commands
// This data gets manipulated throughout the program
type CommandData struct {
	Target             string
	ConnnectionData    ConnectionData
	UserCommand        UserCommand
	AvailableEndpoints map[string]RestEndPoint //key is command name
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
	Parameters map[string]string
}

// RestEndPoint holds endpoint information
type RestEndPoint struct {
	HTTPMethod  string
	URL         string
	CommandName string
	JQFilter    string
	Parameters  []RestAPIParam
}

// RestAPI is used to parse the swagger json response
// first key: url | second key: method (get/post) | value: RestAPIDetail
type RestAPI struct {
	Paths       map[string]map[string]RestAPIDetail `json:"paths"`
	Definitions map[string]DefinitionDetail         `json:"definitions"`
}

// RestAPIDetail provides details about an endpoint
type RestAPIDetail struct {
	CommandName string         `json:"summary"`
	JQFilter    string         `json:"jqFilter"`
	Parameters  []RestAPIParam `json:"parameters"`
}

// DefinitionDetail describes the details of the type definitions
type DefinitionDetail struct {
	Properties map[string]PropertyDetail `json:"properties"`
}

// PropertyDetail describes the details of the properties of type definitions
type PropertyDetail struct {
	Type   string            `json:"type"`
	Ref    string            `json:"$ref"`
	Enum   []string          `json:"enum"`
	Format string            `json:"format"`
	Items  map[string]string `json:"items"`
}

// RestAPIParam contains the information about possible parameters for a call
type RestAPIParam struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	// In describes how params are submitted: "query", "body" or "path"
	In             string            `json:"in"`
	Schema         map[string]string `json:"schema"`
	BodyDefinition map[string]interface{}
}
