package businesslayer

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/arizon-dread/status-checker-api/models"
)

func GetSystemStatus(id int) (models.Systemstatus, error) {
	var err error = nil
	var message string = ""
	system, err := datalayer.GetSystemStatus(id)
	var client *http.Client
	var transport *http.Transport

	if err != nil {
		fmt.Printf("DataLayer returned error: %v when trying to get system from database\n", err)
		//fail
	} else {
		if strings.Contains(system.CallUrl, "https") {
			//cert-check
			message += checkCert(&system)
			//use clientcert in call
			if system.ClientCertID != nil {
				tlsConfig, err := getTlsConfigWithClientCert(system)
				if err == nil {
					transport.TLSClientConfig = tlsConfig
					client.Transport = transport
				} else {
					fmt.Printf("Failed creating tls client config, %v\n", err)
				}
			}
		}
		if system.CallBody != "" {
			//POST
			contentType := getContentType(system.CallBody)

			resp, err := client.Post(system.CallUrl, contentType, strings.NewReader(system.CallBody))
			message += handleResponse(&system, resp, err)

		} else {
			//GET
			resp, err := client.Get(system.CallUrl)
			message += handleResponse(&system, resp, err)
		}
		if message != "" {
			sendAlert(&system, message)
		} else {
			system.Status = "OK"

		}
		createdSys, err := datalayer.SaveSystemStatus(&system)
		if err != nil {
			fmt.Printf("Couldn't persist status to db, %v\n", err)
		}
		return createdSys, err
	}

	return system, errors.New("not found")
}

func getTlsConfigWithClientCert(system models.Systemstatus) (*tls.Config, error) {
	clientCert, getCCErr := getClientCert(*system.ClientCertID)
	if getCCErr != nil {
		fmt.Printf("Could not load certificates from db, %v\n", getCCErr)
	}
	//encode private key as pem structure
	cert, ccDecryptErr := decryptClientCert(clientCert)
	if ccDecryptErr != nil {
		fmt.Printf("Error decrypting cert, %v", ccDecryptErr)
	}

	url, parseErr := url.Parse(system.CallUrl)
	if parseErr != nil {
		fmt.Printf("Could not parse URL to string, %v\n", parseErr)
	}
	serverCerts := getCertFromUrl(*url)
	caCertPool := x509.NewCertPool()
	for _, sCert := range serverCerts {
		caCertPool.AddCert(sCert)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	err := fmt.Errorf("%w + %w + %w", getCCErr, ccDecryptErr, parseErr)
	return tlsConfig, err

}
func DeleteSystem(id int) (bool, error) {
	rowsAffected, err := datalayer.DeleteSystem(id)
	if rowsAffected > 0 && err == nil {
		return true, err
	} else {
		return false, err
	}

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

func checkCert(system *models.Systemstatus) string {
	var message string = ""
	url, err := url.Parse(system.CallUrl)
	if err == nil {
		certs := getCertFromUrl(*url)
		expirationDate := certs[0].NotAfter
		currentDate := time.Now()

		if err == nil {
			alertDays := currentDate.AddDate(0, 0, -system.CertExpirationDays)
			if expirationDate.After(alertDays) {
				message += fmt.Sprintf("Certificate will expire in less than %d days, expiration datetime: %v\n", system.CertExpirationDays, expirationDate)
			} else {
				system.CertStatus = "OK"
			}
		} else {
			fmt.Printf("CertWarningDays from config could not be converted to int, %v", system.CertExpirationDays)
		}
	} else {
		fmt.Printf("url could not be parsed, %v\n", err)
	}
	return message
}

func handleResponse(system *models.Systemstatus, resp *http.Response, err error) string {
	var message string = ""
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
				system.CallStatus = "OK"
				system.LastOKTime = time.Now()
				if system.AlertHasBeenSent {
					sendOKStatus(system)
				}
			} else {
				message += "Response didn't match expected content \n"
			}
		}
	}
	return message
}
