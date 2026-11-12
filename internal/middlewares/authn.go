package middlewares

import (
	"aaa_bekend/internal/lmodel"
	"net/http"
)

func AuthN(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		userCtx := lmodel.UserContext{
            ID:       claims.UserID,
            Username: claims.Username,
            Role:     claims.Role,
        }

		ctx := lmodel.ContextWithUser(r.Context(), userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
