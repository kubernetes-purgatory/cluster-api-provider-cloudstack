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

	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "github.com/aws/cluster-api-provider-cloudstack/api/v1beta1"
)

type ReconcilerMethdod (*CloudStackBaseReconciler) (context.Context) error

// CloudStackZoneReconciler reconciles a CloudStackZone object
type CloudStackZoneReconciler struct {
	CloudStackBaseReconciler
	CSISONetwork *infrav1.CloudStackIsolatedNetwork
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=cloudstackzones,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=cloudstackzones/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=cloudstackzones/finalizers,verbs=update
func (r *CloudStackZoneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Zone", req.Name, "namespace", req.Namespace)
	log.V(1).Info("Reconcile CloudStackZone")

	return r.runWith(ctx, req,
		r.fetchRelatedCRDs,
		r.CheckIfPaused,
		r.reconcileDelete,
		r.reconcile,
		r.patchChangesBackToAPI,
	)

	// // Get CloudStack Zone.
	// csZone := &infrav1.CloudStackZone{}
	// if err := r.Client.Get(ctx, req.NamespacedName, csZone); err != nil {
	// 	if client.IgnoreNotFound(err) == nil {
	// 		log.Info("Zone not found.")
	// 	}
	// 	return ctrl.Result{}, client.IgnoreNotFound(err)
	// }

	// // Get CloudStack cluster that owns the zone.
	// csCluster, err := csCtrlrUtils.GetOwnerCloudStackCluster(ctx, r.Client, csZone.ObjectMeta)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }

	// r.generateIsolatedNetwork(ctx, csZone, csCluster)

	// log.Info("Reconcile CloudStackZone completed successfully.", "spec", csZone.Spec)

	//csZone.Status.Ready = true

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CloudStackZoneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.CloudStackZone{}).
		Complete(r)
}

func (r *CloudStackZoneReconciler) generateIsolatedNetwork(
	ctx context.Context, zone *infrav1.CloudStackZone, csCluster *infrav1.CloudStackCluster) error {

	// csIsoNet := &infrav1.CloudStackIsolatedNetwork{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      zone.Spec.Name,
	// 		Namespace: zone.Namespace,
	// 		// Labels:      internal.ControlPlaneMachineLabelsForCluster(csCluster, csCluster.Name),
	// 		Annotations: map[string]string{},
	// 		OwnerReferences: []metav1.OwnerReference{
	// 			*metav1.NewControllerRef(zone, controlplanev1.GroupVersion.WithKind("CloudStackZone")),
	// 			*metav1.NewControllerRef(csCluster, controlplanev1.GroupVersion.WithKind("CloudStackCluster")),
	// 		},
	// 	},
	// 	Spec: infrav1.CloudStackIsolatedNetworkSpec{Name: zone.Spec.Network.Name},
	// }

	// if err := r.Client.Create(ctx, csIsoNet); err != nil {
	// 	return errors.Wrap(err, "failed to create machine")
	// }
	return nil
}

func (r *CloudStackZoneReconciler) reconcile(context.Context

// Fetches all CRDs relavent to reconciling a CloudStackCluster.
func (r *CloudStackZoneReconciler) fetchRelatedCRDs(
	ctx context.Context, req ctrl.Request) (retRes ctrl.Result, reterr error) {
	r.GetBaseCRDs(ctx, req)
	// Get CloudStackZones.
	// _, reterr = r.fetchIsolatedNetworks(ctx, req)
	// if reterr != nil {
	// 	return reconcile.Result{}, errors.Wrap(reterr, "error encountered fetching CloudStackZone(s)")
	// }
	return ctrl.Result{}, nil
}
