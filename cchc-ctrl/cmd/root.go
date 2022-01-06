package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var force bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cchc-ctrl",
	Short: "Run database management commands",
	Long: `This program runs configuration and management commands for your database.
This includes migrations for the database as well as tasks such as removing 
skipped or failed jobs.

To run this program, you must set the CCHC_DBSTR environment variable to a 
PostgreSQL connection string, such as the following:
	postgres://user:password@hostname:5432/database?sslmode=disable
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "run commands without waiting for user input")
}
