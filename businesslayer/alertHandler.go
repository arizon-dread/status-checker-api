package businesslayer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/arizon-dread/status-checker-api/models"
)

func sendAlert(system *models.Systemstatus, message string) {
	system.Status = fmt.Sprintf("ALERT, %v", message)
	if system.AlertUrl != "" {

		body := ""
		if system.AlertBody == "" {
			body = message
		} else {
			if strings.Contains(system.AlertBody, "$message") {
				body = strings.Replace(system.AlertBody, "$message", message, -1)
			}
		}
		if !system.AlertHasBeenSent {
			sendStatus(system, body)
		}
	}
}

func sendOKStatus(system *models.Systemstatus) {
	body := fmt.Sprintf("System %v status OK", system.Name)
	sendStatus(system, body)
	system.AlertHasBeenSent = false
	datalayer.SaveSystemStatus(system)
}

func sendStatus(system *models.Systemstatus, body string) error {
	contentType := getContentType(system.AlertBody)
	_, err := http.Post(system.AlertUrl, contentType, strings.NewReader(body))
	if err != nil {
		fmt.Printf("error while alerting, %v\n", err)
	}

	if system.AlertEmail != "" {
		//send email
	}

	return err
}
