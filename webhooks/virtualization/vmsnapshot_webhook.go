package virtualization

import (
	"context"
	"encoding/json"

	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

// +kubebuilder:webhook:path=/validate-virtualization-kubesphere-io-v1beta1-vmsnapshot,mutating=false,failurePolicy=fail,sideEffects=None,groups=virtualization.kubesphere.io,resources=vmsnapshots,verbs=create;update,versions=v1beta1,name=vvmsnapshot.kb.io,admissionReviewVersions={v1,v1beta1}

// VMSnapshotWebhook ensures snapshot references are consistent.
type VMSnapshotWebhook struct {
	client.Client
	Decoder *admission.Decoder
}

func (w *VMSnapshotWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	snapshot := &virtualizationv1beta1.VMSnapshot{}
	if err := w.Decoder.Decode(req, snapshot); err != nil {
		return admission.Errored(400, err)
	}

	if snapshot.Spec.SourceRef.Namespace == "" {
		snapshot.Spec.SourceRef.Namespace = req.Namespace
	}

	if snapshot.Spec.SourceRef.Namespace != req.Namespace {
		return admission.Denied("sourceRef namespace must match snapshot namespace")
	}

	if len(snapshot.Spec.IncludedDisks) == 0 {
		return admission.Denied("includedDisks must not be empty to ensure deterministic backups")
	}

	raw, err := json.Marshal(snapshot)
	if err != nil {
		return admission.Errored(500, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, raw)
}

func (w *VMSnapshotWebhook) InjectDecoder(d *admission.Decoder) error {
	w.Decoder = d
	return nil
}
