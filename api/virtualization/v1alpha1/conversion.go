package v1alpha1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	v1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

// ConvertTo converts this VirtualMachine to the Hub version.
func (src *VirtualMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.VirtualMachine)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = v1beta1.VirtualMachineSpec{
		CPU:                   src.Spec.CPU,
		Memory:                src.Spec.Memory,
		CPUModel:              src.Spec.CPUModel,
		DedicatedCPUPlacement: src.Spec.DedicatedCPUPlacement,
		NUMAPolicy:            src.Spec.NUMAPolicy,
		NUMA:                  convertNUMAToBeta(src.Spec.NUMA),
		Hugepages:             src.Spec.Hugepages,
		GPUs:                  convertGPUsToBeta(src.Spec.GPUs),
		Disks:                 convertDisksToBeta(src.Spec.Disks),
		Nets:                  convertNetsToBeta(src.Spec.Nets),
		CloudInit:             convertCloudInitToBeta(src.Spec.CloudInit),
		Console:               convertConsoleToBeta(src.Spec.Console),
		LiveMigration:         convertLiveMigrationToBeta(src.Spec.LiveMigration),
		LivenessProbe:         convertProbeToBeta(src.Spec.LivenessProbe),
		ReadinessProbe:        convertProbeToBeta(src.Spec.ReadinessProbe),
		PowerState:            v1beta1.PowerState(src.Spec.PowerState),
	}
	if src.Spec.KubeVirt != nil {
		dst.Spec.KubeVirt = src.Spec.KubeVirt.DeepCopy()
	}
	dst.Status = v1beta1.VirtualMachineStatus{
		Conditions:     src.Status.Conditions,
		PowerState:     v1beta1.PowerState(src.Status.PowerState),
		Phase:          src.Status.Phase,
		MigrationState: src.Status.MigrationState,
	}
	return nil
}

// ConvertFrom converts from the Hub version to this version.
func (dst *VirtualMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.VirtualMachine)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = VirtualMachineSpec{
		CPU:                   src.Spec.CPU,
		Memory:                src.Spec.Memory,
		CPUModel:              src.Spec.CPUModel,
		DedicatedCPUPlacement: src.Spec.DedicatedCPUPlacement,
		NUMAPolicy:            src.Spec.NUMAPolicy,
		NUMA:                  convertNUMAToAlpha(src.Spec.NUMA),
		Hugepages:             src.Spec.Hugepages,
		GPUs:                  convertGPUsToAlpha(src.Spec.GPUs),
		Disks:                 convertDisksToAlpha(src.Spec.Disks),
		Nets:                  convertNetsToAlpha(src.Spec.Nets),
		CloudInit:             convertCloudInitToAlpha(src.Spec.CloudInit),
		Console:               convertConsoleToAlpha(src.Spec.Console),
		LiveMigration:         convertLiveMigrationToAlpha(src.Spec.LiveMigration),
		PowerState:            PowerState(src.Spec.PowerState),
	}
	dst.Spec.LivenessProbe = convertProbeToAlpha(src.Spec.LivenessProbe)
	dst.Spec.ReadinessProbe = convertProbeToAlpha(src.Spec.ReadinessProbe)
	if src.Spec.KubeVirt != nil {
		dst.Spec.KubeVirt = src.Spec.KubeVirt.DeepCopy()
	}
	dst.Status = VirtualMachineStatus{
		Conditions:     src.Status.Conditions,
		PowerState:     PowerState(src.Status.PowerState),
		Phase:          src.Status.Phase,
		MigrationState: src.Status.MigrationState,
	}
	return nil
}

func convertDisksToBeta(disks []VirtualMachineDisk) []v1beta1.VirtualMachineDisk {
	out := make([]v1beta1.VirtualMachineDisk, len(disks))
	for i, d := range disks {
		out[i] = v1beta1.VirtualMachineDisk{
			Type:      d.Type,
			Bus:       d.Bus,
			Cache:     d.Cache,
			IOThread:  d.IOThread,
			BootOrder: d.BootOrder,
			Hotplug:   d.Hotplug,
			DiskRef:   v1beta1.LocalObjectReference{Name: d.DiskRef.Name},
		}
	}
	return out
}

func convertDisksToAlpha(disks []v1beta1.VirtualMachineDisk) []VirtualMachineDisk {
	out := make([]VirtualMachineDisk, len(disks))
	for i, d := range disks {
		out[i] = VirtualMachineDisk{
			Type:      d.Type,
			Bus:       d.Bus,
			Cache:     d.Cache,
			IOThread:  d.IOThread,
			BootOrder: d.BootOrder,
			Hotplug:   d.Hotplug,
			DiskRef:   LocalObjectReference{Name: d.DiskRef.Name},
		}
	}
	return out
}

func convertNetsToBeta(nets []VirtualMachineNetwork) []v1beta1.VirtualMachineNetwork {
	out := make([]v1beta1.VirtualMachineNetwork, len(nets))
	for i, n := range nets {
		out[i] = v1beta1.VirtualMachineNetwork{
			Type:          n.Type,
			Model:         n.Model,
			Bandwidth:     n.Bandwidth,
			SRIOVResource: n.SRIOVResource,
			Multiqueue:    n.Multiqueue,
		}
		if n.NADRef != nil {
			out[i].NADRef = &v1beta1.NamespacedName{Namespace: n.NADRef.Namespace, Name: n.NADRef.Name}
		}
	}
	return out
}

func convertNetsToAlpha(nets []v1beta1.VirtualMachineNetwork) []VirtualMachineNetwork {
	out := make([]VirtualMachineNetwork, len(nets))
	for i, n := range nets {
		out[i] = VirtualMachineNetwork{
			Type:          n.Type,
			Model:         n.Model,
			Bandwidth:     n.Bandwidth,
			SRIOVResource: n.SRIOVResource,
			Multiqueue:    n.Multiqueue,
		}
		if n.NADRef != nil {
			out[i].NADRef = &NamespacedName{Namespace: n.NADRef.Namespace, Name: n.NADRef.Name}
		}
	}
	return out
}

func convertProbeToBeta(probe *Probe) *v1beta1.Probe {
	if probe == nil {
		return nil
	}
	return &v1beta1.Probe{
		PeriodSeconds:    probe.PeriodSeconds,
		TimeoutSeconds:   probe.TimeoutSeconds,
		FailureThreshold: probe.FailureThreshold,
	}
}

func convertProbeToAlpha(probe *v1beta1.Probe) *Probe {
	if probe == nil {
		return nil
	}
	return &Probe{
		PeriodSeconds:    probe.PeriodSeconds,
		TimeoutSeconds:   probe.TimeoutSeconds,
		FailureThreshold: probe.FailureThreshold,
	}
}

func convertNUMAToBeta(numa *NUMASpec) *v1beta1.NUMASpec {
	if numa == nil {
		return nil
	}
	cells := make([]v1beta1.NUMACell, len(numa.Cells))
	for i, cell := range numa.Cells {
		cells[i] = v1beta1.NUMACell{
			ID:             cell.ID,
			CPUs:           cell.CPUs,
			Memory:         cell.Memory,
			ThreadsPerCore: cell.ThreadsPerCore,
		}
	}
	return &v1beta1.NUMASpec{Cells: cells}
}

func convertNUMAToAlpha(numa *v1beta1.NUMASpec) *NUMASpec {
	if numa == nil {
		return nil
	}
	cells := make([]NUMACell, len(numa.Cells))
	for i, cell := range numa.Cells {
		cells[i] = NUMACell{
			ID:             cell.ID,
			CPUs:           cell.CPUs,
			Memory:         cell.Memory,
			ThreadsPerCore: cell.ThreadsPerCore,
		}
	}
	return &NUMASpec{Cells: cells}
}

func convertGPUsToBeta(gpus []GPUDevice) []v1beta1.GPUDevice {
	out := make([]v1beta1.GPUDevice, len(gpus))
	for i, g := range gpus {
		out[i] = v1beta1.GPUDevice{Name: g.Name, DeviceType: g.DeviceType, ResourceName: g.ResourceName}
	}
	return out
}

func convertGPUsToAlpha(gpus []v1beta1.GPUDevice) []GPUDevice {
	out := make([]GPUDevice, len(gpus))
	for i, g := range gpus {
		out[i] = GPUDevice{Name: g.Name, DeviceType: g.DeviceType, ResourceName: g.ResourceName}
	}
	return out
}

func convertCloudInitToBeta(ci *CloudInitSpec) *v1beta1.CloudInitSpec {
	if ci == nil {
		return nil
	}
	out := &v1beta1.CloudInitSpec{
		UserData:          ci.UserData,
		NetworkData:       ci.NetworkData,
		SSHAuthorizedKeys: append([]string(nil), ci.SSHAuthorizedKeys...),
	}
	if ci.UserDataSecretRef != nil {
		out.UserDataSecretRef = &v1beta1.NamespacedName{Namespace: ci.UserDataSecretRef.Namespace, Name: ci.UserDataSecretRef.Name}
	}
	if ci.NetworkDataSecretRef != nil {
		out.NetworkDataSecretRef = &v1beta1.NamespacedName{Namespace: ci.NetworkDataSecretRef.Namespace, Name: ci.NetworkDataSecretRef.Name}
	}
	return out
}

func convertCloudInitToAlpha(ci *v1beta1.CloudInitSpec) *CloudInitSpec {
	if ci == nil {
		return nil
	}
	out := &CloudInitSpec{
		UserData:          ci.UserData,
		NetworkData:       ci.NetworkData,
		SSHAuthorizedKeys: append([]string(nil), ci.SSHAuthorizedKeys...),
	}
	if ci.UserDataSecretRef != nil {
		out.UserDataSecretRef = &NamespacedName{Namespace: ci.UserDataSecretRef.Namespace, Name: ci.UserDataSecretRef.Name}
	}
	if ci.NetworkDataSecretRef != nil {
		out.NetworkDataSecretRef = &NamespacedName{Namespace: ci.NetworkDataSecretRef.Namespace, Name: ci.NetworkDataSecretRef.Name}
	}
	return out
}

func convertConsoleToBeta(console *ConsoleDevices) *v1beta1.ConsoleDevices {
	if console == nil {
		return nil
	}
	return &v1beta1.ConsoleDevices{VNC: console.VNC, Serial: console.Serial, Type: console.Type}
}

func convertConsoleToAlpha(console *v1beta1.ConsoleDevices) *ConsoleDevices {
	if console == nil {
		return nil
	}
	return &ConsoleDevices{VNC: console.VNC, Serial: console.Serial, Type: console.Type}
}

func convertLiveMigrationToBeta(lm *LiveMigrationSpec) *v1beta1.LiveMigrationSpec {
	if lm == nil {
		return nil
	}
	return &v1beta1.LiveMigrationSpec{
		Enabled:                  lm.Enabled,
		Bandwidth:                lm.Bandwidth,
		CompletionTimeoutSeconds: lm.CompletionTimeoutSeconds,
		AllowPostCopy:            lm.AllowPostCopy,
		AutoConverge:             lm.AutoConverge,
	}
}

func convertLiveMigrationToAlpha(lm *v1beta1.LiveMigrationSpec) *LiveMigrationSpec {
	if lm == nil {
		return nil
	}
	return &LiveMigrationSpec{
		Enabled:                  lm.Enabled,
		Bandwidth:                lm.Bandwidth,
		CompletionTimeoutSeconds: lm.CompletionTimeoutSeconds,
		AllowPostCopy:            lm.AllowPostCopy,
		AutoConverge:             lm.AutoConverge,
	}
}

var _ conversion.Convertible = &VirtualMachine{}
var _ conversion.Convertible = &VirtualDisk{}
var _ conversion.Convertible = &VirtualNet{}
var _ conversion.Convertible = &VMSnapshot{}
var _ conversion.Convertible = &VMTemplate{}

func (src *VirtualDisk) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.VirtualDisk)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = v1beta1.VirtualDiskSpec{
		Backing:      src.Spec.Backing,
		Size:         src.Spec.Size,
		AccessMode:   src.Spec.AccessMode,
		VolumeMode:   src.Spec.VolumeMode,
		StorageClass: src.Spec.StorageClass,
	}
	dst.Status = v1beta1.VirtualDiskStatus{Conditions: src.Status.Conditions}
	return nil
}

func (dst *VirtualDisk) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.VirtualDisk)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = VirtualDiskSpec{
		Backing:      src.Spec.Backing,
		Size:         src.Spec.Size,
		AccessMode:   src.Spec.AccessMode,
		VolumeMode:   src.Spec.VolumeMode,
		StorageClass: src.Spec.StorageClass,
	}
	dst.Status = VirtualDiskStatus{Conditions: src.Status.Conditions}
	return nil
}

func (src *VirtualNet) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.VirtualNet)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = v1beta1.VirtualNetSpec{
		NADTemplate:    src.Spec.NADTemplate,
		BandwidthLimit: src.Spec.BandwidthLimit,
		VLAN:           src.Spec.VLAN,
		SRIOVResource:  src.Spec.SRIOVResource,
	}
	dst.Status = v1beta1.VirtualNetStatus{Conditions: src.Status.Conditions}
	return nil
}

func (dst *VirtualNet) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.VirtualNet)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = VirtualNetSpec{
		NADTemplate:    src.Spec.NADTemplate,
		BandwidthLimit: src.Spec.BandwidthLimit,
		VLAN:           src.Spec.VLAN,
		SRIOVResource:  src.Spec.SRIOVResource,
	}
	dst.Status = VirtualNetStatus{Conditions: src.Status.Conditions}
	return nil
}

func (src *VMSnapshot) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.VMSnapshot)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = v1beta1.VMSnapshotSpec{
		SourceRef:     v1beta1.NamespacedName{Namespace: src.Spec.SourceRef.Namespace, Name: src.Spec.SourceRef.Name},
		IncludedDisks: append([]string{}, src.Spec.IncludedDisks...),
		RetainPolicy:  src.Spec.RetainPolicy,
	}
	dst.Status = v1beta1.VMSnapshotStatus{Conditions: src.Status.Conditions, ReadyToUse: src.Status.ReadyToUse}
	return nil
}

func (dst *VMSnapshot) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.VMSnapshot)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = VMSnapshotSpec{
		SourceRef:     NamespacedName{Namespace: src.Spec.SourceRef.Namespace, Name: src.Spec.SourceRef.Name},
		IncludedDisks: append([]string{}, src.Spec.IncludedDisks...),
		RetainPolicy:  src.Spec.RetainPolicy,
	}
	dst.Status = VMSnapshotStatus{Conditions: src.Status.Conditions, ReadyToUse: src.Status.ReadyToUse}
	return nil
}

func (src *VMTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.VMTemplate)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = v1beta1.VMTemplateSpec{
		Parameters: v1beta1.TemplateParameters{
			CPU:      src.Spec.Parameters.CPU,
			Memory:   src.Spec.Parameters.Memory,
			OS:       src.Spec.Parameters.OS,
			Image:    src.Spec.Parameters.Image,
			Networks: append([]string{}, src.Spec.Parameters.Networks...),
			Disks:    convertTemplateDisksToBeta(src.Spec.Parameters.Disks),
		},
		Constraints: v1beta1.TemplateConstraints{
			MinCPU:    src.Spec.Constraints.MinCPU,
			MaxCPU:    src.Spec.Constraints.MaxCPU,
			MinMemory: src.Spec.Constraints.MinMemory,
			MaxMemory: src.Spec.Constraints.MaxMemory,
		},
		UIHints: src.Spec.UIHints,
	}
	dst.Status = v1beta1.VMTemplateStatus{Conditions: src.Status.Conditions}
	return nil
}

func (dst *VMTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.VMTemplate)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = VMTemplateSpec{
		Parameters: TemplateParameters{
			CPU:      src.Spec.Parameters.CPU,
			Memory:   src.Spec.Parameters.Memory,
			OS:       src.Spec.Parameters.OS,
			Image:    src.Spec.Parameters.Image,
			Networks: append([]string{}, src.Spec.Parameters.Networks...),
			Disks:    convertTemplateDisksToAlpha(src.Spec.Parameters.Disks),
		},
		Constraints: TemplateConstraints{
			MinCPU:    src.Spec.Constraints.MinCPU,
			MaxCPU:    src.Spec.Constraints.MaxCPU,
			MinMemory: src.Spec.Constraints.MinMemory,
			MaxMemory: src.Spec.Constraints.MaxMemory,
		},
		UIHints: src.Spec.UIHints,
	}
	dst.Status = VMTemplateStatus{Conditions: src.Status.Conditions}
	return nil
}

func convertTemplateDisksToBeta(disks []TemplateDisk) []v1beta1.TemplateDisk {
	out := make([]v1beta1.TemplateDisk, len(disks))
	for i, d := range disks {
		out[i] = v1beta1.TemplateDisk{Name: d.Name, Size: d.Size, Type: d.Type}
	}
	return out
}

func convertTemplateDisksToAlpha(disks []v1beta1.TemplateDisk) []TemplateDisk {
	out := make([]TemplateDisk, len(disks))
	for i, d := range disks {
		out[i] = TemplateDisk{Name: d.Name, Size: d.Size, Type: d.Type}
	}
	return out
}
