package businesslayer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
)

func getCertFromUrl(u url.URL) []*x509.Certificate {
	conn, err := tls.Dial("tcp", u.String(), &tls.Config{})

	if err != nil {
		fmt.Printf("error getting cert from remote host, %v\n", err)
	}

	defer conn.Close()

	return conn.ConnectionState().PeerCertificates

}
