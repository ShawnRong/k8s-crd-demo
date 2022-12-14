/*
Copyright 2022.

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
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"shawnrong.github.io/examplecontroller/controllers/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	websitev1alpha1 "shawnrong.github.io/examplecontroller/api/v1alpha1"
)

// WebsiteReconciler reconciles a Website object
type WebsiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=website.shawnrong.github.io,resources=websites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=website.shawnrong.github.io,resources=websites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=website.shawnrong.github.io,resources=websites/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Website object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *WebsiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	//Load Resource by name
	var website websitev1alpha1.Website
	if err := r.Get(ctx, req.NamespacedName, &website); err != nil {
		log.Error(err, "unable to fetch Website")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle deployment
	deployment := utils.NewDeployment(&website)
	if err := controllerutil.SetControllerReference(&website, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if deployment exist, if not, then create
	d := &appv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, d); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, deployment); err != nil {
				log.Error(err, "fail to create deployment")
				return ctrl.Result{}, err
			}
		}
	} else {
		if err := r.Update(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Create service
	service := utils.NewService(&website)
	if err := controllerutil.SetControllerReference(&website, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	s := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      website.Name,
		Namespace: website.Namespace,
	}, s); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, service); err != nil {
				log.Error(err, "fail to create service")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&websitev1alpha1.Website{}).
		Owns(&appv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
