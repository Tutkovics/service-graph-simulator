/*


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
)

// ServicegraphReconciler reconciles a Servicegraph object
type ServicegraphReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=onlab.example.com,resources=servicegraphs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=onlab.example.com,resources=servicegraphs/status,verbs=get;update;patch

func (r *ServicegraphReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("servicegraph", req.NamespacedName)
	reqLogger.Info("Start Loop")

	// your logic here
	servicegraph := &onlabv1.Servicegraph{}
	_ = r.Get(ctx, req.NamespacedName, servicegraph)
	fmt.Printf("\n[SPEC-----------------]\nType: %T --- %+v\n", servicegraph, servicegraph)
	fmt.Printf("\n[ASDASDASD]%+v", servicegraph.Spec)

	return ctrl.Result{}, nil
}

func (r *ServicegraphReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&onlabv1.Servicegraph{}).
		Complete(r)
}
