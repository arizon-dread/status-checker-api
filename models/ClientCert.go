package models 

import(
	"crypto/x509"
	"crypto/pkcs12"

)

type ClientCert struct {

	ID int
	Name		string
	p12			[]byte
	PublicKey	*x509.Certificate 
	PrivateKey  *rsa.PrivateKey //??
	Password	string
}