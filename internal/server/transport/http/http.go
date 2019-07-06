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
	rh     *Handler
}

func NewServer(log zerolog.Logger, conf *server.Config) *Server {
	srv := &Server{
		s:      &fasthttp.Server{},
		bind:   conf.Bind,
		reqttl: conf.RequestTTL * time.Second,
		log:    log,
	}

	srv.rh = NewHandler(log)
	srv.s.Handler = func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			srv.rh.logRequest(srv.rh.autocomplete)(ctx)
		default:
			ctx.Error("unsupported path", fasthttp.StatusNotFound)
		}
	}

	return srv
}

func (srv *Server) Run() error {
	srv.log.Info().Msgf("http server started at: %s", srv.bind)

	return srv.s.ListenAndServe(srv.bind)
}
