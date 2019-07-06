package http

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rustwizard/tpa/internal/server"
	"github.com/valyala/fasthttp"
)

type Server struct {
	s      *fasthttp.Server
	bind   string
	reqttl time.Duration
	log    zerolog.Logger
}

func NewServer(log zerolog.Logger, conf *server.Config) *Server {
	return &Server{
		s: &fasthttp.Server{
			Handler: func(ctx *fasthttp.RequestCtx) {
				switch string(ctx.Path()) {
				case "/":
					LogRequest(autocompleteHandler)(ctx)
				default:
					ctx.Error("unsupported path", fasthttp.StatusNotFound)
				}
			},
		},
		bind:   conf.Bind,
		reqttl: conf.RequestTTL * time.Second,
		log:    log,
	}
}

func (srv *Server) Run() error {
	srv.log.Info().Msgf("http server started at: %s", srv.bind)

	return srv.s.ListenAndServe(srv.bind)
}
