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
var dbstr string
var database *pgxpool.Pool

// getConfig checks for the environment variable or dies trying
func getConfig() {
	str, set := os.LookupEnv("CCHC_DBSTR")
	if !set {
		fmt.Println("The $CCHC_DBSTR environment variable is not set.")
		os.Exit(1)
	}
	strWithApp, err := db.AddApplication(str, "cchc-ctrl")
	if err != nil {
		strWithApp = str
	}

	dbstr = strWithApp
}

// connectDB connects to the database or dies trying
func connectDB(cmd *cobra.Command, args []string) {
	getConfig() // This will set the global variable or fail

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

// getConfirmation gets user confirmation or dies trying
func getConfirmation() {
	fmt.Print("Are you sure you want to proceed? If so, type `yes`: ")
	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		fmt.Println("Confirmation not received")
		shutdown(nil, nil)
		os.Exit(8)
	}
}
