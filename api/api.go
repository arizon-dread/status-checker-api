package api

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
		,c.AbortWithStatus(http.StatusBadRequest)
	}

	createdSys, err := businesslayer.SaveSystemStatus(system)

	if err == nil {
		c.JSON(http.StatusCreated, createdSys)
	}

}

func DeleteSystemStatus(c *gin.Context) {

	//TODO: Implement Delete.
	c.AbortWithStatus(http.NotImplemented)
}

func UploadP12CertAndPass(c *gin.Context) {
	var clientCert models.ClientCert

	err := c.BindJSON(&clientCert)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	clientCert.privkey, clientCert.pubkey, err := pkcs12.Decode
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}
}