package businessLayer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arizon-dread/status-checker-api/models"
)

func sendAlert(system *models.Systemstatus, message string) {
	if system.AlertUrl.String() != "" {
		contentType := getContentType(system.AlertBody)
		body := ""
		if system.AlertBody == "" {
			body = message
		} else {
			if strings.Contains(system.AlertBody, "$message") {
				body = strings.Replace(system.AlertBody, "$message", message, -1)
			}
		}

		_, err := http.Post(system.AlertUrl.String(), contentType, strings.NewReader(body))
		if err != nil {
			fmt.Printf("error while alerting, %v\n", err)
		}

		if system.AlertEmail != "" {
			//send email

		}
	}
}
