package cache

import "github.com/mediocregopher/radix/v3"

type Service struct {
	client radix.Client
}

func NewService(client radix.Client) *Service {
	return &Service{client: client}
}
