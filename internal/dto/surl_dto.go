package dto

type CreateShortURLReq struct {
	OriginalURL string `json:"original_url"`
}

type DeleteShortURLReq struct {
	ShortCode string `json:"short_code"`
}

type GetShortURLReq struct {
	ShortCode string `json:"short_code"`
}

type EditShortURLReq struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

type ListShortURLsReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
