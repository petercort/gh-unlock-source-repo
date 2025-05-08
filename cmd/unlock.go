package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/robandpdx/gh-unlock-source-repo/internal/clients/github"
	ghlog "github.com/robandpdx/gh-unlock-source-repo/pkg/logger"
)

func UnlockRepo() *cobra.Command {
	cmd := &cobra.Command{
		RunE: unlockRepo,
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
	ghlog.Logger.Info("Getting migration ID",
		zap.String("org", org),
		zap.String("repo", repo),
	)
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

	ghlog.Logger.Info("Successfully unlocked repository",
		zap.String("statusCode", strconv.Itoa(response.StatusCode)),
	)

	return nil
}
