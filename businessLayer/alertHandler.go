package businessLayer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arizon-dread/status-checker-api/models"
)

func sendAlert(system models.Systemstatus, message string) {
	if system.AlertUrl.String() != "" {
		contentType := getContentType(system.AlertBody)
		if system.AlertBody == "" {
			system.AlertBody = message
		}
		_, err := http.Post(system.AlertUrl.String(), contentType, strings.NewReader(system.AlertBody))
		if err != nil {
			fmt.Printf("error while alerting, %v\n", err)
		}

		if system.AlertEmail != "" {
			//send email

		}
	}
}
