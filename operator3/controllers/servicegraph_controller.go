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
	// Basic go libraries
	"context"
	"fmt"
	"strings"

	// Manual imports
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	// Imports from framework
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
func (r *ServiceGraphReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("servicegraph", req.NamespacedName)

	// your logic here
	servicegraph := &onlabv2.ServiceGraph{}
	err := r.Get(ctx, req.NamespacedName, servicegraph)

	// printServiceGraph(servicegraph)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Servicegraph resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ServiceGraph resource")
		return ctrl.Result{}, err
	}

	for _, node := range servicegraph.Spec.Nodes {
		// Check if the deployment for the node already exists, if not create a new one
		found := &appsv1.Deployment{}

		err = r.Get(ctx, types.NamespacedName{Name: node.Name, Namespace: "default"}, found)
		if err != nil && errors.IsNotFound(err) {
			//fmt.Printf("######### CREATE: %d node type: %T\n", i, node)
			// Define a new deployment for the node
			dep := r.deploymentForNode(node, servicegraph)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		// Ensure the deployment size is the same as the spec
		size := int32(node.Replicas)
		if *found.Spec.Replicas != size {
			found.Spec.Replicas = &size
			err = r.Update(ctx, found)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		svc := r.serviceForNode(node, servicegraph)
		_ = r.Create(ctx, svc)

	}

	return ctrl.Result{}, nil
}

func (r *ServiceGraphReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&onlabv2.ServiceGraph{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func printServiceGraph(servicegraph *onlabv2.ServiceGraph) {
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
}

func (r *ServiceGraphReconciler) deploymentForNode(node *onlabv2.Node, sg *onlabv2.ServiceGraph) *appsv1.Deployment {
	ls := r.labelsForNode(node)
	replicas := int32(node.Replicas)
	args := r.createCommandForNode(node)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: "default", // TODO: use correct ns
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "tuti/service-graph-simulator:latest",
						Name:    "servicenode",
						Command: args,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(node.ContainerPort),
							Name:          "listen",
						}},
					}},
				},
			},
		},
	}

	fmt.Printf("Created Deployment: %+v", dep)

	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(sg, dep, r.Scheme)
	return dep
}

// labelsForNode returns the labels for selecting the resources
func (r *ServiceGraphReconciler) labelsForNode(node *onlabv2.Node) map[string]string {
	return map[string]string{"app": "servicegraph", "node": node.Name}
}

func (r *ServiceGraphReconciler) serviceForNode(node *onlabv2.Node, sg *onlabv2.ServiceGraph) *corev1.Service {
	ls := r.labelsForNode(node)
	cPort := int32(node.ContainerPort)
	nPort := int32(node.NodePort)
	sType := corev1.ServiceTypeClusterIP

	if nPort != 0 {
		sType = corev1.ServiceTypeNodePort
	}

	port := []corev1.ServicePort{
		corev1.ServicePort{
			Port:     cPort,
			NodePort: nPort,
		},
	}

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports:    port,
			Type:     sType,
		},
	}

	fmt.Printf("\nSERVICE #### Nodeport: %d ######\n%+v\n", node.NodePort, service)

	return service
}

func (r *ServiceGraphReconciler) createCommandForNode(node *onlabv2.Node) []string {
	// eg: -name Backend -delay 90 -port 9999 -cpu 90 -memory 900 -endpoint-url /read -endpoint-cpu 99 -endpoint-delay 192 -endpoint-url /index -endpoint-cpu 22 -endpoint-delay 111 -endpoint-call='"back-end:9898/staus__front-end:9876/health"' -endpoint-call "database:1234/asd?q=user1"
	var cmd []string

	cmd = append(cmd, "/app/main") //start my application

	// application parameters
	cmd = append(cmd, fmt.Sprintf("-name='%s'", node.Name))
	cmd = append(cmd, fmt.Sprintf("-port=%d", node.ContainerPort))
	cmd = append(cmd, fmt.Sprintf("-cpu=%d", node.Resources.CPU))
	cmd = append(cmd, fmt.Sprintf("-memory=%d", node.Resources.Memory))

	for _, ep := range node.Endpoints {
		// convert list of callouts to "call-out1__call-out2..." format
		var tmpArray []string
		for _, ca := range ep.CallOuts {
			tmpArray = append(tmpArray, string(ca.URL))
		}
		callOutParsed := strings.Join(tmpArray, "__")

		cmd = append(cmd, fmt.Sprintf("-endpoint-url='%s'", ep.Path))
		cmd = append(cmd, fmt.Sprintf("-endpoint-delay=%d", ep.Delay))
		cmd = append(cmd, fmt.Sprintf("-endpoint-call='%s'", callOutParsed))
		cmd = append(cmd, fmt.Sprintf("-endpoint-cpu=%d", ep.CPULoad))
	}

	fmt.Printf("Command to run in container: %s", cmd)
	return cmd
}
