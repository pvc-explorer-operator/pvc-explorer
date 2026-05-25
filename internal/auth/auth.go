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

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get,namespace=pvc-explorer-system
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

import (
	"context"
	"errors"
	"slices"
	"strings"

	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleViewer Role = "viewer"

	secretName    = "pvc-explorer-auth"
	configMapName = "pvc-explorer-config"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Authenticator struct {
	reader    client.Reader
	namespace string
}

func NewAuthenticator(r client.Reader, namespace string) *Authenticator {
	return &Authenticator{reader: r, namespace: namespace}
}

func (a *Authenticator) Login(ctx context.Context, username, password string) (Role, error) {
	secret := &corev1.Secret{}
	if err := a.reader.Get(ctx, types.NamespacedName{
		Name:      secretName,
		Namespace: a.namespace,
	}, secret); err != nil {
		return "", err
	}

	hash, ok := secret.Data[username]
	if !ok {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	role := a.resolveRole(ctx, username)
	return role, nil
}

func (a *Authenticator) resolveRole(ctx context.Context, username string) Role {
	cm := &corev1.ConfigMap{}
	if err := a.reader.Get(ctx, types.NamespacedName{
		Name:      configMapName,
		Namespace: a.namespace,
	}, cm); err != nil {
		if username == "admin" {
			return RoleAdmin
		}
		return RoleViewer
	}

	adminUsers := splitTrim(cm.Data["adminUsers"])
	if slices.Contains(adminUsers, username) {
		return RoleAdmin
	}
	return RoleViewer
}

func splitTrim(s string) []string {
	var out []string
	for part := range strings.SplitSeq(s, ",") {
		if t := strings.TrimSpace(part); t != "" {
			out = append(out, t)
		}
	}
	return out
}
