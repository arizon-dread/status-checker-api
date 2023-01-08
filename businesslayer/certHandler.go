package businesslayer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"
)

func getCertFromUrl(u url.URL) []*x509.Certificate {
	url := u.Host
	if !strings.Contains(url, ":") {
		url += ":443"
	}
	conn, err := tls.Dial("tcp", url, &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		fmt.Printf("error getting cert from remote host, %v\n", err)
	}

	defer conn.Close()

	return conn.ConnectionState().PeerCertificates

}
