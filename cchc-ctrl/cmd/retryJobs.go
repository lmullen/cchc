/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// retryJobsCmd represents the retryJobs command
var retryJobsCmd = &cobra.Command{
	Use:   "retry-jobs",
	Short: "Retry skipped and failed jobs",
	Long: `Jobs run which fail or which are skipped are recorded in the database.
This command deletes those jobs so that they can be retried.
`,
	PreRun: connectDB,
	Run: func(cmd *cobra.Command, args []string) {
		if !force {
			fmt.Println("Retrying skipped/failed jobs will delete them from the database.")
			getConfirmation()
		}

		query := `DELETE FROM jobs.fulltext WHERE status = 'skipped' OR status = 'failed';`
		fmt.Println("Deleting jobs might take a long time ...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		tx, err := database.Begin(ctx)
		if err != nil {
			fmt.Printf("Failed to retry skipped/failed jobs with error:\n	%s\n", err)

		}
		defer tx.Rollback(context.TODO())
		_, err = tx.Exec(ctx, query)
		if err != nil {
			fmt.Printf("Failed to retry skipped/failed jobs with error:\n	%s\n", err)
			tx.Rollback(context.TODO())
			shutdown(nil, nil)
			os.Exit(9)
		}

		err = tx.Commit(ctx)
		if err != nil {
			fmt.Printf("Failed to retry skipped/failed jobs with error:\n	%s\n", err)
			shutdown(nil, nil)
			os.Exit(10)
		}

		fmt.Println("Removed skipped/failed jobs successfully")

	},
	PostRun: shutdown,
}

func init() {
	rootCmd.AddCommand(retryJobsCmd)
}
