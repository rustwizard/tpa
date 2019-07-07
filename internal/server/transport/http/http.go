package http

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rustwizard/tpa/internal/server"
	"github.com/valyala/fasthttp"
)

const MIMEApplicationJSON = "application/json"

type Server struct {
	s      *fasthttp.Server
	bind   string
	reqttl time.Duration
	log    zerolog.Logger
	rh     *Handler
}

func NewServer(log zerolog.Logger, conf *server.Config, h *Handler) *Server {
	srv := &Server{
		s:      &fasthttp.Server{},
		bind:   conf.Bind,
		reqttl: conf.RequestTTL,
		log:    log,
	}

	srv.rh = h

	srv.s.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValue("reqttl", srv.reqttl)

		switch string(ctx.Path()) {
		case "/":
			srv.rh.logRequest(fasthttp.TimeoutHandler(srv.rh.autocomplete, srv.reqttl, "request timeout"))(ctx)
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
