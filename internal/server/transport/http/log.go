package http

import (
	"time"

	"github.com/valyala/fasthttp"
)

func (hd *Handler) logRequest(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		h(ctx)
		addr := string(ctx.Request.Header.Peek("X-Real-IP"))
		if len(addr) == 0 {
			addr = string(ctx.Request.Header.Peek("X-Forwarded-For"))
			if len(addr) == 0 {
				addr = ctx.RemoteIP().String()
			}
		}
		hd.log.Info().Bytes("method", ctx.Method()).Uint64("request_id", ctx.ID()).
			Str("addr", addr).Dur("duration", time.Since(start)).Msgf("")
	}
}
