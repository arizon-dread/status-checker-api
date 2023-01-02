package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arizon-dread/status-checker-api/businesslayer"
	"github.com/arizon-dread/status-checker-api/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	err := viper.Unmarshal(&config.Cfg)
	if err != nil {
		fmt.Printf("error when reading config, %v\n", err)
		panic("Failed to read config")
	}
	// config.Cfg.Postgres.PgHost = viper.Get("postgres.pgHost")
	// config.Cfg.Postgres.PgPort = viper.Get("postgres.pgPort")
	// config.Cfg.Postgres.PgDatabase = viper.Get("postgres.pgDatabase")
	// config.Cfg.Postgres.PgUser = viper.Get("postgres.pgUser")
	// config.Cfg.Postgres.PgPassword = viper.Get("postgres.pgPassword")

	router := gin.Default()
	router.GET("/healthz", health)
	router.GET("/systemstatus", systemstatuses)
	router.GET("/systemstatus/:id", systemstatus)
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
