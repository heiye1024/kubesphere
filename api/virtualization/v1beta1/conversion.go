package v1beta1

import "sigs.k8s.io/controller-runtime/pkg/conversion"

func (*VirtualMachine) Hub()     {}
func (*VirtualMachineList) Hub() {}
func (*VirtualDisk) Hub()        {}
func (*VirtualDiskList) Hub()    {}
func (*VirtualNet) Hub()         {}
func (*VirtualNetList) Hub()     {}
func (*VMSnapshot) Hub()         {}
func (*VMSnapshotList) Hub()     {}
func (*VMTemplate) Hub()         {}
func (*VMTemplateList) Hub()     {}

var _ conversion.Hub = &VirtualMachine{}
var _ conversion.Hub = &VirtualDisk{}
var _ conversion.Hub = &VirtualNet{}
var _ conversion.Hub = &VMSnapshot{}
var _ conversion.Hub = &VMTemplate{}
