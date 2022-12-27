package controllers

import (
	"context"
	"fmt"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	webv1alpha1 "tutorial.kubebuilder.io/application-set-operator/api/v1alpha1"
)

// ApplicationSetReconciler reconciles a ApplicationSet object
type ApplicationSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=web.tutorial.kubebuilder.io,resources=applicationsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployment,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=web.tutorial.kubebuilder.io,resources=applicationsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=web.tutorial.kubebuilder.io,resources=applicationsets/finalizers,verbs=update

func (r *ApplicationSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// fetch ApplicationSet resource from cache
	var appSet webv1alpha1.ApplicationSet

	if err := r.Get(ctx, req.NamespacedName, &appSet); err != nil {
		log.Error(err, "unable to fetch ApplicationSet")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// create child resources
	appSetObjects := appSet.ConstructAppSet()
	created := true

	for _, app := range appSetObjects {
		// set owner reference for resources
		if err := ctrl.SetControllerReference(&appSet, app, r.Scheme); err != nil {
			log.Error(err, "unable to set owner reference")
			return ctrl.Result{}, err
		}

		// create and/or update ApplicationSet resources
		if err := r.Create(ctx, app); err != nil {
			created = false
			if apiErrors.IsAlreadyExists(err) {
				if err := r.Update(ctx, app); err != nil {
					if apiErrors.IsInvalid(err) {
						msg := fmt.Sprintf("invalid fields: %s", err)
						log.V(1).Info(msg)
					} else {
						log.Error(err, "unable to update ApplicationSet resources")
						return ctrl.Result{}, err
					}
				}
			}

			log.Error(err, "unable to create app resources")
			return ctrl.Result{}, err
		}
	}

	if created {
		appSet.Status.Created = fmt.Sprintf(
			"%s\n", time.Now().UTC(),
		)

		if err := r.Status().Update(ctx, &appSet); err != nil {
			log.Error(err, "unable to update ApplicationSet status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1alpha1.ApplicationSet{}).
		Complete(r)
}
