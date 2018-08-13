package timer

import (
	"github.com/reef-pi/reef-pi/controller/types"
	"log"
)

type Reminder struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ReminderRunner struct {
	telemetry   types.Telemetry
	title, body string
}

func (r ReminderRunner) Run() {
	_, err := r.telemetry.Alert(r.title, r.body)
	if err != nil {
		log.Println("ERROR: Failed to send reminder. Error:", err)
	}
}
