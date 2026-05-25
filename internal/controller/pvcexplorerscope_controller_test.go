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
	context "context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
)

var _ = Describe("PVCExplorerScope Controller", func() {
	// Utility to create a unique namespace for each test
	createTestNamespace := func() string {
		ns := fmt.Sprintf("test-ns-%s", uuid.NewString())
		nsObj := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
		Expect(k8sClient.Create(context.Background(), nsObj)).To(Succeed())
		return ns
	}
	// Utility to delete a namespace
	deleteTestNamespace := func(ns string) {
		_ = k8sClient.Delete(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}})
	}

	It("creates PVCExplorer for explicit PVC in explicit namespace", func() {
		ns := createTestNamespace()
		DeferCleanup(func() { deleteTestNamespace(ns) })
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "mypvc", Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{"storage": resourceMustParse("1Gi")},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), pvc)).To(Succeed())
		scope := &pvcexplorerv1alpha1.PVCExplorerScope{
			ObjectMeta: metav1.ObjectMeta{Name: "scope-explicit"},
			Spec: pvcexplorerv1alpha1.PVCExplorerScopeSpec{
				Namespaces: pvcexplorerv1alpha1.ScopeNamespacesSpec{Names: []string{ns}},
				Discovery: pvcexplorerv1alpha1.ScopeDiscoverySpec{
					Mode:     pvcexplorerv1alpha1.DiscoveryModeExplicit,
					PVCNames: []string{"mypvc"},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), scope)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(context.Background(), scope) })
		Eventually(func() bool {
			return pvcExplorerExists(ns, "mypvc")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
	})

	It("creates PVCExplorers for all PVCs in namespace in Auto mode", func() {
		ns := createTestNamespace()
		DeferCleanup(func() { deleteTestNamespace(ns) })
		pvc1 := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "pvc1", Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{"storage": resourceMustParse("1Gi")},
				},
			},
		}
		pvc2 := pvc1.DeepCopy()
		pvc2.Name = "pvc2"
		Expect(k8sClient.Create(context.Background(), pvc1)).To(Succeed())
		Expect(k8sClient.Create(context.Background(), pvc2)).To(Succeed())
		scope := &pvcexplorerv1alpha1.PVCExplorerScope{
			ObjectMeta: metav1.ObjectMeta{Name: "scope-auto"},
			Spec: pvcexplorerv1alpha1.PVCExplorerScopeSpec{
				Namespaces: pvcexplorerv1alpha1.ScopeNamespacesSpec{Names: []string{ns}},
				Discovery:  pvcexplorerv1alpha1.ScopeDiscoverySpec{Mode: pvcexplorerv1alpha1.DiscoveryModeAuto},
			},
		}
		Expect(k8sClient.Create(context.Background(), scope)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(context.Background(), scope) })
		Eventually(func() bool {
			return pvcExplorerExists(ns, "pvc1") && pvcExplorerExists(ns, "pvc2")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
	})

	It("excludes PVCs matching excludePVCs glob", func() {
		ns := createTestNamespace()
		DeferCleanup(func() { deleteTestNamespace(ns) })
		pvc1 := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "foo-1", Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{"storage": resourceMustParse("1Gi")},
				},
			},
		}
		pvc2 := pvc1.DeepCopy()
		pvc2.Name = "bar-2"
		Expect(k8sClient.Create(context.Background(), pvc1)).To(Succeed())
		Expect(k8sClient.Create(context.Background(), pvc2)).To(Succeed())
		scope := &pvcexplorerv1alpha1.PVCExplorerScope{
			ObjectMeta: metav1.ObjectMeta{Name: "scope-exclude"},
			Spec: pvcexplorerv1alpha1.PVCExplorerScopeSpec{
				Namespaces: pvcexplorerv1alpha1.ScopeNamespacesSpec{Names: []string{ns}},
				Discovery: pvcexplorerv1alpha1.ScopeDiscoverySpec{
					Mode:        pvcexplorerv1alpha1.DiscoveryModeAuto,
					ExcludePVCs: []string{"foo-*"},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), scope)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(context.Background(), scope) })
		Eventually(func() bool {
			return !pvcExplorerExists(ns, "foo-1") && pvcExplorerExists(ns, "bar-2")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
	})

	It("deletes orphaned PVCExplorers when deletionPolicy is Cleanup", func() {
		ns := createTestNamespace()
		DeferCleanup(func() { deleteTestNamespace(ns) })
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "todelete", Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{"storage": resourceMustParse("1Gi")},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), pvc)).To(Succeed())
		scope := &pvcexplorerv1alpha1.PVCExplorerScope{
			ObjectMeta: metav1.ObjectMeta{Name: "scope-cleanup"},
			Spec: pvcexplorerv1alpha1.PVCExplorerScopeSpec{
				Namespaces: pvcexplorerv1alpha1.ScopeNamespacesSpec{Names: []string{ns}},
				Discovery: pvcexplorerv1alpha1.ScopeDiscoverySpec{
					Mode:     pvcexplorerv1alpha1.DiscoveryModeExplicit,
					PVCNames: []string{"todelete"},
				},
				DeletionPolicy: pvcexplorerv1alpha1.DeletionPolicyCleanup,
			},
		}
		Expect(k8sClient.Create(context.Background(), scope)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(context.Background(), scope) })
		Eventually(func() bool {
			return pvcExplorerExists(ns, "todelete")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
		Eventually(func() error {
			if err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "scope-cleanup"}, scope); err != nil {
				return err
			}
			updated := scope.DeepCopy()
			updated.Spec.Discovery.PVCNames = []string{}
			return k8sClient.Update(context.Background(), updated)
		}, time.Second*5, time.Millisecond*200).Should(Succeed())
		Eventually(func() bool {
			return !pvcExplorerExists(ns, "todelete")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
	})

	It("orphans PVCExplorers when deletionPolicy is Orphan", func() {
		ns := createTestNamespace()
		DeferCleanup(func() { deleteTestNamespace(ns) })
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "orphanme", Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{"storage": resourceMustParse("1Gi")},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), pvc)).To(Succeed())
		scope := &pvcexplorerv1alpha1.PVCExplorerScope{
			ObjectMeta: metav1.ObjectMeta{Name: "scope-orphan"},
			Spec: pvcexplorerv1alpha1.PVCExplorerScopeSpec{
				Namespaces: pvcexplorerv1alpha1.ScopeNamespacesSpec{Names: []string{ns}},
				Discovery: pvcexplorerv1alpha1.ScopeDiscoverySpec{
					Mode:     pvcexplorerv1alpha1.DiscoveryModeExplicit,
					PVCNames: []string{"orphanme"},
				},
				DeletionPolicy: pvcexplorerv1alpha1.DeletionPolicyOrphan,
			},
		}
		Expect(k8sClient.Create(context.Background(), scope)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(context.Background(), scope) })
		Eventually(func() bool {
			return pvcExplorerExists(ns, "orphanme")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
		Eventually(func() error {
			if err := k8sClient.Get(context.Background(), types.NamespacedName{Name: "scope-orphan"}, scope); err != nil {
				return err
			}
			updated := scope.DeepCopy()
			updated.Spec.Discovery.PVCNames = []string{}
			return k8sClient.Update(context.Background(), updated)
		}, time.Second*5, time.Millisecond*200).Should(Succeed())
		Consistently(func() bool {
			return pvcExplorerExists(ns, "orphanme")
		}, time.Second*3, time.Millisecond*200).Should(BeTrue())
	})

	It("selects namespaces by labelSelector", func() {
		labelKey := "test-label"
		labelValue := "scope"
		ns := createTestNamespace()
		DeferCleanup(func() { deleteTestNamespace(ns) })
		// Patch label
		patch := clientMergePatch(map[string]any{"metadata": map[string]any{"labels": map[string]any{labelKey: labelValue}}})
		Expect(k8sClient.Patch(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, patch)).To(Succeed())
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: "lblpvc", Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{"storage": resourceMustParse("1Gi")},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), pvc)).To(Succeed())
		selector := metav1.LabelSelector{MatchLabels: map[string]string{labelKey: labelValue}}
		scope := &pvcexplorerv1alpha1.PVCExplorerScope{
			ObjectMeta: metav1.ObjectMeta{Name: "scope-labelsel"},
			Spec: pvcexplorerv1alpha1.PVCExplorerScopeSpec{
				Namespaces: pvcexplorerv1alpha1.ScopeNamespacesSpec{LabelSelector: &selector},
				Discovery: pvcexplorerv1alpha1.ScopeDiscoverySpec{
					Mode: pvcexplorerv1alpha1.DiscoveryModeAuto,
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), scope)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(context.Background(), scope) })
		Eventually(func() bool {
			return pvcExplorerExists(ns, "lblpvc")
		}, time.Second*10, time.Millisecond*200).Should(BeTrue())
	})
})

// Helper: does a PVCExplorer exist for ns/pvc?
func pvcExplorerExists(ns, pvc string) bool {
	list := &pvcexplorerv1alpha1.PVCExplorerList{}
	err := k8sClient.List(context.Background(), list)
	if err != nil {
		return false
	}
	for _, e := range list.Items {
		if e.Namespace == ns && e.Name == pvc {
			return true
		}
	}
	return false
}

// Helper: parse resource quantity
func resourceMustParse(s string) resource.Quantity {
	return resource.MustParse(s)
}

// Helper: create a merge patch
func clientMergePatch(obj map[string]any) client.Patch {
	b, _ := json.Marshal(obj)
	return client.RawPatch(types.MergePatchType, b)
}
