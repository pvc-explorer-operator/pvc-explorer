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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pvcv1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
	"github.com/pvc-explorer-operator/pvc-explorer/internal/scaler"
)

const keyName = "name"

type RestHandler struct {
	client          client.Client
	scaler          *scaler.Scaler
	broadcaster     *Broadcaster
	systemNamespace string
	version         string
}

func NewRestHandler(c client.Client, s *scaler.Scaler, b *Broadcaster, ns string, version string) *RestHandler {
	return &RestHandler{client: c, scaler: s, broadcaster: b, systemNamespace: ns, version: version}
}

func (h *RestHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/scopes", h.listScopes)
	mux.HandleFunc("GET /api/v1/scopes/{name}", h.getScope)
	mux.HandleFunc("POST /api/v1/scopes", h.createScope)
	mux.HandleFunc("PUT /api/v1/scopes/{name}", h.updateScope)
	mux.HandleFunc("DELETE /api/v1/scopes/{name}", h.deleteScope)

	mux.HandleFunc("GET /api/v1/explorers", h.listExplorers)
	mux.HandleFunc("GET /api/v1/explorers/{ns}/{name}", h.getExplorer)
	mux.HandleFunc("POST /api/v1/explorers", h.createExplorer)
	mux.HandleFunc("PUT /api/v1/explorers/{ns}/{name}", h.updateExplorer)
	mux.HandleFunc("DELETE /api/v1/explorers/{ns}/{name}", h.deleteExplorer)

	mux.HandleFunc("POST /api/v1/explorers/{ns}/{name}/wake", h.wakeExplorer)
	mux.HandleFunc("POST /api/v1/explorers/{ns}/{name}/sleep", h.sleepExplorer)
	mux.HandleFunc("POST /api/v1/explorers/{ns}/{name}/heartbeat", h.heartbeat)
	mux.HandleFunc("/api/v1/explorers/", h.proxyDispatch)

	mux.HandleFunc("GET /api/v1/labels", h.listLabels)
	mux.HandleFunc("GET /api/v1/namespaces", h.listNamespaces)
	mux.HandleFunc("GET /api/v1/namespaces/{ns}/pvcs", h.listPVCs)
	mux.HandleFunc("GET /api/v1/theme", h.getTheme)
	mux.HandleFunc("GET /api/version", h.getVersion)
}

func (h *RestHandler) getVersion(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, h.version)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *RestHandler) listScopes(w http.ResponseWriter, r *http.Request) {
	var list pvcv1.PVCExplorerScopeList
	if err := h.client.List(r.Context(), &list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, list.Items)
}

func (h *RestHandler) getScope(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var scope pvcv1.PVCExplorerScope
	if err := h.client.Get(r.Context(), types.NamespacedName{Name: name}, &scope); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, scope)
}

func (h *RestHandler) createScope(w http.ResponseWriter, r *http.Request) {
	var scope pvcv1.PVCExplorerScope
	if err := json.NewDecoder(r.Body).Decode(&scope); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.client.Create(r.Context(), &scope); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.broadcaster.Publish(string(EventScopeUpdated), scope)
	writeJSON(w, http.StatusCreated, scope)
}

func (h *RestHandler) updateScope(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var incoming pvcv1.PVCExplorerScope
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var existing pvcv1.PVCExplorerScope
	if err := h.client.Get(r.Context(), types.NamespacedName{Name: name}, &existing); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	patch := client.MergeFrom(existing.DeepCopy())
	existing.Spec = incoming.Spec
	existing.Labels = incoming.Labels
	if err := h.client.Patch(r.Context(), &existing, patch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.broadcaster.Publish(string(EventScopeUpdated), existing)
	writeJSON(w, http.StatusOK, existing)
}

func (h *RestHandler) deleteScope(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var scope pvcv1.PVCExplorerScope
	if err := h.client.Get(r.Context(), types.NamespacedName{Name: name}, &scope); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.client.Delete(r.Context(), &scope); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.broadcaster.Publish(string(EventScopeDeleted), map[string]string{keyName: name})
	w.WriteHeader(http.StatusNoContent)
}

func (h *RestHandler) listExplorers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	opts := []client.ListOption{}
	if ns := q.Get("namespace"); ns != "" {
		opts = append(opts, client.InNamespace(ns))
	}
	labelFilters := map[string]string{}
	if scope := q.Get("scope"); scope != "" {
		labelFilters["pvcexplorer.io/scope"] = scope
	}
	if mode := q.Get("mode"); mode != "" {
		labelFilters["pvcexplorer.io/mode"] = strings.ToLower(mode)
	}
	if mountState := q.Get("mountState"); mountState != "" {
		labelFilters["pvcexplorer.io/mount-state"] = mountState
	}
	if accessMode := q.Get("accessMode"); accessMode != "" {
		labelFilters["pvcexplorer.io/access-mode"] = strings.ToLower(accessMode)
	}
	if len(labelFilters) > 0 {
		opts = append(opts, client.MatchingLabels(labelFilters))
	}

	var list pvcv1.PVCExplorerList
	if err := h.client.List(r.Context(), &list, opts...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := list.Items
	if search := q.Get("search"); search != "" {
		search = strings.ToLower(search)
		var filtered []pvcv1.PVCExplorer
		for _, e := range items {
			if strings.Contains(strings.ToLower(e.Name), search) ||
				strings.Contains(strings.ToLower(e.Namespace), search) ||
				strings.Contains(strings.ToLower(e.Spec.PVCName), search) {
				filtered = append(filtered, e)
			}
		}
		items = filtered
	}
	if phase := q.Get("phase"); phase != "" {
		var filtered []pvcv1.PVCExplorer
		for _, e := range items {
			if strings.EqualFold(string(e.Status.Phase), phase) {
				filtered = append(filtered, e)
			}
		}
		items = filtered
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *RestHandler) getExplorer(w http.ResponseWriter, r *http.Request) {
	ns, name := r.PathValue("ns"), r.PathValue("name")
	var explorer pvcv1.PVCExplorer
	if err := h.client.Get(r.Context(), types.NamespacedName{Namespace: ns, Name: name}, &explorer); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, explorer)
}

func (h *RestHandler) createExplorer(w http.ResponseWriter, r *http.Request) {
	var explorer pvcv1.PVCExplorer
	if err := json.NewDecoder(r.Body).Decode(&explorer); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.client.Create(r.Context(), &explorer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.broadcaster.Publish(string(EventExplorerUpdated), explorer)
	writeJSON(w, http.StatusCreated, explorer)
}

func (h *RestHandler) updateExplorer(w http.ResponseWriter, r *http.Request) {
	ns, name := r.PathValue("ns"), r.PathValue("name")
	var incoming pvcv1.PVCExplorer
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var existing pvcv1.PVCExplorer
	if err := h.client.Get(r.Context(), types.NamespacedName{Namespace: ns, Name: name}, &existing); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	patch := client.MergeFrom(existing.DeepCopy())
	existing.Spec = incoming.Spec
	existing.Labels = incoming.Labels
	if err := h.client.Patch(r.Context(), &existing, patch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.broadcaster.Publish(string(EventExplorerUpdated), existing)
	writeJSON(w, http.StatusOK, existing)
}

func (h *RestHandler) deleteExplorer(w http.ResponseWriter, r *http.Request) {
	ns, name := r.PathValue("ns"), r.PathValue("name")
	var explorer pvcv1.PVCExplorer
	if err := h.client.Get(r.Context(), types.NamespacedName{Namespace: ns, Name: name}, &explorer); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.client.Delete(r.Context(), &explorer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.broadcaster.Publish(string(EventExplorerDeleted), map[string]string{"namespace": ns, keyName: name})
	w.WriteHeader(http.StatusNoContent)
}

func (h *RestHandler) wakeExplorer(w http.ResponseWriter, r *http.Request) {
	ns, name := r.PathValue("ns"), r.PathValue("name")
	_ = h.broadcaster.Publish(string(EventAgentWaking), map[string]string{"namespace": ns, keyName: name})
	if err := h.scaler.WakeAgent(r.Context(), ns, name); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "waking"})
}

func (h *RestHandler) sleepExplorer(w http.ResponseWriter, r *http.Request) {
	ns, name := r.PathValue("ns"), r.PathValue("name")
	if err := h.scaler.SleepAgent(r.Context(), ns, name); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "sleeping"})
}

type heartbeatResponse struct {
	RemainingSeconds int64  `json:"remainingSeconds"`
	IdleTimeout      string `json:"idleTimeout"`
	Phase            string `json:"phase"`
}

func (h *RestHandler) heartbeat(w http.ResponseWriter, r *http.Request) {
	ns, name := r.PathValue("ns"), r.PathValue("name")
	remaining, err := h.scaler.ResetIdleTimer(r.Context(), ns, name)
	if err != nil {
		if strings.Contains(err.Error(), "agent not running") {
			var explorer pvcv1.PVCExplorer
			_ = h.client.Get(r.Context(), types.NamespacedName{Namespace: ns, Name: name}, &explorer)
			writeJSON(w, http.StatusConflict, map[string]string{
				"error": "agent not running",
				"phase": string(explorer.Status.Phase),
			})
			return
		}
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var explorer pvcv1.PVCExplorer
	_ = h.client.Get(r.Context(), types.NamespacedName{Namespace: ns, Name: name}, &explorer)
	writeJSON(w, http.StatusOK, heartbeatResponse{
		RemainingSeconds: int64(remaining.Seconds()),
		IdleTimeout:      explorer.Spec.Scaling.IdleTimeout,
		Phase:            string(explorer.Status.Phase),
	})
}

func (h *RestHandler) proxyDispatch(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/api/v1/explorers/"), "/", 3)
	if len(parts) < 3 || parts[2] == "" {
		http.NotFound(w, r)
		return
	}
	ns, name, tail := parts[0], parts[1], parts[2]
	if !strings.HasPrefix(tail, "proxy/") {
		http.NotFound(w, r)
		return
	}

	var explorer pvcv1.PVCExplorer
	if err := h.client.Get(r.Context(), types.NamespacedName{Namespace: ns, Name: name}, &explorer); err != nil {
		if apierrors.IsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if explorer.Status.AgentEndpoint == "" {
		http.Error(w, "agent endpoint not available", http.StatusServiceUnavailable)
		return
	}

	target, err := url.Parse(explorer.Status.AgentEndpoint)
	if err != nil {
		http.Error(w, "invalid agent endpoint", http.StatusInternalServerError)
		return
	}

	_, _ = h.scaler.ResetIdleTimer(r.Context(), ns, name)

	// Read the agent authentication token from the shared Secret.
	agentToken, err := h.readAgentToken(r.Context(), ns, name)
	if err != nil {
		http.Error(w, "failed to read agent token", http.StatusInternalServerError)
		return
	}

	proxyPath := "/" + strings.TrimPrefix(tail, "proxy/")
	r2 := r.Clone(r.Context())
	r2.URL.Path = proxyPath
	r2.URL.RawPath = ""
	r2.Host = target.Host
	r2.Header.Set("Authorization", "Bearer "+agentToken)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r2)
}

func (h *RestHandler) listLabels(w http.ResponseWriter, r *http.Request) {
	var list pvcv1.PVCExplorerList
	if err := h.client.List(r.Context(), &list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	labels := map[string]map[string]struct{}{}
	for _, e := range list.Items {
		for k, v := range e.Labels {
			if strings.HasPrefix(k, "pvcexplorer.io/") {
				continue
			}
			if labels[k] == nil {
				labels[k] = map[string]struct{}{}
			}
			labels[k][v] = struct{}{}
		}
	}

	result := map[string][]string{}
	for k, vs := range labels {
		for v := range vs {
			result[k] = append(result[k], v)
		}
	}
	writeJSON(w, http.StatusOK, result)
}

type namespaceInfo struct {
	Namespace     string `json:"namespace"`
	PVCCount      int    `json:"pvcCount"`
	ExplorerCount int    `json:"explorerCount"`
}

func (h *RestHandler) listNamespaces(w http.ResponseWriter, r *http.Request) {
	var explorerList pvcv1.PVCExplorerList
	if err := h.client.List(r.Context(), &explorerList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nsMap := map[string]*namespaceInfo{}
	for _, e := range explorerList.Items {
		ns := e.Namespace
		if nsMap[ns] == nil {
			nsMap[ns] = &namespaceInfo{Namespace: ns}
		}
		nsMap[ns].ExplorerCount++
	}

	for ns, info := range nsMap {
		var pvcList corev1.PersistentVolumeClaimList
		if err := h.client.List(r.Context(), &pvcList, client.InNamespace(ns)); err == nil {
			info.PVCCount = len(pvcList.Items)
		}
	}

	result := make([]namespaceInfo, 0, len(nsMap))
	for _, v := range nsMap {
		result = append(result, *v)
	}
	writeJSON(w, http.StatusOK, result)
}

type pvcInfo struct {
	Name        string                              `json:"name"`
	Namespace   string                              `json:"namespace"`
	Phase       string                              `json:"phase"`
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Consumers   []pvcv1.ConsumerInfo                `json:"consumers,omitempty"`
}

func (h *RestHandler) listPVCs(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("ns")

	var pvcList corev1.PersistentVolumeClaimList
	if err := h.client.List(r.Context(), &pvcList, client.InNamespace(ns)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var explorerList pvcv1.PVCExplorerList
	_ = h.client.List(r.Context(), &explorerList, client.InNamespace(ns))
	consumerMap := map[string][]pvcv1.ConsumerInfo{}
	for _, e := range explorerList.Items {
		consumerMap[e.Spec.PVCName] = e.Status.Mount.Consumers
	}

	result := make([]pvcInfo, 0, len(pvcList.Items))
	for _, pvc := range pvcList.Items {
		result = append(result, pvcInfo{
			Name:        pvc.Name,
			Namespace:   pvc.Namespace,
			Phase:       string(pvc.Status.Phase),
			AccessModes: pvc.Spec.AccessModes,
			Consumers:   consumerMap[pvc.Name],
		})
	}
	writeJSON(w, http.StatusOK, result)
}

var _ = metav1.Now

type themeResponse struct {
	AppName         string `json:"appName"`
	LogoURL         string `json:"logoUrl,omitempty"`
	PrimaryColor    string `json:"primaryColor,omitempty"`
	DarkModeDefault bool   `json:"darkModeDefault"`
}

func (h *RestHandler) getTheme(w http.ResponseWriter, r *http.Request) {
	defaults := themeResponse{
		AppName:         "PVC Explorer",
		DarkModeDefault: false,
	}

	var cm corev1.ConfigMap
	err := h.client.Get(r.Context(), types.NamespacedName{
		Namespace: h.systemNamespace,
		Name:      "pvc-explorer-theme",
	}, &cm)
	if err != nil {
		writeJSON(w, http.StatusOK, defaults)
		return
	}

	raw, ok := cm.Data["theme.json"]
	if !ok {
		writeJSON(w, http.StatusOK, defaults)
		return
	}

	var theme themeResponse
	if err := json.Unmarshal([]byte(raw), &theme); err != nil {
		writeJSON(w, http.StatusOK, defaults)
		return
	}

	if theme.AppName == "" {
		theme.AppName = defaults.AppName
	}

	writeJSON(w, http.StatusOK, theme)
}

const agentTokenSecretSuffix = "-agent-token"
const agentTokenSecretKey = "token"

// readAgentToken reads the operator-generated bearer token from the
// agent authentication Secret in the explorer's namespace.
func (h *RestHandler) readAgentToken(ctx context.Context, ns, explorerName string) (string, error) {
	secret := &corev1.Secret{}
	if err := h.client.Get(ctx, types.NamespacedName{
		Namespace: ns,
		Name:      explorerName + agentTokenSecretSuffix,
	}, secret); err != nil {
		return "", err
	}
	token, ok := secret.Data[agentTokenSecretKey]
	if !ok {
		return "", fmt.Errorf("secret %q has no %q key", explorerName+agentTokenSecretSuffix, agentTokenSecretKey)
	}
	return string(token), nil
}
