package businesslayer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/arizon-dread/status-checker-api/config"
	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/arizon-dread/status-checker-api/models"
	"golang.org/x/crypto/pkcs12"
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

func SaveCertificate(cf models.CertUploadForm) (int, error) {

	c, convertErr := formToModel(cf)
	if convertErr != nil {
		fmt.Printf("error converting form data to certmodel, %v\n", convertErr)
		return 0, convertErr
	}
	encPass, err := encryptPassword(c.Password)
	if err != nil {
		fmt.Printf("error encrypting password, %v\n", err)
		return 0, err
	}
	c.Password = encPass

	return datalayer.SaveClientCert(&c)
}

func GetCertList() (map[int]string, error) {
	crts, err := datalayer.GetCertList()
	m := make(map[int]string)
	if err != nil {
		fmt.Printf("Could not get certlist from database, %v", err)
		return m, err
	}
	for _, c := range crts {
		m[c.ID] = c.Name
	}

	return m, err

}
func VerifyCertificate(cert models.CertUploadForm) bool {
	idx := cert.ID != nil
	cx, err := datalayer.CertExists(cert.Name, cert.ID)
	//if id is sent in the form, we want to update the cert, otherwise verify == false.
	if err == nil {
		if cx && !idx {
			return false
		}
		if !cx && idx {
			return false
		}
	}
	cc, err := formToModel(cert)
	if err != nil {
		return false
	}
	_, err = decryptClientCert(cc)
	return err == nil
}

func DeleteClientCert(id int) error {
	err := datalayer.DeleteClientCert(id)

	if err != nil {
		fmt.Printf("error deleting cert, %v\n", err)
	}
	return err
}

func getClientCert(id int) (models.ClientCert, error) {
	clientCert, dlErr := datalayer.GetClientCert(id)
	if dlErr != nil {
		fmt.Printf("Could not get clientCert from database, %v", dlErr)
	}
	pw, decryptErr := decryptPassword(clientCert.Password)
	if decryptErr != nil {
		fmt.Printf("Failed decrypting password, %v", decryptErr)
	}
	err := fmt.Errorf("%w + %w", dlErr, decryptErr)
	clientCert.Password = pw
	return clientCert, err
}

func formToModel(form models.CertUploadForm) (models.ClientCert, error) {
	var cc models.ClientCert
	cc.Name = form.Name
	cc.Password = form.Password
	file, err := form.P12.Open()
	if err != nil {
		fmt.Printf("Failed opening the p12 file, %v\n", err)
		return models.ClientCert{}, err
	}
	p12Container, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed reading file stream, %v\n", err)
		return models.ClientCert{}, err
	}
	cc.P12 = p12Container
	if form.ID != nil {
		cc.ID = *form.ID
	}
	return cc, err
}
func encryptPassword(clearTxt string) (string, error) {
	var err error = nil
	gcmInstance, err := getGcmInstance()
	//Create a nonce which must be unique.
	nonce := make([]byte, gcmInstance.NonceSize())
	// read the byte length into the nonce
	_, _ = io.ReadFull(rand.Reader, nonce)
	cryptoResult := string(gcmInstance.Seal(nonce, nonce, []byte(clearTxt), nil))
	result := base64.RawStdEncoding.EncodeToString([]byte(cryptoResult))

	return result, err

}

func decryptPassword(encPwd string) (string, error) {
	cipher, b64Err := base64.StdEncoding.DecodeString(encPwd)
	if b64Err != nil {
		fmt.Printf("error decoding base64 string")
	}
	gcmInstance, gcmErr := getGcmInstance()
	if gcmErr != nil {
		fmt.Printf("error getting gcmInstance, %v\n", gcmErr)
	}
	nonceSize := gcmInstance.NonceSize()
	//slice the nonce off the text
	nonce, cipherText := cipher[:nonceSize], cipher[nonceSize:]

	clearPwd, decryptErr := gcmInstance.Open(nil, nonce, cipherText, nil)
	if decryptErr != nil {
		fmt.Printf("Could not decrypt password, %v", decryptErr)
	}
	var err error
	if gcmErr != nil || decryptErr != nil || b64Err != nil {
		err = fmt.Errorf("%w + %w + %w", gcmErr, decryptErr, b64Err)
	}

	return string(clearPwd), err

}
func getGcmInstance() (cipher.AEAD, error) {
	cfg := config.GetInstance()
	var err error = nil
	// Create aesBlock with encryption key
	aesBlock, aesErr := aes.NewCipher([]byte(cfg.General.EncryptionKey))

	if aesErr != nil {
		fmt.Printf("Error creating cipher, %v", aesErr)
	}
	// create a gcmInstance to use for encryption
	gcmInstance, gcmErr := cipher.NewGCM(aesBlock)
	if gcmErr != nil {
		fmt.Printf("Error creating GCMInstance, %v\n", gcmErr)
	}
	if aesErr != nil || gcmErr != nil {
		err = fmt.Errorf("%w + %w", aesErr, gcmErr)
	}

	return gcmInstance, err
}

func decryptClientCert(cc models.ClientCert) (tls.Certificate, error) {

	blocks, err := pkcs12.ToPEM(cc.P12, cc.Password)
	if err != nil {
		fmt.Printf("Could not unpack P12, %v\n", err)
		return tls.Certificate{}, err
	}
	var pubKey *pem.Block
	var privKey *pem.Block
	for _, b := range blocks {
		if b.Type == "CERTIFICATE" {
			pubKey = b
		}
		if b.Type == "PRIVATE KEY" {
			privKey = b
		}
	}
	cert, err := tls.X509KeyPair(pem.EncodeToMemory(pubKey), pem.EncodeToMemory(privKey))
	if err != nil {
		fmt.Printf("Could not load x509KeyPair, %v\n", err)
	}
	return cert, err
}
