package main

import (
	"os"
	"fmt"
	"flag"

	"reme/controllers"

	"github.com/TheCreeper/go-notify"
)

func main() {

	asDaemon := flag.Bool("daemon", false, "Start as daemon.")
	flag.Parse()

	if *asDaemon {
		controllers.StartDaemon()
		return
	}

	ntfs := notify.NewNotification("Test from after return.", "Nothing special going on here.")

	if _, err := ntfs.Show(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not show notification: %v\n", err)
		return
	}

	newEvent, err := controllers.GetNewEvent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get new event: %v\n", err)
		return
	}

	eventData, file, err := controllers.WriteEvent(&newEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing new event: %v\n", err)
		return
	}

	defer file.Close()
	for _, event := range eventData.Events {
		fmt.Println(event.Subject)
	}
}
