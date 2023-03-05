package models

import (
	"time"
)

type Systemstatus struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	CallStatus         string    `json:"callStatus"`
	CallUrl            string    `json:"callUrl"`
	CallBody           string    `json:"callBody"`
	HttpMethod         string    `json:"httpMethod"`
	CertStatus         string    `json:"certStatus"`
	CertExpirationDays int       `json:"certExpirationDays"`
	Message            string    `json:"message"`
	ResponseMatch      string    `json:"responseMatch"`
	AlertBody          string    `json:"alertBody"`
	AlertUrl           string    `json:"alertUrl"`
	AlertEmail         string    `json:"alertEmail"`
	AlertHasBeenSent   bool      `json:"alertHasBeenSent"`
	Status             string    `json:"status"`
	LastOKTime         time.Time `json:"lastOkTime"`
	LastFailTime       time.Time `json:"lastFailTime"`
	ClientCert_ID      *int      `json:"clientCertId"`
}
