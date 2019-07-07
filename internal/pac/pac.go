package pac

import (
	"context"

	"github.com/rs/zerolog"
)

// Service implements main logic with interaction the remote API
type Service struct {
	log           zerolog.Logger
	RemoteAPIPath string
}

type Request struct {
	URI string
}

type Case struct {
	Da string `json:"da" valid:"-"`
	Tv string `json:"tv" valid:"-"`
	Vi string `json:"vi" valid:"-"`
	Pr string `json:"pr" valid:"-"`
	Ro string `json:"ro" valid:"-"`
}

type Coordinate struct {
	Lon float64 `json:"lon" valid:"-"`
	Lat float64 `json:"lat" valid:"-"`
}

type Response struct {
	Code            string      `json:"code" valid:"-"`
	MainAirportName string      `json:"main_airport_name" valid:"-"`
	CountryCases    string      `json:"country_cases" valid:"-"`
	IndexStrings    []string    `json:"index_strings" valid:"-"`
	Weight          int         `json:"weight" valid:"-"`
	Cases           *Case       `json:"cases" valid:"-"`
	CountryName     string      `json:"country_name" valid:"-"`
	Type            string      `json:"type" valid:"-"`
	Coordinates     *Coordinate `json:"coordinates" valid:"-"`
	Name            string      `json:"name" valid:"-"`
	StateCode       int         `json:"state_code" valid:"-"`
}

func NewService(log zerolog.Logger, APIPath string) *Service {
	return &Service{log: log, RemoteAPIPath: APIPath}
}

func (s *Service) Do(ctx context.Context, req *Request) (*Response, error) {
	return nil, nil
}
