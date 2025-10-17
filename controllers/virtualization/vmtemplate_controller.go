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

const templateFinalizer = "virtualization.kubesphere.io/vmtemplate-cleanup"

// VMTemplateReconciler ensures template inventory stays healthy.
type VMTemplateReconciler struct {
	client.Client
	Catalog CatalogBackend
	Clock   func() time.Time
}

//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=vmtemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=vmtemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=vmtemplates/finalizers,verbs=update

func (r *VMTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	template := &virtualizationv1beta1.VMTemplate{}
	if err := r.Get(ctx, req.NamespacedName, template); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if template.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(template, templateFinalizer) {
			controllerutil.AddFinalizer(template, templateFinalizer)
			if err := r.Update(ctx, template); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(template, templateFinalizer) {
			if err := r.Catalog.RemoveTemplate(ctx, template); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(template, templateFinalizer)
			if err := r.Update(ctx, template); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.Catalog.Sync(ctx, template); err != nil {
		metav1.SetStatusCondition(&template.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			Reason:             "SyncFailed",
			Message:            err.Error(),
			LastTransitionTime: metav1.NewTime(r.now()),
		})
		_ = r.Status().Update(ctx, template)
		return ctrl.Result{}, err
	}

	metav1.SetStatusCondition(&template.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "Synced",
		Message:            "Template synchronized",
		LastTransitionTime: metav1.NewTime(r.now()),
	})
	if err := r.Status().Update(ctx, template); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VMTemplateReconciler) now() time.Time {
	if r.Clock != nil {
		return r.Clock()
	}
	return time.Now()
}

func (r *VMTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&virtualizationv1beta1.VMTemplate{}).
		Complete(r)
}

// CatalogBackend stores templates for UI consumption.
type CatalogBackend interface {
	Sync(context.Context, *virtualizationv1beta1.VMTemplate) error
	RemoveTemplate(context.Context, *virtualizationv1beta1.VMTemplate) error
}
