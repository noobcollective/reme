package main

import (
	"encoding/json"
	"os"
	"fmt"
	"flag"

	"reme/controllers"
	"reme/entities"
)

var filePath = flag.String("ef", "", "Custom path to events file.")

func main() {
	asDaemon := flag.Bool("daemon", false, "Start as daemon.")
	flag.Parse()

	if err := checkOrCreateFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error checking or creating file: <%v>\n", err)
		return
	}

	fmt.Println(*filePath)
	controllers.FilePath = *filePath

	if *asDaemon {
		controllers.StartDaemon()
		return
	}

	newEvent, err := controllers.GetNewEvent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get new event: %v\n", err)
		return
	}

	if _, err := controllers.WriteEvent(&newEvent); err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing new event: %v\n", err)
		return
	}
}


// Checks if events file already exists.
// If not events file will be created in `~/.local/share/reme` folder.
func checkOrCreateFile() error {
	// Events file configured via arguments.
	if *filePath != "" {
		return nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filesDir := userHome + "/.local/share/reme"
	tmpFilePath := filesDir + "/events.json"
	if _, err := os.Stat(tmpFilePath); err == nil {
		*filePath = tmpFilePath
		return nil
	}

	if err := os.MkdirAll(filesDir, 0777); err != nil {
		return err
	}

	*filePath = tmpFilePath
	emptyObj := entities.Events{Events: make([]entities.Event, 0)}
	emptyData, err := json.Marshal(emptyObj)
	if err != nil {
		return err
	}

	if err := os.WriteFile(*filePath, []byte(emptyData), 0644); err != nil {
		return err
	}

	return nil
}
