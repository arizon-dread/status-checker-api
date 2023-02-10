package api

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
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

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusCreated, createdSys)

}

func DeleteSystemStatus(c *gin.Context) {

	//TODO: Implement Delete.
	c.AbortWithStatus(http.StatusNotImplemented)
}

/*
 * upload cert pass and cert-name and get a Location header for where to upload the file with PUT.
 */

func UploadP12Pass(c *gin.Context) {
	var clientCert models.ClientCert

	err := c.BindJSON(&clientCert)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	id, saveErr := businesslayer.SaveCertificate(clientCert)
	if saveErr != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	var url = fmt.Sprintf("%v://%v%v/%v", c.Request.URL.Scheme, c.Request.URL.Host, c.Request.URL.Path, id)
	c.Header("Location", url)
	c.JSON(http.StatusAccepted, id)
}

/*
 * upload file with put, add header X-FILENAME
 * the file should be a form file.
 */
func UploadCertFile(c *gin.Context) {
	file, err := c.FormFile(c.GetHeader("X-FILENAME"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	fileContent, err := file.Open()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	clientCert, err := businesslayer.GetCertificate(id)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
	}

	privateKey, publicKey, err := pkcs12.Decode(fileBytes, clientCert.Password)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
	} else {
		clientCert.PrivateKey = privateKey.(*rsa.PrivateKey)
		clientCert.PublicKey = publicKey
		clientCert.P12 = fileBytes

		_, err := businesslayer.SaveCertificate(clientCert)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusCreated, nil)
		}
	}
}
