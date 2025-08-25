package dto

type AuthResponseDto struct {
	Token  string `json:"token"`
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Jti    string `json:"jti"`
}

type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
