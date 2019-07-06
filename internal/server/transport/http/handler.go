package http

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func autocompleteHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "RequestURI is %q", ctx.RequestURI())
}
