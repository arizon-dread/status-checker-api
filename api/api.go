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

// since execution will run throughout each func, equal comparison is the most logic operator when checking for errors.
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

	if err == nil {
		createdSys, err := businesslayer.SaveSystemStatus(system)
		if err == nil {
			c.JSON(http.StatusCreated, createdSys)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func DeleteSystemStatus(c *gin.Context) {

	//TODO: Implement Delete.
	c.AbortWithStatus(http.StatusNotImplemented)
}

func GetCertList(c *gin.Context) {
	certlist, err := businesslayer.GetCertList()
	if err == nil {
		c.JSON(http.StatusOK, certlist)
	} else {
		fmt.Printf("error getting certlist, %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func UploadP12(c *gin.Context) {
	var cc models.CertUploadForm
	bindErr := c.ShouldBind(&cc)
	if bindErr == nil {
		valid := businesslayer.VerifyCertificate(cc)
		if valid {
			_, saveErr := businesslayer.SaveCertificate(cc)
			if saveErr == nil {
				if bindErr == nil && valid && saveErr == nil {
					c.JSON(http.StatusCreated, cc.Name)
				}
			} else {
				fmt.Printf("error saving certificate, %v\n", saveErr)
				c.AbortWithError(http.StatusInternalServerError, saveErr)
			}
		} else {
			fmt.Printf("cert could not be validated\n")
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("unable to decrypt P12"))
		}
	} else {
		fmt.Printf("could not bind input form to form model\n")
		c.AbortWithError(http.StatusBadRequest, bindErr)
	}
}

func DeleteClientCert(c *gin.Context) {
	var id int = 0
	bindErr := c.ShouldBind(&id)
	if bindErr == nil {
		err := businesslayer.DeleteClientCert(id)
		if err == nil {
			c.JSON(http.StatusOK, id)
		}
	}
}
