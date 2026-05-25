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

package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	pvcv1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
	"github.com/pvc-explorer-operator/pvc-explorer/internal/api"
	"github.com/pvc-explorer-operator/pvc-explorer/internal/scaler"
)

func restScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = corev1.AddToScheme(s)
	_ = pvcv1.AddToScheme(s)
	return s
}

func newRestMux(t *testing.T, objs ...runtime.Object) *http.ServeMux {
	t.Helper()
	c := fake.NewClientBuilder().WithScheme(restScheme()).WithRuntimeObjects(objs...).Build()
	b := api.NewBroadcaster()
	sc := scaler.New(c, b)
	h := api.NewRestHandler(c, sc, b, "pvc-explorer-system", "test-version")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func doRequest(t *testing.T, mux *http.ServeMux, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func TestListScopes_Empty(t *testing.T) {
	mux := newRestMux(t)
	w := doRequest(t, mux, http.MethodGet, "/api/v1/scopes", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateAndGetScope(t *testing.T) {
	mux := newRestMux(t)
	scope := pvcv1.PVCExplorerScope{
		ObjectMeta: metav1.ObjectMeta{Name: "prod"},
	}
	w := doRequest(t, mux, http.MethodPost, "/api/v1/scopes", scope)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	w2 := doRequest(t, mux, http.MethodGet, "/api/v1/scopes/prod", nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var got pvcv1.PVCExplorerScope
	_ = json.NewDecoder(w2.Body).Decode(&got)
	if got.Name != "prod" {
		t.Errorf("unexpected name: %s", got.Name)
	}
}

func TestGetScope_NotFound(t *testing.T) {
	mux := newRestMux(t)
	w := doRequest(t, mux, http.MethodGet, "/api/v1/scopes/missing", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteScope(t *testing.T) {
	existing := &pvcv1.PVCExplorerScope{
		ObjectMeta: metav1.ObjectMeta{Name: "prod"},
	}
	mux := newRestMux(t, existing)
	w := doRequest(t, mux, http.MethodDelete, "/api/v1/scopes/prod", nil)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestListExplorers_FilterByNamespace(t *testing.T) {
	e1 := &pvcv1.PVCExplorer{ObjectMeta: metav1.ObjectMeta{Name: "e1", Namespace: "ns-a"}}
	e2 := &pvcv1.PVCExplorer{ObjectMeta: metav1.ObjectMeta{Name: "e2", Namespace: "ns-b"}}
	mux := newRestMux(t, e1, e2)

	w := doRequest(t, mux, http.MethodGet, "/api/v1/explorers?namespace=ns-a", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var items []pvcv1.PVCExplorer
	_ = json.NewDecoder(w.Body).Decode(&items)
	if len(items) != 1 || items[0].Name != "e1" {
		t.Errorf("expected [e1], got %+v", items)
	}
}

func TestListExplorers_SearchFilter(t *testing.T) {
	e1 := &pvcv1.PVCExplorer{
		ObjectMeta: metav1.ObjectMeta{Name: "my-data", Namespace: "ns-a"},
		Spec:       pvcv1.PVCExplorerSpec{PVCName: "my-data"},
	}
	e2 := &pvcv1.PVCExplorer{
		ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "ns-a"},
		Spec:       pvcv1.PVCExplorerSpec{PVCName: "other"},
	}
	mux := newRestMux(t, e1, e2)

	w := doRequest(t, mux, http.MethodGet, "/api/v1/explorers?search=my", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var items []pvcv1.PVCExplorer
	_ = json.NewDecoder(w.Body).Decode(&items)
	if len(items) != 1 || items[0].Name != "my-data" {
		t.Errorf("expected [my-data], got %+v", items)
	}
}

func TestCreateExplorer(t *testing.T) {
	mux := newRestMux(t)
	explorer := pvcv1.PVCExplorer{
		ObjectMeta: metav1.ObjectMeta{Name: "my-pvc", Namespace: "default"},
		Spec:       pvcv1.PVCExplorerSpec{PVCName: "my-pvc"},
	}
	w := doRequest(t, mux, http.MethodPost, "/api/v1/explorers", explorer)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteExplorer_NotFound(t *testing.T) {
	mux := newRestMux(t)
	w := doRequest(t, mux, http.MethodDelete, "/api/v1/explorers/default/missing", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestListLabels(t *testing.T) {
	e := &pvcv1.PVCExplorer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "e1",
			Namespace: "default",
			Labels: map[string]string{
				"team": "aoi",
				"env":  "prod",
			},
		},
	}
	mux := newRestMux(t, e)
	w := doRequest(t, mux, http.MethodGet, "/api/v1/labels", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result map[string][]string
	_ = json.NewDecoder(w.Body).Decode(&result)
	if _, ok := result["team"]; !ok {
		t.Errorf("expected 'team' label key in response: %v", result)
	}
}

func TestListNamespaces(t *testing.T) {
	e := &pvcv1.PVCExplorer{ObjectMeta: metav1.ObjectMeta{Name: "e1", Namespace: "ns-a"}}
	mux := newRestMux(t, e)
	w := doRequest(t, mux, http.MethodGet, "/api/v1/namespaces", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListPVCs(t *testing.T) {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "my-pvc", Namespace: "default"},
	}
	mux := newRestMux(t, pvc)
	w := doRequest(t, mux, http.MethodGet, "/api/v1/namespaces/default/pvcs", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
