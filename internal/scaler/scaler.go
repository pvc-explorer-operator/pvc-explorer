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

package scaler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	pvcv1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
)

const (
	AnnotationIdleDeadline = "pvcexplorer.io/idle-deadline"
	defaultIdleTimeout     = 10 * time.Minute
	idleWatchInterval      = 5 * time.Second
	warningThreshold       = 60 * time.Second
)

type Broadcaster interface {
	Publish(eventType string, payload any) error
}

type idleWarnState struct {
	warned bool
}

type Scaler struct {
	Client      client.Client
	broadcaster Broadcaster

	mu        sync.Mutex
	warnState map[string]*idleWarnState
}

func New(c client.Client, b Broadcaster) *Scaler {
	return &Scaler{
		Client:      c,
		broadcaster: b,
		warnState:   make(map[string]*idleWarnState),
	}
}

func (s *Scaler) WakeAgent(ctx context.Context, namespace, name string) error {
	explorer := &pvcv1.PVCExplorer{}
	if err := s.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, explorer); err != nil {
		return err
	}

	patch := client.MergeFrom(explorer.DeepCopy())
	explorer.Spec.Mode = pvcv1.ExplorerModeDeployment
	if explorer.Annotations == nil {
		explorer.Annotations = map[string]string{}
	}
	explorer.Annotations[AnnotationIdleDeadline] = time.Now().Add(idleTimeout(explorer)).UTC().Format(time.RFC3339)
	if err := s.Client.Patch(ctx, explorer, patch); err != nil {
		return err
	}
	s.resetWarnState(namespace, name)
	return nil
}

func (s *Scaler) SleepAgent(ctx context.Context, namespace, name string) error {
	explorer := &pvcv1.PVCExplorer{}
	if err := s.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, explorer); err != nil {
		return err
	}

	patch := client.MergeFrom(explorer.DeepCopy())
	explorer.Spec.Mode = pvcv1.ExplorerModeScaledToZero
	if explorer.Annotations != nil {
		delete(explorer.Annotations, AnnotationIdleDeadline)
	}
	return s.Client.Patch(ctx, explorer, patch)
}

func (s *Scaler) ResetIdleTimer(ctx context.Context, namespace, name string) (remaining time.Duration, err error) {
	explorer := &pvcv1.PVCExplorer{}
	if err := s.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, explorer); err != nil {
		return 0, err
	}

	if explorer.Status.Phase != pvcv1.ExplorerPhaseRunning {
		return 0, fmt.Errorf("agent not running: phase=%s", explorer.Status.Phase)
	}

	timeout := idleTimeout(explorer)
	patch := client.MergeFrom(explorer.DeepCopy())
	if explorer.Annotations == nil {
		explorer.Annotations = map[string]string{}
	}
	explorer.Annotations[AnnotationIdleDeadline] = time.Now().Add(timeout).UTC().Format(time.RFC3339)
	if err := s.Client.Patch(ctx, explorer, patch); err != nil {
		return 0, err
	}

	s.resetWarnState(namespace, name)

	if s.broadcaster != nil {
		_ = s.broadcaster.Publish("idle.tick", map[string]any{
			"namespace":        namespace,
			"name":             name,
			"remainingSeconds": int64(timeout.Seconds()),
			"idleTimeout":      explorer.Spec.Scaling.IdleTimeout,
		})
	}

	return timeout, nil
}

// RunIdleWatcher polls running explorers every 5 seconds and emits idle.tick/warning/expired events.
func (s *Scaler) RunIdleWatcher(ctx context.Context) {
	log := logf.FromContext(ctx).WithName("idle-watcher")
	ticker := time.NewTicker(idleWatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tickAll(ctx, log)
		}
	}
}

func (s *Scaler) tickAll(ctx context.Context, log interface{ Info(string, ...any) }) {
	var list pvcv1.PVCExplorerList
	if err := s.Client.List(ctx, &list); err != nil {
		return
	}

	for i := range list.Items {
		explorer := &list.Items[i]
		if explorer.Status.Phase != pvcv1.ExplorerPhaseRunning {
			continue
		}

		deadline, ok := IdleDeadline(explorer)
		if !ok {
			continue
		}

		remaining := time.Until(deadline)
		ns, name := explorer.Namespace, explorer.Name
		timeout := idleTimeout(explorer)

		if remaining <= 0 {
			if err := s.SleepAgent(ctx, ns, name); err == nil {
				if s.broadcaster != nil {
					_ = s.broadcaster.Publish("idle.expired", map[string]any{
						"namespace": ns,
						"name":      name,
						"expiredAt": time.Now().UTC().Format(time.RFC3339),
					})
				}
				s.deleteWarnState(ns, name)
			}
			continue
		}

		if s.broadcaster != nil {
			_ = s.broadcaster.Publish("idle.tick", map[string]any{
				"namespace":        ns,
				"name":             name,
				"remainingSeconds": int64(remaining.Seconds()),
				"idleTimeout":      timeout.String(),
			})
		}

		if remaining <= warningThreshold {
			state := s.getOrCreateWarnState(ns, name)
			if !state.warned {
				state.warned = true
				if s.broadcaster != nil {
					_ = s.broadcaster.Publish("idle.warning", map[string]any{
						"namespace":        ns,
						"name":             name,
						"remainingSeconds": int64(remaining.Seconds()),
						"warningThreshold": int64(warningThreshold.Seconds()),
					})
				}
			}
		}
	}
}

func (s *Scaler) getOrCreateWarnState(ns, name string) *idleWarnState {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := ns + "/" + name
	if st, ok := s.warnState[key]; ok {
		return st
	}
	st := &idleWarnState{}
	s.warnState[key] = st
	return st
}

func (s *Scaler) resetWarnState(ns, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := ns + "/" + name
	if st, ok := s.warnState[key]; ok {
		st.warned = false
	}
}

func (s *Scaler) deleteWarnState(ns, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.warnState, ns+"/"+name)
}

func idleTimeout(explorer *pvcv1.PVCExplorer) time.Duration {
	if explorer.Spec.Scaling.IdleTimeout != "" {
		if d, err := time.ParseDuration(explorer.Spec.Scaling.IdleTimeout); err == nil && d > 0 {
			return d
		}
	}
	return defaultIdleTimeout
}

func IdleDeadline(explorer *pvcv1.PVCExplorer) (time.Time, bool) {
	if explorer.Annotations == nil {
		return time.Time{}, false
	}
	raw, ok := explorer.Annotations[AnnotationIdleDeadline]
	if !ok {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
