package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arizon-dread/status-checker-api/businesslayer"
	"github.com/arizon-dread/status-checker-api/config"
	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/arizon-dread/status-checker-api/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	readConfig()

	err := datalayer.PerformMigrations()
	if err != nil {
		panic("migrations failed")
	}
	router := gin.Default()
	router.GET("/healthz", health)
	router.GET("/systemstatus", systemstatuses)
	router.GET("/systemstatus/:id", systemstatus)
	router.POST("/systemstatus", saveSystemStatus)

	router.Run(":8080")
}
func readConfig() {

	cfg := config.GetInstance()

	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		breakOnNoConfig(err)
	}
	//all keys that are read from config will get overwritten by their env equivalents, as long as they exist in config (empty or not)
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		breakOnNoConfig(err)
	}

}

func breakOnNoConfig(err error) {
	fmt.Printf("error when reading config, %v\n", err)
	panic("Failed to read config")
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func systemstatuses(c *gin.Context) {

	status, err := businesslayer.GetSystemStatuses()

	if err == nil {
		c.JSON(http.StatusOK, status)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func systemstatus(c *gin.Context) {
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

func saveSystemStatus(c *gin.Context) {
	var system models.Systemstatus

	c.BindJSON(&system)

	createdSys, err := businesslayer.SaveSystemStatus(system)

	if err == nil {
		c.JSON(http.StatusCreated, createdSys)
	}

}
