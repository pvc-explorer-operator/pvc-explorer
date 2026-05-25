/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/pvc-explorer-operator/pvc-explorer/internal/auth"
)

type contextKey string

const (
	ctxRole     contextKey = "role"
	ctxUsername contextKey = "username"
)

type permission struct {
	method  string
	prefix  string
	minRole auth.Role
}

var routePermissions = []permission{
	{"POST", "/api/v1/auth/", auth.Role("")},
	{"GET", "/api/v1/health", auth.Role("")},

	{"POST", "/api/v1/explorers/", auth.RoleAdmin},
	{"PUT", "/api/v1/explorers/", auth.RoleAdmin},
	{"DELETE", "/api/v1/explorers/", auth.RoleAdmin},
	{"POST", "/api/v1/scopes", auth.RoleAdmin},
	{"PUT", "/api/v1/scopes/", auth.RoleAdmin},
	{"DELETE", "/api/v1/scopes/", auth.RoleAdmin},
}

type AuthMiddleware struct {
	sessions *auth.SessionStore
}

func NewAuthMiddleware(sessions *auth.SessionStore) *AuthMiddleware {
	return &AuthMiddleware{sessions: sessions}
}

func (a *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRoute(r) {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(auth.SessionCookieName)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		username, role, ok := a.sessions.Get(cookie.Value)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !isAllowed(r.Method, r.URL.Path, role) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), ctxRole, role)
		ctx = context.WithValue(ctx, ctxUsername, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) WrapFunc(fn http.HandlerFunc) http.HandlerFunc {
	return a.Wrap(fn).(http.HandlerFunc)
}

func isPublicRoute(r *http.Request) bool {
	return r.URL.Path == "/api/v1/health" ||
		r.URL.Path == "/api/v1/theme" ||
		strings.HasPrefix(r.URL.Path, "/api/v1/auth/") ||
		isStaticAsset(r.URL.Path)
}

func isStaticAsset(path string) bool {
	// Pass all non-API paths through — they are SPA routes served as index.html
	return !strings.HasPrefix(path, "/api/") ||
		path == "/" ||
		path == "/index.html" ||
		strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/favicon")
}

func isAllowed(method, path string, role auth.Role) bool {
	for _, p := range routePermissions {
		if p.minRole == "" {
			continue
		}
		if p.method == method && strings.HasPrefix(path, p.prefix) {
			return role == p.minRole || (p.minRole == auth.RoleViewer)
		}
	}
	if strings.HasPrefix(path, "/api/v1/explorers/") && strings.Contains(path, "/proxy/api/") {
		if method != http.MethodGet {
			return role == auth.RoleAdmin
		}
	}
	return true
}

func RoleFromContext(ctx context.Context) auth.Role {
	r, _ := ctx.Value(ctxRole).(auth.Role)
	return r
}

func UsernameFromContext(ctx context.Context) string {
	u, _ := ctx.Value(ctxUsername).(string)
	return u
}
