package http

import "fmt"

type AutocompleteRequest struct {
	Term   string   `json:"term" valid:"required"`
	Locale string   `json:"locale" valid:"required"`
	Types  []string `json:"types" valid:"-"`
}

func (a AutocompleteRequest) String() string {
	return fmt.Sprintf("term=%s, locale=%s, types=%q", a.Term, a.Locale, a.Types)
}
