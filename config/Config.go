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
	AlertSMTP
	//General  General
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Database string
	Password string
}

type AlertSMTP struct {
	Server   string
	Port     int
	User     string
	Password string
}

// type General struct {
// 	CertWarningDays string
// }
