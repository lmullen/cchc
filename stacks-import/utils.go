package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// Check that there is no problem with reading the files passed in
func checkPathsToBatches(paths []string) error {
	for _, path := range paths {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return err
		} else if filepath.Ext(path) != ".json" {
			err := fmt.Errorf("Error: %s is not a .json file", path)
			return err
		} else if err != nil {
			return err
		}
	}
	return nil
}

func extractDateString(j json.RawMessage) string {
	s := string(j)
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	return re.FindString(s)
}

func year(date string) (int, error) {
	if len(date) >= 4 {
		date = date[0:4]
	}
	year, err := strconv.Atoi(date)
	return year, err
}
