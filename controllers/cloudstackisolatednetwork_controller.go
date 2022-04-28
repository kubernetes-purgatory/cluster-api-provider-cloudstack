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

	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "github.com/aws/cluster-api-provider-cloudstack/api/v1beta1"
	csCtrlrUtils "github.com/aws/cluster-api-provider-cloudstack/controllers/utils"
	"github.com/aws/cluster-api-provider-cloudstack/pkg/cloud"
)

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=cloudstackisolatednetworks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=cloudstackisolatednetworks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=cloudstackisolatednetworks/finalizers,verbs=update

// CloudStackIsoNetReconciler reconciles a CloudStackZone object
type CloudStackIsoNetReconciler struct {
	csCtrlrUtils.ReconcilerBase
}

// CloudStackZoneReconciliationRunner is a ReconciliationRunner with extensions specific to CloudStackCluster reconciliation.
type CloudStackIsoNetReconciliationRunner struct {
	csCtrlrUtils.ReconciliationRunner
	Zone                  *infrav1.CloudStackZone
	ReconciliationSubject *infrav1.CloudStackIsolatedNetwork
	CSUser                cloud.Client
}

// Initialize a new CloudStackIsoNet reconciliation runner with concrete types and initialized member fields.
func NewCSIsoNetReconciliationRunner() *CloudStackIsoNetReconciliationRunner {
	// Set concrete type and init pointers.
	r := &CloudStackIsoNetReconciliationRunner{ReconciliationSubject: &infrav1.CloudStackIsolatedNetwork{}}
	r.Zone = &infrav1.CloudStackZone{}
	// Setup the base runner. Initializes pointers and links reconciliation methods.
	r.ReconciliationRunner = csCtrlrUtils.NewRunner(r, r.ReconciliationSubject)
	return r
}

func (reconciler *CloudStackIsoNetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, retErr error) {
	return NewCSIsoNetReconciliationRunner().
		UsingBaseReconciler(reconciler.ReconcilerBase).
		ForRequest(req).
		WithRequestCtx(ctx).
		RunBaseReconciliationStages()
}

func (r *CloudStackIsoNetReconciliationRunner) Reconcile() (retRes ctrl.Result, retErr error) {
	if res, err := r.RequeueIfMissingBaseCRDs(); r.ShouldReturn(res, err) {
		return res, err
	}
	if res, err := r.GetParent(r.ReconciliationSubject, r.Zone)(); r.ShouldReturn(res, err) {
		return res, err
	}
	controllerutil.AddFinalizer(r.ReconciliationSubject, infrav1.IsolatedNetworkFinalizer)

	// Setup isolated network, endpoint, egress, and load balancing.
	// Set endpoint of CloudStackCluster if it is not currently set. (uses patcher to do so)
	if csClusterPatcher, err := patch.NewHelper(r.CSCluster, r.Client); err != nil {
		return r.ReturnWrappedError(retErr, "error encountered while setting up CloudStackCluster patcher")
	} else if err := r.CS.GetOrCreateIsolatedNetwork(r.Zone, r.ReconciliationSubject, r.CSCluster); err != nil {
		return ctrl.Result{}, err
	} else if err := csClusterPatcher.Patch(r.RequestCtx, r.CSCluster); err != nil {
		return r.ReturnWrappedError(err, "error encountered when patching endpoint update to CloudStackCluster")
	}

	r.ReconciliationSubject.Status.Ready = true
	return ctrl.Result{}, nil
}

func (r *CloudStackIsoNetReconciliationRunner) ReconcileDelete() (retRes ctrl.Result, retErr error) {
	r.CS.DisposeIsoNetResources(r.Zone, r.CSCluster)
	controllerutil.RemoveFinalizer(r.ReconciliationSubject, infrav1.IsolatedNetworkFinalizer)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (reconciler *CloudStackIsoNetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.CloudStackIsolatedNetwork{}).
		Complete(reconciler)
}
