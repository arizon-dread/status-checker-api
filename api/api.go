package api

import (
	"crypto/rsa"
	"fmt"
	"strconv"

	"net/http"

	"github.com/arizon-dread/status-checker-api/businesslayer"
	"github.com/arizon-dread/status-checker-api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/pkcs12"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func Systemstatuses(c *gin.Context) {

	status, err := businesslayer.GetSystemStatuses()

	if err == nil {
		c.JSON(http.StatusOK, status)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func Systemstatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		status, err := businesslayer.GetSystemStatus(id)
		if err == nil {
			c.JSON(http.StatusOK, status)
		} else {
			fmt.Printf("err: %v", err)
			if err.Error() == "NotFound" {
				c.AbortWithStatus(http.StatusNotFound)
			}
			c.AbortWithStatus(http.StatusInternalServerError)
		}

	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func SaveSystemStatus(c *gin.Context) {
	var system models.Systemstatus

	err := c.BindJSON(&system)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	createdSys, err := businesslayer.SaveSystemStatus(system)

	if err == nil {
		c.JSON(http.StatusCreated, createdSys)
	}

}

func DeleteSystemStatus(c *gin.Context) {

	//TODO: Implement Delete.
	c.AbortWithStatus(http.StatusNotImplemented)
}

func UploadP12CertAndPass(c *gin.Context) {
	var clientCert models.ClientCert

	err := c.BindJSON(&clientCert)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	privateKey, publicKey, err := pkcs12.Decode(clientCert.P12, clientCert.Password)
	clientCert.PrivateKey = privateKey.(*rsa.PrivateKey)
	clientCert.PublicKey = publicKey
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}
}
