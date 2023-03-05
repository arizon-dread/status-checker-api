package businesslayer

import (
	"fmt"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/arizon-dread/status-checker-api/config"
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
		err = sendEmail(system, body)
	}

	return err
}

func sendEmail(system *models.Systemstatus, body string) error {
	conf := config.GetInstance()

	auth := smtp.PlainAuth("", conf.AlertSMTP.User, conf.AlertSMTP.Password, conf.AlertSMTP.Server)
	recipients := strings.Split(system.AlertEmail, ";")

	err := smtp.SendMail(fmt.Sprintf("%v:%d", conf.AlertSMTP.Server, conf.AlertSMTP.Port), auth,
		fmt.Sprintf("%v@%v", conf.AlertSMTP.User, conf.AlertSMTP.Server), recipients,
		[]byte(system.AlertBody))
	if err != nil {
		fmt.Printf("error sending email, %v\n", err)
	}
	return err
}
