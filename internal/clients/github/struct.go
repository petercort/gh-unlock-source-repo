package github

type UnlockRepoInput struct {
	MigrationId string `json:"migrationId"`
	OrgName     string `json:"orgName"`
	RepoName    string `json:"repoName"`
}

type UnlockRepoResponse struct {
	StatusCode int `json:"statusCode"`
}

type Migration struct {
	Id    int64 `json:"id"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Guid             string `json:"guid"`
	State            string `json:"state"`
	LockRepositories bool   `json:"lock_repositories"`
	Repositories     []struct {
		Id       int64  `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repositories"`
}
