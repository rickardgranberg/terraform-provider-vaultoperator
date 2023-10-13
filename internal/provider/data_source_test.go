package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccDataSourceInitVar = fmt.Sprintf("data.%[1]s.test", resInit)
var testAccDataSourceInit = fmt.Sprintf(`
provider "%[1]s" {
}

data "%[2]s" "test" {
}
`, provider, resInit)

func TestAccDataSourceInit(t *testing.T) {
	startVault(t, false)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccDataSourceInitVar, argInitialized, "false"),
				),
			},
		},
	})
}

func TestAccDataSourceInitComplete(t *testing.T) {
	startVault(t, false)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInit,
			},
			{
				Config: testAccDataSourceInit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccDataSourceInitVar, argInitialized, "true"),
				),
			},
		},
	})
}
