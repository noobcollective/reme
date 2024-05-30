package controllers

import (
	"encoding/json"
	"os"
	"time"

	"reme/entities"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slices"
)

var FilePath string


// Reads all events from configured file.
func ReadEvents() (entities.Events, error) {
	var data entities.Events

	fileContent, err := os.ReadFile(FilePath)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(fileContent, &data); err != nil {
		return data, nil
	}
	return data, nil
}


// Filter todays events from the whole events.
// Returns pointer to events.
func GetTodaysEvents() (*entities.Events, error) {
	data, err := ReadEvents()
	if err != nil {
		return nil, err
	}
	today := time.Now()

	for index, event := range data.Events {
		eventTime, err := time.Parse(time.RFC3339, event.Time)

		if err != nil {
			return nil, err
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
func WriteEvent(currEvent *entities.Event) (*entities.Events, error) {
	data, err := ReadEvents()
	if err != nil {
		return nil, err
	}

	var newEvent bool = false
	for idx, event := range data.Events {
		if event.ID != currEvent.ID {
			continue
		}

		data.Events[idx] = *currEvent
		newEvent = true
	} 

	if !newEvent {
		data.Events = append(data.Events, *currEvent)
	}

	newData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(FilePath, newData, 0644); err != nil {
		return nil, err
	}

	return &data, nil
}


// Continuously watches the configured file and
// notifies listeners for file writes.
func WatchFile(fileChanged chan fsnotify.Event, chanError chan error, fileDone chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		chanError <- err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Write {
					watcher.Add(event.Name)
				}
				fileChanged <- event

			case err = <-watcher.Errors:
				chanError <- err
			}
		}
	}()

	// Add the file to watcher.
	if err := watcher.Add(FilePath); err != nil {
		chanError <- err
		watcher.Close()
	}

	<-fileDone
	watcher.Close()
}
