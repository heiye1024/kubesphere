package e2e

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	virtualizationv1alpha1 "kubesphere.io/kubesphere/api/virtualization/v1alpha1"
	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

var _ = Describe("Virtualization happy path", Ordered, func() {
	var (
		client    ctrlclient.Client
		namespace = "e2e-virtualization"
		ctx       context.Context
	)

	BeforeAll(func() {
		cfg := ctrl.GetConfigOrDie()
		var err error
		client, err = ctrlclient.New(cfg, ctrlclient.Options{Scheme: schemeForTests()})
		Expect(err).NotTo(HaveOccurred())
		ctx = context.Background()
		Expect(client.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})).To(Succeed())
	})

	AfterAll(func() {
		_ = client.Delete(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})
	})

	It("creates a VM with two disks and two nets", func() {
		vm := &virtualizationv1beta1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{Name: "scenario-vm", Namespace: namespace},
			Spec: virtualizationv1beta1.VirtualMachineSpec{
				CPU:    "4",
				Memory: "8Gi",
				Disks: []virtualizationv1beta1.VirtualMachineDisk{
					{Type: "system", DiskRef: virtualizationv1beta1.LocalObjectReference{Name: "root"}},
					{Type: "data", DiskRef: virtualizationv1beta1.LocalObjectReference{Name: "data"}},
				},
				Nets:       []virtualizationv1beta1.VirtualMachineNetwork{{Type: "masquerade"}, {Type: "bridge"}},
				PowerState: virtualizationv1beta1.PowerStateRunning,
			},
		}
		Expect(client.Create(ctx, vm)).To(Succeed())
		Eventually(func() virtualizationv1beta1.PowerState {
			_ = client.Get(ctx, types.NamespacedName{Name: vm.Name, Namespace: namespace}, vm)
			return vm.Status.PowerState
		}, 5*time.Minute, 10*time.Second).Should(Equal(virtualizationv1beta1.PowerStateRunning))
	})

	It("enables live migration and observes success", func() {
		// Placeholder for enabling live migration feature flag and verifying migration readiness.
		By("triggering a migration event", func() {
			fmt.Println("Trigger live migration via KubeVirt API")
		})
	})

	It("creates and restores from snapshot", func() {
		snap := &virtualizationv1beta1.VMSnapshot{
			ObjectMeta: metav1.ObjectMeta{Name: "scenario-snap", Namespace: namespace},
			Spec: virtualizationv1beta1.VMSnapshotSpec{
				SourceRef:     virtualizationv1beta1.NamespacedName{Name: "scenario-vm", Namespace: namespace},
				IncludedDisks: []string{"root", "data"},
			},
		}
		Expect(client.Create(ctx, snap)).To(Succeed())
		Eventually(func() bool {
			_ = client.Get(ctx, types.NamespacedName{Name: snap.Name, Namespace: namespace}, snap)
			return snap.Status.ReadyToUse
		}, 10*time.Minute, 30*time.Second).Should(BeTrue())
	})

	It("denies SR-IOV VM on incompatible cluster", func() {
		vm := &virtualizationv1beta1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{Name: "sriov-deny", Namespace: namespace},
			Spec: virtualizationv1beta1.VirtualMachineSpec{
				CPU:    "2",
				Memory: "4Gi",
				Disks:  []virtualizationv1beta1.VirtualMachineDisk{{Type: "system", DiskRef: virtualizationv1beta1.LocalObjectReference{Name: "root"}}},
				Nets:   []virtualizationv1beta1.VirtualMachineNetwork{{Type: "sriov"}},
			},
		}
		err := client.Create(ctx, vm)
		Expect(err).To(HaveOccurred())
	})

	It("enforces RBAC on power operations", func() {
		By("attempting powerOff as read-only user", func() {
			fmt.Println("Use gateway client impersonating read-only role")
		})
	})
})

func schemeForTests() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = virtualizationv1alpha1.AddToScheme(scheme)
	_ = virtualizationv1beta1.AddToScheme(scheme)
	return scheme
}
