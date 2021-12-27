package main

// import (
// 	"context"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"syscall"

// 	log "github.com/sirupsen/logrus"
// )

// var app = &App{}

// func main() {
// 	// Check for interrupt signals and stop gracefully
// 	ctx, cancel := context.WithCancel(context.Background())
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
// 	defer func() {
// 		signal.Stop(quit)
// 		cancel()
// 	}()
// 	go func() {
// 		select {
// 		case <-quit:
// 			log.Info("Received signal to quit the process ...")
// 			cancel()
// 		case <-ctx.Done():
// 		}
// 	}()

// 	err := app.Init(ctx)
// 	if err != nil {
// 		log.WithError(err).Fatal("Error initializing application")
// 	}
// 	defer app.Shutdown()

// 	wg := sync.WaitGroup{}

// 	wg.Add(1)
// 	go func() {
// 		enqueueForQuotations(ctx)
// 		wg.Done()
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		enqueueForLanguages(ctx)
// 		wg.Done()
// 	}()

// 	wg.Wait()

// }
