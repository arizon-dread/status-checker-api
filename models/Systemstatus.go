package models

import (
	"crypto/x509"
	"net/url"
)

type Systemstatus struct {
	Id            int
	CallStatus    string
	CallUrl       url.URL
	HttpMethod    string
	Cert          []x509.Certificate
	CertStatus    string
	CallBody      string
	Message       string
	ResponseMatch string
	AlertBody     string
	AlertUrl      url.URL
	AlertEmail    string
	Status        string
}
