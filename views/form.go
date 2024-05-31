package views

import (
	"errors"
	"fmt"
	"strconv"
	reg "regexp"

	"reme/controllers"
	"reme/entities"
	"reme/tools"

	"github.com/charmbracelet/huh"
	"github.com/rs/xid"
)


// Runs the form an creates the event.
func RunForm() error {
	var err error
	var subject, jsonTime  string

	// Get basic infos
	if err = runForm(&subject, &jsonTime); err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}

		return err
	}

	if err = createEvent(subject, jsonTime); err != nil {
        return err
    }

	fmt.Println("Event created 🎉")
	return nil
}


// Runs a huh? form to ask the user for event details.
// Parses a json formatted time from input and assigns it.
func runForm(subject *string, jsonTime *string) error {
	var eventType string

	// Create a form with needed input fields.
	form := huh.NewForm(
		huh.NewGroup(
			// Ask for option to create new event.
			huh.NewSelect[string]().Title("What would you like to do?").
				Options(
					huh.NewOption("Set a timer", "timer"),
					huh.NewOption("Set an appointment", "appointment"),
				).Value(&eventType),

			huh.NewInput().Title("Subject").Value(subject).Validate(validateEmpty),
		),

		// Informations for a Timer
		huh.NewGroup(
			huh.NewInput().Key("hours").Title("Hours").Placeholder("HH").Validate(validateTime),

			huh.NewInput().Key("minutes").Title("Minutes").Placeholder("MM").Validate(validateTime),
		).WithHideFunc(func() bool { return eventType != "timer" }),

		// Informations for an Appointment
		huh.NewGroup(
			huh.NewInput().Key("date").Title("Date").Placeholder("YYYY-MM-DD").Validate(validateDate),

			huh.NewInput().Key("hour").Title("Hour").Placeholder("HH").Validate(validateEmpty),

			huh.NewInput().Key("minute").Title("Minute").Placeholder("MM").Validate(validateEmpty),
		).WithHideFunc(func() bool { return eventType != "appointment" }),
	)

	var err error
	if err = form.Run(); err != nil {
		return err
	}

	// Assign jsonTime from given input based on event type.
	switch (eventType) {
		case "appointment":
			timeStr := form.GetString("hour") + ":" + form.GetString("minute")
			*jsonTime, err = tools.GetFixedJsonTime(form.GetString("date"), timeStr)

		case "timer":
			var hours, minutes uint64
			hours, err = strconv.ParseUint(form.GetString("hours"), 10, 64)
			minutes, err = strconv.ParseUint(form.GetString("minutes"), 10, 64)
			*jsonTime = tools.GetRelativeJsonTime(hours, minutes)
	}

	if err != nil {
		return err
	}

	return nil
}


func validateEmpty(value string) error {
	if value == "" {
		return errors.New("Please set a value for this field.")
	}

	return nil
}


func validateDate(value string) error {
	if err := validateEmpty(value); err != nil {
        return err
    }

	if match, _ := reg.MatchString(`\d{4}-\d{1,2}-\d{1,2}`, value); !match {
		return errors.New("Please provide date in the format of `YYYY-MM-DD`.")
	}

	return nil
}


func validateTime(value string) error {
	if err := validateEmpty(value); err != nil {
        return err
    }

	if match, _ := reg.MatchString(`\d{1,2}`, value); !match {
		return errors.New("Please provide time in the format of `dd`.")
	}

	return nil
}


func createEvent(subject string, jsonTime string) error {
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
