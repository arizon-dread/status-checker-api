package datalayer

import "github.com/arizon-dread/status-checker-api/models"

func SaveClientCert(c *models.ClientCert) error {
	var err error = nil
	db, err := getDbConn()

	if err == nil {
		tx := db.Save(&c)
		err = tx.Error
	}

	return err
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
