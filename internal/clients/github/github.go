package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/robandpdx/gh-unlock-source-repo/internal/clients"
	ghlog "github.com/robandpdx/gh-unlock-source-repo/pkg/logger"

	"go.uber.org/zap"
)

func UnlockRepo(input UnlockRepoInput) (*UnlockRepoResponse, error) {
	// Get environment variables
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubHost := os.Getenv("GITHUB_API_ENDPOINT")

	if githubHost == "" {
		githubHost = "api.github.com"
	}

	// Initialize GitHub client with proper headers
	githubClient := clients.NewGitHubClient(githubToken)
	client, err := githubClient.GitHubAuth()

	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %v", err)
	}

	url := fmt.Sprintf("https://%s/orgs/%s/migrations/%s/repos/%s/lock", githubHost, input.OrgName, input.MigrationId, input.RepoName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GraphQL request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			ghlog.Logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	var unlockResponse UnlockRepoResponse

	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	} else {
		// new UnlockRepoResponse to hold status code
		unlockResponse = UnlockRepoResponse{
			StatusCode: resp.StatusCode,
		}
	}

	return &unlockResponse, nil
}

func GetMigrationId(orgName string, repoName string) (string, error) {
	// Get environment variables
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubHost := os.Getenv("GITHUB_API_ENDPOINT")

	if githubHost == "" {
		githubHost = "api.github.com"
	}

	// Initialize GitHub client with proper headers
	githubClient := clients.NewGitHubClient(githubToken)
	client, err := githubClient.GitHubAuth()

	if err != nil {
		return "", fmt.Errorf("failed to create GitHub client: %v", err)
	}

	url := fmt.Sprintf("https://%s/orgs/%s/migrations", githubHost, orgName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Client().Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make GraphQL request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			ghlog.Logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var migrations []Migration

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &migrations); err != nil {
		ghlog.Logger.Error("Failed to decode response", zap.Error(err))
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Check if the response contains any errors
	if len(migrations) == 0 {
		return "", fmt.Errorf("no migrations found for org: %s", orgName)
	}

	// Loop through the migrations and find the one that contains the repo
	for _, migration := range migrations {
		for _, repo := range migration.Repositories {
			if repo.FullName == fmt.Sprintf("%s/%s", orgName, repoName) {
				return fmt.Sprintf("%d", migration.Id), nil
			}
		}
	}

	// no migration found for the repo, return an error
	return "", fmt.Errorf("no migration found for repo: %s/%s", orgName, repoName)
}
