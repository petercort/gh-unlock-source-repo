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

	// Get the migation ID
	ghlog.Logger.Info("Getting migration ID")
	migrationId, err := github.GetMigrationId(org, repo)
	if err != nil {
		ghlog.Logger.Error("failed to get migration ID", zap.Error(err))
		return err
	}

	// unlock the repository
	ghlog.Logger.Info("Unlocking repository",
		zap.String("org", org),
		zap.String("repo", repo),
		zap.String("migrationId", migrationId),
	)

	input := github.UnlockRepoInput{
		MigrationId: migrationId,
		OrgName:     org,
		RepoName:    repo,
	}

	response, err := github.UnlockRepo(input)

	// log the satus code
	ghlog.Logger.Info("Unlocking repository status code", zap.Int("statusCode", response.StatusCode))

	return nil
}
