package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Flood-project/backend-flood/internal/token"
)

func UserUnathorized(response http.ResponseWriter) {
	response.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(response).Encode(map[string]string{
		"error:": "unauthorized",
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

			_, err := tokenManager.ValidateToken(accessTokenString)
			if err != nil {
				UserUnathorized(response)
				return
			}

			next.ServeHTTP(response, request)
		})
	}
}