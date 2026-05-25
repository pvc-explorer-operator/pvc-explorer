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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"maps"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
	"github.com/pvc-explorer-operator/pvc-explorer/internal/consumer"
)

type agentBroadcaster interface {
	Publish(eventType string, payload any) error
}

const (
	agentRequeueAfter      = 30 * time.Second
	defaultAgentImage      = "ghcr.io/pvc-explorer-operator/pvc-explorer-agent:latest"
	defaultAgentPort       = int32(8081)
	agentContainerName     = "agent"
	agentPVCVolumeName     = "pvc-data"
	agentPVCMountPath      = "/data"
	agentTokenSecretSuffix = "-agent-token"
	agentTokenSecretKey    = "token"

	annotationRestartAt = "pvcexplorer.io/restart-at"
	condReady           = "Ready"
	strategyDirect      = "Direct"
	labelApp            = "app"
)

// AgentReconciler reconciles a PVCExplorer object.
//
// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=pvcexplorer.io,resources=pvcexplorers/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=statefulsets;daemonsets,verbs=get;list;watch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

type PVCExplorerReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ConsumerIndex *consumer.Index
	Broadcaster   agentBroadcaster
}

type mountDecision struct {
	readOnly        bool
	forceRWDeferred bool
	targetNode      string
	strategy        string
	accessMode      string
	pending         bool
}

func (r *PVCExplorerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	explorer := &pvcexplorerv1alpha1.PVCExplorer{}
	if err := r.Get(ctx, req.NamespacedName, explorer); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	pvcKey := types.NamespacedName{Name: explorer.Spec.PVCName, Namespace: explorer.Namespace}
	if err := r.Get(ctx, pvcKey, pvc); err != nil {
		if apierrors.IsNotFound(err) {
			return r.setFailed(ctx, explorer, "PVCNotFound",
				fmt.Sprintf("PVC %q not found in namespace %q", explorer.Spec.PVCName, explorer.Namespace))
		}
		return ctrl.Result{}, fmt.Errorf("fetching PVC %s: %w", explorer.Spec.PVCName, err)
	}
	if pvc.Status.Phase != corev1.ClaimBound {
		if pvc.Status.Phase == corev1.ClaimPending {
			patch := client.MergeFrom(explorer.DeepCopy())
			explorer.Status.Phase = pvcexplorerv1alpha1.ExplorerPhasePending
			explorer.Status.ObservedGeneration = explorer.Generation
			apimeta.SetStatusCondition(&explorer.Status.Conditions, metav1.Condition{
				Type:               condReady,
				Status:             metav1.ConditionFalse,
				Reason:             "PVCNotBound",
				Message:            fmt.Sprintf("PVC %q is Pending, waiting for volume to bind", explorer.Spec.PVCName),
				ObservedGeneration: explorer.Generation,
			})
			if err := r.Status().Patch(ctx, explorer, patch); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: agentRequeueAfter}, nil
		}
		return r.setFailed(ctx, explorer, "PVCNotBound",
			fmt.Sprintf("PVC %q is in phase %q, expected Bound", explorer.Spec.PVCName, pvc.Status.Phase))
	}

	consumers, err := consumer.Detect(ctx, r.Client, explorer.Namespace, explorer.Spec.PVCName)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("detecting consumers: %w", err)
	}
	if r.ConsumerIndex != nil {
		if idx := r.ConsumerIndex.Get(explorer.Namespace, explorer.Spec.PVCName); len(idx) > 0 {
			consumers = idx
		}
	}
	consumers = excludeSelf(consumers, explorer.Name)

	dec := r.mountPolicy(explorer, pvc, consumers)

	if dec.pending {
		return r.setPending(ctx, explorer, consumers)
	}

	if err := r.reconcileAgentTokenSecret(ctx, explorer); err != nil {
		return ctrl.Result{}, fmt.Errorf("reconciling agent token Secret: %w", err)
	}

	if err := r.reconcileDeployment(ctx, explorer, dec); err != nil {
		return ctrl.Result{}, fmt.Errorf("reconciling Deployment: %w", err)
	}
	if err := r.reconcileService(ctx, explorer); err != nil {
		return ctrl.Result{}, fmt.Errorf("reconciling Service: %w", err)
	}
	if err := r.reconcileAgentRBAC(ctx, explorer); err != nil {
		return ctrl.Result{}, fmt.Errorf("reconciling agent RBAC: %w", err)
	}
	if err := r.syncStatus(ctx, explorer, dec, consumers); err != nil {
		return ctrl.Result{}, fmt.Errorf("syncing status: %w", err)
	}

	if explorer.Status.Phase == pvcexplorerv1alpha1.ExplorerPhaseRunning {
		if deadline, ok := idleDeadline(explorer); ok && time.Now().After(deadline) {
			patch := client.MergeFrom(explorer.DeepCopy())
			explorer.Spec.Mode = pvcexplorerv1alpha1.ExplorerModeScaledToZero
			delete(explorer.Annotations, "pvcexplorer.io/idle-deadline")
			if err := r.Patch(ctx, explorer, patch); err != nil {
				return ctrl.Result{}, fmt.Errorf("sleeping idle agent: %w", err)
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	log.Info("Reconciled PVCExplorer",
		"name", explorer.Name,
		"namespace", explorer.Namespace,
		"phase", explorer.Status.Phase,
		"readOnly", dec.readOnly,
		"consumers", len(consumers),
	)
	return ctrl.Result{RequeueAfter: agentRequeueAfter}, nil
}

func (r *PVCExplorerReconciler) mountPolicy(
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	pvc *corev1.PersistentVolumeClaim,
	consumers []pvcexplorerv1alpha1.ConsumerInfo,
) mountDecision {
	accessMode := ""
	if len(pvc.Spec.AccessModes) > 0 {
		accessMode = string(pvc.Spec.AccessModes[0])
	}
	dec := mountDecision{accessMode: accessMode}

	isROX := len(pvc.Spec.AccessModes) > 0 && pvc.Spec.AccessModes[0] == corev1.ReadOnlyMany
	if isROX {
		dec.readOnly = true
		dec.strategy = strategyDirect
		return dec
	}

	isRWO := len(pvc.Spec.AccessModes) > 0 && pvc.Spec.AccessModes[0] == corev1.ReadWriteOnce

	if len(consumers) == 0 {
		dec.readOnly = !explorer.Spec.ForceRW
		dec.strategy = strategyDirect
		return dec
	}

	dec.readOnly = true
	dec.forceRWDeferred = explorer.Spec.ForceRW

	if isRWO {
		nodes := uniqueNodes(consumers)
		if len(nodes) > 1 {
			dec.pending = true
			return dec
		}
		if len(nodes) == 1 {
			dec.targetNode = nodes[0]
			dec.strategy = "NodeAffinity"
			return dec
		}
	}

	dec.strategy = "Direct"
	return dec
}

func uniqueNodes(consumers []pvcexplorerv1alpha1.ConsumerInfo) []string {
	seen := map[string]struct{}{}
	var nodes []string
	for _, c := range consumers {
		if c.NodeName != "" {
			if _, ok := seen[c.NodeName]; !ok {
				seen[c.NodeName] = struct{}{}
				nodes = append(nodes, c.NodeName)
			}
		}
	}
	return nodes
}

func (r *PVCExplorerReconciler) reconcileDeployment(
	ctx context.Context,
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	dec mountDecision,
) error {
	desired := r.buildDeployment(explorer, dec)
	if err := controllerutil.SetControllerReference(explorer, desired, r.Scheme); err != nil {
		return err
	}

	existing := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, existing)
	if apierrors.IsNotFound(err) {
		logf.FromContext(ctx).Info("Creating Deployment", "name", desired.Name, "namespace", desired.Namespace)
		return r.Create(ctx, desired)
	}
	if err != nil {
		return err
	}

	affinityChanged := nodeAffinityChanged(existing, desired)
	patch := client.MergeFrom(existing.DeepCopy())
	existing.Spec.Replicas = desired.Spec.Replicas
	existing.Spec.Template.Spec.Containers = desired.Spec.Template.Spec.Containers
	existing.Spec.Template.Spec.Volumes = desired.Spec.Template.Spec.Volumes
	existing.Spec.Template.Spec.Affinity = desired.Spec.Template.Spec.Affinity
	if affinityChanged {
		if existing.Spec.Template.Annotations == nil {
			existing.Spec.Template.Annotations = map[string]string{}
		}
		existing.Spec.Template.Annotations[annotationRestartAt] = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return r.Patch(ctx, existing, patch)
}

func nodeAffinityChanged(existing, desired *appsv1.Deployment) bool {
	return extractAffinityNode(existing) != extractAffinityNode(desired)
}

func extractAffinityNode(deploy *appsv1.Deployment) string {
	aff := deploy.Spec.Template.Spec.Affinity
	if aff == nil || aff.NodeAffinity == nil {
		return ""
	}
	req := aff.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	if req == nil || len(req.NodeSelectorTerms) == 0 {
		return ""
	}
	exprs := req.NodeSelectorTerms[0].MatchExpressions
	if len(exprs) == 0 || len(exprs[0].Values) == 0 {
		return ""
	}
	return exprs[0].Values[0]
}

func (r *PVCExplorerReconciler) buildDeployment(
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	dec mountDecision,
) *appsv1.Deployment {
	replicas := r.desiredReplicas(explorer)
	image := explorer.Spec.Image
	if image == "" {
		image = defaultAgentImage
	}
	port := explorer.Spec.Port
	if port == 0 {
		port = defaultAgentPort
	}

	podLabels := map[string]string{
		managedByLabelKey: managedByValue,
		labelApp:          explorer.Name,
	}
	maps.Copy(podLabels, explorer.Spec.ExplorerLabels)

	container := corev1.Container{
		Name:  agentContainerName,
		Image: image,
		Args: []string{
			"--root", agentPVCMountPath,
			"--pvc", explorer.Spec.PVCName,
		},
		Ports: []corev1.ContainerPort{{Name: "http", ContainerPort: port, Protocol: corev1.ProtocolTCP}},
		Env: []corev1.EnvVar{
			{Name: "POD_NAMESPACE", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"},
			}},
			{Name: "POD_NAME", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"},
			}},
			{
				Name: "AUTH_TOKEN",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: tokenSecretName(explorer.Name),
						},
						Key: agentTokenSecretKey,
					},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      agentPVCVolumeName,
				MountPath: agentPVCMountPath,
				SubPath:   explorer.Spec.SubPath,
				ReadOnly:  dec.readOnly,
			},
		},
		Resources: explorer.Spec.Resources,
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      explorer.Name,
			Namespace: explorer.Namespace,
			Labels: map[string]string{
				managedByLabelKey: managedByValue,
				scopeLabelKey:     explorer.Labels[scopeLabelKey],
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{labelApp: explorer.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: podLabels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{container},
					Volumes: []corev1.Volume{
						{
							Name: agentPVCVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: explorer.Spec.PVCName,
									ReadOnly:  dec.readOnly,
								},
							},
						},
					},
				},
			},
		},
	}

	if dec.targetNode != "" {
		deploy.Spec.Template.Spec.Affinity = nodeAffinityForNode(dec.targetNode)
	}

	return deploy
}

func nodeAffinityForNode(nodeName string) *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "kubernetes.io/hostname",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{nodeName},
							},
						},
					},
				},
			},
		},
	}
}

func (r *PVCExplorerReconciler) desiredReplicas(explorer *pvcexplorerv1alpha1.PVCExplorer) int32 {
	if explorer.Spec.Mode == pvcexplorerv1alpha1.ExplorerModeDeployment {
		return 1
	}
	return 0
}

func (r *PVCExplorerReconciler) reconcileAgentRBAC(
	ctx context.Context,
	explorer *pvcexplorerv1alpha1.PVCExplorer,
) error {
	roleName := explorer.Name + "-pvcwatch"
	ns := explorer.Namespace

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{Name: roleName, Namespace: ns},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{""}, Resources: []string{"pods"}, Verbs: []string{"list"}},
		},
	}
	if err := controllerutil.SetControllerReference(explorer, role, r.Scheme); err != nil {
		return err
	}
	existing := &rbacv1.Role{}
	err := r.Get(ctx, types.NamespacedName{Name: roleName, Namespace: ns}, existing)
	if apierrors.IsNotFound(err) {
		return r.Create(ctx, role)
	}
	if err != nil {
		return err
	}

	binding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: roleName, Namespace: ns},
		RoleRef:    rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "Role", Name: roleName},
		Subjects:   []rbacv1.Subject{{Kind: "ServiceAccount", Name: "default", Namespace: ns}},
	}
	if err := controllerutil.SetControllerReference(explorer, binding, r.Scheme); err != nil {
		return err
	}
	existingBinding := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: roleName, Namespace: ns}, existingBinding)
	if apierrors.IsNotFound(err) {
		return r.Create(ctx, binding)
	}
	return err
}

// tokenSecretName returns the deterministic name for the agent authentication Secret.
func tokenSecretName(explorerName string) string {
	return explorerName + agentTokenSecretSuffix
}

// reconcileAgentTokenSecret creates or deletes the agent authentication Secret
// based on the desired number of agent replicas:
//   - replicas == 0 (scaled to zero): delete the Secret — no pod means no token.
//   - replicas >  0 (running or waking): ensure the Secret exists with a fresh random token.
//
// The Secret carries an owner reference so it is cleaned up on PVCExplorer deletion.
func (r *PVCExplorerReconciler) reconcileAgentTokenSecret(ctx context.Context, explorer *pvcexplorerv1alpha1.PVCExplorer) error {
	name := tokenSecretName(explorer.Name)
	ns := explorer.Namespace

	if r.desiredReplicas(explorer) == 0 {
		// Agent is scaling to zero — delete the token Secret.
		existing := &corev1.Secret{}
		err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, existing)
		if apierrors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
		logf.FromContext(ctx).Info("Deleting agent token Secret", "name", name, "namespace", ns)
		return r.Delete(ctx, existing)
	}

	// Agent is (or will be) running — ensure the token Secret exists.
	existing := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, existing)
	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    map[string]string{managedByLabelKey: managedByValue},
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			agentTokenSecretKey: token,
		},
	}
	if err := controllerutil.SetControllerReference(explorer, secret, r.Scheme); err != nil {
		return err
	}
	logf.FromContext(ctx).Info("Creating agent token Secret", "name", name, "namespace", ns)
	return r.Create(ctx, secret)
}

// generateToken returns a cryptographically random hex-encoded token (64 hex chars = 32 bytes).
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *PVCExplorerReconciler) reconcileService(
	ctx context.Context,
	explorer *pvcexplorerv1alpha1.PVCExplorer,
) error {
	desired := r.buildService(explorer)
	if err := controllerutil.SetControllerReference(explorer, desired, r.Scheme); err != nil {
		return err
	}

	existing := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, existing)
	if apierrors.IsNotFound(err) {
		logf.FromContext(ctx).Info("Creating Service", "name", desired.Name, "namespace", desired.Namespace)
		return r.Create(ctx, desired)
	}
	if err != nil {
		return err
	}

	patch := client.MergeFrom(existing.DeepCopy())
	existing.Spec.Selector = desired.Spec.Selector
	existing.Spec.Ports = desired.Spec.Ports
	return r.Patch(ctx, existing, patch)
}

func (r *PVCExplorerReconciler) buildService(explorer *pvcexplorerv1alpha1.PVCExplorer) *corev1.Service {
	port := explorer.Spec.Port
	if port == 0 {
		port = defaultAgentPort
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      explorer.Name,
			Namespace: explorer.Namespace,
			Labels:    map[string]string{managedByLabelKey: managedByValue},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{labelApp: explorer.Name},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       port,
					TargetPort: intstr.FromString("http"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (r *PVCExplorerReconciler) syncStatus(
	ctx context.Context,
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	dec mountDecision,
	consumers []pvcexplorerv1alpha1.ConsumerInfo,
) error {
	deploy := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: explorer.Name, Namespace: explorer.Namespace}, deploy)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	patch := client.MergeFrom(explorer.DeepCopy())
	prevPhase := explorer.Status.Phase

	port := explorer.Spec.Port
	if port == 0 {
		port = defaultAgentPort
	}
	explorer.Status.AgentEndpoint = fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", explorer.Name, explorer.Namespace, port)
	explorer.Status.Mode = explorer.Spec.Mode
	if explorer.Status.Mode == "" {
		explorer.Status.Mode = pvcexplorerv1alpha1.ExplorerModeScaledToZero
	}
	explorer.Status.ObservedGeneration = explorer.Generation
	explorer.Status.Mount = pvcexplorerv1alpha1.MountStatus{
		AccessMode:      dec.accessMode,
		Strategy:        dec.strategy,
		ReadOnly:        dec.readOnly,
		ForceRWDeferred: dec.forceRWDeferred,
		TargetNode:      dec.targetNode,
		Consumers:       consumers,
	}
	explorer.Status.Phase = r.derivePhase(explorer, deploy)

	apimeta.SetStatusCondition(&explorer.Status.Conditions, metav1.Condition{
		Type:               condReady,
		Status:             metav1.ConditionTrue,
		Reason:             string(explorer.Status.Phase),
		Message:            fmt.Sprintf("Agent %s", explorer.Status.Phase),
		ObservedGeneration: explorer.Generation,
	})

	if dec.forceRWDeferred {
		apimeta.SetStatusCondition(&explorer.Status.Conditions, metav1.Condition{
			Type:               "ForceRWDeferred",
			Status:             metav1.ConditionTrue,
			Reason:             "ConsumerActive",
			Message:            fmt.Sprintf("Mounted ReadOnly due to %d active consumer(s); will remount ReadWrite on release", len(consumers)),
			ObservedGeneration: explorer.Generation,
		})
	} else {
		apimeta.SetStatusCondition(&explorer.Status.Conditions, metav1.Condition{
			Type:               "ForceRWDeferred",
			Status:             metav1.ConditionFalse,
			Reason:             "NoConsumers",
			Message:            "No active consumers; mount policy applied as configured",
			ObservedGeneration: explorer.Generation,
		})
	}

	if err := r.Status().Patch(ctx, explorer, patch); err != nil {
		return err
	}

	if r.Broadcaster != nil &&
		prevPhase != pvcexplorerv1alpha1.ExplorerPhaseRunning &&
		explorer.Status.Phase == pvcexplorerv1alpha1.ExplorerPhaseRunning {
		_ = r.Broadcaster.Publish("agent.ready", map[string]string{
			"namespace": explorer.Namespace,
			"name":      explorer.Name,
		})
	}

	return nil
}

func (r *PVCExplorerReconciler) derivePhase(
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	deploy *appsv1.Deployment,
) pvcexplorerv1alpha1.ExplorerPhase {
	if deploy == nil || deploy.Name == "" {
		return pvcexplorerv1alpha1.ExplorerPhasePending
	}
	desired := r.desiredReplicas(explorer)
	if desired == 0 {
		if deploy.Status.ReadyReplicas == 0 {
			return pvcexplorerv1alpha1.ExplorerPhaseScaledToZero
		}
		return pvcexplorerv1alpha1.ExplorerPhaseWaking
	}
	if deploy.Status.ReadyReplicas >= 1 {
		return pvcexplorerv1alpha1.ExplorerPhaseRunning
	}
	if deploy.Status.UnavailableReplicas > 0 {
		return pvcexplorerv1alpha1.ExplorerPhaseFailed
	}
	return pvcexplorerv1alpha1.ExplorerPhaseWaking
}

func (r *PVCExplorerReconciler) setPending(
	ctx context.Context,
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	consumers []pvcexplorerv1alpha1.ConsumerInfo,
) (ctrl.Result, error) {
	logf.FromContext(ctx).Info("PVCExplorer Pending due to multi-node RWO consumers",
		"name", explorer.Name, "consumers", len(consumers))
	patch := client.MergeFrom(explorer.DeepCopy())
	explorer.Status.Phase = pvcexplorerv1alpha1.ExplorerPhasePending
	explorer.Status.ObservedGeneration = explorer.Generation
	explorer.Status.Mount.Consumers = consumers
	apimeta.SetStatusCondition(&explorer.Status.Conditions, metav1.Condition{
		Type:               condReady,
		Status:             metav1.ConditionFalse,
		Reason:             "MultiNodeConflict",
		Message:            fmt.Sprintf("RWO PVC mounted on %d nodes; waiting for conflict to resolve", len(consumers)),
		ObservedGeneration: explorer.Generation,
	})
	if err := r.Status().Patch(ctx, explorer, patch); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: agentRequeueAfter}, nil
}

func (r *PVCExplorerReconciler) setFailed(
	ctx context.Context,
	explorer *pvcexplorerv1alpha1.PVCExplorer,
	reason, message string,
) (ctrl.Result, error) {
	logf.FromContext(ctx).Info("Setting PVCExplorer Failed", "reason", reason, "message", message)
	patch := client.MergeFrom(explorer.DeepCopy())
	explorer.Status.Phase = pvcexplorerv1alpha1.ExplorerPhaseFailed
	explorer.Status.ObservedGeneration = explorer.Generation
	apimeta.SetStatusCondition(&explorer.Status.Conditions, metav1.Condition{
		Type:               condReady,
		Status:             metav1.ConditionFalse,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: explorer.Generation,
	})
	if err := r.Status().Patch(ctx, explorer, patch); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: agentRequeueAfter}, nil
}

func (r *PVCExplorerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	podPredicate := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return podMountsPVC(e.Object) && r.isRegisteredNamespace(e.Object.GetNamespace())
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldPod, ok1 := e.ObjectOld.(*corev1.Pod)
			newPod, ok2 := e.ObjectNew.(*corev1.Pod)
			if !ok1 || !ok2 {
				return false
			}
			return oldPod.Status.Phase != newPod.Status.Phase &&
				podMountsPVC(e.ObjectNew) &&
				r.isRegisteredNamespace(e.ObjectNew.GetNamespace())
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return podMountsPVC(e.Object) && r.isRegisteredNamespace(e.Object.GetNamespace())
		},
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&pvcexplorerv1alpha1.PVCExplorer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Watches(
			&corev1.Pod{},
			handler.EnqueueRequestsFromMapFunc(r.podToPVCExplorers),
			builder.WithPredicates(podPredicate),
		).
		WithOptions(controller.Options{
			RateLimiter: workqueue.NewTypedItemFastSlowRateLimiter[reconcile.Request](
				100*time.Millisecond,
				10*time.Second,
				5,
			),
		}).
		Named("agent").
		Complete(r)
}

func podMountsPVC(obj client.Object) bool {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return false
	}
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim != nil {
			return true
		}
	}
	return false
}

func (r *PVCExplorerReconciler) isRegisteredNamespace(ns string) bool {
	explorerList := &pvcexplorerv1alpha1.PVCExplorerList{}
	if err := r.List(context.Background(), explorerList, client.InNamespace(ns)); err != nil {
		return false
	}
	return len(explorerList.Items) > 0
}

func (r *PVCExplorerReconciler) podToPVCExplorers(ctx context.Context, obj client.Object) []reconcile.Request {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil
	}

	if r.ConsumerIndex != nil {
		if err := r.ConsumerIndex.Sync(ctx, r.Client, pod.Namespace); err != nil {
			logf.FromContext(ctx).Error(err, "Failed to sync consumer index", "namespace", pod.Namespace)
		}
	}

	pvcNames := map[string]struct{}{}
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim != nil {
			pvcNames[v.PersistentVolumeClaim.ClaimName] = struct{}{}
		}
	}
	var requests []reconcile.Request
	for pvcName := range pvcNames {
		key := types.NamespacedName{Name: pvcName, Namespace: pod.Namespace}
		explorer := &pvcexplorerv1alpha1.PVCExplorer{}
		if err := r.Get(ctx, key, explorer); err == nil {
			requests = append(requests, reconcile.Request{NamespacedName: key})
		}
	}
	return requests
}

func idleDeadline(explorer *pvcexplorerv1alpha1.PVCExplorer) (time.Time, bool) {
	if explorer.Annotations == nil {
		return time.Time{}, false
	}
	raw, ok := explorer.Annotations["pvcexplorer.io/idle-deadline"]
	if !ok {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func excludeSelf(consumers []pvcexplorerv1alpha1.ConsumerInfo, explorerName string) []pvcexplorerv1alpha1.ConsumerInfo {
	filtered := consumers[:0:0]
	for _, c := range consumers {
		if c.OwnerKind == "Deployment" && c.OwnerName == explorerName {
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}
