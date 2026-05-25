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
	"maps"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
)

// indexKey identifies a PVC within a namespace.
type indexKey struct {
	namespace string
	pvcName   string
}

type consumerBroadcaster interface {
	Publish(eventType string, payload any) error
}

// Index is a thread-safe in-memory cache of active consumers per PVC.
// It is updated by the pod watch handler and read by the agent reconciler,
// eliminating per-reconcile full pod list scans.
type Index struct {
	mu          sync.RWMutex
	store       map[indexKey][]pvcexplorerv1alpha1.ConsumerInfo
	broadcaster consumerBroadcaster
}

// NewIndex creates an empty Index.
func NewIndex() *Index {
	return &Index{store: make(map[indexKey][]pvcexplorerv1alpha1.ConsumerInfo)}
}

func NewIndexWithBroadcaster(b consumerBroadcaster) *Index {
	return &Index{
		store:       make(map[indexKey][]pvcexplorerv1alpha1.ConsumerInfo),
		broadcaster: b,
	}
}

// Get returns the current consumers for a PVC.
func (idx *Index) Get(namespace, pvcName string) []pvcexplorerv1alpha1.ConsumerInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	v := idx.store[indexKey{namespace, pvcName}]
	if v == nil {
		return nil
	}
	out := make([]pvcexplorerv1alpha1.ConsumerInfo, len(v))
	copy(out, v)
	return out
}

// Sync re-scans all pods in namespace and refreshes the index for every PVC
// they reference. It is called by the pod watch handler on each relevant event.
func (idx *Index) Sync(ctx context.Context, c client.Client, namespace string) error {
	podList := &corev1.PodList{}
	if err := c.List(ctx, podList, client.InNamespace(namespace)); err != nil {
		return err
	}

	updated := map[indexKey][]pvcexplorerv1alpha1.ConsumerInfo{}
	for i := range podList.Items {
		pod := &podList.Items[i]
		if pod.Status.Phase != corev1.PodRunning && pod.Status.Phase != corev1.PodPending {
			continue
		}
		for _, v := range pod.Spec.Volumes {
			if v.PersistentVolumeClaim == nil {
				continue
			}
			pvcName := v.PersistentVolumeClaim.ClaimName
			key := indexKey{namespace, pvcName}
			info := pvcexplorerv1alpha1.ConsumerInfo{
				PodName:       pod.Name,
				NodeName:      pod.Spec.NodeName,
				MountReadOnly: isMountedReadOnly(pod, pvcName),
			}
			info.OwnerKind, info.OwnerName = resolveOwnerChain(ctx, c, pod)
			updated[key] = append(updated[key], info)
		}
	}

	idx.mu.Lock()
	var attached, detached []consumerEvent
	if idx.broadcaster != nil {
		allKeys := make(map[indexKey]struct{})
		for k := range idx.store {
			if k.namespace == namespace {
				allKeys[k] = struct{}{}
			}
		}
		for k := range updated {
			allKeys[k] = struct{}{}
		}
		for k := range allKeys {
			attached = append(attached, diffAdded(k, idx.store[k], updated[k])...)
			detached = append(detached, diffRemoved(k, idx.store[k], updated[k])...)
		}
	}
	for k := range idx.store {
		if k.namespace == namespace {
			delete(idx.store, k)
		}
	}
	maps.Copy(idx.store, updated)
	idx.mu.Unlock()

	if idx.broadcaster != nil {
		for _, ev := range attached {
			_ = idx.broadcaster.Publish("consumer.attached", consumerPayload(ev))
		}
		for _, ev := range detached {
			_ = idx.broadcaster.Publish("consumer.detached", consumerPayload(ev))
		}
	}

	return nil
}

type consumerEvent struct {
	key  indexKey
	info pvcexplorerv1alpha1.ConsumerInfo
}

func consumerPayload(ev consumerEvent) map[string]any {
	return map[string]any{
		"namespace": ev.key.namespace,
		"pvcName":   ev.key.pvcName,
		"podName":   ev.info.PodName,
		"nodeName":  ev.info.NodeName,
		"ownerKind": ev.info.OwnerKind,
		"ownerName": ev.info.OwnerName,
		"readOnly":  ev.info.MountReadOnly,
	}
}

func diffAdded(k indexKey, old, neu []pvcexplorerv1alpha1.ConsumerInfo) []consumerEvent {
	oldSet := podNameSet(old)
	var out []consumerEvent
	for _, info := range neu {
		if !oldSet[info.PodName] {
			out = append(out, consumerEvent{key: k, info: info})
		}
	}
	return out
}

func diffRemoved(k indexKey, old, neu []pvcexplorerv1alpha1.ConsumerInfo) []consumerEvent {
	newSet := podNameSet(neu)
	var out []consumerEvent
	for _, info := range old {
		if !newSet[info.PodName] {
			out = append(out, consumerEvent{key: k, info: info})
		}
	}
	return out
}

func podNameSet(consumers []pvcexplorerv1alpha1.ConsumerInfo) map[string]bool {
	s := make(map[string]bool, len(consumers))
	for _, c := range consumers {
		s[c.PodName] = true
	}
	return s
}
