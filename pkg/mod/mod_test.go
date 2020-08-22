package mod

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mod", func() {
	Context("Repo", func() {
		It("should be able to read the repo of this module", func() {
			Ω(repoOf("../../go.mod")).Should(Equal("bakito/goverage-badge"))
		})
	})
})
