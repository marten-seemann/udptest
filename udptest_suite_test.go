package udptest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUdptest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Udptest Suite")
}
