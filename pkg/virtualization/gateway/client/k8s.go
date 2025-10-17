package client

import (
	"context"
	"errors"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	virtualizationv1alpha1 "kubesphere.io/kubesphere/api/virtualization/v1alpha1"
	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
	"kubesphere.io/kubesphere/pkg/virtualization/gateway/dto"
)

// QueryOptions holds filters for list requests.
type QueryOptions struct {
	Namespace string
	Clusters  []string
}

// VirtualizationStore abstracts Kubernetes CRUD for virtualization resources.
type VirtualizationStore interface {
	ListVMs(context.Context, QueryOptions) (*virtualizationv1beta1.VirtualMachineList, error)
	CreateVM(context.Context, *virtualizationv1beta1.VirtualMachine) error
	UpdatePowerState(context.Context, string, string, virtualizationv1beta1.PowerState) error
	MigrateVM(context.Context, string, string) error
	Console(context.Context, string, string) (string, error)
	ListDisks(context.Context, QueryOptions) (*virtualizationv1beta1.VirtualDiskList, error)
	CreateDisk(context.Context, *virtualizationv1beta1.VirtualDisk) error
	ListNetworks(context.Context, QueryOptions) (*virtualizationv1beta1.VirtualNetList, error)
	CreateNetwork(context.Context, *virtualizationv1beta1.VirtualNet) error
	ListSnapshots(context.Context, QueryOptions) (*virtualizationv1beta1.VMSnapshotList, error)
	CreateSnapshot(context.Context, *virtualizationv1beta1.VMSnapshot) error
	ListTemplates(context.Context, QueryOptions) (*virtualizationv1beta1.VMTemplateList, error)
	CreateTemplate(context.Context, *virtualizationv1beta1.VMTemplate) error
}

// Store implements VirtualizationStore.
type Store struct {
	client ctrlclient.Client
}

var errForbidden = errors.New("forbidden")

// IsForbidden checks whether the error originated from RBAC denial.
func IsForbidden(err error) bool {
	return errors.Is(err, errForbidden)
}

func NewDynamicClient(cfg *rest.Config) (ctrlclient.Client, error) {
	scheme := runtime.NewScheme()
	_ = virtualizationv1alpha1.AddToScheme(scheme)
	_ = virtualizationv1beta1.AddToScheme(scheme)
	return ctrlclient.New(cfg, ctrlclient.Options{Scheme: scheme})
}

func NewVirtualizationStore(c ctrlclient.Client) VirtualizationStore {
	return &Store{client: c}
}

func (s *Store) ListVMs(ctx context.Context, opts QueryOptions) (*virtualizationv1beta1.VirtualMachineList, error) {
	return s.listVirtualMachines(ctx, opts)
}

func (s *Store) CreateVM(ctx context.Context, vm *virtualizationv1beta1.VirtualMachine) error {
	return s.client.Create(ctx, vm)
}

func (s *Store) UpdatePowerState(ctx context.Context, namespace, name string, state virtualizationv1beta1.PowerState) error {
	vm := &virtualizationv1beta1.VirtualMachine{}
	if err := s.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, vm); err != nil {
		if apierrors.IsForbidden(err) {
			return errForbidden
		}
		return err
	}
	vm.Spec.PowerState = state
	if err := s.client.Update(ctx, vm); err != nil {
		if apierrors.IsForbidden(err) {
			return errForbidden
		}
		return err
	}
	return nil
}

func (s *Store) MigrateVM(ctx context.Context, namespace, name string) error {
	// Placeholder for invoking KubeVirt migration
	return nil
}

func (s *Store) Console(ctx context.Context, namespace, name string) (string, error) {
	return fmt.Sprintf("/console/%s/%s", namespace, name), nil
}

func (s *Store) ListDisks(ctx context.Context, opts QueryOptions) (*virtualizationv1beta1.VirtualDiskList, error) {
	list := &virtualizationv1beta1.VirtualDiskList{}
	if err := s.client.List(ctx, list, buildListOptions(opts)...); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) CreateDisk(ctx context.Context, disk *virtualizationv1beta1.VirtualDisk) error {
	return s.client.Create(ctx, disk)
}

func (s *Store) ListNetworks(ctx context.Context, opts QueryOptions) (*virtualizationv1beta1.VirtualNetList, error) {
	list := &virtualizationv1beta1.VirtualNetList{}
	if err := s.client.List(ctx, list, buildListOptions(opts)...); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) CreateNetwork(ctx context.Context, netObj *virtualizationv1beta1.VirtualNet) error {
	return s.client.Create(ctx, netObj)
}

func (s *Store) ListSnapshots(ctx context.Context, opts QueryOptions) (*virtualizationv1beta1.VMSnapshotList, error) {
	list := &virtualizationv1beta1.VMSnapshotList{}
	if err := s.client.List(ctx, list, buildListOptions(opts)...); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) CreateSnapshot(ctx context.Context, snapshot *virtualizationv1beta1.VMSnapshot) error {
	return s.client.Create(ctx, snapshot)
}

func (s *Store) ListTemplates(ctx context.Context, opts QueryOptions) (*virtualizationv1beta1.VMTemplateList, error) {
	list := &virtualizationv1beta1.VMTemplateList{}
	if err := s.client.List(ctx, list, buildListOptions(opts)...); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) listVirtualMachines(ctx context.Context, opts QueryOptions) (*virtualizationv1beta1.VirtualMachineList, error) {
	clusters := normalizeClusters(opts.Clusters)
	if len(clusters) <= 1 {
		list := &virtualizationv1beta1.VirtualMachineList{}
		if err := s.client.List(ctx, list, buildListOptions(opts)...); err != nil {
			return nil, err
		}
		return list, nil
	}

	aggregated := &virtualizationv1beta1.VirtualMachineList{}
	for _, cluster := range clusters {
		list := &virtualizationv1beta1.VirtualMachineList{}
		localOpts := QueryOptions{Namespace: opts.Namespace, Clusters: []string{cluster}}
		if err := s.client.List(ctx, list, buildListOptions(localOpts)...); err != nil {
			return nil, err
		}
		aggregated.Items = append(aggregated.Items, list.Items...)
	}
	return aggregated, nil
}

func buildListOptions(opts QueryOptions) []ctrlclient.ListOption {
	options := []ctrlclient.ListOption{}
	if opts.Namespace != "" {
		options = append(options, ctrlclient.InNamespace(opts.Namespace))
	}
	clusters := normalizeClusters(opts.Clusters)
	if len(clusters) == 1 && clusters[0] != "" && clusters[0] != "all" {
		options = append(options, ctrlclient.MatchingLabels{dto.LabelCluster: clusters[0]})
	}
	return options
}

func normalizeClusters(clusters []string) []string {
	if len(clusters) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	normalized := make([]string, 0, len(clusters))
	for _, cluster := range clusters {
		if cluster == "" {
			continue
		}
		if cluster == "all" {
			// "all" means no filter.
			return []string{"all"}
		}
		if _, ok := seen[cluster]; ok {
			continue
		}
		seen[cluster] = struct{}{}
		normalized = append(normalized, cluster)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func (s *Store) CreateTemplate(ctx context.Context, template *virtualizationv1beta1.VMTemplate) error {
	return s.client.Create(ctx, template)
}
