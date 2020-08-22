package shield

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shield", func() {
	Context("SeverityMap.color", func() {
		var (
			sm SeverityMap
		)

		Context("default", func() {
			BeforeEach(func() {
				sm = defaultSeverity
			})
			It("should return red", func() {
				c := sm.color(0)
				Ω(c).Should(Equal("red"))
			})
			It("should return orange", func() {
				c := sm.color(20)
				Ω(c).Should(Equal("orange"))
			})
			It("should return yellow", func() {
				c := sm.color(30)
				Ω(c).Should(Equal("yellow"))
			})
			It("should return yellowgreen", func() {
				c := sm.color(40)
				Ω(c).Should(Equal("yellowgreen"))
			})
			It("should return green", func() {
				c := sm.color(55)
				Ω(c).Should(Equal("green"))
			})
			It("should return brighgreen", func() {
				c := sm.color(70)
				Ω(c).Should(Equal("brighgreen"))
			})
		})
		It("should return an empty string", func() {
			sm = SeverityMap{}
			c := sm.color(1)
			Ω(c).Should(BeEmpty())
		})
		It("should use color from default severity if from is the same for two colors", func() {
			sm = SeverityMap{
				"red":   0,
				"green": 0,
			}
			c := sm.color(0)
			Ω(c).Should(Equal("red"))
		})
		It("should use color from default severity if from is the same for two colors and one color is not known by the defaults", func() {
			sm = SeverityMap{
				"red":  0,
				"blue": 0,
			}
			c := sm.color(0)
			Ω(c).Should(Equal("red"))
		})
	})
})
