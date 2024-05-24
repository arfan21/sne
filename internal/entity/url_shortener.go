package entity

import "time"

type URLShortener struct {
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt int64     `json:"expired_at"`
}

func (URLShortener) TableName() string {
	return "url_shortener"
}
