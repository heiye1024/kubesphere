package virtualization

import (
	"context"
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

// +kubebuilder:webhook:path=/validate-virtualization-kubesphere-io-v1beta1-virtualmachine,mutating=false,failurePolicy=fail,sideEffects=None,groups=virtualization.kubesphere.io,resources=virtualmachines,verbs=create;update,versions=v1beta1,name=vvirtualmachine.kb.io,admissionReviewVersions={v1,v1beta1}
// +kubebuilder:webhook:path=/mutate-virtualization-kubesphere-io-v1beta1-virtualmachine,mutating=true,failurePolicy=fail,sideEffects=None,groups=virtualization.kubesphere.io,resources=virtualmachines,verbs=create;update,versions=v1beta1,name=mvirtualmachine.kb.io,admissionReviewVersions={v1,v1beta1}

// VirtualMachineWebhook validates VirtualMachine resources against cluster capabilities.
type VirtualMachineWebhook struct {
	client.Client
	Decoder *admission.Decoder
}

func (w *VirtualMachineWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	vm := &virtualizationv1beta1.VirtualMachine{}
	if err := w.Decoder.Decode(req, vm); err != nil {
		return admission.Errored(400, err)
	}

	if err := w.defaultVirtualMachine(vm); err != nil {
		return admission.Errored(400, err)
	}

	if req.AdmissionRequest.Operation == admissionv1.Create || req.AdmissionRequest.Operation == admissionv1.Update {
		if err := w.validateSRIOV(ctx, vm); err != nil {
			return admission.Denied(err.Error())
		}
		if err := w.validateHugePages(ctx, vm); err != nil {
			return admission.Denied(err.Error())
		}
	}

	marshaled, err := jsonMarshal(vm)
	if err != nil {
		return admission.Errored(500, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}

func (w *VirtualMachineWebhook) defaultVirtualMachine(vm *virtualizationv1beta1.VirtualMachine) error {
	if vm.Spec.NUMAPolicy == "" {
		vm.Spec.NUMAPolicy = "none"
	}
	if vm.Spec.PowerState == "" {
		vm.Spec.PowerState = virtualizationv1beta1.PowerStateRunning
	}
	return nil
}

func (w *VirtualMachineWebhook) validateSRIOV(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine) error {
	for _, net := range vm.Spec.Nets {
		if net.Type != "sriov" {
			continue
		}
		if net.NADRef == nil {
			return fmt.Errorf("SR-IOV network requires nadRef")
		}
		nodeList := &corev1.NodeList{}
		if err := w.List(ctx, nodeList, client.MatchingLabels{"sriov.capable": "true"}); err != nil {
			return err
		}
		if len(nodeList.Items) == 0 {
			return fmt.Errorf("no sriov.capable=true nodes available for SR-IOV network")
		}
	}
	return nil
}

func (w *VirtualMachineWebhook) validateHugePages(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine) error {
	if vm.Spec.Hugepages == "" {
		return nil
	}
	nodes := &corev1.NodeList{}
	if err := w.List(ctx, nodes); err != nil {
		return err
	}
	for _, node := range nodes.Items {
		if quantity, ok := node.Status.Allocatable[corev1.ResourceName("hugepages-"+vm.Spec.Hugepages)]; ok {
			if !quantity.IsZero() {
				return nil
			}
		}
	}
	return fmt.Errorf("no nodes advertise hugepages %s", vm.Spec.Hugepages)
}

func (w *VirtualMachineWebhook) InjectDecoder(d *admission.Decoder) error {
	w.Decoder = d
	return nil
}

// jsonMarshal is a small indirection for testability.
var jsonMarshal = func(obj runtime.Object) ([]byte, error) {
	return json.Marshal(obj)
}
