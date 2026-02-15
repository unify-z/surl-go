package admin_dto

type AdminBanUserReq struct {
	UserID uint `json:"user_id"`
}

type AdminUnbanUserReq struct {
	UserID uint `json:"user_id"`
}

type AdminDeleteUserReq struct {
	ID uint `json:"id"`
}
