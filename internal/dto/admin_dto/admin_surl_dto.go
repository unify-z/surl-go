package admin_dto

type AdminDeleteShortURLReq struct {
	ID uint
}

type AdminListShortURLsReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type AdminEditShortURLReq struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}
