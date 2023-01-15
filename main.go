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
	router.GET("/healthz", api.Health)
	router.GET("/systemstatus", api.Systemstatuses)
	router.GET("/systemstatus/:id", api.Systemstatus)
	router.POST("/systemstatus", api.SaveSystemStatus)

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



