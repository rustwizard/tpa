package pac

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vmihailenco/msgpack"

	"github.com/valyala/fasthttp"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

// Service implements main logic with interaction the remote API
type Service struct {
	log           zerolog.Logger
	RemoteAPIPath string
	cachesvc      *Cache
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
	Collection []*RespEntity
}

func (c *Case) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(c.Da, c.Pr, c.Ro, c.Tv, c.Vi)
}

func (c *Case) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&c.Da, &c.Pr, &c.Ro, &c.Tv, &c.Vi)
}

func (c *Coordinate) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(c.Lat, c.Lon)
}

func (c *Coordinate) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&c.Lat, &c.Lon)
}

func (r *Response) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(r.Collection)
}

func (r *Response) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&r.Collection)
}

func (re *RespEntity) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(re.Code, re.MainAirportName, re.CountryCases, re.IndexStrings, re.Weight, re.Cases,
		re.CountryName, re.Type, re.Coordinates, re.Name, re.StateCode, re.CityName)
}

func (re *RespEntity) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&re.Code, &re.MainAirportName, &re.CountryCases, &re.IndexStrings, &re.Weight, &re.Cases,
		&re.CountryName, &re.Type, &re.Coordinates, &re.Name, &re.StateCode, &re.CityName)
}

func NewService(log zerolog.Logger, APIPath string, cachesvc *Cache) *Service {
	return &Service{
		log:           log,
		RemoteAPIPath: APIPath,
		cachesvc:      cachesvc,
	}
}

func (s *Service) Do(ctx context.Context, req *Request) (*Response, error) {
	var resp Response
	resp.Collection = make([]*RespEntity, 0)
	if len(req.URI) == 0 {
		return &resp, errors.New("pac: request has empty URI")
	}

	ttl := ctx.Value("reqttl").(time.Duration)

	url := s.RemoteAPIPath + req.URI
	s.log.Debug().Str("url", url).Msg("")

	cacheresp, err := s.cachesvc.GetResponse(req.URI)
	if err != nil && err != ErrNotFound {
		s.log.Debug().Err(err).Msg("get from cache")
		return &resp, err
	}

	if cacheresp != nil && len(cacheresp.Collection) > 0 {
		s.log.Debug().Msg("got response from cache!")
		return cacheresp, nil
	}

	r := fasthttp.AcquireRequest()
	rp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(r)
	defer fasthttp.ReleaseResponse(rp)

	r.SetRequestURI(url)

	s.log.Debug().Str("request header", r.Header.String()).
		Float64("reqttl", ttl.Seconds()).Msg("")

	client := &fasthttp.Client{}
	if err := client.DoTimeout(r, rp, ttl); err != nil {
		return &resp, err
	}

	s.log.Debug().RawJSON("response body", rp.Body()).Msg("")

	if err := json.Unmarshal(rp.Body(), &resp.Collection); err != nil {
		return &resp, err
	}

	if err := s.cachesvc.SetResponse(req.URI, &resp); err != nil {
		s.log.Debug().Err(err).Msg("set to cache")
		return &resp, err
	}

	return &resp, nil
}
