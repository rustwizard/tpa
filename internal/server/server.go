package server

import (
	"time"
)

type Config struct {
	Bind          string        `mapstructure:"SERVER_BIND"`
	RequestTTL    time.Duration `mapstructure:"SERVER_REQUEST_TTL"`
	Transport     string        `mapstructure:"SERVER_TRANSPORT"`
	RemoteAPIPath string        `mapstructure:"REMOTE_API_PATH"`
}

type Runner interface {
	Run() error
}
