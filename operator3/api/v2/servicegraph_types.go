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

package v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceGraphSpec defines the desired state of ServiceGraph
type ServiceGraphSpec struct {
	// +kubebuilder:validation:Minimum=0
	// Size is the size of the deployment
	//Size int32 `json:"size"`

	// kubebuilder:validation:Type=string
	// Name for the deployment
	//Name string `json:"name"`

	// +kubebuilder:validation:Required
	// Nodes contain the configs to service
	Nodes []*Node `json:"nodes,omitempty"`
}

// ServiceGraphStatus defines the observed state of ServiceGraph
type ServiceGraphStatus struct {
	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ServiceGraph is the Schema for the servicegraphs API
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

// Here come my structs

// Node struct. Contains specification of a node
type Node struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="[a-z-]*"
	// Name of the node
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Number of replica to run
	Replicas uint `json:"replicas"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Container port to open and listen
	ContainerPort uint `json:"port"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// If the service will listen on node port
	NodePort uint `json:"nodePort"`

	// +kubebuilder:validation:Required
	// Resource to consume
	Resources Resource `json:"resources"`

	// +kubebuilder:validation:Required
	// Setup and configure endpoints to node
	Endpoints []Endpoint `json:"endpoints"`
}

// Resource structure
type Resource struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Memory to use (kB)
	Memory uint `json:"memory"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// CPU to use (mCPU)
	CPU uint `json:"cpu"`
}

// Endpoint structure and behaviour
type Endpoint struct {

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// Regex for pattern eg: /index
	// +kubebuilder:validation:Pattern="/[a-z]*"
	// Path to listen and answer
	Path string `json:"path"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// CPU usage until request
	CPULoad uint `json:"cpuLoad"`

	//TODO: Memory load

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// Delay / process time to "serve" request (ms)
	Delay uint `json:"delay"`

	// +kubebuilder:validation:Optional
	// Regex for pattern eg: db-user:890/read?from=table#site
	// __+kubebuilder:validation:Pattern="[a-z-]*:[0-9]*/[a-z?=#]*"
	// Ask further services
	CallOuts []string `json:"callouts,omitempty"`
}
