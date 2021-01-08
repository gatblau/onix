package build

import (
	"github.com/gatblau/onix/artisan/core"
	"github.com/gatblau/onix/artisan/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestSample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Package Builder Suite")
}

var _ = Describe("Signing and Verifying a Package", func() {
	Context("when creating and signing a package with a private PGP key", func() {
		builder := NewBuilder()
		builder.Build(".", "", "", core.ParseName("artisan"), "linux", false, false, "root_rsa_key.pgp")
		// time.Sleep(time.Second * 2)
		It("should open successfully with the corresponding public PGP key", func() {
			local := registry.NewLocalRegistry()
			local.Open(core.ParseName("artisan"), "", false, "test", "root_rsa_pub.pgp", false)
		})
		os.RemoveAll("images")
		os.RemoveAll("version")
		os.RemoveAll("test")
	})
})
