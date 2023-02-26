package datalayer

import (
	"errors"
	"fmt"

	"github.com/arizon-dread/status-checker-api/config"
	"github.com/arizon-dread/status-checker-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PerformMigrations() error {
	var err error = nil
	db, err := getDbConn()

	if err == nil {
		systemErr := db.AutoMigrate(models.Systemstatus{})
		certErr := db.AutoMigrate(models.ClientCert{})
		if systemErr != nil || certErr != nil {
			err = errors.Join(fmt.Errorf("%w %w", systemErr, certErr))
		}
	}

	return err
}

func getDbConn() (*gorm.DB, error) {
	cfg := config.GetInstance()
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Stockholm",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to get db connection: %v\n", err)
	}
	return db, err

}
