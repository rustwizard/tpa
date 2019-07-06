package http

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"

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
	h.log.Info().Str("request", "autocomplete").Uint64("request_id", ctx.ID()).
		Bytes("body", ctx.PostBody()).Msg("")

	var ac AutocompleteRequest

	if err := json.Unmarshal(ctx.PostBody(), &ac); err != nil {
		h.log.Err(err).Uint64("request_id", ctx.ID()).Msg("parse json")
		h.response(ctx, fasthttp.StatusBadRequest, []byte{})
		return
	}

	if _, err := govalidator.ValidateStruct(ac); err != nil {
		h.log.Err(err).Uint64("request_id", ctx.ID()).Msg("validate request")
		h.response(ctx, fasthttp.StatusBadRequest, []byte{})
		return
	}

	h.response(ctx, fasthttp.StatusOK, []byte{})
}

func (h *Handler) response(ctx *fasthttp.RequestCtx, status int, res []byte) {
	ctx.SetStatusCode(status)
	ctx.SetContentType(MIMEApplicationJSON)
	ctx.SetBody(res)
}
