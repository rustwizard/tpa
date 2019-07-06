package server

import (
	"time"
)

type Config struct {
	Bind       string        `mapstructure:"SERVER_BIND"`
	RequestTTL time.Duration `mapstructure:"SERVER_REQUEST_TTL"`
	Transport  string        `mapstructure:"SERVER_TRANSPORT"`
}

type Runner interface {
	Run() error
}
