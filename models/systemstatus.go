package models

import (
	"crypto/x509"
	"net/url"
)

type Systemstatus struct {
	id         int
	callStatus string
	url        url.URL
	cert       []x509.Certificate
	certStatus string
	message    string
}
