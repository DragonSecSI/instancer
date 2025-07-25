package auth

import (
	"net/http"
)

type AuthAdmin struct {
	IsAdmin func(r *http.Request, token string) bool
}

func authAdminIsAdmin(r *http.Request, token string) bool {
	t := r.Header.Get("Authorization")
	return t == token
}
