package dataLayer

import "github.com/arizon-dread/status-checker-api/models"

func getAllSystemStatuses() (models.Systemstatus, error) {
	connStr := "systatus:muchs4f3!@localhost:5432:/systemstatus"
	var statuses []models.Systemstatus
}
