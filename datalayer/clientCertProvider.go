package datalayer

import (
	"fmt"

	"github.com/arizon-dread/status-checker-api/models"
	"gorm.io/gorm"
)

func SaveClientCert(c *models.ClientCert) (int, error) {
	var err error = nil
	db, err := getDbConn()

	if err == nil {
		tx := db.Save(&c)
		err = tx.Error
	}

	return c.ID, err
}
func GetClientCert(id int) (models.ClientCert, error) {
	var err error = nil
	var cert models.ClientCert
	db, err := getDbConn()

	if err == nil {
		tx := db.First(&cert, "id = ?", id)

		if tx.Error != nil {
			err = tx.Error
		}
	}
	return cert, err
}
func CertExists(name string, id *int) (bool, error) {
	var err error = nil
	var c []models.ClientCert
	db, err := getDbConn()
	if err == nil {
		var tx *gorm.DB
		if id == nil {
			if name == "" {
				return false, fmt.Errorf("neither id or name was supplied, cannot determine if cert exists. This should never happen")
			}
			tx = db.Find(&c, "name = ?", name)

		} else {
			tx = db.Find(&c, "id = ?", id)
		}
		if tx.Error != nil {
			err = tx.Error
		}

	}
	return len(c) > 0, err
}

func GetCertList() ([]models.ClientCert, error) {
	var err error = nil
	var certs []models.ClientCert
	db, err := getDbConn()
	if err == nil {
		tx := db.Find(&certs)
		if tx.Error != nil {
			err = tx.Error
		}
	}
	return certs, err
}

func DeleteClientCert(id int) error {
	var err error = nil
	var ss models.Systemstatus
	var noSS int64
	db, err := getDbConn()
	if err == nil {
		tx := db.Find(&ss, "client_cert_id = ?", id).Count(&noSS)
		var i int = 0
		if noSS == int64(i) {
			tx = db.Delete(&models.ClientCert{}, id)
		} else {
			err = fmt.Errorf("could not delete cert since it has relations in systemstatus objects")
		}
		if tx.Error != nil {
			err = tx.Error
		}

	}
	return err
}
