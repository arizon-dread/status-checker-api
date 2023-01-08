package config

import (
	"sync"
)

var cfg *Config

var lock = &sync.Mutex{}

// This is the right way to create a singleton, source: https://refactoring.guru/design-patterns/singleton/go/example
func GetInstance() *Config {
	if cfg == nil {
		lock.Lock()
		defer lock.Unlock()
		if cfg == nil {
			cfg = &Config{}
		}
	}
	return cfg
}

type Config struct {
	Postgres Postgres
	//General  General
}

type Postgres struct {
	PgHost     string
	PgPort     string
	PgUser     string
	PgDatabase string
	PgPassword string
}

// type General struct {
// 	CertWarningDays string
// }
