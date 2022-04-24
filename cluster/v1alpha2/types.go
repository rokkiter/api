package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	ClusterReadyCondition   = "Ready"
	ClusterSynchroCondition = "ClusterSynchroInitialized"

	SyncStatusPending = "Pending"
	SyncStatusSyncing = "Syncing"
	SyncStatusStop    = "Stop"
	SyncStatusUnknown = "Unknown"

	InvalidConfigConditionReason = "InvalidConfig"
	InitialFailedConditionReason = "InitialFailed"
	RunningConditionReason       = "Running"

	HealthyReason            = "Healthy"
	PendingReason            = "Pending"
	UnhealthyReason          = "Unhealthy"
	NotReachableReason       = "NotReachable"
	ClusterSynchroStopReason = "ClusterSynchroStop"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Cluster"
// +kubebuilder:printcolumn:name="APIServer",type=string,JSONPath=".spec.apiserver"
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=".status.version"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=".status.conditions[?(@.type == 'Ready')].reason"
type PediaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec ClusterSpec `json:"spec,omitempty"`

	// +optional
	Status ClusterStatus `json:"status,omitempty"`
}

type ClusterSpec struct {
	// +optional
	Kubeconfig []byte `json:"kubeconfig,omitempty"`

	// +optional
	APIServer string `json:"apiserver,omitempty"`

	// +optional
	TokenData []byte `json:"tokenData,omitempty"`

	// +optional
	CAData []byte `json:"caData,omitempty"`

	// +optional
	CertData []byte `json:"certData,omitempty"`

	// +optional
	KeyData []byte `json:"keyData,omitempty"`

	// +required
	SyncResources []ClusterGroupResources `json:"syncResources"`

	// +optional
	SyncAllCustomResources bool `json:"syncAllCustomResources,omitempty"`
}

type ClusterGroupResources struct {
	Group string `json:"group"`

	// +optional
	Versions []string `json:"versions,omitempty"`

	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Resources []string `json:"resources"`
}

type ClusterStatus struct {
	// +required
	// +kubebuilder:validation:Required
	Version string `json:"version,omitempty"`

	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	SyncResources []ClusterGroupResourcesStatus `json:"syncResources,omitempty"`
}

type ClusterGroupResourcesStatus struct {
	// +required
	// +kubebuilder:validation:Required
	Group string `json:"group"`

	// +required
	// +kubebuilder:validation:Required
	Resources []ClusterResourceStatus `json:"resources"`
}

type ClusterResourceStatus struct {
	// +required
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +required
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// +required
	// +kubebuilder:validation:Required
	Namespaced bool `json:"namespaced"`

	// +required
	// +kubebuilder:validation:Required
	SyncConditions []ClusterResourceSyncCondition `json:"syncConditions"`
}

type ClusterResourceSyncCondition struct {
	// +required
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// optional
	SyncVersion string `json:"syncVersion,omitempty"`

	// optional
	SyncResource string `json:"syncResource,omitempty"`

	// optional
	StorageVersion string `json:"storageVersion,omitempty"`

	// optional
	StorageResource string `json:"storrageResource,omitempty"`

	// +required
	// +kubebuilder:validation:Required
	Status string `json:"status"`

	// optional
	Reason string `json:"reason,omitempty"`

	// optional
	Message string `json:"message,omitempty"`

	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

func (cond ClusterResourceSyncCondition) SyncGVR(resource schema.GroupResource) schema.GroupVersionResource {
	if cond.Version == "" || cond.SyncVersion == "" {
		return schema.GroupVersionResource{}
	}

	if cond.SyncResource != "" {
		resource = schema.ParseGroupResource(cond.SyncResource)
	}
	if cond.SyncVersion != "" {
		return resource.WithVersion(cond.StorageVersion)
	}
	return resource.WithVersion(cond.Version)
}

func (cond ClusterResourceSyncCondition) StorageGVR(resource schema.GroupResource) schema.GroupVersionResource {
	if cond.Version == "" || cond.StorageVersion == "" {
		return schema.GroupVersionResource{}
	}

	if cond.StorageResource != "" {
		return schema.ParseGroupResource(cond.StorageResource).WithVersion(cond.StorageVersion)
	}
	return resource.WithVersion(cond.StorageVersion)
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PediaClusterList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PediaCluster `json:"items"`
}
