package main

import (
	"net/http"

	"github.com/arizon-dread/status-checker-api/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInconfig()

	config.Cfg.PgHost = viper.Get("pgHost")
	config.Cfg.PgPort = viper.Get("pgPort")
	config.Cfg.PgDatabase = viper.Get("pgDatabase")
	config.Cfg.PgUser = viper.Get("pgUser")
	config.Cfg.PgPassword = viper.Get("pgPassword")

	router := gin.Default()
	router.GET("/healthz", health)
	router.GET("/systemstatuses", systemstatuses)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, "Healthy")
}

func systemstatuses(c *gin.Context) {
	//businesslayer.
}
