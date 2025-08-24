package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	jwt.RegisteredClaims
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Jti    string `json: "jti"`
}
