package businesslayer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/arizon-dread/status-checker-api/models"
)

func GetSystemStatus(id int) (models.Systemstatus, error) {
	var err error = nil
	var message string = ""
	system, err := datalayer.GetSystemStatus(id)

	if err != nil {
		fmt.Printf("DataLayer returned error: %v when trying to get system from database\n", err)
		//fail
	} else {
		if strings.Contains(system.CallUrl.Scheme, "https") {
			certs := getCertFromUrl(system.CallUrl)

			expirationDate := certs[0].NotAfter
			currentDate := time.Now()
			alertDays := currentDate.AddDate(0, 0, -20)
			if expirationDate.After(alertDays) {
				message += fmt.Sprintf("Certificate will expire in less than 20 days, expiration time: %v\n", expirationDate)
			} else {
				system.CertStatus = "OK"
			}

		}
		if system.CallBody != "" {
			//POST
			contentType := getContentType(system.CallBody)

			resp, err := http.Post(system.CallUrl.String(), contentType, strings.NewReader(system.CallBody))
			if err != nil || resp.StatusCode > 399 {
				message += fmt.Sprintf("Failed posting to endpoint: %v, error was: %v\n", system.CallUrl.String(), err)
				fmt.Print(message)
			} else {
				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					message += fmt.Sprintf("Failed to read response from endpoint: %v, response was: %v\n", system.CallUrl.String(), err)
					fmt.Print(message)
				} else {
					if strings.Contains(string(body[:]), system.ResponseMatch) {
						system.Status = "OK"
					} else {
						message += "Response didn't match expected content"
					}
				}
			}
		} else {
			//GET
			resp, err := http.Get(system.CallUrl.String())
			if err != nil || resp.StatusCode > 399 {
				message += fmt.Sprintf("Failed sending GET-request to endpoint: %v, error was: %v\n", system.CallUrl.String(), err)
				fmt.Print(message)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				message += "Failed reading body of response\n"
			}
			if strings.Contains(string(body[:]), system.ResponseMatch) {
				system.Status = "OK"
			} else {
				message += "Response didn't match expected content\n"
			}
		}
		if message != "" {
			sendAlert(&system, message)
		}
	}
	err = datalayer.SaveSystemStatus(&system)
	if err != nil {
		fmt.Printf("Couldn't persist status to db, %v\n", err)
	}
	return system, err
}

func GetSystemStatuses() ([]models.Systemstatus, error) {
	var err error = nil
	statuses, err := datalayer.GetAllSystemStatuses()

	if err != nil {
		fmt.Printf("couldn't get systemStatuses from db, %v\n", err)
	}
	return statuses, err
}
