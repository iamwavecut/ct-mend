// Package config Application config
package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultTimeout = 5 * time.Second
)

type (
	// Key for context value.
	Key struct{}

	// TLS Web-server config.
	TLS struct {
		Addr string `env:"TLS_ADDR"`
	}

	// Storage Database config.
	Storage struct {
		Type string `env:"STORAGE_TYPE"`
		Addr string `env:"STORAGE_ADDR"`
	}

	// Config Application config.
	Config struct {
		TLS             TLS
		Storage         Storage
		AppLogLevel     log.Level     `env:"LOG_LEVEL" envDefault:"trace"`
		GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT" envDefault:"10s"`
	}
)
