package http

import (
	"fmt"
	"net/url"

	"github.com/rustwizard/tpa/internal/pac"
)

type AutocompleteRequest struct {
	Term   string   `json:"term" valid:"required"`
	Locale string   `json:"locale" valid:"required"`
	Types  []string `json:"types" valid:"-"`
}

func (a AutocompleteRequest) String() string {
	return fmt.Sprintf("term=%s, locale=%s, types=%q", a.Term, a.Locale, a.Types)
}

func makePacRequest(a *AutocompleteRequest) *pac.Request {
	pacreq := &pac.Request{}
	pacreq.URI = fmt.Sprintf("term=%s&locale=%s", url.QueryEscape(a.Term), a.Locale)
	for _, v := range a.Types {
		pacreq.URI += fmt.Sprintf("&%s=%s", url.QueryEscape("types[]"), v)
	}
	return pacreq
}
