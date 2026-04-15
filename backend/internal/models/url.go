package models

type URL struct {
	LongURL   string `json:"long_url"`
	ShortCode string `json:"short_code"`
	Clicks    int    `json:"clicks"`
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

type URLListResponse struct {
	URLs []URL `json:"urls"`
}
