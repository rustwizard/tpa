package pac

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rustwizard/tpa/internal/cache"

	"github.com/valyala/fasthttp"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

// Service implements main logic with interaction the remote API
type Service struct {
	log           zerolog.Logger
	RemoteAPIPath string
	cachesvc      *cache.Service
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
type RespEntity struct {
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
	StateCode       string      `json:"state_code" valid:"-"`
	CityName        string      `json:"city_name" valid:"-"`
}

type Response struct {
	Collection []RespEntity
}

func NewService(log zerolog.Logger, APIPath string, cachesvc *cache.Service) *Service {
	return &Service{log: log, RemoteAPIPath: APIPath, cachesvc: cachesvc}
}

func (s *Service) Do(ctx context.Context, req *Request) (*Response, error) {
	var resp Response
	resp.Collection = make([]RespEntity, 0)
	if len(req.URI) == 0 {
		return &resp, errors.New("pac: request has empty URI")
	}

	ttl := ctx.Value("reqttl").(time.Duration)

	url := s.RemoteAPIPath + req.URI
	s.log.Debug().Str("url", url).Msg("")

	r := fasthttp.AcquireRequest()
	rp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(r)
	defer fasthttp.ReleaseResponse(rp)

	r.SetRequestURI(url)

	s.log.Debug().Str("request header", r.Header.String()).
		Dur("reqttl", ttl).Msg("")

	client := &fasthttp.Client{}
	if err := client.DoTimeout(r, rp, ttl); err != nil {
		return &resp, err
	}

	s.log.Debug().RawJSON("response body", rp.Body()).Msg("")

	if err := json.Unmarshal(rp.Body(), &resp.Collection); err != nil {
		return &resp, err
	}

	return &resp, nil
}
