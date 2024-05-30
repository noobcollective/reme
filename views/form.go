package views

import (
	"fmt"
	"reme/controllers"
	"reme/entities"
	"reme/tools"
	"strconv"
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/rs/xid"
)

// Basics
var (
    eventType string
	subject   string
	jsonTime  string
)


func RunForms() error {
	var err error
	// Get basic infos
	if err = runBasicsForm(); err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}

		return err
	}

	switch(eventType) {
		case "timer":
			err = runTimerForm()
	
		case "appointment":
			err = runAppointmentForm()
		
		default:
			return nil
	}

	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}

		return err
	}

	if err = createEvent(); err != nil {
        return err
    }

	fmt.Println("Event created 🎉")

	return nil
}


// Runs a huh? form to ask the user for basic event details.
func runBasicsForm() error {
	form := huh.NewForm(
		huh.NewGroup(
			// Ask what to do.
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
				huh.NewOption("Set a timer", "timer"),
				huh.NewOption("Set an appointment", "appointment"),
				).
				Value(&eventType), // store the chosen option in the "burger" variable

			// Ask for subject.
			huh.NewInput().Title("Subject").Value(&subject).Validate(validateEmpty),
		),
	)

	return form.Run()
}


// Runs a huh? form to ask the user for timer dates.
func runTimerForm() error {
	var strHours, strMinutes string

	form := huh.NewForm(
		huh.NewGroup(
			// Ask for subject.
			huh.NewInput().Title("Hours").Placeholder("HH").Value(&strHours).Validate(validateEmpty),
			huh.NewInput().Title("Minutes").Placeholder("MM").Value(&strMinutes).Validate(validateEmpty),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	hours, err := strconv.ParseUint(strHours, 10, 64)
	minutes, err := strconv.ParseUint(strMinutes, 10, 64)
	if err != nil {
		return err
	}

	jsonTime = tools.GetRelativeJsonTime(hours, minutes)

	return nil
}

// Runs a huh? form to ask the user for timer dates.
func runAppointmentForm() error {
	var err error
	var eventDate, eventTime string

	form := huh.NewForm(
		huh.NewGroup(
			// Ask for subject.
			huh.NewInput().Title("Date").Placeholder("YYYY-MM-DD").Value(&eventDate).Validate(validateEmpty),
			huh.NewInput().Title("Time").Placeholder("HH:MM").Value(&eventTime).Validate(validateEmpty),
		),
	)

	if err = form.Run(); err != nil {
		return err
	}

	jsonTime, err = tools.GetFixedJsonTime(eventDate, eventTime)
	if err != nil {
		return err
	}

	return nil
}


func validateEmpty(value string) error {
	if value != "" {
		return nil
	}

	return errors.New("Please set a value for this field.")
}


func createEvent() error {
	newEvent := entities.Event{
		ID: xid.New().String(),
		Time: jsonTime,
		Subject: subject,
		AlreadyDispatched: false,
	}

	if _, err := controllers.WriteEvent(&newEvent); err != nil {
		return err
	}

	return nil
}
