package virtualization

import (
	"context"
	"encoding/json"

	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

// +kubebuilder:webhook:path=/validate-virtualization-kubesphere-io-v1beta1-vmtemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=virtualization.kubesphere.io,resources=vmtemplates,verbs=create;update,versions=v1beta1,name=vvmtemplate.kb.io,admissionReviewVersions={v1,v1beta1}

// VMTemplateWebhook enforces template constraints.
type VMTemplateWebhook struct {
	client.Client
	Decoder *admission.Decoder
}

func (w *VMTemplateWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	template := &virtualizationv1beta1.VMTemplate{}
	if err := w.Decoder.Decode(req, template); err != nil {
		return admission.Errored(400, err)
	}

	if template.Spec.Parameters.CPU == "" || template.Spec.Parameters.Memory == "" {
		return admission.Denied("parameters.cpu and parameters.memory must be specified")
	}

	if template.Spec.Constraints.MaxCPU != "" && template.Spec.Constraints.MinCPU != "" {
		if template.Spec.Constraints.MaxCPU < template.Spec.Constraints.MinCPU {
			return admission.Denied("constraints.maxCPU must be >= minCPU")
		}
	}

	raw, err := json.Marshal(template)
	if err != nil {
		return admission.Errored(500, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, raw)
}

func (w *VMTemplateWebhook) InjectDecoder(d *admission.Decoder) error {
	w.Decoder = d
	return nil
}
