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
						"vaultoperator_init.foo", "secret_shares", regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(
						"vaultoperator_init.foo", "secret_threshold", regexp.MustCompile("3")),
				),
			},
		},
	})
}

const testAccResourceInit = `
provider "vaultoperator" {
	vault_url = "http://localhost:8200"
}
resource "vaultoperator_init" "foo" {
	secret_shares    = 5
	secret_threshold = 3
}
`
