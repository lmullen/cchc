package main

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var db *pgxpool.Pool

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path/to/a/batch.json ...> \n", os.Args[0])
		flag.PrintDefaults()
	}
	debug := flag.BoolP("debug", "d", false, "Turn on debugging messages")
	help := flag.Bool("help", false, "help")
	flag.Parse()
	batches := flag.Args()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if len(batches) == 0 {
		log.Error("Provide paths to Stacks JSON files")
		flag.Usage()
		os.Exit(1)
	}

	err := checkPathsToBatches(batches)
	if err != nil {
		log.Println("Error with arguments: ", err)
		os.Exit(1)
	}

	db, err = DBConnect()
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	defer db.Close()

	for _, b := range batches {
		processBatch(b)
	}

}
