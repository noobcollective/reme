package main

import (
	"fmt"
	"flag"

	"reme/common"
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
		panic(err)
	}

	newEvent, err := controllers.GetNewEvent()
	common.CheckErr(err)

	eventData, file := controllers.WriteEvent(&newEvent)

	defer file.Close()
	
	for _, event := range eventData.Events {
		fmt.Println(event.Subject)
	}
}
