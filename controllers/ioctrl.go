package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"reme/common"
	"reme/models"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slices"
)

// FIXME: Possibility to configure where file is.
const filename = "events.json"

// Reads all events from configured file.
func ReadEvents() (models.Events, *os.File) {

	var data models.Events
	
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	common.CheckErr(err)
	log.Printf("Successfully opened file %v.\n", filename)

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprint("File error: %v", err))
	}

	json.Unmarshal(byteValue, &data)

	return data, file
}


// Filter todays events from the whole events.
// Returns pointer to events.
func GetTodaysEvents() (*models.Events) {
	data, _ := ReadEvents()
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

	return &data
}


// Writes an event to configured file.
func WriteEvent(event *models.Event) (*models.Events, *os.File) {

	data, file := ReadEvents()

	defer file.Close()

	data.Events = append(data.Events, *event)

	newData, err := json.MarshalIndent(data, "", "	")

	if err != nil {
		panic(fmt.Sprintf("JSON error: %v", err))
	}

	_, writeErr := file.WriteAt(newData, 0)
	if writeErr != nil {
		panic(fmt.Sprintf("File error: %v", err))
	}

	return &data, file
}

// Continuously watches the configured file and
// notifies listeners for file writes.
func WatchFile(fileChanged chan fsnotify.Event, fileError chan error, fileDone chan bool) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fileError <- err
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
					fileError <- err
			}
		}
	}()

	// Add the file to watcher.
	if err := watcher.Add("events.json"); err != nil {
		fileError <- err
		watcher.Close()
	}

	<-fileDone
	watcher.Close()
}
