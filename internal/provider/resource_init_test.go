package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccResourceInitVar = fmt.Sprintf("%[1]s.test", resInit)
var testAccResourceInit = fmt.Sprintf(`
provider "%[1]s" {
}

resource "%[2]s" "test" {
	secret_shares    = 5
	secret_threshold = 3
}
`, provider, resInit)

func TestAccResourceInit(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(testAccResourceInitVar, argSecretShares, regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argSecretThreshold, regexp.MustCompile("3")),
					resource.TestCheckResourceAttrSet(testAccResourceInitVar, argRootToken),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argRootToken, regexp.MustCompile(`s\.[A-Za-z0-9]+`)),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeys+".#", regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeys+".1", regexp.MustCompile("[a-z0-9]+")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeysBase64+".#", regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeysBase64+".1", regexp.MustCompile("[A-Za-z0-9]+")),
				),
			},
		},
	})
}
