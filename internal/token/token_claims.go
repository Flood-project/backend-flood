package token

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	jwt.RegisteredClaims
	IdUser    int32  `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	UserGroup int32  `json:"id_user_group" db:"id_user_group"`
	Type      string `json:"type"`
}

type RefreshTokenCustomClaims struct {
	IdUser int32
	Type   string
	jwt.RegisteredClaims
}
