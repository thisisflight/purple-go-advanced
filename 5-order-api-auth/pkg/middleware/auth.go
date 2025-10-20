package middleware

import (
	"context"
	"net/http"
	"purple/links/configs"
	"purple/links/pkg/jwt"
	"strings"
)

type Key string

const (
	ContextEmailKey Key = "ContextEmailKey"
)

func writeUnauthorizedHeader(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func AuthMiddleware(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorizedHeader(w)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthorizedHeader(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.PhoneNumber)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
