package config

var Cfg Config

type Config struct {
	Postgres Postgres
	General  General
}

type Postgres struct {
	PgHost     string
	PgPort     string
	PgUser     string
	PgDatabase string
	PgPassword string
}

type General struct {
	CertWarningDays string
}
