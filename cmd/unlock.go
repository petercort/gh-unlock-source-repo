package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/robandpdx/gh-unlock-repo/internal/clients/github"
	ghlog "github.com/robandpdx/gh-unlock-repo/pkg/logger"
)

func UnlockRepo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gh-unlock-repo --org <org-name> --repo <repo-name>",
		Short: "Unlock a repository in GitHub",
		Long:  `Unlock the specified repository in GitHub.`,

		Example: `  gh glx unlock-repo --org my-org --repo my-repo`,
		RunE:    unlockRepo,
	}

	return cmd
}

func unlockRepo(cmd *cobra.Command, args []string) error {
	ghlog.Logger.Info("Reading input values for generating pre-signed URL")

	org, err := cmd.Flags().GetString("org")
	if err != nil {
		ghlog.Logger.Error("failed to get org flag", zap.Error(err))
		return err
	}

	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		ghlog.Logger.Error("failed to get repo flag", zap.Error(err))
		return err
	}

	ghlog.Logger.Info("Unlocking repository", zap.String("org", org), zap.String("repo", repo))

	input := github.UnlockRepoInput{
		MigrationId: "test",
		OrgName:     org,
		RepoName:    repo,
	}

	response, err := github.UnlockRepo(input)

	// log the satus code
	ghlog.Logger.Info("Unlocking repository status code", zap.Int("statusCode", response.StatusCode))

	// =========================================

	// ctx := context.Background()

	// // Initialize GitHub GraphQL client (replace YOUR_GITHUB_TOKEN with actual token handling)
	// src := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	// )
	// httpClient := oauth2.NewClient(ctx, src)
	// client := githubv4.NewClient(httpClient)

	// var mutation struct {
	// 	UnlockImportedRepositories struct {
	// 		UnlockedRepositories []struct {
	// 			Name string
	// 		}
	// 	} `graphql:"unlockImportedRepositories(input: {owner: $owner, repositories: [$repo]})"`
	// }

	// variables := map[string]interface{}{
	// 	"owner": githubv4.String(org),
	// 	"repo":  githubv4.String(repo),
	// }

	// err = client.Mutate(ctx, &mutation, variables)
	// if err != nil {
	// 	ghlog.Logger.Error("failed to unlock repository", zap.Error(err))
	// 	return err
	// }

	// if len(mutation.UnlockImportedRepositories.UnlockedRepositories) == 0 {
	// 	errMsg := fmt.Sprintf("no repositories unlocked for org: %s, repo: %s", org, repo)
	// 	ghlog.Logger.Error(errMsg)
	// 	return fmt.Errorf(errMsg)
	// }

	// ghlog.Logger.Info("Successfully unlocked repository", zap.String("org", org), zap.String("repo", repo))

	// =========================================
	return nil
}
