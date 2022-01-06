package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lmullen/cchc/common/db"
	"github.com/spf13/cobra"
)

// database is a global variable for the database
var database *pgxpool.Pool

// connectDB connects to the database or dies trying
func connectDB(cmd *cobra.Command, args []string) {
	dbstr, set := os.LookupEnv("CCHC_DBSTR")
	if !set {
		fmt.Println("The $CCHC_DBSTR environment variable is not set.")
		os.Exit(1)
	}

	ctx, cancel := timeout()
	defer cancel()
	conn, err := db.Connect(ctx, dbstr, "cchc-ctrl")
	if err != nil {
		fmt.Println("Failed to connect to the database with this error:")
		fmt.Printf("	%s\n", err)
		os.Exit(2)
	}
	database = conn
}

// shutdown makes sure we close the database
func shutdown(cmd *cobra.Command, args []string) {
	database.Close()
}

// timeout returns a timeout context
func timeout() (context.Context, context.CancelFunc) {
	time := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), time)
	return ctx, cancel
}
