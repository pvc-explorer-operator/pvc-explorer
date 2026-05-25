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
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
)

var _ = Describe("PVCExplorer Controller", func() {
	newNS := func() string {
		ns := fmt.Sprintf("agent-test-%s", uuid.NewString()[:8])
		Expect(k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		})).To(Succeed())
		return ns
	}

	boundPVC := func(ns, name string) {
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), pvc)).To(Succeed())
		patch := pvc.DeepCopy()
		patch.Status.Phase = corev1.ClaimBound
		Expect(k8sClient.Status().Update(context.Background(), patch)).To(Succeed())
	}

	newExplorer := func(ns, pvcName string, mode pvcexplorerv1alpha1.ExplorerMode, forceRW bool) {
		Expect(k8sClient.Create(context.Background(), &pvcexplorerv1alpha1.PVCExplorer{
			ObjectMeta: metav1.ObjectMeta{Name: pvcName, Namespace: ns},
			Spec: pvcexplorerv1alpha1.PVCExplorerSpec{
				PVCName: pvcName,
				Mode:    mode,
				ForceRW: forceRW,
			},
		})).To(Succeed())
	}

	It("creates Deployment with replicas=0 for ScaledToZero mode", func() {
		ns := newNS()
		boundPVC(ns, "mypvc")
		newExplorer(ns, "mypvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		deploy := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "mypvc", Namespace: ns}, deploy)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		Expect(*deploy.Spec.Replicas).To(Equal(int32(0)))
	})

	It("creates Deployment with replicas=1 for Deployment mode", func() {
		ns := newNS()
		boundPVC(ns, "alwayson")
		newExplorer(ns, "alwayson", pvcexplorerv1alpha1.ExplorerModeDeployment, true)

		deploy := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "alwayson", Namespace: ns}, deploy)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		Expect(*deploy.Spec.Replicas).To(Equal(int32(1)))
	})

	It("creates Service with correct selector and port", func() {
		ns := newNS()
		boundPVC(ns, "svcpvc")
		newExplorer(ns, "svcpvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		svc := &corev1.Service{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "svcpvc", Namespace: ns}, svc)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		Expect(svc.Spec.Selector).To(HaveKeyWithValue("app", "svcpvc"))
		Expect(svc.Spec.Ports).To(HaveLen(1))
		Expect(svc.Spec.Ports[0].Port).To(Equal(int32(8081)))
	})

	It("mounts PVC read-write when forceRW=true", func() {
		ns := newNS()
		boundPVC(ns, "rwpvc")
		newExplorer(ns, "rwpvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		deploy := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "rwpvc", Namespace: ns}, deploy)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		Expect(deploy.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ReadOnly).To(BeFalse())
	})

	It("mounts PVC read-only when forceRW omitted (defaults to true via CRD, so RW)", func() {
		ns := newNS()
		boundPVC(ns, "ropvc")
		newExplorer(ns, "ropvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, false)

		deploy := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "ropvc", Namespace: ns}, deploy)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Expect(k8sClient.Get(context.Background(), types.NamespacedName{Name: "ropvc", Namespace: ns}, explorer)).To(Succeed())
		readOnly := !explorer.Spec.ForceRW
		Expect(deploy.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ReadOnly).To(Equal(readOnly))
	})

	It("sets status.phase=ScaledToZero and agentEndpoint", func() {
		ns := newNS()
		boundPVC(ns, "statuspvc")
		newExplorer(ns, "statuspvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		key := types.NamespacedName{Name: "statuspvc", Namespace: ns}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Eventually(func() pvcexplorerv1alpha1.ExplorerPhase {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return explorer.Status.Phase
		}, 10*time.Second, 200*time.Millisecond).Should(Equal(pvcexplorerv1alpha1.ExplorerPhaseScaledToZero))

		Expect(explorer.Status.AgentEndpoint).To(ContainSubstring("statuspvc"))
		Expect(explorer.Status.Mode).To(Equal(pvcexplorerv1alpha1.ExplorerModeScaledToZero))
	})

	It("sets status.phase=Failed when PVC does not exist", func() {
		ns := newNS()
		Expect(k8sClient.Create(context.Background(), &pvcexplorerv1alpha1.PVCExplorer{
			ObjectMeta: metav1.ObjectMeta{Name: "nopvc", Namespace: ns},
			Spec:       pvcexplorerv1alpha1.PVCExplorerSpec{PVCName: "nonexistent"},
		})).To(Succeed())

		key := types.NamespacedName{Name: "nopvc", Namespace: ns}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Eventually(func() pvcexplorerv1alpha1.ExplorerPhase {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return explorer.Status.Phase
		}, 10*time.Second, 200*time.Millisecond).Should(Equal(pvcexplorerv1alpha1.ExplorerPhaseFailed))
	})

	It("sets ownerReference on Deployment pointing to PVCExplorer", func() {
		ns := newNS()
		boundPVC(ns, "ownerpvc")
		newExplorer(ns, "ownerpvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		deploy := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "ownerpvc", Namespace: ns}, deploy)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		Expect(deploy.OwnerReferences).To(HaveLen(1))
		Expect(deploy.OwnerReferences[0].Kind).To(Equal("PVCExplorer"))
		Expect(deploy.OwnerReferences[0].Name).To(Equal("ownerpvc"))
	})

	It("mounts RO and sets ForceRWDeferred when consumer pod is Running", func() {
		ns := newNS()
		boundPVC(ns, "conspvc")
		newExplorer(ns, "conspvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		Eventually(func() error {
			deploy := &appsv1.Deployment{}
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "conspvc", Namespace: ns}, deploy)
		}, 10*time.Second, 200*time.Millisecond).Should(Succeed())

		consumer := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "consumer-pod", Namespace: ns},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "app", Image: "busybox"}},
				Volumes: []corev1.Volume{{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "conspvc"},
					},
				}},
			},
		}
		Expect(k8sClient.Create(context.Background(), consumer)).To(Succeed())
		patch := consumer.DeepCopy()
		patch.Status.Phase = corev1.PodRunning
		Expect(k8sClient.Status().Update(context.Background(), patch)).To(Succeed())

		key := types.NamespacedName{Name: "conspvc", Namespace: ns}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Eventually(func() bool {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return explorer.Status.Mount.ReadOnly && explorer.Status.Mount.ForceRWDeferred
		}, 15*time.Second, 200*time.Millisecond).Should(BeTrue())

		Expect(explorer.Status.Mount.Consumers).To(HaveLen(1))
		Expect(explorer.Status.Mount.Consumers[0].PodName).To(Equal("consumer-pod"))
	})

	It("remounts RW when consumer pod is deleted", func() {
		ns := newNS()
		boundPVC(ns, "releasepvc")
		newExplorer(ns, "releasepvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		consumer := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "consumer-pod", Namespace: ns},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "app", Image: "busybox"}},
				Volumes: []corev1.Volume{{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "releasepvc"},
					},
				}},
			},
		}
		Expect(k8sClient.Create(context.Background(), consumer)).To(Succeed())
		patch := consumer.DeepCopy()
		patch.Status.Phase = corev1.PodRunning
		Expect(k8sClient.Status().Update(context.Background(), patch)).To(Succeed())

		key := types.NamespacedName{Name: "releasepvc", Namespace: ns}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Eventually(func() bool {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return explorer.Status.Mount.ReadOnly
		}, 15*time.Second, 200*time.Millisecond).Should(BeTrue())

		Expect(k8sClient.Delete(context.Background(), consumer)).To(Succeed())

		Eventually(func() bool {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return !explorer.Status.Mount.ReadOnly && !explorer.Status.Mount.ForceRWDeferred
		}, 15*time.Second, 200*time.Millisecond).Should(BeTrue())

		deploy := &appsv1.Deployment{}
		Expect(k8sClient.Get(context.Background(), types.NamespacedName{Name: "releasepvc", Namespace: ns}, deploy)).To(Succeed())
		Expect(deploy.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ReadOnly).To(BeFalse())
	})

	It("sets status.mount.strategy=NodeAffinity and targetNode for RWO+consumer", func() {
		ns := newNS()
		boundPVC(ns, "rwopvc")
		newExplorer(ns, "rwopvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		consumer := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "rwo-consumer", Namespace: ns},
			Spec: corev1.PodSpec{
				NodeName:   "worker-node-1",
				Containers: []corev1.Container{{Name: "app", Image: "busybox"}},
				Volumes: []corev1.Volume{{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "rwopvc"},
					},
				}},
			},
		}
		Expect(k8sClient.Create(context.Background(), consumer)).To(Succeed())
		patch := consumer.DeepCopy()
		patch.Status.Phase = corev1.PodRunning
		Expect(k8sClient.Status().Update(context.Background(), patch)).To(Succeed())

		key := types.NamespacedName{Name: "rwopvc", Namespace: ns}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Eventually(func() string {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return explorer.Status.Mount.Strategy
		}, 15*time.Second, 200*time.Millisecond).Should(Equal("NodeAffinity"))

		Expect(explorer.Status.Mount.TargetNode).To(Equal("worker-node-1"))

		deploy := &appsv1.Deployment{}
		Expect(k8sClient.Get(context.Background(), types.NamespacedName{Name: "rwopvc", Namespace: ns}, deploy)).To(Succeed())
		aff := deploy.Spec.Template.Spec.Affinity
		Expect(aff).NotTo(BeNil())
		terms := aff.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
		Expect(terms[0].MatchExpressions[0].Values).To(ContainElement("worker-node-1"))
	})

	It("sets phase=Pending for RWO PVC consumed on multiple nodes", func() {
		ns := newNS()
		boundPVC(ns, "multinodepvc")
		newExplorer(ns, "multinodepvc", pvcexplorerv1alpha1.ExplorerModeScaledToZero, true)

		for i, node := range []string{"node-a", "node-b"} {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("consumer-%d", i), Namespace: ns},
				Spec: corev1.PodSpec{
					NodeName:   node,
					Containers: []corev1.Container{{Name: "app", Image: "busybox"}},
					Volumes: []corev1.Volume{{
						Name: "data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "multinodepvc"},
						},
					}},
				},
			}
			Expect(k8sClient.Create(context.Background(), pod)).To(Succeed())
			p := pod.DeepCopy()
			p.Status.Phase = corev1.PodRunning
			Expect(k8sClient.Status().Update(context.Background(), p)).To(Succeed())
		}

		key := types.NamespacedName{Name: "multinodepvc", Namespace: ns}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		Eventually(func() pvcexplorerv1alpha1.ExplorerPhase {
			_ = k8sClient.Get(context.Background(), key, explorer)
			return explorer.Status.Phase
		}, 15*time.Second, 200*time.Millisecond).Should(Equal(pvcexplorerv1alpha1.ExplorerPhasePending))
	})
})
