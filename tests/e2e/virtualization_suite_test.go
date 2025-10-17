package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVirtualization(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtualization Suite")
}
