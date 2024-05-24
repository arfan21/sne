package model

type CreateURLShortedRequest struct {
	LongUrl string `json:"long_url" validate:"required,url"`
}

type CreateURLShortedResponse struct {
	ShortUrl string `json:"short_url"`
}
