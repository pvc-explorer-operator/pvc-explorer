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
	"testing"

	"github.com/pvc-explorer-operator/pvc-explorer/internal/auth"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testConfigMapName = "pvc-explorer-config"
	testPolicyCSV     = `g, "admin-group", admin`
	testOIDCEnabled   = "oidc.enabled"
	testTrue          = "true"
	testOIDCIssuer    = "oidc.issuer"
)

func TestMapGroupsToRole(t *testing.T) {
	tests := []struct {
		name     string
		groups   []string
		rbac     *auth.RBACConfig
		expected auth.Role
	}{
		{
			name:   "admin group matches",
			groups: []string{"admin-group", "other-group"},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleViewer,
				PolicyCSV:   testPolicyCSV,
			},
			expected: auth.RoleAdmin,
		},
		{
			name:   "user group matches",
			groups: []string{"user-group"},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleViewer,
				PolicyCSV: testPolicyCSV + `
g, "user-group", user`,
			},
			expected: auth.RoleUser,
		},
		{
			name:   "no match returns default",
			groups: []string{"unknown-group"},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleViewer,
				PolicyCSV:   testPolicyCSV,
			},
			expected: auth.RoleViewer,
		},
		{
			name:   "default role is user",
			groups: []string{"unknown-group"},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleUser,
				PolicyCSV:   testPolicyCSV,
			},
			expected: auth.RoleUser,
		},
		{
			name:   "empty groups returns default",
			groups: []string{},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleViewer,
				PolicyCSV:   testPolicyCSV,
			},
			expected: auth.RoleViewer,
		},
		{
			name:   "first matching rule wins",
			groups: []string{"group-a", "group-b"},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleViewer,
				PolicyCSV: `g, "group-a", user
g, "group-b", admin`,
			},
			expected: auth.RoleUser,
		},
		{
			name:   "azure ad GUID format",
			groups: []string{"d1f2e3a4-b5c6-7890-abcd-ef1234567890"},
			rbac: &auth.RBACConfig{
				DefaultRole: auth.RoleViewer,
				PolicyCSV:   `g, "d1f2e3a4-b5c6-7890-abcd-ef1234567890", admin`,
			},
			expected: auth.RoleAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := auth.MapGroupsToRole(tt.groups, tt.rbac)
			if result != tt.expected {
				t.Errorf("expected role %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestLoadOIDCConfig_Disabled(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testConfigMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{
			testOIDCEnabled: "false",
		},
	}
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).WithObjects(cm).Build()

	cfg, err := auth.LoadOIDCConfig(context.Background(), c, testNamespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Fatal("expected nil config when OIDC is disabled")
	}
}

func TestLoadOIDCConfig_MissingConfigMap(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).Build()

	cfg, err := auth.LoadOIDCConfig(context.Background(), c, testNamespace)
	if err == nil {
		t.Fatal("expected error when ConfigMap is missing")
	}
	if cfg != nil {
		t.Fatal("expected nil config")
	}
}

func TestLoadOIDCConfig_MissingFields(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testConfigMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{
			testOIDCEnabled: testTrue,
			testOIDCIssuer:  "https://example.com",
		},
	}
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).WithObjects(cm).Build()

	_, err := auth.LoadOIDCConfig(context.Background(), c, testNamespace)
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func TestLoadOIDCConfig_ValidConfig(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testConfigMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{
			testOIDCEnabled:     testTrue,
			testOIDCIssuer:      "https://login.microsoftonline.com/tenant/v2.0",
			"oidc.clientID":     "test-client-id",
			"oidc.clientSecret": "$oidc.clientSecret",
			"oidc.redirectURI":  "https://app.example.com/callback",
			"oidc.scopes":       "openid,profile,email,groups",
			"oidc.groupClaim":   "groups",
		},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "oidc",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"clientSecret": []byte("test-secret"),
		},
	}
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).WithObjects(cm, secret).Build()

	cfg, err := auth.LoadOIDCConfig(context.Background(), c, testNamespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.Issuer != "https://login.microsoftonline.com/tenant/v2.0" {
		t.Errorf("unexpected issuer: %s", cfg.Issuer)
	}
	if cfg.ClientID != "test-client-id" {
		t.Errorf("unexpected clientID: %s", cfg.ClientID)
	}
	if cfg.ClientSecret != "test-secret" {
		t.Errorf("unexpected clientSecret: %s", cfg.ClientSecret)
	}
	if cfg.RedirectURI != "https://app.example.com/callback" {
		t.Errorf("unexpected redirectURI: %s", cfg.RedirectURI)
	}
	if len(cfg.Scopes) != 4 {
		t.Errorf("expected 4 scopes, got %d", len(cfg.Scopes))
	}
	if cfg.GroupClaim != "groups" {
		t.Errorf("unexpected groupClaim: %s", cfg.GroupClaim)
	}
	if cfg.SkipTLSVerify {
		t.Error("expected SkipTLSVerify to be false by default")
	}
}

func TestLoadOIDCConfig_SkipTLSVerify(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testConfigMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{
			testOIDCEnabled:      testTrue,
			testOIDCIssuer:       "https://localhost:5556",
			"oidc.clientID":      "test-client",
			"oidc.clientSecret":  "$oidc.clientSecret",
			"oidc.redirectURI":   "https://localhost:8080/callback",
			"oidc.skipTLSVerify": testTrue,
		},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "oidc",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"clientSecret": []byte("secret"),
		},
	}
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).WithObjects(cm, secret).Build()

	cfg, err := auth.LoadOIDCConfig(context.Background(), c, testNamespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.SkipTLSVerify {
		t.Error("expected SkipTLSVerify to be true")
	}
}

func TestLoadRBACConfig_Defaults(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).Build()

	rbac := auth.LoadRBACConfig(context.Background(), c, testNamespace)
	if rbac.DefaultRole != auth.RoleViewer {
		t.Errorf("expected default role viewer, got %s", rbac.DefaultRole)
	}
}

func TestLoadRBACConfig_CustomDefaults(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pvc-explorer-rbac",
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"policy.default": "user",
			"policy.csv":     testPolicyCSV,
		},
	}
	c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).WithObjects(cm).Build()

	rbac := auth.LoadRBACConfig(context.Background(), c, testNamespace)
	if rbac.DefaultRole != auth.RoleUser {
		t.Errorf("expected default role user, got %s", rbac.DefaultRole)
	}
	if rbac.PolicyCSV != testPolicyCSV {
		t.Errorf("unexpected policy CSV: %s", rbac.PolicyCSV)
	}
}

func TestIsOIDCConfigured(t *testing.T) {
	tests := []struct {
		name     string
		cmData   map[string]string
		expected bool
	}{
		{
			name:     "enabled",
			cmData:   map[string]string{testOIDCEnabled: testTrue},
			expected: true,
		},
		{
			name:     "enabled case insensitive",
			cmData:   map[string]string{testOIDCEnabled: "True"},
			expected: true,
		},
		{
			name:     "disabled",
			cmData:   map[string]string{testOIDCEnabled: "false"},
			expected: false,
		},
		{
			name:     "missing key",
			cmData:   map[string]string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testConfigMapName,
					Namespace: testNamespace,
				},
				Data: tt.cmData,
			}
			c := fake.NewClientBuilder().WithScheme(oidcTestScheme()).WithObjects(cm).Build()

			result := auth.IsOIDCConfigured(context.Background(), c, testNamespace)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func oidcTestScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = corev1.AddToScheme(s)
	return s
}
