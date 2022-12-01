package provider

import (
	"fmt"
	"regexp"
	"testing"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
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
	startVault(t)

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

func testAccCheckDecryptable(t *testing.T, pgpKey *crypto.Key) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		bytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}

		pgpMessage := crypto.NewPGPMessage(bytes)

		keyRing, err := crypto.NewKeyRing(pgpKey)
		if err != nil {
			return err
		}

		_, err = keyRing.Decrypt(pgpMessage, nil, 0)
		if err != nil {
			return err
		}

		return nil
	}
}

func TestAccResourceInitPgp(t *testing.T) {
	pgpKeys := make([]*crypto.Key, 6)
	publicKeys := make([]string, 6)

	for i := range pgpKeys {
		pgpKey, err := crypto.GenerateKey(
			"Rickard Granberg",
			"rickardg@outlook.com",
			"x25519",
			0,
		)

		if err != nil {
			t.Fatal(err)
		}

		publicKeyBytes, err := pgpKey.GetPublicKey()
		if err != nil {
			t.Fatal(err)
		}

		publicKey := base64.StdEncoding.EncodeToString(publicKeyBytes)

		pgpKeys[i] = pgpKey
		publicKeys[i] = publicKey
	}

	testAccResourceInitPgp := fmt.Sprintf(`
provider "%[1]s" {
}

resource "%[2]s" "test" {
	secret_shares      = 5
	secret_threshold   = 3
	root_token_pgp_key = "%[3]s"
	pgp_keys           = ["%[4]s", "%[5]s", "%[6]s", "%[7]s", "%[8]s"]
}`,
		provider,
		resInit,
		publicKeys[0],
		publicKeys[1],
		publicKeys[2],
		publicKeys[3],
		publicKeys[4],
		publicKeys[5],
	)

	startVault(t)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInitPgp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(testAccResourceInitVar, argSecretShares, regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argSecretThreshold, regexp.MustCompile("3")),
					resource.TestCheckResourceAttrSet(testAccResourceInitVar, argRootToken),
					resource.TestCheckResourceAttrWith(testAccResourceInitVar, argRootToken, testAccCheckDecryptable(t, pgpKeys[0])),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeys+".#", regexp.MustCompile("5")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeys+".1", regexp.MustCompile("[a-z0-9]+")),
					resource.TestMatchResourceAttr(testAccResourceInitVar, argKeysBase64+".#", regexp.MustCompile("5")),
					resource.TestCheckResourceAttrWith(testAccResourceInitVar, argKeysBase64+".0", testAccCheckDecryptable(t, pgpKeys[1])),
					resource.TestCheckResourceAttrWith(testAccResourceInitVar, argKeysBase64+".1", testAccCheckDecryptable(t, pgpKeys[2])),
					resource.TestCheckResourceAttrWith(testAccResourceInitVar, argKeysBase64+".2", testAccCheckDecryptable(t, pgpKeys[3])),
					resource.TestCheckResourceAttrWith(testAccResourceInitVar, argKeysBase64+".3", testAccCheckDecryptable(t, pgpKeys[4])),
					resource.TestCheckResourceAttrWith(testAccResourceInitVar, argKeysBase64+".4", testAccCheckDecryptable(t, pgpKeys[5])),
				),
			},
		},
	})
}
