package auth

import (
	"net/http"
)

type AuthAdmin struct {
	IsAdmin func(r *http.Request, token string) bool
}

func authAdminIsAdmin(r *http.Request, token string) bool {
	t := r.Header.Get("Authorization")
	if t == token {
		return true
	}

	t = r.URL.Query().Get("token")
	return t == token
}
