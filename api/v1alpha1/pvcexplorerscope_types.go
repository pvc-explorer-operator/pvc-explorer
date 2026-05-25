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

// DiscoveryMode controls how PVCs are discovered within the registered namespaces.
// +kubebuilder:validation:Enum=Auto;Explicit
type DiscoveryMode string

const (
	// DiscoveryModeAuto discovers all PVCs in the registered namespaces,
	// subject to excludePVCs patterns.
	DiscoveryModeAuto DiscoveryMode = "Auto"
	// DiscoveryModeExplicit only manages PVCs listed in pvcNames.
	DiscoveryModeExplicit DiscoveryMode = "Explicit"
)

// DeletionPolicy controls what happens to PVCExplorer CRs when the scope is deleted.
// +kubebuilder:validation:Enum=Cleanup;Orphan
type DeletionPolicy string

const (
	// DeletionPolicyCleanup deletes all owned PVCExplorer CRs (and their agent
	// resources) when the scope is deleted. PVCs are never touched.
	DeletionPolicyCleanup DeletionPolicy = "Cleanup"
	// DeletionPolicyOrphan leaves PVCExplorer CRs in place as standalone resources.
	DeletionPolicyOrphan DeletionPolicy = "Orphan"
)

// FallbackPolicy controls agent behaviour when a PVC is in use on multiple nodes.
// +kubebuilder:validation:Enum=Pending
type FallbackPolicy string

const (
	// FallbackPolicyPending keeps the agent in Pending phase until the conflict resolves.
	FallbackPolicyPending FallbackPolicy = "Pending"
)

// ScopeNamespacesSpec selects which namespaces are registered in this scope.
type ScopeNamespacesSpec struct {
	// names is a static list of namespace names to register.
	// +optional
	Names []string `json:"names,omitempty"`

	// labelSelector registers namespaces that match this label selector.
	// Combined with names using OR logic.
	// +optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
}

// ScopeDiscoverySpec controls PVC discovery within registered namespaces.
type ScopeDiscoverySpec struct {
	// mode is Auto (all PVCs) or Explicit (only pvcNames).
	// +kubebuilder:default=Auto
	Mode DiscoveryMode `json:"mode"`

	// pvcNames is only used when mode is Explicit.
	// +optional
	PVCNames []string `json:"pvcNames,omitempty"`

	// excludePVCs is a list of PVC name patterns (glob syntax) to exclude even
	// in Auto mode.
	// +optional
	ExcludePVCs []string `json:"excludePVCs,omitempty"`
}

// ScopeMountStrategySpec describes the default mount strategy for agents in this scope.
type ScopeMountStrategySpec struct {
	// allowNodeAffinity permits the controller to set nodeAffinity on agent pods
	// when the PVC is RWO and consumed by another pod. Defaults to true.
	// +kubebuilder:default=true
	// +optional
	AllowNodeAffinity bool `json:"allowNodeAffinity,omitempty"`

	// fallbackOnConflict controls agent behaviour when a RWO PVC has consumers on
	// multiple nodes. Currently only Pending is supported.
	// +kubebuilder:default=Pending
	// +optional
	FallbackOnConflict FallbackPolicy `json:"fallbackOnConflict,omitempty"`
}

// ScopeScalingSpec defines the default scaling parameters for agents in this scope.
type ScopeScalingSpec struct {
	// idleTimeout is the duration an agent may remain idle before being scaled to zero.
	// +kubebuilder:default="10m"
	// +optional
	IdleTimeout string `json:"idleTimeout,omitempty"`

	// startupTimeout is the maximum time to wait for the agent pod to pass its
	// /healthz probe after being woken.
	// +kubebuilder:default="60s"
	// +optional
	StartupTimeout string `json:"startupTimeout,omitempty"`
}

// ScopeDefaultsSpec provides default values applied to every PVCExplorer created
// by this scope.
type ScopeDefaultsSpec struct {
	// mode is the default deployment mode for agents.
	// +kubebuilder:default=ScaledToZero
	// +optional
	Mode ExplorerMode `json:"mode,omitempty"`

	// image is the default container image for agent pods.
	// +optional
	Image string `json:"image,omitempty"`

	// forceRW, when true, mounts the PVC read-write when no consumers are present.
	// Automatically overridden to read-only when consumers are detected.
	// +kubebuilder:default=true
	// +optional
	ForceRW bool `json:"forceRW,omitempty"`

	// scaling holds the default idle/startup timeouts.
	// +optional
	Scaling ScopeScalingSpec `json:"scaling,omitempty"`

	// mountStrategy holds the default mount strategy options.
	// +optional
	MountStrategy ScopeMountStrategySpec `json:"mountStrategy,omitempty"`

	// resources defines default compute resource requests/limits for agent pods.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// PVCExplorerScopeSpec defines the desired state of PVCExplorerScope.
type PVCExplorerScopeSpec struct {
	// namespaces selects which namespaces are registered in this scope.
	// +required
	Namespaces ScopeNamespacesSpec `json:"namespaces"`

	// discovery controls how PVCs are discovered within the registered namespaces.
	// +optional
	Discovery ScopeDiscoverySpec `json:"discovery,omitempty"`

	// deletionPolicy controls what happens to owned PVCExplorer CRs when this scope
	// is deleted. PVCs are never touched.
	// +kubebuilder:default=Cleanup
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// defaults provides default values applied to every PVCExplorer created by this scope.
	// +optional
	Defaults ScopeDefaultsSpec `json:"defaults,omitempty"`
}

// PVCExplorerScopeStatus defines the observed state of PVCExplorerScope.
type PVCExplorerScopeStatus struct {
	// namespaceCount is the number of namespaces currently registered by this scope.
	// +optional
	NamespaceCount int `json:"namespaceCount,omitempty"`

	// explorerCount is the number of PVCExplorer CRs owned by this scope.
	// +optional
	ExplorerCount int `json:"explorerCount,omitempty"`

	// observedGeneration is the .metadata.generation that the status was last
	// reconciled against.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// conditions represent the current state of the PVCExplorerScope resource.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=pvcs
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Namespaces",type=integer,JSONPath=`.status.namespaceCount`
// +kubebuilder:printcolumn:name="Explorers",type=integer,JSONPath=`.status.explorerCount`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// PVCExplorerScope is the Schema for the pvcexplorerscopes API.
// A scope registers one or more namespaces for PVC exploration and provides
// default configuration for all PVCExplorer agents it creates.
type PVCExplorerScope struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of PVCExplorerScope.
	// +required
	Spec PVCExplorerScopeSpec `json:"spec"`

	// status defines the observed state of PVCExplorerScope.
	// +optional
	Status PVCExplorerScopeStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// PVCExplorerScopeList contains a list of PVCExplorerScope.
type PVCExplorerScopeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []PVCExplorerScope `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PVCExplorerScope{}, &PVCExplorerScopeList{})
}
