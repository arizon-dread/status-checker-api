package dataLayer

import (
	"fmt"

	"github.com/arizon-dread/status-checker-api/config"
	"github.com/arizon-dread/status-checker-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetAllSystemStatuses() ([]models.Systemstatus, error) {
	var err error = nil

	var systemstatuses []models.Systemstatus
	db, err := getDbConn()

	if err == nil {
		results := db.Find(&systemstatuses)

		if results.Error != nil {
			err = results.Error
		}
	}

	return systemstatuses, err
}
func GetSystemStatus(id int) (models.Systemstatus, error) {
	var systemstatus models.Systemstatus
	db, err := getDbConn()
	if err == nil {
		result := db.First(&systemstatus, "id = ?", id)

		if result.Error != nil {
			err = result.Error
		}
	}
	return systemstatus, err
}
func SaveSystemStatus(system models.Systemstatus) error {
	var err error = nil
	//persist data
	db, err := getDbConn()
	if err == nil {
		db.Save(system)
	}

	return err
}

func getDbConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Stockholm",
		config.Cfg.PgHost, config.Cfg.PgUser, config.Cfg.PgPassword, config.Cfg.PgDatabase, config.Cfg.PgPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to get db connection: %v\n", err)
	}
	return db, err

}
