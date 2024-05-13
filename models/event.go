package models

type Events struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID string `json:"id"`
	Time string `json:"time"`
	Subject string `json:"subject"`
	AlreadyDispatched bool `json:"alreadyDispatched"`
}
