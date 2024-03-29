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

// make it mockable
var dlGetSystemStatus = datalayer.GetSystemStatus

func GetSystemStatus(id int) (models.Systemstatus, error) {
	var err error = nil
	system, err := dlGetSystemStatus(id)
	seconds := 30
	client := &http.Client{Timeout: time.Duration(seconds) * time.Second}
	var transport *http.Transport = http.DefaultTransport.(*http.Transport)

	if err != nil {
		fmt.Printf("DataLayer returned error: %v when trying to get system from database\n", err)
		//fail
	} else {
		if strings.Contains(system.CallUrl, "https") {
			//cert-check
			certErr := checkCert(&system)
			//use clientcert in call
			if certErr == nil {
				if system.ClientCertID != nil {
					tlsConfig, tlsErr := getTlsConfigWithClientCert(system)
					if tlsErr == nil {
						transport.TLSClientConfig = tlsConfig
						client.Transport = transport
					} else {
						fmt.Printf("Failed creating tls client config, %v\n", tlsErr)
					}
				}
			}
		}
		var callErr error = nil
		if system.CallBody != "" {
			//POST
			contentType := getContentType(system.CallBody)

			resp, err := client.Post(system.CallUrl, contentType, strings.NewReader(system.CallBody))
			callErr = blHandleResponse(&system, resp, err)

		} else {
			//GET
			resp, err := client.Get(system.CallUrl)
			callErr = blHandleResponse(&system, resp, err)
		}
		if callErr != nil {
			sendAlert(&system, fmt.Sprintf("%v", callErr))
		} else {
			system.Status = "OK"
		}
		createdSys, err := dlSaveSystemStatus(&system)
		if err != nil {
			fmt.Printf("Couldn't persist status to db, %v\n", err)
		}
		return createdSys, err
	}

	return system, errors.New("NotFound")
}

func getTlsConfigWithClientCert(system models.Systemstatus) (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	var ccDecryptErr error
	var parseErr error
	clientCert, getCCErr := getClientCert(*system.ClientCertID)

	if getCCErr != nil {
		fmt.Printf("Could not load certificates from db, %v\n", getCCErr)
	} else {
		// encode private key as pem structure
		cert, ccDecryptErr := decryptClientCert(clientCert)
		if ccDecryptErr != nil {
			fmt.Printf("Error decrypting cert, %v", ccDecryptErr)
		} else {
			url, parseErr := url.Parse(system.CallUrl)
			if parseErr != nil {
				fmt.Printf("Could not parse URL to string, %v\n", parseErr)
			} else {
				serverCerts := getCertFromUrl(*url)
				caCertPool := x509.NewCertPool()
				for _, sCert := range serverCerts {
					caCertPool.AddCert(sCert)
				}

				tlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      caCertPool,
					CipherSuites: []uint16{
						tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					},
					PreferServerCipherSuites: true,
					MinVersion:               tls.VersionTLS12,
					MaxVersion:               tls.VersionTLS12,
				}
			}
		}
	}

	err := errors.Join(getCCErr, ccDecryptErr, parseErr)
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

func checkCert(system *models.Systemstatus) error {
	var checkErr error
	url, parseErr := url.Parse(system.CallUrl)
	if parseErr == nil {
		certs := getCertFromUrl(*url)
		expirationDate := certs[0].NotAfter
		currentDate := time.Now()

		alertDays := currentDate.AddDate(0, 0, system.CertExpirationDays)
		if !expirationDate.After(alertDays) {
			checkErr = fmt.Errorf("certificate will expire in less than %d days, expiration datetime: %v", system.CertExpirationDays, expirationDate)
		} else {
			system.CertStatus = "OK"
		}
	}
	err := errors.Join(checkErr, parseErr)

	return err
}

// var to make it mockable.
// handleResponse: Checks error code, reads body, checks if it matches the ResponseMatch prop of the models.SystemStatus struct.
// returns: a string of anything that went wrong.
var blHandleResponse = handleResponse

func handleResponse(system *models.Systemstatus, resp *http.Response, err error) error {
	var respErr, readErr, matchErr error
	if err != nil || resp.StatusCode > 399 {
		respErr = fmt.Errorf("failed posting to endpoint: %v", system.CallUrl)
		fmt.Print(respErr)
	} else {
		body, readErr := io.ReadAll(resp.Body)

		if readErr != nil {
			fmt.Print(readErr)
		} else {
			if strings.Contains(string(body[:]), system.ResponseMatch) {
				system.CallStatus = "OK"
				system.LastOKTime = time.Now()
				if system.AlertHasBeenSent {
					sendOKStatus(system)
				}
			} else {
				matchErr = fmt.Errorf("response didn't match expected content")
			}
		}
	}
	return errors.Join(respErr, readErr, matchErr)
}
