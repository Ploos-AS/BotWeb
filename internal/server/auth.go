package server

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

type AuthConfig struct {
	ViewerToken   string
	OperatorToken string
}

func (c AuthConfig) enabled() bool { return c.ViewerToken != "" || c.OperatorToken != "" }

func (c AuthConfig) role(token string) string {
	if c.OperatorToken != "" && secureEqual(token, c.OperatorToken) {
		return "operator"
	}
	if c.ViewerToken != "" && secureEqual(token, c.ViewerToken) {
		return "viewer"
	}
	return ""
}

func secureEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func bearerToken(r *http.Request) string {
	v := r.Header.Get("Authorization")
	if !strings.HasPrefix(v, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(v, "Bearer "))
}

func authMiddleware(c AuthConfig, next http.Handler) http.Handler {
	if !c.enabled() {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		role := c.role(bearerToken(r))
		if role == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead && role != "operator" {
			http.Error(w, "operator role required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
