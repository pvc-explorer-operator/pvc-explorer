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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExplorerMode controls whether the agent runs continuously or scales to zero.
// +kubebuilder:validation:Enum=ScaledToZero;Deployment
type ExplorerMode string

const (
	// ExplorerModeScaledToZero keeps the agent at replicas=0 until a user
	// requests access. This is the default and recommended mode.
	ExplorerModeScaledToZero ExplorerMode = "ScaledToZero"
	// ExplorerModeDeployment keeps the agent always running (replicas=1).
	ExplorerModeDeployment ExplorerMode = "Deployment"
)

// ExplorerPhase is the high-level lifecycle phase of the explorer agent.
// +kubebuilder:validation:Enum=ScaledToZero;Waking;Running;Failed;Pending
type ExplorerPhase string

const (
	ExplorerPhaseScaledToZero ExplorerPhase = "ScaledToZero"
	ExplorerPhaseWaking       ExplorerPhase = "Waking"
	ExplorerPhaseRunning      ExplorerPhase = "Running"
	ExplorerPhaseFailed       ExplorerPhase = "Failed"
	ExplorerPhasePending      ExplorerPhase = "Pending"
)

// ScalingProvider selects the scale-to-zero implementation.
// +kubebuilder:validation:Enum=auto;native;knative
type ScalingProvider string

const (
	ScalingProviderAuto    ScalingProvider = "auto"
	ScalingProviderNative  ScalingProvider = "native"
	ScalingProviderKnative ScalingProvider = "knative"
)

// ScalingSpec configures scale-to-zero behaviour for this agent.
type ScalingSpec struct {
	// provider selects the scaling backend. "auto" uses Knative if available,
	// otherwise falls back to native Deployment scaling.
	// +kubebuilder:default=auto
	// +optional
	Provider ScalingProvider `json:"provider,omitempty"`

	// idleTimeout is the duration of inactivity before the agent is scaled to zero.
	// +kubebuilder:default="10m"
	// +optional
	IdleTimeout string `json:"idleTimeout,omitempty"`

	// startupTimeout is the maximum time to wait for the agent pod to pass its
	// /healthz probe after being woken.
	// +kubebuilder:default="60s"
	// +optional
	StartupTimeout string `json:"startupTimeout,omitempty"`
}

// MountStrategySpec configures how the agent handles PVC mounting edge cases.
type MountStrategySpec struct {
	// autoDetect enables automatic consumer detection. When consumers are found
	// the agent is always mounted read-only, regardless of forceRW.
	// +kubebuilder:default=true
	// +optional
	AutoDetect bool `json:"autoDetect,omitempty"`

	// allowNodeAffinity permits the controller to add nodeAffinity to the agent
	// pod when the PVC is RWO and a consumer is running on a specific node.
	// +kubebuilder:default=true
	// +optional
	AllowNodeAffinity bool `json:"allowNodeAffinity,omitempty"`

	// fallbackOnConflict controls behaviour when a RWO PVC has consumers on
	// multiple nodes simultaneously.
	// +kubebuilder:default=Pending
	// +optional
	FallbackOnConflict FallbackPolicy `json:"fallbackOnConflict,omitempty"`
}

// PVCExplorerSpec defines the desired state of PVCExplorer.
type PVCExplorerSpec struct {
	// pvcName is the name of the PersistentVolumeClaim to explore.
	// The PVC must exist in the same namespace as this PVCExplorer.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	PVCName string `json:"pvcName"`

	// subPath mounts only a sub-directory of the PVC rather than its root.
	// +optional
	SubPath string `json:"subPath,omitempty"`

	// mode controls whether the agent scales to zero when idle.
	// +kubebuilder:default=ScaledToZero
	// +optional
	Mode ExplorerMode `json:"mode,omitempty"`

	// image is the container image for the agent pod.
	// +optional
	Image string `json:"image,omitempty"`

	// port is the HTTP port the agent listens on.
	// +kubebuilder:default=8081
	// +optional
	Port int32 `json:"port,omitempty"`

	// forceRW, when true, mounts the PVC read-write when no consumers are present.
	// Automatically overridden to read-only when consumers are detected.
	// +kubebuilder:default=true
	// +optional
	ForceRW bool `json:"forceRW,omitempty"`

	// explorerLabels are extra labels propagated onto the agent Deployment and pods.
	// +optional
	ExplorerLabels map[string]string `json:"explorerLabels,omitempty"`

	// scaling configures scale-to-zero behaviour.
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`

	// mountStrategy configures how mount edge cases are handled.
	// +optional
	MountStrategy MountStrategySpec `json:"mountStrategy,omitempty"`

	// resources specifies compute resource requests and limits for the agent pod.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// ConsumerInfo describes a workload pod that is currently mounting the tracked PVC.
type ConsumerInfo struct {
	// podName is the name of the consumer pod.
	PodName string `json:"podName"`

	// ownerKind is the kind of the top-level owner (Deployment, StatefulSet, Job, etc.).
	// +optional
	OwnerKind string `json:"ownerKind,omitempty"`

	// ownerName is the name of the top-level owner.
	// +optional
	OwnerName string `json:"ownerName,omitempty"`

	// nodeName is the node the consumer pod is scheduled on.
	// +optional
	NodeName string `json:"nodeName,omitempty"`

	// mountReadOnly indicates whether the consumer mounts the PVC read-only.
	// +optional
	MountReadOnly bool `json:"mountReadOnly,omitempty"`
}

// MountStatus describes the current mount state of the agent.
type MountStatus struct {
	// accessMode is the access mode of the underlying PVC (ReadWriteOnce, etc.).
	// +optional
	AccessMode string `json:"accessMode,omitempty"`

	// strategy is the mount strategy in use (Direct, NodeAffinity, etc.).
	// +optional
	Strategy string `json:"strategy,omitempty"`

	// readOnly indicates whether the agent currently has the PVC mounted read-only.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`

	// forceRWDeferred is true when the agent is currently read-only due to active
	// consumers, but will be remounted read-write once they release the PVC.
	// +optional
	ForceRWDeferred bool `json:"forceRWDeferred,omitempty"`

	// targetNode is the node the agent pod is (or must be) scheduled on. Non-empty
	// only when nodeAffinity is in use for RWO PVCs.
	// +optional
	TargetNode string `json:"targetNode,omitempty"`

	// consumers lists workload pods that are currently mounting the tracked PVC.
	// +optional
	Consumers []ConsumerInfo `json:"consumers,omitempty"`
}

// AgentStatus captures metadata reported by the running agent pod.
type AgentStatus struct {
	// version is the Unix-timestamp build version string from the agent binary.
	// +optional
	Version string `json:"version,omitempty"`

	// cluster is the cluster name reported by the agent.
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// pvcWatchEnabled indicates whether the agent has PVC-event watching enabled.
	// +optional
	PVCWatchEnabled bool `json:"pvcWatchEnabled,omitempty"`
}

// PVCExplorerStatus defines the observed state of PVCExplorer.
type PVCExplorerStatus struct {
	// phase is the high-level lifecycle phase of the agent.
	// +optional
	Phase ExplorerPhase `json:"phase,omitempty"`

	// mode mirrors spec.mode as observed.
	// +optional
	Mode ExplorerMode `json:"mode,omitempty"`

	// agentEndpoint is the in-cluster HTTP address of the agent service.
	// +optional
	AgentEndpoint string `json:"agentEndpoint,omitempty"`

	// mount describes the current PVC mount state.
	// +optional
	Mount MountStatus `json:"mount,omitempty"`

	// agent captures metadata reported by the running agent pod.
	// +optional
	Agent AgentStatus `json:"agent,omitempty"`

	// lastHealthCheck is the timestamp of the last successful /healthz probe.
	// Empty when the agent is not running.
	// +optional
	LastHealthCheck *metav1.Time `json:"lastHealthCheck,omitempty"`

	// observedGeneration is the .metadata.generation that the status was last
	// reconciled against.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// conditions represent the current state of the PVCExplorer resource.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=pvcexp
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="PVC",type=string,JSONPath=`.spec.pvcName`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="ReadOnly",type=boolean,JSONPath=`.status.mount.readOnly`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// PVCExplorer is the Schema for the pvcexplorers API.
// Each PVCExplorer manages a single agent pod that provides HTTP/WebSocket
// access to a PersistentVolumeClaim.
type PVCExplorer struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of PVCExplorer.
	// +required
	Spec PVCExplorerSpec `json:"spec"`

	// status defines the observed state of PVCExplorer.
	// +optional
	Status PVCExplorerStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// PVCExplorerList contains a list of PVCExplorer.
type PVCExplorerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []PVCExplorer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PVCExplorer{}, &PVCExplorerList{})
}
