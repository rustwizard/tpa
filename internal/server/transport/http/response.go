package http

import (
	"fmt"

	"github.com/rustwizard/tpa/internal/pac"
)

type AutocompleteResponse struct {
	Slug     string `json:"slug" valid:"required"`
	Subtitle string `json:"subtitle" valid:"required"`
	Title    string `json:"title" valid:"required"`
}

func (a AutocompleteResponse) String() string {
	return fmt.Sprintf("slug=%s, subtitle=%s, title=%q", a.Slug, a.Subtitle, a.Title)
}

func umarshalPacResponse(pr *pac.Response) *AutocompleteResponse {
	return nil
}
