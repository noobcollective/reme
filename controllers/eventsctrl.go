package controllers

import (
	"fmt"
	"log"
	"time"

	"reme/models"

	"github.com/fsnotify/fsnotify"
	"github.com/TheCreeper/go-notify"
)

// Listens to channel, if there were new events added to events file.
func WatchForNewEvents(fileChange chan fsnotify.Event, noEvents chan bool, fileError chan error) {

	// Get todays events.
	data, _, err := ReadEvents()
	if err != nil {
		log.Fatal("Error reading events: ", err)
	}

	if len(data.Events) == 0 {
		log.Println("No events for today.")
		return
	}

	log.Printf("Todays events are: %v", data.Events)

	// Listen to channel.
	for {
		select {
			case event := <-fileChange:
				switch event.Op {
					case fsnotify.Write:
						log.Println("New write to file.")
						updateData(&data)
						log.Printf("New events are: %v\n", data)

					default:
						log.Printf("Other operations: %v", event)
				}

			case <-noEvents:
				log.Println("No more events.")

			case <-time.Tick(5 * time.Second):
				for idx := range data.Events {
					checkIfNow(&data.Events[idx], fileError)
				}
		}
	}
}


// Update data with new events.
func updateData(data *models.Events) {
	newData, _, err := ReadEvents()
	if err != nil {
		log.Fatal("Error reading events: ", err)
	}

	if ( (len(data.Events) == len(newData.Events)) ||
			(len(data.Events) == 0 || len(newData.Events) == 0) ) {
		return
	}

	data.Events = append(data.Events, newData.Events[len(newData.Events) - 1])
}

// Checks if the diff from now to the event time is between +/- 5 seconds.
// If so, notify listeners and mark event as passed.
func checkIfNow(event *models.Event, fileError chan error) {
	if ( event.AlreadyDispatched ) {
		return
	}

	eventTime, err := time.Parse(time.RFC3339, event.Time)

	if err != nil {
		log.Println("Date parsing error... skipping this event.")
	}

	// Event is out of threshold
	if ( !time.Now().After( eventTime ) ) {
		return
	}

	event.AlreadyDispatched = true
	ntfs := notify.NewNotification( "REME Notification", fmt.Sprintf("Event %s passed.", event.Subject) )

	if _, err := ntfs.Show(); err != nil {
		panic(err)
	}
}
