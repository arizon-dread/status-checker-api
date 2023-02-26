package api

import (
	"fmt"
	"strconv"

	"net/http"

	"github.com/arizon-dread/status-checker-api/businesslayer"
	"github.com/arizon-dread/status-checker-api/models"
	"github.com/gin-gonic/gin"
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

func UploadP12(c *gin.Context) {
	var cc models.CertUploadForm
	bindErr := c.ShouldBind(&cc)
	if bindErr != nil {
		fmt.Printf("could not bind input form to form model")
		c.AbortWithError(http.StatusBadRequest, bindErr)
	}
	valid := businesslayer.VerifyCertificate(cc)
	if !valid {
		fmt.Printf("cert could not be validated")
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("unable to decrypt P12"))
	}
	_, err := businesslayer.SaveCertificate(cc)
	if err != nil {
		fmt.Printf("error saving certificate, %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, cc.Name)
}
