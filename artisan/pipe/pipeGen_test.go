/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-2021 by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/
package pipe

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestSample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pipeline Generator Suite")
}

var _ = Describe("Loading flow from file", func() {
	var (
		err       error
		generator *Generator
	)
	generator, err = NewGeneratorFromPath("test-flow.yaml")

	Context("when creating a new generator from a file", func() {
		It("should not return an error", func() {
			Expect(err).Should(BeNil())
		})
		It("should have a valid flow", func() {
			if generator != nil {
				Expect(generator.flow).ShouldNot(Equal(nil))
			}
		})
	})
})
