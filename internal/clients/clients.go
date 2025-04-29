package clients

import (
	"fmt"

	"github.com/google/go-github/v69/github"

	"github.com/robandpdx/gh-unlock-source-repo/pkg/logger"
)

type GitHubClient interface {
	GitHubAuth() (*github.Client, error)
}

type GitHubClientImpl struct {
	githubPAT string
}

func NewGitHubClient(pat string) GitHubClient {
	return &GitHubClientImpl{
		githubPAT: pat,
	}
}

func (g *GitHubClientImpl) GitHubAuth() (*github.Client, error) {
	if g.githubPAT == "" {
		logger.Logger.Error("GITHUB_TOKEN is not set")
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	client := github.NewClient(nil).WithAuthToken(g.githubPAT)
	if client == nil {
		logger.Logger.Error("Failed to create GitHub client")
		return nil, fmt.Errorf("failed to initialize GitHub client")
	}
	return client, nil
}
