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

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const (
	SessionCookieName = "pvc_session"
	sessionTTL        = 8 * time.Hour
)

type sessionEntry struct {
	Username string
	Role     Role
	expiry   time.Time
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]sessionEntry
}

func NewSessionStore() *SessionStore {
	s := &SessionStore{sessions: make(map[string]sessionEntry)}
	go s.reap()
	return s
}

func (s *SessionStore) Create(username string, role Role) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.sessions[token] = sessionEntry{
		Username: username,
		Role:     role,
		expiry:   time.Now().Add(sessionTTL),
	}
	s.mu.Unlock()
	return token, nil
}

func (s *SessionStore) Get(token string) (username string, role Role, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, exists := s.sessions[token]
	if !exists || time.Now().After(e.expiry) {
		delete(s.sessions, token)
		return "", "", false
	}
	return e.Username, e.Role, true
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func (s *SessionStore) reap() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for token, e := range s.sessions {
			if now.After(e.expiry) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
