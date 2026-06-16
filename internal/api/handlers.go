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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/pvc-explorer-operator/pvc-explorer/internal/auth"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

type Handler struct {
	auth     *auth.Authenticator
	sessions *auth.SessionStore
	oidc     *auth.OIDCProvider
}

func NewHandler(authenticator *auth.Authenticator, sessions *auth.SessionStore, oidcProvider *auth.OIDCProvider) *Handler {
	return &Handler{auth: authenticator, sessions: sessions, oidc: oidcProvider}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
	mux.HandleFunc("POST /api/v1/auth/logout", h.handleLogout)
	mux.HandleFunc("GET /api/v1/auth/me", h.handleMe)
	mux.HandleFunc("GET /api/v1/auth/config", h.handleAuthConfig)
	mux.HandleFunc("GET /api/v1/health", h.handleHealth)

	if h.oidc != nil {
		mux.HandleFunc("GET /api/v1/auth/oidc/start", h.handleOIDCStart)
		mux.HandleFunc("GET /api/v1/auth/oidc/callback", h.handleOIDCCallback)
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Role     string   `json:"role"`
	Username string   `json:"username"`
	Email    string   `json:"email,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Subject  string   `json:"subject,omitempty"`
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	log := ctrllog.FromContext(r.Context())

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	role, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			log.Info("Login failed", "method", "local", "username", req.Username, "reason", "invalid credentials")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Error(err, "Login error", "method", "local", "username", req.Username)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	token, err := h.sessions.Create(req.Username, role, "", nil, "")
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((8 * time.Hour).Seconds()),
	})

	log.Info("Login successful", "method", "local", "username", req.Username, "role", role)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(loginResponse{
		Role:     string(role),
		Username: req.Username,
	})
}

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	entry, ok := h.sessions.Get(cookie.Value)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(loginResponse{
		Role:     string(entry.Role),
		Username: entry.Username,
		Email:    entry.Email,
		Groups:   entry.Groups,
		Subject:  entry.Subject,
	})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err == nil {
		h.sessions.Delete(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   true,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type authConfigResponse struct {
	OIDCEnabled bool `json:"oidcEnabled"`
}

func (h *Handler) handleAuthConfig(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authConfigResponse{
		OIDCEnabled: h.oidc != nil,
	})
}

func (h *Handler) handleOIDCStart(w http.ResponseWriter, r *http.Request) {
	if h.oidc == nil {
		http.Error(w, "OIDC not configured", http.StatusNotFound)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oidc_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	http.Redirect(w, r, h.oidc.AuthCodeURL(state), http.StatusFound)
}

func (h *Handler) handleOIDCCallback(w http.ResponseWriter, r *http.Request) {
	log := ctrllog.FromContext(r.Context())

	if h.oidc == nil {
		http.Error(w, "OIDC not configured", http.StatusNotFound)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie("oidc_state")
	if err != nil {
		log.Info("OIDC callback failed", "reason", "missing state cookie")
		http.Error(w, "Missing state cookie", http.StatusBadRequest)
		return
	}
	stateParam := r.URL.Query().Get("state")
	if stateParam == "" || stateParam != stateCookie.Value {
		log.Info("OIDC callback failed", "reason", "invalid state parameter")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	clearCookie := &http.Cookie{
		Name:     "oidc_state",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	}
	http.SetCookie(w, clearCookie)

	oauth2Token, err := h.oidc.Exchange(r.Context(), code)
	if err != nil {
		log.Error(err, "OIDC token exchange failed")
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Info("OIDC callback failed", "reason", "missing id_token in token response")
		http.Error(w, "Missing id_token in token response", http.StatusInternalServerError)
		return
	}

	claims, err := h.oidc.VerifyToken(r.Context(), rawIDToken)
	if err != nil {
		log.Error(err, "OIDC token verification failed")
		http.Error(w, "Failed to verify token", http.StatusUnauthorized)
		return
	}

	groups := claims.Groups
	role := h.oidc.MapGroupsToRole(groups)

	username := claims.Name
	if username == "" {
		username = claims.Email
	}
	if username == "" {
		username = claims.Subject
	}

	token, err := h.sessions.Create(username, role, claims.Email, claims.Groups, claims.Subject)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((8 * time.Hour).Seconds()),
	})

	log.Info("Login successful", "method", "oidc", "username", username, "email", claims.Email, "subject", claims.Subject, "groups", groups, "role", role)

	http.Redirect(w, r, "/", http.StatusFound)
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
