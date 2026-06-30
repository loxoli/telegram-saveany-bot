package api

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/krau/SaveAny-Bot/config"
)

// tokenContextKey 用於在 context 中儲存 token
type tokenContextKey struct{}

// AuthMiddleware 回傳認證中間件
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := config.C().API

			// 從請求標頭取得 token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, http.StatusUnauthorized, "unauthorized", "missing authorization header")
				return
			}

			// 提取 Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid authorization header format")
				return
			}

			token := parts[1]

			// 驗證 token
			if subtle.ConstantTimeCompare([]byte(token), []byte(cfg.Token)) != 1 {
				WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}

			// 將 token 加入 context
			ctx := context.WithValue(r.Context(), tokenContextKey{}, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
