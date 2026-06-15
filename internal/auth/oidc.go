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
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	oidcConfigMapName = "pvc-explorer-config"
	rbacConfigMapName = "pvc-explorer-rbac"

	oidcKeyEnabled        = "oidc.enabled"
	oidcKeyIssuer         = "oidc.issuer"
	oidcKeyExternalIssuer = "oidc.externalIssuer"
	oidcKeyClientID       = "oidc.clientID"
	oidcKeyClientSecret   = "oidc.clientSecret"
	oidcKeyRedirectURI    = "oidc.redirectURI"
	oidcKeyScopes         = "oidc.scopes"
	oidcKeyGroupClaim     = "oidc.groupClaim"
	oidcKeySkipTLS        = "oidc.skipTLSVerify"

	rbacKeyPolicyDefault = "policy.default"
	rbacKeyPolicyCSV     = "policy.csv"
)

type OIDCConfig struct {
	Enabled        bool
	Issuer         string
	ExternalIssuer string
	ClientID       string
	ClientSecret   string
	RedirectURI    string
	Scopes         []string
	GroupClaim     string
	SkipTLSVerify  bool
}

type RBACConfig struct {
	DefaultRole Role
	PolicyCSV   string
}

type OIDCProvider struct {
	Config           *OIDCConfig
	RBAC             *RBACConfig
	Provider         *oidc.Provider
	OAuth2Config     *oauth2.Config
	Verifier         *oidc.IDTokenVerifier
	OAuth2HTTPClient *http.Client
}

type OIDCClaims struct {
	Subject string   `json:"sub"`
	Email   string   `json:"email"`
	Name    string   `json:"name"`
	Groups  []string `json:"groups"`
}

func LoadOIDCConfig(ctx context.Context, reader client.Reader, namespace string) (*OIDCConfig, error) {
	cm := &corev1.ConfigMap{}
	if err := reader.Get(ctx, types.NamespacedName{
		Name:      oidcConfigMapName,
		Namespace: namespace,
	}, cm); err != nil {
		return nil, fmt.Errorf("failed to get OIDC config: %w", err)
	}

	enabled := strings.EqualFold(cm.Data[oidcKeyEnabled], "true")
	if !enabled {
		return nil, nil
	}

	issuer := cm.Data[oidcKeyIssuer]
	clientID := cm.Data[oidcKeyClientID]
	redirectURI := cm.Data[oidcKeyRedirectURI]

	if issuer == "" || clientID == "" || redirectURI == "" {
		return nil, fmt.Errorf("oidc config missing required fields (issuer, clientID, redirectURI)")
	}

	secretRef := cm.Data[oidcKeyClientSecret]
	clientSecret := ""
	if secretRef != "" {
		secret, err := readSecretRef(ctx, reader, namespace, secretRef)
		if err != nil {
			return nil, fmt.Errorf("failed to read OIDC client secret: %w", err)
		}
		clientSecret = secret
	}

	scopes := []string{oidc.ScopeOpenID, "profile", "email"}
	if s := cm.Data[oidcKeyScopes]; s != "" {
		scopes = splitTrim(s)
	}

	groupClaim := "groups"
	if g := cm.Data[oidcKeyGroupClaim]; g != "" {
		groupClaim = g
	}

	return &OIDCConfig{
		Enabled:        true,
		Issuer:         issuer,
		ExternalIssuer: cm.Data[oidcKeyExternalIssuer],
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RedirectURI:    redirectURI,
		Scopes:         scopes,
		GroupClaim:     groupClaim,
		SkipTLSVerify:  strings.EqualFold(cm.Data[oidcKeySkipTLS], "true"),
	}, nil
}

func LoadRBACConfig(ctx context.Context, reader client.Reader, namespace string) *RBACConfig {
	cm := &corev1.ConfigMap{}
	if err := reader.Get(ctx, types.NamespacedName{
		Name:      rbacConfigMapName,
		Namespace: namespace,
	}, cm); err != nil {
		return &RBACConfig{DefaultRole: RoleViewer}
	}

	defaultRole := RoleViewer
	if d := cm.Data[rbacKeyPolicyDefault]; d != "" {
		defaultRole = Role(d)
	}

	return &RBACConfig{
		DefaultRole: defaultRole,
		PolicyCSV:   cm.Data[rbacKeyPolicyCSV],
	}
}

func NewOIDCProvider(ctx context.Context, cfg *OIDCConfig, rbac *RBACConfig) (*OIDCProvider, error) {
	var provider *oidc.Provider
	var err error

	if cfg.SkipTLSVerify {
		httpCtx := oidc.ClientContext(ctx, &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // dev only
			},
		})
		provider, err = oidc.NewProvider(httpCtx, cfg.Issuer)
	} else {
		provider, err = oidc.NewProvider(ctx, cfg.Issuer)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OIDC discovery document: %w", err)
	}

	oauth2HTTPClient := &http.Client{}
	if cfg.SkipTLSVerify {
		oauth2HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // dev only
			},
		}
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURI,
		Scopes:       cfg.Scopes,
	}

	verifier := provider.Verifier(&oidc.Config{
		SkipClientIDCheck:          false,
		SkipIssuerCheck:            false,
		SkipExpiryCheck:            false,
		InsecureSkipSignatureCheck: false,
		ClientID:                   cfg.ClientID,
	})

	return &OIDCProvider{
		Config:           cfg,
		RBAC:             rbac,
		OAuth2Config:     oauth2Config,
		Provider:         provider,
		Verifier:         verifier,
		OAuth2HTTPClient: oauth2HTTPClient,
	}, nil
}

func (p *OIDCProvider) AuthCodeURL(state string) string {
	authURL := p.OAuth2Config.AuthCodeURL(state, oidc.Nonce(state))

	if p.Config.ExternalIssuer != "" {
		if parsed, err := url.Parse(authURL); err == nil {
			if ext, err := url.Parse(p.Config.ExternalIssuer); err == nil {
				parsed.Host = ext.Host
				parsed.Scheme = ext.Scheme
				authURL = parsed.String()
			}
		}
	}

	return authURL
}

func (p *OIDCProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if p.OAuth2HTTPClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, p.OAuth2HTTPClient)
	}
	return p.OAuth2Config.Exchange(ctx, code)
}

func (p *OIDCProvider) VerifyToken(ctx context.Context, rawIDToken string) (*OIDCClaims, error) {
	token, err := p.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var claims OIDCClaims
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	return &claims, nil
}

func (p *OIDCProvider) MapGroupsToRole(groups []string) Role {
	return MapGroupsToRole(groups, p.RBAC)
}

func MapGroupsToRole(groups []string, rbac *RBACConfig) Role {
	rules := parsePolicyCSV(rbac.PolicyCSV)
	groupSet := make(map[string]bool, len(groups))
	for _, g := range groups {
		groupSet[g] = true
	}

	for _, rule := range rules {
		if rule.Type != "g" {
			continue
		}
		if groupSet[rule.Subject] {
			return rule.Role
		}
	}

	return rbac.DefaultRole
}

type policyRule struct {
	Type    string
	Subject string
	Role    Role
}

func parsePolicyCSV(csv string) []policyRule {
	var rules []policyRule
	for line := range strings.SplitSeq(csv, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 3)
		if len(parts) != 3 {
			continue
		}
		rules = append(rules, policyRule{
			Type:    strings.TrimSpace(parts[0]),
			Subject: strings.Trim(strings.TrimSpace(parts[1]), `"'`),
			Role:    Role(strings.TrimSpace(parts[2])),
		})
	}
	return rules
}

func readSecretRef(ctx context.Context, reader client.Reader, namespace, ref string) (string, error) {
	ref = strings.TrimPrefix(ref, "$")
	parts := strings.SplitN(ref, ".", 2)
	secretName := "pvc-explorer-oidc"
	key := "clientSecret"
	if len(parts) == 2 {
		secretName = parts[0]
		key = parts[1]
	}

	secret := &corev1.Secret{}
	if err := reader.Get(ctx, types.NamespacedName{
		Name:      secretName,
		Namespace: namespace,
	}, secret); err != nil {
		return "", err
	}

	val, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret %s", key, secretName)
	}

	return string(val), nil
}

func IsOIDCConfigured(ctx context.Context, reader client.Reader, namespace string) bool {
	cm := &corev1.ConfigMap{}
	if err := reader.Get(ctx, types.NamespacedName{
		Name:      oidcConfigMapName,
		Namespace: namespace,
	}, cm); err != nil {
		return false
	}
	return strings.EqualFold(cm.Data[oidcKeyEnabled], "true")
}
