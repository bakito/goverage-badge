package coverage_test

import (
	"github.com/bakito/goverage-badge/pkg/coverage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coverage", func() {
	Context("Calculate", func() {
		It("should be 57.9%", func() {
			Ω(coverage.Calculate("../../testdata/coverage.out")).
				Should(BeNumerically("~", 57.9, 0.01))
		})
	})
})
