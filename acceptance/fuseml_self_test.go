package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fuseml maintenance operations", func() {
	Describe("info", func() {
		It("prints information about the FuseML components and platform", func() {
			info, err := Fuseml("info", "")
			Expect(err).ToNot(HaveOccurred())
			Expect(info).To(MatchRegexp("Platform: k3s"))
			Expect(info).To(MatchRegexp("Kubernetes Version: v1.20"))
		})
	})
})
