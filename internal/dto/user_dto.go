package dto

type RegisterReq struct {
	Username        string `json:"username"`
	PasswordMD5     string `json:"password_md5"`
	Email           string `json:"email"`
	EmailVerifyCode string `json:"email_verify_code"`
}

type SendEmailVerifyCodeReq struct {
	Email string `json:"email"`
}

type LoginReq struct {
	Username    string `json:"username"`
	PasswordMD5 string `json:"password_md5"`
}

type EditProfileReq struct {
	Username        string `json:"username,omitempty"`
	PasswordMD5     string `json:"password_md5,omitempty"`
	CodeID          string `json:"code_id,omitempty"`
	EmailVerifyCode string `json:"email_verify_code,omitempty"`
}
