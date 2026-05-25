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

package controller

import (
	"context"
	"fmt"
	"path"
	"slices"
	"time"

	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
)

// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorerscopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorerscopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorerscopes/finalizers,verbs=update
// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch

const (
	scopeRequeueAfter = 30 * time.Second
	scopeLabelKey     = "pvcexplorer.io/scope"
	managedByLabelKey = "pvcexplorer.io/managed-by"
	managedByValue    = "pvc-explorer"
	scopeFinalizer    = "pvcexplorer.io/scope-cleanup"
)

type PVCExplorerScopeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PVCExplorerScopeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	scope := &pvcexplorerv1alpha1.PVCExplorerScope{}
	if err := r.Get(ctx, req.NamespacedName, scope); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !scope.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, scope)
	}

	if !controllerutil.ContainsFinalizer(scope, scopeFinalizer) {
		patch := client.MergeFrom(scope.DeepCopy())
		controllerutil.AddFinalizer(scope, scopeFinalizer)
		if err := r.Patch(ctx, scope, patch); err != nil {
			return ctrl.Result{}, err
		}
	}

	namespaces, err := r.resolveNamespaces(ctx, scope)
	if err != nil {
		log.Error(err, "Could not resolve namespaces")
		return ctrl.Result{RequeueAfter: scopeRequeueAfter}, nil
	}

	pvcMap, err := r.discoverPVCs(ctx, scope, namespaces)
	if err != nil {
		log.Error(err, "Could not discover PVCs")
		return ctrl.Result{RequeueAfter: scopeRequeueAfter}, nil
	}

	totalExplorers := 0
	for _, pvcs := range pvcMap {
		totalExplorers += len(pvcs)
	}

	if syncErr := r.syncPVCExplorers(ctx, scope, pvcMap); syncErr != nil {
		log.Error(syncErr, "Could not sync PVCExplorer CRs")
		r.setReadyCondition(ctx, scope, metav1.ConditionFalse, "SyncFailed", syncErr.Error())
		return ctrl.Result{RequeueAfter: scopeRequeueAfter}, nil
	}

	patch := client.MergeFrom(scope.DeepCopy())
	scope.Status.NamespaceCount = len(namespaces)
	scope.Status.ExplorerCount = totalExplorers
	scope.Status.ObservedGeneration = scope.Generation
	apimeta.SetStatusCondition(&scope.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "Reconciled",
		Message:            fmt.Sprintf("Reconciled %d namespaces, %d explorers", len(namespaces), totalExplorers),
		ObservedGeneration: scope.Generation,
	})
	if err := r.Status().Patch(ctx, scope, patch); err != nil {
		log.Error(err, "Could not update PVCExplorerScope status")
	}

	log.Info("Reconciled PVCExplorerScope", "name", scope.Name, "namespaces", len(namespaces), "explorers", totalExplorers)
	return ctrl.Result{RequeueAfter: scopeRequeueAfter}, nil
}

func (r *PVCExplorerScopeReconciler) resolveNamespaces(ctx context.Context, scope *pvcexplorerv1alpha1.PVCExplorerScope) ([]corev1.Namespace, error) {
	nameSet := map[string]struct{}{}
	for _, ns := range scope.Spec.Namespaces.Names {
		nameSet[ns] = struct{}{}
	}
	if scope.Spec.Namespaces.LabelSelector != nil {
		sel, err := metav1.LabelSelectorAsSelector(scope.Spec.Namespaces.LabelSelector)
		if err != nil {
			return nil, err
		}
		nsList := &corev1.NamespaceList{}
		if err := r.List(ctx, nsList, &client.ListOptions{LabelSelector: sel}); err != nil {
			return nil, err
		}
		for _, ns := range nsList.Items {
			nameSet[ns.Name] = struct{}{}
		}
	}
	var result []corev1.Namespace
	for ns := range nameSet {
		obj := corev1.Namespace{}
		if err := r.Get(ctx, types.NamespacedName{Name: ns}, &obj); err == nil {
			result = append(result, obj)
		}
	}
	return result, nil
}

func (r *PVCExplorerScopeReconciler) discoverPVCs(ctx context.Context, scope *pvcexplorerv1alpha1.PVCExplorerScope, namespaces []corev1.Namespace) (map[string][]corev1.PersistentVolumeClaim, error) {
	pvcMap := make(map[string][]corev1.PersistentVolumeClaim)
	mode := scope.Spec.Discovery.Mode
	explicitSet := map[string]struct{}{}
	if mode == pvcexplorerv1alpha1.DiscoveryModeExplicit {
		for _, n := range scope.Spec.Discovery.PVCNames {
			explicitSet[n] = struct{}{}
		}
	}

	for _, ns := range namespaces {
		pvcList := &corev1.PersistentVolumeClaimList{}
		if err := r.List(ctx, pvcList, client.InNamespace(ns.Name)); err != nil {
			return nil, err
		}
		var filtered []corev1.PersistentVolumeClaim
		for _, pvc := range pvcList.Items {
			if mode == pvcexplorerv1alpha1.DiscoveryModeExplicit {
				if _, ok := explicitSet[pvc.Name]; !ok {
					continue
				}
			}
			if matchesAnyGlob(pvc.Name, scope.Spec.Discovery.ExcludePVCs) {
				continue
			}
			filtered = append(filtered, pvc)
		}
		if len(filtered) > 0 {
			pvcMap[ns.Name] = filtered
		}
	}
	return pvcMap, nil
}

func matchesAnyGlob(name string, patterns []string) bool {
	for _, pat := range patterns {
		if matched, err := path.Match(pat, name); err == nil && matched {
			return true
		}
	}
	return false
}

func (r *PVCExplorerScopeReconciler) syncPVCExplorers(ctx context.Context, scope *pvcexplorerv1alpha1.PVCExplorerScope, pvcMap map[string][]corev1.PersistentVolumeClaim) error {
	log := logf.FromContext(ctx)

	allExplorers := &pvcexplorerv1alpha1.PVCExplorerList{}
	if err := r.List(ctx, allExplorers, client.MatchingLabels{scopeLabelKey: scope.Name}); err != nil {
		return err
	}
	existing := map[string]*pvcexplorerv1alpha1.PVCExplorer{}
	for i := range allExplorers.Items {
		e := &allExplorers.Items[i]
		existing[e.Namespace+"/"+e.Name] = e
	}

	desired := map[string]corev1.PersistentVolumeClaim{}
	for ns, pvcs := range pvcMap {
		for _, pvc := range pvcs {
			desired[ns+"/"+pvc.Name] = pvc
		}
	}

	for key, pvc := range desired {
		if _, found := existing[key]; !found {
			explorer := r.buildExplorer(scope, pvc)
			if err := r.Create(ctx, explorer); err != nil {
				return fmt.Errorf("creating PVCExplorer %s: %w", key, err)
			}
			log.Info("Created PVCExplorer", "name", pvc.Name, "namespace", pvc.Namespace)
		}
	}

	if scope.Spec.DeletionPolicy == pvcexplorerv1alpha1.DeletionPolicyCleanup {
		for key, e := range existing {
			if _, ok := desired[key]; !ok {
				if err := r.Delete(ctx, e); err != nil {
					return fmt.Errorf("deleting PVCExplorer %s: %w", key, err)
				}
				log.Info("Deleted PVCExplorer", "name", e.Name, "namespace", e.Namespace)
			}
		}
	}
	return nil
}

func (r *PVCExplorerScopeReconciler) buildExplorer(scope *pvcexplorerv1alpha1.PVCExplorerScope, pvc corev1.PersistentVolumeClaim) *pvcexplorerv1alpha1.PVCExplorer {
	d := scope.Spec.Defaults
	return &pvcexplorerv1alpha1.PVCExplorer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvc.Name,
			Namespace: pvc.Namespace,
			Labels: map[string]string{
				scopeLabelKey:     scope.Name,
				managedByLabelKey: managedByValue,
			},
		},
		Spec: pvcexplorerv1alpha1.PVCExplorerSpec{
			PVCName:       pvc.Name,
			Mode:          d.Mode,
			Image:         d.Image,
			ForceRW:       d.ForceRW,
			Scaling:       pvcexplorerv1alpha1.ScalingSpec{IdleTimeout: d.Scaling.IdleTimeout, StartupTimeout: d.Scaling.StartupTimeout},
			MountStrategy: pvcexplorerv1alpha1.MountStrategySpec{AllowNodeAffinity: d.MountStrategy.AllowNodeAffinity, FallbackOnConflict: d.MountStrategy.FallbackOnConflict},
			Resources:     d.Resources,
		},
	}
}

func (r *PVCExplorerScopeReconciler) setReadyCondition(ctx context.Context, scope *pvcexplorerv1alpha1.PVCExplorerScope, status metav1.ConditionStatus, reason, msg string) {
	patch := client.MergeFrom(scope.DeepCopy())
	apimeta.SetStatusCondition(&scope.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             status,
		Reason:             reason,
		Message:            msg,
		ObservedGeneration: scope.Generation,
	})
	_ = r.Status().Patch(ctx, scope, patch)
}

func (r *PVCExplorerScopeReconciler) reconcileDelete(ctx context.Context, scope *pvcexplorerv1alpha1.PVCExplorerScope) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if scope.Spec.DeletionPolicy == pvcexplorerv1alpha1.DeletionPolicyCleanup {
		explorers := &pvcexplorerv1alpha1.PVCExplorerList{}
		if err := r.List(ctx, explorers, client.MatchingLabels{scopeLabelKey: scope.Name}); err != nil {
			return ctrl.Result{}, err
		}
		for i := range explorers.Items {
			e := &explorers.Items[i]
			if e.DeletionTimestamp.IsZero() {
				if err := r.Delete(ctx, e); err != nil {
					return ctrl.Result{}, fmt.Errorf("deleting PVCExplorer %s/%s: %w", e.Namespace, e.Name, err)
				}
				log.Info("Deleted PVCExplorer during scope cleanup", "name", e.Name, "namespace", e.Namespace)
			}
		}
		if len(explorers.Items) > 0 {
			return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
		}
	}

	patch := client.MergeFrom(scope.DeepCopy())
	controllerutil.RemoveFinalizer(scope, scopeFinalizer)
	return ctrl.Result{}, r.Patch(ctx, scope, patch)
}

func (r *PVCExplorerScopeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pvcexplorerv1alpha1.PVCExplorerScope{}).
		Watches(
			&corev1.Namespace{},
			handler.EnqueueRequestsFromMapFunc(r.namespaceToScopes),
		).
		Watches(
			&corev1.PersistentVolumeClaim{},
			handler.EnqueueRequestsFromMapFunc(r.pvcToScopes),
		).
		Named("scope").
		Complete(r)
}

func (r *PVCExplorerScopeReconciler) namespaceToScopes(ctx context.Context, obj client.Object) []reconcile.Request {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return nil
	}
	scopeList := &pvcexplorerv1alpha1.PVCExplorerScopeList{}
	if err := r.List(ctx, scopeList); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, scope := range scopeList.Items {
		if selectsNamespace(&scope, ns) {
			reqs = append(reqs, reconcile.Request{NamespacedName: types.NamespacedName{Name: scope.Name}})
		}
	}
	return reqs
}

func (r *PVCExplorerScopeReconciler) pvcToScopes(ctx context.Context, obj client.Object) []reconcile.Request {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil
	}
	scopeList := &pvcexplorerv1alpha1.PVCExplorerScopeList{}
	if err := r.List(ctx, scopeList); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, scope := range scopeList.Items {
		if selectsNamespace(&scope, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: pvc.Namespace}}) {
			reqs = append(reqs, reconcile.Request{NamespacedName: types.NamespacedName{Name: scope.Name}})
		}
	}
	return reqs
}

func selectsNamespace(scope *pvcexplorerv1alpha1.PVCExplorerScope, ns *corev1.Namespace) bool {
	if slices.Contains(scope.Spec.Namespaces.Names, ns.Name) {
		return true
	}
	if scope.Spec.Namespaces.LabelSelector != nil {
		sel, err := metav1.LabelSelectorAsSelector(scope.Spec.Namespaces.LabelSelector)
		if err == nil && sel.Matches(labels.Set(ns.Labels)) {
			return true
		}
	}
	return false
}
