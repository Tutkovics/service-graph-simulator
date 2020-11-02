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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	onlabv2 "project.msc/m/v2/api/v2"
)

// ServiceGraphReconciler reconciles a ServiceGraph object
type ServiceGraphReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=onlab.project.msc,resources=servicegraphs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=onlab.project.msc,resources=servicegraphs/status,verbs=get;update;patch

func (r *ServiceGraphReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("servicegraph", req.NamespacedName)

	// your logic here
	servicegraph := &onlabv2.ServiceGraph{}
	_ = r.Get(ctx, req.NamespacedName, servicegraph)

	fmt.Printf("[RECONCILE]\tName: %s\n", servicegraph.Name)
	fmt.Printf("[RECONCILE]\t%+v\n", servicegraph)
	for i, node := range servicegraph.Spec.Nodes {
		fmt.Printf("\t\t--- %d ---\n", i)
		fmt.Printf("\t\tname: %s\n", node.Name)
		fmt.Printf("\t\tport: %d\n", node.ContainerPort)
		fmt.Printf("\t\tnode port: %d\n", node.NodePort)
		fmt.Printf("\t\tresource: %d (kB), %d (mCPU)\n", node.Resources.Memory, node.Resources.CPU)
		fmt.Printf("\t\t# of endpoints: %d\n", len(node.Endpoints))
		for _, ep := range node.Endpoints {
			fmt.Printf("\t\t\tpath: %s\n", ep.Path)
			fmt.Printf("\t\t\tcpu: %d\n", ep.CPULoad)
			fmt.Printf("\t\t\tdelay: %d\n", ep.Delay)
			for _, co := range ep.CallOuts {
				fmt.Printf("\t\t\t\tcall out: %s\n", co)
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *ServiceGraphReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&onlabv2.ServiceGraph{}).
		Complete(r)
}
