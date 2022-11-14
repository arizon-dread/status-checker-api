package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", health)
	router.GET("/systemstatuses", systemstatuses)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func systemstatuses(c *gin.Context) {
	businesslayer.
}