package pkgutil

type HTTPResponse struct {
	Code   int `json:"code" example:"200"`
	Data   any `json:"data,omitempty" `
	Errors any `json:"errors,omitempty" `
}

type PaginationResponse struct {
	TotalData int `json:"total_data" example:"1"`
	TotalPage int `json:"total_page" example:"1"`
	Page      int `json:"page" example:"1"`
	Limit     int `json:"limit" example:"10"`
	Data      any `json:"data,omitempty" `
}

type ErrValidationResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
