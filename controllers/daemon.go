package controllers

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"syscall"
	"context"
	"os/signal"

	"github.com/fsnotify/fsnotify"
)

func StartDaemon() {

	// Initialize context and sigchannel
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Defer the cancelation of the context, until it's done.
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
				case s := <-signalChan:
					switch s {
						case syscall.SIGINT, syscall.SIGTERM:
							log.Printf("Got SIGINT/SIGTERM, exiting.")
							cancel()
							os.Exit(1)

						case syscall.SIGHUP:
							log.Printf("Got SIGHUP, reloading config now.")
					}
				case <-ctx.Done():
					log.Printf("Done.")
					os.Exit(1)
			}
		}
	}()

	if err := controlDaemon(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func controlDaemon(ctx context.Context, out io.Writer) (error) {

	fileDone := make(chan bool)
	chanError := make(chan error)
	fileChange := make(chan fsnotify.Event)
	noEvents := make(chan bool)

	go WatchFile(fileChange, chanError, fileDone)
	go WatchForNewEvents(fileChange, noEvents, chanError)

	for {
		select {
			case <-ctx.Done():
				fileDone <- true
				log.Println("Done, exiting now.")
				os.Exit(1)

			case err := <-chanError:
				log.Fatal("Channel faile with an error: ", err)

			case <-time.Tick(10 * time.Second):
				log.Println("Still running...")
		}
	}
}
