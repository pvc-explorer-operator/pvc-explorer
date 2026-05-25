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

package consumer

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
)

const (
	KindReplicaSet = "ReplicaSet"
	KindJob        = "Job"
)

// Detect returns ConsumerInfo for every Running/Pending pod in namespace that mounts pvcName.
// Owner chain is fully resolved: Pod→RS→Deployment, Pod→Job→CronJob.
func Detect(ctx context.Context, c client.Client, namespace, pvcName string) ([]pvcexplorerv1alpha1.ConsumerInfo, error) {
	podList := &corev1.PodList{}
	if err := c.List(ctx, podList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	var consumers []pvcexplorerv1alpha1.ConsumerInfo
	for i := range podList.Items {
		pod := &podList.Items[i]
		if pod.Status.Phase != corev1.PodRunning && pod.Status.Phase != corev1.PodPending {
			continue
		}
		if !MountsPVC(pod, pvcName) {
			continue
		}
		info := pvcexplorerv1alpha1.ConsumerInfo{
			PodName:       pod.Name,
			NodeName:      pod.Spec.NodeName,
			MountReadOnly: isMountedReadOnly(pod, pvcName),
		}
		info.OwnerKind, info.OwnerName = resolveOwnerChain(ctx, c, pod)
		consumers = append(consumers, info)
	}
	return consumers, nil
}

// MountsPVC reports whether the pod has a volume backed by pvcName.
func MountsPVC(pod *corev1.Pod, pvcName string) bool {
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim != nil && v.PersistentVolumeClaim.ClaimName == pvcName {
			return true
		}
	}
	return false
}

func isMountedReadOnly(pod *corev1.Pod, pvcName string) bool {
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim != nil && v.PersistentVolumeClaim.ClaimName == pvcName {
			return v.PersistentVolumeClaim.ReadOnly
		}
	}
	return false
}

// resolveOwnerChain walks Pod→RS→Deployment and Pod→Job→CronJob via API lookups.
func resolveOwnerChain(ctx context.Context, c client.Client, pod *corev1.Pod) (kind, name string) {
	if len(pod.OwnerReferences) == 0 {
		return "Pod", pod.Name
	}
	ref := pod.OwnerReferences[0]
	ns := pod.Namespace

	switch ref.Kind {
	case KindReplicaSet:
		rs := &appsv1.ReplicaSet{}
		if err := c.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: ns}, rs); err != nil {
			return KindReplicaSet, ref.Name
		}
		if len(rs.OwnerReferences) > 0 && rs.OwnerReferences[0].Kind == "Deployment" {
			return "Deployment", rs.OwnerReferences[0].Name
		}
		return KindReplicaSet, ref.Name

	case "StatefulSet", "DaemonSet":
		return ref.Kind, ref.Name

	case KindJob:
		job := &batchv1.Job{}
		if err := c.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: ns}, job); err != nil {
			return KindJob, ref.Name
		}
		if len(job.OwnerReferences) > 0 && job.OwnerReferences[0].Kind == "CronJob" {
			return "CronJob", job.OwnerReferences[0].Name
		}
		return KindJob, ref.Name

	default:
		return ref.Kind, ref.Name
	}
}
