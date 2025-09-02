package token

import (
	"errors"
	"log"
	"time"
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateToken(account account_user.Account) (string, error)
	ValidateToken(tokenString string) (string, error)
}

type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) TokenManager {
	return &JWT{
		secretKey: secretKey,
	}
}

func (jwtToken *JWT) GenerateToken(account account_user.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": account.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.NewString(), //unique token
	})

	tokenString, err := token.SignedString([]byte(jwtToken.secretKey))
	if err != nil {
		return "", errors.New("erro ao gerar token")
	}

	return tokenString, nil
}

func (jwtToken *JWT) ValidateToken(tokenString string) (string, error) {
	jwtClaims := jwt.MapClaims{}

	tokenVerified, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (any, error) {
		return []byte(jwtToken.secretKey), nil
	})
	if err != nil {
		log.Println("erro jwt.Parse: ", err)
		return "", errors.New("erro ao validar token")
	}

	if !tokenVerified.Valid {
		return "", errors.New("token inválido")
	}

	claims, ok := tokenVerified.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("claims: ", claims)
		return "", errors.New("erro ao receber claims do token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("email não encontrado no token")
	}

	return email, nil
}
