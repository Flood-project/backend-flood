package util

import (
	"github.com/Flood-project/backend-flood/internal/token"
)

func ExtractClaims(customClaims *token.CustomClaims) *token.CustomClaims {
	return customClaims
}