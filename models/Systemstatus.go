package models

import (
	"crypto/x509"
	"time"
)

type Systemstatus struct {
	ID            int       `json:"id"`
	CallStatus    string    `json:"callStatus"`
	CallUrl       string    `json:"callUrl"`
	HttpMethod    string    `json:"httpMethod"`
	Cert          []Cert    `json:"-"`
	CertStatus    string    `json:"certStatus"`
	CallBody      string    `json:"callBody"`
	Message       string    `json:"message"`
	ResponseMatch string    `json:"responseMatch"`
	AlertBody     string    `json:"alertBody"`
	AlertUrl      string    `json:"alertUrl"`
	AlertEmail    string    `json:"alertEmail"`
	Status        string    `json:"status"`
	LastOKTime    time.Time `json:"lastOkTime"`
	LastFailTime  time.Time `json:"lastFailTime"`
}
type Cert struct {
	ID          int
	Certificate x509.Certificate
}
