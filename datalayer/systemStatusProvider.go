package datalayer

import (
	"fmt"

	"github.com/arizon-dread/status-checker-api/models"
)

type system interface {
	UpdateSystem(models.Systemstatus) error
}

func GetAllSystemStatuses() ([]models.Systemstatus, error) {
	var err error = nil

	var systemstatuses []models.Systemstatus
	db, err := getDbConn()

	if err == nil {
		tx := db.Find(&systemstatuses)

		if tx.Error != nil {
			err = tx.Error
		}
	}

	return systemstatuses, err
}
func GetSystemStatus(id int) (models.Systemstatus, error) {
	var systemstatus models.Systemstatus
	db, err := getDbConn()
	if err == nil {
		tx := db.First(&systemstatus, "id = ?", id)

		if tx.Error != nil {
			err = tx.Error
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
		tx := db.Save(&system)
		if tx.Error != nil {
			err = tx.Error
		}
		tx = db.First(&createdSys, &system)
		if tx.Error != nil {
			err = tx.Error
		}
	}

	return createdSys, err
}

func DeleteSystem(id int) (int, error) {
	var err error = nil

	db, err := getDbConn()
	if err == nil {
		tx := db.Delete(&models.Systemstatus{}, id)
		if tx.Error == nil {
			if tx.RowsAffected == 0 {
				err = fmt.Errorf("rowsAffected was 0, probably because the id did not exist")
			}
		}
		return int(tx.RowsAffected), err
	}
	return 0, err
}

func UpdateSystem(system *models.Systemstatus) error {
	var err error = nil
	db, err := getDbConn()
	if err == nil {
		tx := db.Save(system)
		if tx.Error != nil {
			err = tx.Error
		}
	}
	return err
}
