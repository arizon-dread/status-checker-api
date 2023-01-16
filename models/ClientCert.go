package models

import (
	"crypto/rsa"
	"crypto/x509"
)

type ClientCert struct {
	ID         int
	Name       string
	P12        []byte
	PublicKey  *x509.Certificate
	PrivateKey *rsa.PrivateKey //??
	Password   string
}
