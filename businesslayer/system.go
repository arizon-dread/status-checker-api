package businesslayer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/arizon-dread/status-checker-api/config"
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
		if strings.Contains(system.CallUrl, "https") {
			//cert-check
			url, err := url.Parse(system.CallUrl)
			if err == nil {
				certs := getCertFromUrl(*url)
				cfg := config.GetInstance()
				expirationDate := certs[0].NotAfter
				currentDate := time.Now()
				certWarningDays, err := strconv.Atoi(cfg.General.CertWarningDays)
				if err == nil {
					alertDays := currentDate.AddDate(0, 0, -certWarningDays)
					if expirationDate.After(alertDays) {
						message += fmt.Sprintf("Certificate will expire in less than 20 days, expiration time: %v\n", expirationDate)
					} else {
						system.CertStatus = "OK"
					}
				} else {
					fmt.Printf("CertWarningDays from config could not be converted to int, %v", cfg.General.CertWarningDays)
				}
			} else {
				fmt.Printf("url could not be parsed, %v\n", err)
			}
		}
		if system.CallBody != "" {
			//POST
			contentType := getContentType(system.CallBody)

			resp, err := http.Post(system.CallUrl, contentType, strings.NewReader(system.CallBody))
			if err != nil || resp.StatusCode > 399 {
				message += fmt.Sprintf("Failed posting to endpoint: %v, error was: %v\n", system.CallUrl, err)
				fmt.Print(message)
			} else {
				body, err := io.ReadAll(resp.Body)

				if err != nil {
					message += fmt.Sprintf("Failed to read response from endpoint: %v, response was: %v\n", system.CallUrl, err)
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
			resp, err := http.Get(system.CallUrl)
			if err != nil || resp.StatusCode > 399 {
				message += fmt.Sprintf("Failed sending GET-request to endpoint: %v, error was: %v\n", system.CallUrl, err)
				fmt.Print(message)
			}
			body, err := io.ReadAll(resp.Body)
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
	createdSys, err := datalayer.SaveSystemStatus(&system)
	if err != nil {
		fmt.Printf("Couldn't persist status to db, %v\n", err)
	}
	return createdSys, err
}

func GetSystemStatuses() ([]models.Systemstatus, error) {
	var err error = nil
	statuses, err := datalayer.GetAllSystemStatuses()

	if err != nil {
		fmt.Printf("couldn't get systemStatuses from db, %v\n", err)
	}
	return statuses, err
}

func SaveSystemStatus(system models.Systemstatus) (models.Systemstatus, error) {

	system, err := datalayer.SaveSystemStatus(&system)
	return system, err
}
