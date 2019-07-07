package http

import (
	"fmt"

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

func marshalPacRequest(a *AutocompleteRequest) *pac.Request {
	return &pac.Request{}
}
