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

type ClusterManagementResults struct {
	StatusCode string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	MemberStatus []MemberStatus `json:"memberStatus"`
	Results []IndividualClusterManagementResult `json:"result"`
}

type IndividualClusterManagementResult struct {
	Config map[string]interface{} `json:"config"`
	RuntimeInfo []map[string]interface{} `json:"runtimeInfo"`
}

type MemberStatus struct {
	ServerName string
	Success bool
	Message string
}

type PostJson struct {
	name string `json:"name"`
	_type string `json:"type"`
}
