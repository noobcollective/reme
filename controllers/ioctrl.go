package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"reme/models"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slices"
)

// FIXME: Possibility to configure where file is.
const filename = "events.json"

// Reads all events from configured file.
func ReadEvents() (models.Events, *os.File, error) {

	var data models.Events
	
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return data, nil, err
	}
	log.Printf("Successfully opened file %v.\n", filename)

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Sprint("File error: %v", err))
	}

	json.Unmarshal(byteValue, &data)

	return data, file, nil
}


// Filter todays events from the whole events.
// Returns pointer to events.
func GetTodaysEvents() (*models.Events, error) {
	data, _, err := ReadEvents()
	if err != nil {
		return nil, err
	}
	today := time.Now()

	for index, event := range data.Events {
		eventTime, err := time.Parse(time.RFC3339, event.Time)

		if err != nil {
			panic(fmt.Sprintf("Time error: %v", err))
		}

		if (eventTime.Day() == today.Day() &&
			eventTime.Month() == today.Month() &&
			eventTime.Year() == today.Year()) {
				continue
		}

		data.Events = slices.Delete(data.Events, index, index + 1)
	}

	return &data, nil
}


// Writes an event to configured file.
func WriteEvent(event *models.Event) (*models.Events, *os.File, error) {
	data, file, err := ReadEvents()
	if err != nil {
		return nil, nil, err
	}
	

	defer file.Close()
	data.Events = append(data.Events, *event)
	newData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return nil, nil, err
	}

	if _, err := file.WriteAt(newData, 0); err != nil {
		return nil, nil, err
	}

	return &data, file, nil
}

// Continuously watches the configured file and
// notifies listeners for file writes.
func WatchFile(fileChanged chan fsnotify.Event, chanError chan error, fileDone chan bool) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		chanError <- err
	}

	defer watcher.Close()

	go func() {
		for {
			select {
				case event := <-watcher.Events:
					fileChanged <- event
					if event.Op != fsnotify.Write {
						watcher.Add(event.Name)
					}

				case err = <-watcher.Errors:
					chanError <- err
			}
		}
	}()

	// Add the file to watcher.
	if err := watcher.Add("events.json"); err != nil {
		chanError <- err
		watcher.Close()
	}

	<-fileDone
	watcher.Close()
}
