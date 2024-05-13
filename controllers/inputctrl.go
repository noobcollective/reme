package controllers

import (
	"os"
	"fmt"
	"time"
	"bufio"

	"reme/models"
	"github.com/rs/xid"
)

func GetNewEvent() (models.Event, error) {

	var chosen string
	var subject string
	now := time.Now()
	_, offset := now.Zone()

	allowed := map[string] func(*string, int) (models.Event, error) {
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


func setTimerData(subject *string, offset int) (models.Event, error) {

	var hours uint
	var minutes uint

	fmt.Println("Hours: ")
	fmt.Scanf("%d", &hours)

	fmt.Println("Minutes: ")
	fmt.Scanf("%d", &minutes)

	timeIn := time.Now().Add(time.Hour * time.Duration(hours) + time.Minute * time.Duration(minutes))
	jsonTime := timeIn.In(time.Local).Format(time.RFC3339)

	return models.Event{
		ID: xid.New().String(),
		Time: jsonTime,
		Subject: *subject,
		AlreadyDispatched: false,
	}, nil
}

func setPointData(subject *string, offset int) (models.Event, error) {

	var on string
	var at string

	fmt.Println("On: ")
	fmt.Scan(&on)

	fmt.Println("At: ")
	fmt.Scan(&at)

	timestring := fmt.Sprintf("%v %v:00", on, at)

	//seconds := (hours / 3600) + (minutes / 60)
	timeIn, err := time.Parse("2006-01-02 15:04:05", timestring)
	if err != nil {
		return models.Event{}, err
	}

	jsonTime := timeIn.In(time.Local).Format(time.RFC3339)

	return models.Event {
		ID: xid.New().String(),
		Time: jsonTime,
		Subject: *subject,
		AlreadyDispatched: false,
	}, nil
}
