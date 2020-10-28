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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//Random comment
type endpoint struct {
	NodeName     string `json:"nodename,omitempty"`
	EndpointName string `json:"name,omitempty"`
	EndpointPath string `json:"path,omitempty"`
	CPU          string `json:"cpuload,omitempty"`
	Delay        uint   `json:"delay,omitempty"`
	CallOut      string `json:"callout,omitempty"`
}

// ServiceGraphSpec defines the desired state of ServiceGraph
type ServiceGraphSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Array for nodes in service-graph
	Nodes []*node `json:"nodes,omitempty"`

	// API endpoints for each node
	APIEndpoints []*endpoint `json:"apiendpoints,omitempty"`
}

// MemcachedStatus defines the observed state of Memcached
type node struct {
	Name          string   `json:"name,omitempty"`
	ContainerPort uint     `json:"containerport,omitempty"`
	NodePort      uint     `json:"nodeport,omitempty"`
	Replicas      uint     `json:"replicas,omitempty"`
	Resources     resource `json:"resources,omitempty"`
	//APIEndpoints  *[]endpoint `json:"apiendpoints,omitempty"` Nem működött, ha ide volt beágyazva
	// api/v1/zz_generated.deepcopy.go:96:10: (*in).DeepCopyInto undefined (type *node has no field or method DeepCopyInto)

}

// MemcachedStatus defines the observed state of Memcached
type resource struct {
	Requests memcpu `json:"requests,omitempty"`
	Limits   memcpu `json:"limit,omitempty"`
}

//Random comment
type memcpu struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

// ServiceGraphStatus defines the observed state of ServiceGraph
type ServiceGraphStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//Nodes []node `json:"node"`
}

// +kubebuilder:object:root=true

// ServiceGraph is the Schema for the servicegraphs API
// +kubebuilder:subresource:status
type ServiceGraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceGraphSpec   `json:"spec,omitempty"`
	Status ServiceGraphStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServiceGraphList contains a list of ServiceGraph
type ServiceGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceGraph `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceGraph{}, &ServiceGraphList{})
}
