package pac

import "github.com/rs/zerolog"

// Service implements main logic with interaction the remote API
type Service struct {
	log           zerolog.Logger
	RemoteAPIPath string
}
