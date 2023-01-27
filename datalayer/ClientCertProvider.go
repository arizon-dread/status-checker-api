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
