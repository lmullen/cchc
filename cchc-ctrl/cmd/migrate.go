package cmd

import (
	"fmt"
	"os"

	"github.com/lmullen/cchc/common/db"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database to the current schema",
	Long: `Migration will bring the database from its current state to 
the current schema for the application.`,
	PreRun: connectDB,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := timeout()
		defer cancel()
		err := db.MigrateUp(ctx, dbstr)
		if err != nil {
			fmt.Println("Failed to run migrations with this error:")
			fmt.Printf("	%s\n", err)
			shutdown(nil, nil)
			os.Exit(3)
		}
		fmt.Println("Migrated the database successfully")
	},
	PostRun: shutdown,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
