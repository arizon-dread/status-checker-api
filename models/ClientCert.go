package models

import (
	"crypto/rsa"
	"crypto/x509"
)

type ClientCert struct {
	ID   *int   `json:"id"`
	Name string `json:"name"`
	//P12        []byte            `json:"p12"`
	PublicKey  *x509.Certificate `json:"-"`
	PrivateKey *rsa.PrivateKey   `json:"-"` //??
	Password   string            `json:"password"`
}
