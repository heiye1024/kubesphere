package main

import (
	"context"
	"flag"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	virtualizationv1alpha1 "kubesphere.io/kubesphere/api/virtualization/v1alpha1"
	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
	"kubesphere.io/kubesphere/controllers/virtualization"
)

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	cfg := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = virtualizationv1alpha1.AddToScheme(scheme)
	_ = virtualizationv1beta1.AddToScheme(scheme)

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{MetricsBindAddress: metricsAddr, Scheme: scheme})
	if err != nil {
		log.Fatalf("failed to create manager: %v", err)
	}

	vmReconciler := &virtualization.VirtualMachineReconciler{
		Client:         mgr.GetClient(),
		KubeVirtClient: noopVMBackend{},
		CDIClient:      noopDataVolume{},
		MultusClient:   noopNAD{},
	}
	if err := vmReconciler.SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create VM controller: %v", err)
	}
	snapshotReconciler := &virtualization.VMSnapshotReconciler{Client: mgr.GetClient(), Snapshotter: noopSnapshot{}}
	if err := snapshotReconciler.SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create snapshot controller: %v", err)
	}
	templateReconciler := &virtualization.VMTemplateReconciler{Client: mgr.GetClient(), Catalog: noopCatalog{}}
	if err := templateReconciler.SetupWithManager(mgr); err != nil {
		log.Fatalf("unable to create template controller: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Fatalf("unable to set up health check: %v", err)
	}

	log.Println("starting virtualization controller manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatalf("manager exited: %v", err)
	}
}

type noopVMBackend struct{}

func (noopVMBackend) EnsureVM(context.Context, *virtualizationv1beta1.VirtualMachine) error {
	return nil
}
func (noopVMBackend) PowerOn(context.Context, *virtualizationv1beta1.VirtualMachine) error {
	return nil
}
func (noopVMBackend) PowerOff(context.Context, *virtualizationv1beta1.VirtualMachine) error {
	return nil
}
func (noopVMBackend) Cleanup(context.Context, *virtualizationv1beta1.VirtualMachine) error {
	return nil
}

type noopDataVolume struct{}

func (noopDataVolume) EnsureDataVolume(context.Context, *virtualizationv1beta1.VirtualMachine, virtualizationv1beta1.VirtualMachineDisk) error {
	return nil
}
func (noopDataVolume) DeleteOwnedVolumes(context.Context, *virtualizationv1beta1.VirtualMachine) error {
	return nil
}

type noopNAD struct{}

func (noopNAD) ValidateSRIOVNetwork(context.Context, string, virtualizationv1beta1.VirtualMachineNetwork) error {
	return nil
}

type noopSnapshot struct{}

func (noopSnapshot) Sync(context.Context, *virtualizationv1beta1.VMSnapshot) error { return nil }
func (noopSnapshot) DeleteSnapshot(context.Context, *virtualizationv1beta1.VMSnapshot) error {
	return nil
}

type noopCatalog struct{}

func (noopCatalog) Sync(context.Context, *virtualizationv1beta1.VMTemplate) error { return nil }
func (noopCatalog) RemoveTemplate(context.Context, *virtualizationv1beta1.VMTemplate) error {
	return nil
}
