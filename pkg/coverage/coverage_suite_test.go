package coverage_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMod(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mod Suite")
}
