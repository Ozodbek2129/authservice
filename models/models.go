package models

type UserInfo struct {
	Id       string
	UserName string
	Password string
	FullName string
	UserType string
}

type RefreshToken struct {
	UserId    string
	Token     string
	ExpiresAt int64
}

type Request struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}
