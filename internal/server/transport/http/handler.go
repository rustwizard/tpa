package http

import (
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
	h.response(ctx, fasthttp.StatusOK, ctx.Request.Host())
}

func (h *Handler) response(ctx *fasthttp.RequestCtx, status int, res []byte) {
	ctx.SetStatusCode(status)
	ctx.SetContentType(MIMEApplicationJSON)
	ctx.SetBody(res)
}
