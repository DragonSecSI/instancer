package middleware

import (
	"context"
	"net/http"

	"github.com/DragonSecSI/instancer/backend/pkg/helpers"
	"gorm.io/gorm"
)

func AuthAdminMiddleware(token string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !helpers.Auth.Admin.IsAdmin(r, token) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "admin", true)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func AuthUserMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			team, err := helpers.Auth.Token.GetTeam(db, r)
			if err != nil || team == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "team", team)
			ctx = context.WithValue(ctx, "admin", false)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func AuthAnyMiddleware(db *gorm.DB, token string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !helpers.Auth.Admin.IsAdmin(r, token) {
				team, err := helpers.Auth.Token.GetTeam(db, r)
				if err != nil || team == nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), "team", team)
				ctx = context.WithValue(ctx, "admin", false)
				r = r.WithContext(ctx)
			} else {
				ctx := context.WithValue(r.Context(), "admin", true)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
