package main

// import (
// 	"io"
// 	"os"

// 	log "github.com/sirupsen/logrus"
// )

// func checkFile(path string) error {
// 	file, err := os.Open(path)
// 	defer file.Close()
// 	if err != nil {
// 		return err
// 	}
// 	b, err := io.ReadAll(file)
// 	if err != nil {
// 		return err
// 	}
// 	log.WithField("file-path", path).WithField("file-contents", string(b)).Debug("Checked file")
// 	return nil
// }
