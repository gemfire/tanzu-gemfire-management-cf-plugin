package main
type BasicPlugin struct{}

type ServiceKeyUsers struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type ServiceKeyUrls struct {
	Management string `json:"management"`
	Gfsh string `json:"gfsh"`
}

type ServiceKey struct {
	Urls  ServiceKeyUrls    `json:"urls"`
	Users []ServiceKeyUsers `json:"users"`
}

type RestAPICall struct {
	command string
	parameters map[string]string
}

type IndividualEndpoint struct {
	HttpMethod string 	`json:"httpMethod"`
	Url string 			`json:"url"`
	CommandCall string 	`json:"summary"`
}


type SwaggerInfo struct {
	Paths map[string]map[string]FurtherEndpointDetails `json:"paths"`
}

type FurtherEndpointDetails struct {
	Summary string `json:"summary"`
}


