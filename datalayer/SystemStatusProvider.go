package datalayer

import (
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
		err = db.AutoMigrate(models.Systemstatus{})
	}

	return err
}

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
func SaveSystemStatus(system *models.Systemstatus) (models.Systemstatus, error) {
	var err error = nil
	//persist data
	db, err := getDbConn()
	var createdSys models.Systemstatus
	if err == nil {
		db.Save(&system)
		db.First(&createdSys, &system)
	}

	return createdSys, err
}

func getDbConn() (*gorm.DB, error) {
	cfg := config.GetInstance()
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Stockholm",
		cfg.Postgres.PgHost, cfg.Postgres.PgUser, "supers3cret", cfg.Postgres.PgDatabase, cfg.Postgres.PgPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to get db connection: %v\n", err)
	}
	return db, err

}
