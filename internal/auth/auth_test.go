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

package auth_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/pvc-explorer-operator/pvc-explorer/internal/auth"
)

const testNamespace = "pvc-explorer-system"

func testScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = corev1.AddToScheme(s)
	return s
}

func bcryptHash(t *testing.T, password string) []byte {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func authSecret(users map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pvc-explorer-auth",
			Namespace: testNamespace,
		},
		Data: users,
	}
}

func TestLogin_ValidAdmin(t *testing.T) {
	secret := authSecret(map[string][]byte{
		"admin": bcryptHash(t, "secret"),
	})
	c := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(secret).Build()
	a := auth.NewAuthenticator(c, testNamespace)

	role, err := a.Login(context.Background(), "admin", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != auth.RoleAdmin {
		t.Errorf("expected admin role, got %s", role)
	}
}

func TestLogin_ValidViewer(t *testing.T) {
	secret := authSecret(map[string][]byte{
		"alice": bcryptHash(t, "pass"),
	})
	c := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(secret).Build()
	a := auth.NewAuthenticator(c, testNamespace)

	role, err := a.Login(context.Background(), "alice", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != auth.RoleViewer {
		t.Errorf("expected viewer role, got %s", role)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	secret := authSecret(map[string][]byte{
		"admin": bcryptHash(t, "correct"),
	})
	c := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(secret).Build()
	a := auth.NewAuthenticator(c, testNamespace)

	_, err := a.Login(context.Background(), "admin", "wrong")
	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	secret := authSecret(map[string][]byte{
		"admin": bcryptHash(t, "pass"),
	})
	c := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(secret).Build()
	a := auth.NewAuthenticator(c, testNamespace)

	_, err := a.Login(context.Background(), "nobody", "pass")
	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_ConfigMapAdminList(t *testing.T) {
	secret := authSecret(map[string][]byte{
		"ops-lead": bcryptHash(t, "pass"),
	})
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pvc-explorer-config",
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"adminUsers": "admin, ops-lead",
		},
	}
	c := fake.NewClientBuilder().WithScheme(testScheme()).WithObjects(secret, cm).Build()
	a := auth.NewAuthenticator(c, testNamespace)

	role, err := a.Login(context.Background(), "ops-lead", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != auth.RoleAdmin {
		t.Errorf("expected admin role via ConfigMap, got %s", role)
	}
}

func TestSessionStore_CreateAndGet(t *testing.T) {
	s := auth.NewSessionStore()
	token, err := s.Create("alice", auth.RoleViewer)
	if err != nil {
		t.Fatal(err)
	}

	username, role, ok := s.Get(token)
	if !ok {
		t.Fatal("expected session to exist")
	}
	if username != "alice" || role != auth.RoleViewer {
		t.Errorf("unexpected session: %s / %s", username, role)
	}
}

func TestSessionStore_DeleteRemovesSession(t *testing.T) {
	s := auth.NewSessionStore()
	token, _ := s.Create("bob", auth.RoleAdmin)
	s.Delete(token)

	_, _, ok := s.Get(token)
	if ok {
		t.Fatal("expected session to be deleted")
	}
}

func TestSessionStore_UnknownTokenReturnsFalse(t *testing.T) {
	s := auth.NewSessionStore()
	_, _, ok := s.Get("nonexistent-token")
	if ok {
		t.Fatal("expected false for unknown token")
	}
}
