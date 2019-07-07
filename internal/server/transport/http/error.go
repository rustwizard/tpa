package http

var (
	ErrParseJSON = Error{
		Code:    1000,
		Message: "Parse JSON",
	}
	ErrValidateRequest = Error{
		Code:    1001,
		Message: "Validate request",
	}
	ErrProcessRequest = Error{
		Code:    1002,
		Message: "Process request",
	}
)

type Error struct {
	Code    int    `json:"err_code"`
	Message string `json:"message"`
}
