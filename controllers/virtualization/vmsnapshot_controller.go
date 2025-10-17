package virtualization

import (
	"context"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

const snapshotFinalizer = "virtualization.kubesphere.io/vmsnapshot-cleanup"

// VMSnapshotReconciler coordinates VM snapshots with CSI snapshot API.
type VMSnapshotReconciler struct {
	client.Client
	Snapshotter SnapshotBackend
	Clock       func() time.Time
}

//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=vmsnapshots,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=vmsnapshots/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=vmsnapshots/finalizers,verbs=update
//+kubebuilder:rbac:groups=snapshot.storage.k8s.io,resources=volumesnapshots;volumesnapshotcontents,verbs=*

func (r *VMSnapshotReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	snapshot := &virtualizationv1beta1.VMSnapshot{}
	if err := r.Get(ctx, req.NamespacedName, snapshot); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if snapshot.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(snapshot, snapshotFinalizer) {
			controllerutil.AddFinalizer(snapshot, snapshotFinalizer)
			if err := r.Update(ctx, snapshot); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(snapshot, snapshotFinalizer) {
			if err := r.Snapshotter.DeleteSnapshot(ctx, snapshot); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(snapshot, snapshotFinalizer)
			if err := r.Update(ctx, snapshot); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.Snapshotter.Sync(ctx, snapshot); err != nil {
		metav1.SetStatusCondition(&snapshot.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			Reason:             "SyncFailed",
			Message:            err.Error(),
			LastTransitionTime: metav1.NewTime(r.now()),
		})
		_ = r.Status().Update(ctx, snapshot)
		return ctrl.Result{}, err
	}

	metav1.SetStatusCondition(&snapshot.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "Synced",
		Message:            "Snapshot reconciled",
		LastTransitionTime: metav1.NewTime(r.now()),
	})
	snapshot.Status.ReadyToUse = true
	if err := r.Status().Update(ctx, snapshot); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VMSnapshotReconciler) now() time.Time {
	if r.Clock != nil {
		return r.Clock()
	}
	return time.Now()
}

func (r *VMSnapshotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&virtualizationv1beta1.VMSnapshot{}).
		Complete(r)
}

// SnapshotBackend abstracts CSI snapshot operations.
type SnapshotBackend interface {
	Sync(context.Context, *virtualizationv1beta1.VMSnapshot) error
	DeleteSnapshot(context.Context, *virtualizationv1beta1.VMSnapshot) error
}
