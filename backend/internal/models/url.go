package models

import "time"


type URL struct {
	LongURL string `json:"long_url"`
	ShortCode string `json:"short_code"`
	Clicks int `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	Code string 	`json:"code"`
}

type URLListResponse struct {
	URLs []URL `json:"urls"`
}