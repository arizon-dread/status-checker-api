package datalayer

import (
	"github.com/arizon-dread/status-checker-api/models"
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
