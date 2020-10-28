/*
Copyright 2020 Tutkovics Andras.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Random comment
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type endpoint struct {
	// Random Comment
	NodeName string `json:"nodename,omitempty"`
	// Random Comment
	EndpointName string `json:"name,omitempty"`
	// Random Comment
	EndpointPath string `json:"path,omitempty"`
	// Random Comment
	CPU string `json:"cpuload,omitempty"`
	// Random Comment
	Delay uint `json:"delay,omitempty"`
	// Random Comment
	CallOut string `json:"callout,omitempty"`
}

// ServiceGraphSpec defines the desired state of ServiceGraph
// +k8s:openapi-gen=true
type ServiceGraphSpec struct {

	// Array for nodes in service-graph
	// +listType=atomic
	Nodes []*node `json:"nodes,omitempty"`

	Random int `json:"random"`

	// API endpoints for each node
	// +listType=atomic
	APIEndpoints []*endpoint `json:"apiendpoints,omitempty"`
}

// MemcachedStatus defines the observed state of Memcached
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type node struct {
	// Random Comment
	Name string `json:"name,omitempty"`
	// Random Comment
	ContainerPort uint `json:"containerport,omitempty"`
	// Random Comment
	NodePort uint `json:"nodeport,omitempty"`
	// Random Comment
	Replicas uint `json:"replicas,omitempty"`
	// Random Comment
	Resources resource `json:"resources,omitempty"`
	//APIEndpoints  *[]endpoint `json:"apiendpoints,omitempty"` Nem működött, ha ide volt beágyazva
	// api/v1/zz_generated.deepcopy.go:96:10: (*in).DeepCopyInto undefined (type *node has no field or method DeepCopyInto)

}

// MemcachedStatus defines the observed state of Memcached
// +k8s:openapi-gen=true
type resource struct {
	// Random Comment
	Requests memcpu `json:"requests,omitempty"`
	// Random Comment
	Limits memcpu `json:"limit,omitempty"`
}

//Random comment
// +k8s:openapi-gen=true
type memcpu struct {
	// Random Comment
	Memory string `json:"memory,omitempty"`
	// Random Comment
	CPU string `json:"cpu,omitempty"`
}

// ServiceGraphStatus defines the observed state of ServiceGraph
type ServiceGraphStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//Nodes []node `json:"node"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceGraph is the Schema for the servicegraphs API
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=servicegraph,scope=Cluster
type ServiceGraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceGraphSpec   `json:"spec,omitempty"`
	Status ServiceGraphStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceGraphList contains a list of ServiceGraph
type ServiceGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceGraph `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceGraph{}, &ServiceGraphList{})
}
