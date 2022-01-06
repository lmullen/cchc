package cmd

import (
	"fmt"
	"os"

	"github.com/lmullen/cchc/common/db"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database (deletes all data)",
	Long: `This command deletes all data in the database and migrates it to the 
current schema. This WILL result in data loss; use with caution. This command
require interactive user confirmation, unless it is run with the --force/-f flag.
`,
	PreRun: connectDB,
	Run: func(cmd *cobra.Command, args []string) {
		if !force {
			fmt.Println("Reseting the database will delete all your data.")
			getConfirmation()
		}

		ctx, cancel := timeout()
		defer cancel()
		err := db.MigrateDown(ctx, dbstr)
		if err != nil {
			fmt.Println("Failed to reset the database with this error:")
			fmt.Printf("	%s\n", err)
			os.Exit(5)
		}

		err = db.MigrateUp(ctx, dbstr)
		if err != nil {
			fmt.Println("Failed to reset the database with this error:")
			fmt.Printf("	%s\n", err)
			shutdown(nil, nil)
			os.Exit(6)
		}

		fmt.Println("Reset the database successfully")
	},
	PostRun: shutdown,
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
