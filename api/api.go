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
	err := c.ShouldBind(&cc)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	businesslayer.SaveCertificate(cc)

}
