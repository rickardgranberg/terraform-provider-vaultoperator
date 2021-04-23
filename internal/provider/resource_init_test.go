package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceInit(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"vaultoperator_init.foo", argSecretShares, regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(
						"vaultoperator_init.foo", argSecretThreshold, regexp.MustCompile("3")),
					resource.TestMatchResourceAttr(
						"vaultoperator_init.foo", argRootToken, regexp.MustCompile(`s\.[A-Za-z0-9]+`)),
					resource.TestCheckResourceAttrSet("vaultoperator_init.foo", argRootToken),
					resource.TestMatchResourceAttr("vaultoperator_init.foo", argKeys+".#", regexp.MustCompile("5")),
					resource.TestMatchResourceAttr("vaultoperator_init.foo", argKeys+".1", regexp.MustCompile("[a-z0-9]+")),
					resource.TestMatchResourceAttr("vaultoperator_init.foo", argKeysBase64+".#", regexp.MustCompile("5")),
					resource.TestMatchResourceAttr("vaultoperator_init.foo", argKeysBase64+".1", regexp.MustCompile("[A-Za-z0-9]+")),
				),
			},
		},
	})
}

const testAccResourceInit = `
provider "vaultoperator" {
}

resource "vaultoperator_init" "foo" {
	secret_shares    = 5
	secret_threshold = 3
}
`
