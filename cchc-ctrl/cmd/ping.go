package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Check connection to the database",
	Long: `Checks the connection to the database by seeing whether the 
$CCHC_DBSTR environment variable is set and whether you can succesfully 
connect to the database.`,
	PreRun: connectDB,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := timeout()
		defer cancel()
		err := database.Ping(ctx)
		if err != nil {
			fmt.Println("Failed to connect to and ping the database:")
			fmt.Printf("	%s\n", err)
			return
		}
		fmt.Println("Successfully connected to and pinged the database")
	},
	PostRun: shutdown,
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
