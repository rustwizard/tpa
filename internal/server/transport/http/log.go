package http

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/valyala/fasthttp"
)

var l = log.With().Str("pkg", "http").Logger()

func LogRequest(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		l.Info().Bytes("method", ctx.Method()).Bytes("path", ctx.RequestURI()).
			Str("request", string(ctx.Request.Body())).Msg("")
		h(ctx)
		addr := string(ctx.Request.Header.Peek("X-Real-IP"))
		if len(addr) == 0 {
			addr = string(ctx.Request.Header.Peek("X-Forwarded-For"))
			if len(addr) == 0 {
				addr = ctx.RemoteIP().String()
			}
		}
		l.Info().Bytes("method", ctx.Method()).Bytes("path", ctx.RequestURI()).
			Str("addr", addr).Dur("duration", time.Since(start)).Msgf("")
	}
}
