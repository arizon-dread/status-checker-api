package businessLayer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/arizon-dread/status-checker-api/dataLayer"
	"github.com/arizon-dread/status-checker-api/models"
)

func GetSystemStatus(id int) (models.Systemstatus, error) {
	var err error = nil

	system, err := dataLayer.GetSystemStatus(id)

	if err != nil {
		fmt.Printf("GetAllSystemStatuses returned error: %v", err)
		//fail
	} else {
		if strings.Contains(system.CallUrl.Scheme, "https") {
			certs := getCertFromUrl(system.CallUrl)

			expirationDate := certs[0].NotAfter
			currentDate := time.Now()
			alertDays := currentDate.AddDate(0, 0, -20)
			if expirationDate.After(alertDays) {
				sendAlert(system, "Certificate will expire in less than 20 days")
			} else {
				system.CertStatus = "OK"
			}

		}
		if system.CallBody != "" {
			//post
			contentType := getContentType(system.CallBody)
			resp, err := http.Post(system.CallUrl.String(), contentType, strings.NewReader(system.CallBody))
			if err != nil || resp.StatusCode > 399 {
				message := fmt.Sprintf("Failed posting to endpoint: %v, error was: %v\n", system.AlertUrl.String(), err)
				sendAlert(system, message)
				fmt.Printf(message)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				message := fmt.Sprintf("Failed to read response from endpoint: %v, response was: %v\n", system.AlertUrl.String(), err)
				sendAlert(system, message)
				fmt.Printf(message)
			}
			if strings.Contains(string(body[:]), system.ResponseMatch) {
				system.Status = "OK"
			} else {
				sendAlert(system, "Response didn't match expected content")
			}
		} else {
			resp, err := http.Get(system.CallUrl.String())
			if err != nil || resp.StatusCode > 399 {
				message := fmt.Sprintf("Failed sending GET-request to endpoint: %v, error was: %v\n", system.AlertUrl.String(), err)
				sendAlert(system, message)
				fmt.Printf(message)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if strings.Contains(string(body[:]), system.ResponseMatch) {
				system.Status = "OK"
			} else {
				sendAlert(system, "Response didn't match expected content")
			}
		}
	}
	err = dataLayer.SaveSystemStatus(system)
	if err != nil {
		fmt.Printf("Couldn't persist status to db, %v\n", err)
	}
	return system, err
}

func GetSystemStatuses() ([]models.Systemstatus, error) {
	var err error = nil
	statuses, err := dataLayer.GetAllSystemStatuses()

	if err != nil {
		fmt.Printf("couldn't get systemStatuses from db, %v\n", err)
	}
	return statuses, err
}
