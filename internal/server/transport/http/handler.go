package http

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rustwizard/tpa/internal/pac"

	"github.com/asaskevich/govalidator"

	"github.com/rs/zerolog"

	"github.com/valyala/fasthttp"
)

type Handler struct {
	log zerolog.Logger
	pac *pac.Service
}

func NewHandler(log zerolog.Logger, pac *pac.Service) *Handler {
	return &Handler{log: log, pac: pac}
}

func (h *Handler) autocomplete(ctx *fasthttp.RequestCtx) {
	h.log.Info().Str("request", "autocomplete").Uint64("request_id", ctx.ID()).
		Bytes("body", ctx.PostBody()).Msg("")

	var areq AutocompleteRequest

	if err := json.Unmarshal(ctx.PostBody(), &areq); err != nil {
		h.log.Err(err).Uint64("request_id", ctx.ID()).Msg("parse json")
		h.responseError(ctx, fasthttp.StatusBadRequest, ErrParseJSON)
		return
	}

	if _, err := govalidator.ValidateStruct(areq); err != nil {
		h.log.Err(err).Uint64("request_id", ctx.ID()).Msg("validate request")
		h.responseError(ctx, fasthttp.StatusBadRequest, ErrValidateRequest)
		return
	}

	h.log.Debug().Str("request", areq.String()).Msg("")

	ttl := ctx.Value("reqttl")
	rctx, cancel := context.WithTimeout(ctx, ttl.(time.Duration))
	defer cancel()

	pr, err := h.pac.Do(rctx, makePacRequest(&areq))
	if err != nil {
		h.log.Err(err).Uint64("request_id", ctx.ID()).Msg("process request")
		h.responseError(ctx, fasthttp.StatusBadGateway, ErrProcessRequest)
		return

	}

	aresp := makeAutocompleteResponse(pr)

	b, err := json.Marshal(aresp.Collection)
	if err != nil {
		h.log.Err(err).Uint64("request_id", ctx.ID()).Msg("parse json")
		h.responseError(ctx, fasthttp.StatusBadGateway, ErrParseJSON)
		return
	}

	h.responseOK(ctx, b)
}

func (h *Handler) responseOK(ctx *fasthttp.RequestCtx, res []byte) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Success(MIMEApplicationJSON, res)
}

func (h *Handler) responseError(ctx *fasthttp.RequestCtx, status int, e Error) {
	ctx.SetStatusCode(status)
	ctx.SetContentType(MIMEApplicationJSON)
	res, err := json.Marshal(e)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}
	ctx.SetBody(res)
}
