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

	onlabv1 "example.com/m/v2/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	//corev1 "k8s.io/api/core/v1"
)

// ServiceGraphReconciler reconciles a ServiceGraph object
type ServiceGraphReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=onlab.my.domain,resources=servicegraphs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=onlab.my.domain,resources=servicegraphs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;pat$
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
func (r *ServiceGraphReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	reqLogger := r.Log.WithValues("servicegraph", req.NamespacedName)
	reqLogger.Info("Reconciling servicegraph")

	// your logic here
	servicegraph := &onlabv1.ServiceGraph{}
	ctx := context.TODO()
	_ = r.Get(ctx, req.NamespacedName, servicegraph)
	fmt.Printf("\n[REQUEST]\t%+v\n", req)
	fmt.Printf("\n[NODES]\t%+v\n", (servicegraph.ObjectMeta.Annotations))
	fmt.Printf("\n[SPEC-----------------]\nNodes: %+v\n", servicegraph.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]["specs"])

	for node := range servicegraph.Spec.Nodes {
		deployment := &appsv1.Deployment{}
		_ = r.Client.Get(context.TODO(), req.NamespacedName, deployment)
		fmt.Printf("\n[NODE]:\t%+v\n", node)
		reqLogger.Info("-------------------------------")
		fmt.Printf("\n[DEPLOYMENT]:\t%+v\n", deployment)
		//if err != nil && errors.IsNotFound(err) {
		// Define a new Deployment
		// dep := r.deploymentForMemcached(memcached)
		// reqLogger.Info("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		// err = r.Client.Create(context.TODO(), dep)
		// if err != nil {
		// 	reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		// 	return ctrl.Result{}, err
		// }
		// Deployment created successfully - return and requeue
		// NOTE: that the requeue is made with the purpose to provide the deployment object for the next step to ensure the deployment size is the same as the spec.
		// Also, you could GET the deployment object again instead of requeue if you wish. See more over it here: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Reco$
		// 	return reconcile.Result{Requeue: true}, nil
		// } else if err != nil {
		// 	reqLogger.Error(err, "Failed to get Deployment.")
		// 	return ctrl.Result{}, err
		// }
	}

	return ctrl.Result{}, nil
}

func (r *ServiceGraphReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&onlabv1.ServiceGraph{}).
		Complete(r)
}

// deploymentForMemcached returns a memcached Deployment object
// func (r *ServiceGraphReconciler) deploymentForMemcached(m *onlabv1.ServiceGraph) *appsv1.Deployment {
// 	//ls := labelsForMemcached(m.Name)
// 	replicas := m.Spec.Nodes
// 	image := "tuti/service-greph-simulator:latest"
// 	appName := "Asd" //m.Spec.Name

// 	log := r.Log.WithValues("memcached", "createDeployment")
// 	// for _, node := range m.Spec.MeshNodes {
// 	// 	s := fmt.Sprint("Name: ", node.Name, " --> ", node.Port)
// 	// 	log.Info(s)
// 	// }

// 	// dep := &appsv1.Deployment{
// 	// 	ObjectMeta: v1.ObjectMeta{
// 	// 		Name:      m.Name,
// 	// 		Namespace: m.Namespace,
// 	// 	},
// 	// 	Spec: appsv1.DeploymentSpec{
// 	// 		Replicas: &replicas,
// 	// 		Selector: &v1.LabelSelector{
// 	// 			MatchLabels: ls,
// 	// 		},
// 	// 		Template: corev1.PodTemplateSpec{
// 	// 			ObjectMeta: v1.ObjectMeta{
// 	// 				Labels: ls,
// 	// 			},
// 	// 			Spec: corev1.PodSpec{
// 	// 				Containers: []corev1.Container{{
// 	// 					Image:   image, //"memcached:1.4.36-alpine",
// 	// 					Name:    appName,
// 	// 					Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
// 	// 					Ports: []corev1.ContainerPort{{
// 	// 						ContainerPort: 11211,
// 	// 						Name:          "memcached",
// 	// 					}},
// 	// 				}},
// 	// 			},
// 	// 		},
// 	// 	},
// 	// }

// 	// // Set Memcached instance as the owner of the Deployment.
// 	// ctrl.SetControllerReference(m, dep, r.Scheme) //todo check how to get the schema

// 	return dep
// }
