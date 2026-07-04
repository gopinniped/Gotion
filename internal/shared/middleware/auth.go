package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gopinniped/gotion/internal/shared/errs"
	"github.com/gopinniped/gotion/internal/shared/mapper"
	"github.com/gopinniped/gotion/pkg/token"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

func Auth(maker *token.TokenMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "missing authorization header"))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "invalid authorization format"))
				return
			}

			claims, err := maker.Validate(parts[1])
			if err != nil {
				mapper.SendError(w, errs.New(errs.TypeUnauthenticated, "invalid or expired token"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}

func GetUsername(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(UsernameKey).(string)
	return name, ok
}
