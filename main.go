package main

import (
	"fmt"
	"strings"

	"github.com/arizon-dread/status-checker-api/api"
	"github.com/arizon-dread/status-checker-api/config"
	"github.com/arizon-dread/status-checker-api/datalayer"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	readConfig()

	err := datalayer.PerformMigrations()
	if err != nil {
		fmt.Printf("Migrations failed: %v\n", err)
		panic("migrations failed")
	}
	router := gin.Default()
	router.GET("/healthz", api.Health)
	router.GET("/systemstatus", api.Systemstatuses)
	router.GET("/systemstatus/:id", api.Systemstatus)
	router.POST("/systemstatus", api.SaveSystemStatus)
	router.DELETE("/systemstatus/:id", api.DeleteSystemStatus)
	router.POST("/clientcert", api.UploadP12)
	router.GET("/clientcerts", api.GetCertList)
	router.DELETE("/clientcert/:id", api.DeleteClientCert)

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
