package controllers

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"reme/entities"
	"reme/tools"

	"github.com/rs/xid"
)

// Prompts for info and creates a new event.
func GetNewEvent() (entities.Event, error) {
	var chosen string
	var subject string
	now := time.Now()
	_, offset := now.Zone()

	allowed := map[string] func(*string, int) (entities.Event, error) {
		"t": setTimerData,
		"p": setPointData,
	}

	fmt.Println("Hey you. Press 't' to set a timer or 'p' to set an appointment.")
	fmt.Scan(&chosen)

	if _, ok := allowed[chosen]; ! ok {
		fmt.Printf("Type %v not allowed. Please use either 't' to set a timer or 'p' to set an appointment:\n", chosen)
		fmt.Scan(&chosen)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Subject: ")
	if scanner.Scan() {
		subject = scanner.Text()
	}

	return allowed[chosen](&subject, offset)
}


// Sets the event with relative time given.
func setTimerData(subject *string, offset int) (entities.Event, error) {
	var hours uint64
	var minutes uint64

	fmt.Println("Hours: ")
	fmt.Scanf("%d", &hours)

	fmt.Println("Minutes: ")
	fmt.Scanf("%d", &minutes)

	jsonTime := tools.GetRelativeJsonTime(hours, minutes)
	return entities.Event{
		ID: xid.New().String(),
		Time: jsonTime,
		Subject: *subject,
		AlreadyDispatched: false,
	}, nil
}


// Sets the event with fixed date given.
func setPointData(subject *string, offset int) (entities.Event, error) {
	var on string
	var at string

	fmt.Println("On: ")
	fmt.Scan(&on)

	fmt.Println("At: ")
	fmt.Scan(&at)

	jsonTime, err := tools.GetFixedJsonTime(on, at)
	if err != nil {
		return entities.Event{}, err
	}

	return entities.Event {
		ID: xid.New().String(),
		Time: jsonTime,
		Subject: *subject,
		AlreadyDispatched: false,
	}, nil
}
