package http

import (
	"github.com/rustwizard/tpa/internal/pac"
)

type ResponseEntity struct {
	Slug     string `json:"slug" valid:"required"`
	Subtitle string `json:"subtitle" valid:"required"`
	Title    string `json:"title" valid:"required"`
}

type AutocompleteResponse struct {
	Collection []*ResponseEntity
}

func makeAutocompleteResponse(pr *pac.Response) *AutocompleteResponse {
	acr := &AutocompleteResponse{}
	acr.Collection = make([]*ResponseEntity, 0)
	for _, v := range pr.Collection {
		switch v.Type {
		case "city":
			acr.Collection = append(acr.Collection, &ResponseEntity{
				Slug:     v.Code,
				Subtitle: v.CountryName,
				Title:    v.Name,
			})
		case "airport":
			acr.Collection = append(acr.Collection, &ResponseEntity{
				Slug:     v.Code,
				Subtitle: v.CityName,
				Title:    v.Name,
			})
		case "country":
			acr.Collection = append(acr.Collection, &ResponseEntity{
				Slug:     v.Code,
				Subtitle: v.Name,
				Title:    v.Name,
			})
		}
	}
	return acr
}
