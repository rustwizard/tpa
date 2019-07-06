package http

type AutocompleteRequest struct {
	Term   string   `json:"term" valid:"required"`
	Locale string   `json:"locale" valid:"required"`
	Types  []string `json:"types" valid:"-"`
}
