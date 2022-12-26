package config

type Config struct {
	PgHost     string `yaml:"pgHost" env: "PGHOST" env-default:"localhost"`
	PgPort     string `yaml:"pgPort" env: "PGPORT" env-default:"5432"`
	PgUser     string `yaml:"pgUser" env: "PGUSER" env-default:"statusUser"`
	PgDatabase string `yaml:"pgDatabase" env:"PGDB" env-default:"status`
	PgPassword string `yaml:"pgPassword" env:"PGPASSWORD"`
}

var Cfg Config
