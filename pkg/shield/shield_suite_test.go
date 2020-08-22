package shield

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShield(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shield Suite")
}
