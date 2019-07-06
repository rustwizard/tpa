package http

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/valyala/fasthttp"
)

type Handler struct {
	log zerolog.Logger
}

func NewHandler(log zerolog.Logger) *Handler {
	return &Handler{log: log}
}

func (h *Handler) autocomplete(ctx *fasthttp.RequestCtx) {
	h.log.Info().Str("request", "autocomplete").Uint64("request_id", ctx.ID()).Msg("")
	fmt.Fprintf(ctx, "RequestURI is %q", ctx.RequestURI())
}
