package businesslayer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arizon-dread/status-checker-api/models"
)

func sendAlert(system *models.Systemstatus, message string) {
	system.Status = fmt.Sprintf("ALERT, %v", message)
	if system.AlertUrl != "" {
		contentType := getContentType(system.AlertBody)
		body := ""
		if system.AlertBody == "" {
			body = message
		} else {
			if strings.Contains(system.AlertBody, "$message") {
				body = strings.Replace(system.AlertBody, "$message", message, -1)
			}
		}

		_, err := http.Post(system.AlertUrl, contentType, strings.NewReader(body))
		if err != nil {
			fmt.Printf("error while alerting, %v\n", err)
		}

	}
	if system.AlertEmail != "" {
		//send email

	}
}
