package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Flood-project/backend-flood/internal/token"
	"github.com/Flood-project/backend-flood/internal/token/util"
)

func UserUnathorized(response http.ResponseWriter) {
	response.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(response).Encode(map[string]string{
		"erro": "sem permissão",
	})
}

func CheckAuthentication(tokenManager token.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			authHeader := request.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				UserUnathorized(response)
				return
			}

			accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := tokenManager.ValidateToken(accessTokenString)
			if err != nil {
				UserUnathorized(response)
				return
			}

			if token.Type != "access" {
				UserUnathorized(response)
				return
			}


			claims := util.ExtractClaims(token)

			ctx := context.WithValue(request.Context(), "user_id", claims.IdUser)

			next.ServeHTTP(response, request.WithContext(ctx))
		})
	}
}