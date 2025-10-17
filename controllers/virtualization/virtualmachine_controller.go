package virtualization

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

const (
	vmFinalizer = "virtualization.kubesphere.io/vm-cleanup"
)

// VirtualMachineReconciler orchestrates VirtualMachine lifecycle and power states.
type VirtualMachineReconciler struct {
	client.Client
	KubeVirtClient VMBackendClient
	CDIClient      DataVolumeClient
	MultusClient   NADClient
	Clock          func() time.Time
}

func (r *VirtualMachineReconciler) defaultClock() time.Time {
	if r.Clock != nil {
		return r.Clock()
	}
	return time.Now()
}

//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=virtualmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=virtualmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=virtualization.kubesphere.io,resources=virtualmachines/finalizers,verbs=update
//+kubebuilder:rbac:groups=kubevirt.io,resources=virtualmachines;virtualmachineinstances,verbs=*
//+kubebuilder:rbac:groups=cdi.kubevirt.io,resources=datavolumes,verbs=*
//+kubebuilder:rbac:groups=k8s.cni.cncf.io,resources=network-attachment-definitions,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *VirtualMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	vm := &virtualizationv1beta1.VirtualMachine{}
	if err := r.Get(ctx, req.NamespacedName, vm); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if vm.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(vm, vmFinalizer) {
			controllerutil.AddFinalizer(vm, vmFinalizer)
			if err := r.Update(ctx, vm); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(vm, vmFinalizer) {
			if err := r.teardownVM(ctx, vm); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(vm, vmFinalizer)
			if err := r.Update(ctx, vm); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.ensureBackingResources(ctx, vm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.syncPowerState(ctx, vm); err != nil {
		r.setCondition(ctx, vm, "Ready", metav1.ConditionFalse, "PowerSyncFailed", err.Error())
		return ctrl.Result{}, err
	}

	r.setCondition(ctx, vm, "Ready", metav1.ConditionTrue, "PowerSynced", "VirtualMachine power reconciled")
	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *VirtualMachineReconciler) ensureBackingResources(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine) error {
	for _, disk := range vm.Spec.Disks {
		switch disk.Type {
		case "system", "data":
			if err := r.CDIClient.EnsureDataVolume(ctx, vm, disk); err != nil {
				return err
			}
		case "ephemeral":
			// No-op for ephemeral disk
		default:
			return fmt.Errorf("unsupported disk type %s", disk.Type)
		}
	}

	for _, net := range vm.Spec.Nets {
		if net.Type == "sriov" {
			if err := r.MultusClient.ValidateSRIOVNetwork(ctx, vm.Namespace, net); err != nil {
				return err
			}
		}
	}

	return r.KubeVirtClient.EnsureVM(ctx, vm)
}

func (r *VirtualMachineReconciler) syncPowerState(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine) error {
	switch vm.Spec.PowerState {
	case virtualizationv1beta1.PowerStateRunning:
		if err := r.KubeVirtClient.PowerOn(ctx, vm); err != nil {
			return err
		}
	case virtualizationv1beta1.PowerStateStopped:
		if err := r.KubeVirtClient.PowerOff(ctx, vm); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown power state %s", vm.Spec.PowerState)
	}

	vm.Status.PowerState = vm.Spec.PowerState
	return r.Status().Update(ctx, vm)
}

func (r *VirtualMachineReconciler) setCondition(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine, cType string, status metav1.ConditionStatus, reason, msg string) {
	condition := metav1.Condition{
		Type:               cType,
		Status:             status,
		Reason:             reason,
		Message:            msg,
		LastTransitionTime: metav1.NewTime(r.defaultClock()),
	}
	metav1.SetStatusCondition(&vm.Status.Conditions, condition)
	_ = r.Status().Update(ctx, vm)
}

func (r *VirtualMachineReconciler) teardownVM(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine) error {
	if err := r.KubeVirtClient.Cleanup(ctx, vm); err != nil {
		return err
	}
	return r.CDIClient.DeleteOwnedVolumes(ctx, vm)
}

func (r *VirtualMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&virtualizationv1beta1.VirtualMachine{}).
		Complete(r)
}

// VMBackendClient abstracts KubeVirt interactions.
type VMBackendClient interface {
	EnsureVM(context.Context, *virtualizationv1beta1.VirtualMachine) error
	PowerOn(context.Context, *virtualizationv1beta1.VirtualMachine) error
	PowerOff(context.Context, *virtualizationv1beta1.VirtualMachine) error
	Cleanup(context.Context, *virtualizationv1beta1.VirtualMachine) error
}

// DataVolumeClient abstracts CDI interactions.
type DataVolumeClient interface {
	EnsureDataVolume(context.Context, *virtualizationv1beta1.VirtualMachine, virtualizationv1beta1.VirtualMachineDisk) error
	DeleteOwnedVolumes(context.Context, *virtualizationv1beta1.VirtualMachine) error
}

// NADClient validates NAD resources.
type NADClient interface {
	ValidateSRIOVNetwork(context.Context, string, virtualizationv1beta1.VirtualMachineNetwork) error
}
