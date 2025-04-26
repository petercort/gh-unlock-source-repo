package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/robandpdx/gh-unlock-repo/internal/clients"
	ghlog "github.com/robandpdx/gh-unlock-repo/pkg/logger"

	"go.uber.org/zap"
)

func UnlockRepo(input UnlockRepoInput) (*UnlockRepoResponse, error) {
	ghlog.Logger.Info("Unlocking repository",
		zap.String("orgName", input.OrgName),
		zap.String("repoName", input.RepoName))

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

	mutation := `
	mutation unlockImportedRepositories(
			$migrationId: ID!
			$org: String!
			$repo: String!
	) {
			unlockImportedRepositories(
					input: {
						migrationId: $migrationId,
						org:         "${input.OrgName}",
						repo:        "${input.RepoName}"
					}
			) {
					migration {
						guid
						id
						state
					}
					unlockedRepositories {
						nameWithOwner
					}
			}
	}`

	requestBody := map[string]interface{}{
		"query": mutation,
		"variables": map[string]interface{}{
			"migrationId": input.MigrationId,
			"org":         input.OrgName,
			"repo":        input.RepoName,
		},
		"operationName": "unlockImportedRepositories",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		ghlog.Logger.Error("Failed to marshal request body", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	url := fmt.Sprintf("https://%s/graphql", githubHost)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		ghlog.Logger.Error("Failed to create HTTP request", zap.Error(err))
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gh-glx-migrator")
	req.Header.Set("GraphQL-Features", "octoshift_gl_exporter")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Client().Do(req)
	if err != nil {
		ghlog.Logger.Error("Failed to make GraphQL request", zap.Error(err))
		return nil, fmt.Errorf("failed to make GraphQL request: %v", err)
	}

	// if the response status is not 200, show error message
	// and the response body and return an error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, func() string {
			body, _ := io.ReadAll(resp.Body)
			return string(body)
		}())
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			ghlog.Logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ghlog.Logger.Error("Failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	ghlog.Logger.Debug("Raw response", zap.String("body", string(body)))

	var response struct {
		Data   UnlockRepoResponse `json:"data"`
		Errors []struct {
			Message string   `json:"message"`
			Type    string   `json:"type"`
			Path    []string `json:"path"`
		} `json:"errors,omitempty"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		ghlog.Logger.Error("Failed to decode response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(response.Errors) > 0 {
		errMsg := response.Errors[0].Message
		ghlog.Logger.Error("GraphQL mutation returned an error",
			zap.String("error", errMsg),
			zap.String("type", response.Errors[0].Type),
			zap.Strings("path", response.Errors[0].Path))
		return nil, fmt.Errorf("GraphQL error: %s", errMsg)
	}

	if response.Data.StatusCode != 200 {
		errMsg := fmt.Sprintf("failed to unlock repository: %s", response.Errors[0].Message)
		ghlog.Logger.Error("Failed to unlock repository",
			zap.String("error", errMsg),
			zap.String("type", response.Errors[0].Type),
			zap.Strings("path", response.Errors[0].Path))
		return nil, fmt.Errorf("failed to unlock repository: %s", errMsg)
	}

	ghlog.Logger.Info("Successfully unlocked repository",
		zap.String("orgName", input.OrgName),
		zap.String("repoName", input.RepoName),
		zap.String("migrationId", input.MigrationId),
		zap.String("statusCode", fmt.Sprintf("%d", response.Data.StatusCode)))

	return &response.Data, nil
}
