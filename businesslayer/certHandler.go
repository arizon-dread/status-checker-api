package businesslayer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"

	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/arizon-dread/status-checker-api/models"
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

func SaveCertificate(c models.ClientCert) (*int, error) {

	return datalayer.SaveClientCert(&c)
}

func GetCertificate(id int) (models.ClientCert, error) {
	clientCert, err := datalayer.GetClientCert(id)
	if err != nil {
		fmt.Printf("Could not get clientCert from database, %v", err)
	}

	return clientCert, err
}
