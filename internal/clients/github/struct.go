package github

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type UnlockRepoInput struct {
	MigrationId string `json:"migrationId"`
	OrgName     string `json:"orgName"`
	RepoName    string `json:"repoName"`
}

type UnlockRepoResponse struct {
	StatusCode int `json:"statusCode"`
}
