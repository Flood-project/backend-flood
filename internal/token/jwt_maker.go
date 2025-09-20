package token

import (
	"errors"
	"fmt"
	"time"
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/golang-jwt/jwt/v5"
	//"github.com/google/uuid"
)

type TokenManager interface {
	GenerateToken(account account_user.Account) (string, string, error)
	ValidateToken(tokenString string) (*CustomClaims, error)
	ValidateRefreshToken(token string) (*RefreshTokenCustomClaims, error)
}

type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) TokenManager {
	return &JWT{
		secretKey: secretKey,
	}
}

func (jwtToken *JWT) GenerateToken(account account_user.Account) (string, string, error) {
	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%d", account.Id_account),
		},
		IdUser:    account.Id_account,
		Email:     account.Email,
		UserGroup: account.IdUserGroup,
		Type: "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtToken.secretKey))
	if err != nil {
		return "", "", errors.New("erro ao gerar token")
	}

	// refreshToken
	refreshClaims := &RefreshTokenCustomClaims{
		IdUser: account.Id_account,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshJwt, err := refreshToken.SignedString([]byte(jwtToken.secretKey))
	if err != nil {
		return "", "", errors.New("erro ao gerar refresh token")
	}

	return tokenString, refreshJwt, nil
}

func (jwtToken *JWT) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtToken.secretKey), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("erro ao receber claims do token")
	}

	return claims, nil
}

func (jwtToken *JWT) ValidateRefreshToken(tokenString string) (*RefreshTokenCustomClaims, error) {
	refreshToken, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{},
	func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtToken.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !refreshToken.Valid {
		return nil, errors.New("token inválido")
	}

	RefreshClaims, ok := refreshToken.Claims.(*RefreshTokenCustomClaims)
	if !ok {
		return nil, errors.New("erro ao receber claims do token")
	}

	return RefreshClaims, nil
}
