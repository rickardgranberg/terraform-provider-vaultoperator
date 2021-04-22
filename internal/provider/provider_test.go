package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"vaultoperator": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ schema.Provider = *New("dev")()
}

func TestProvider_configure_url(t *testing.T) {
	ctx := context.TODO()

	rc := terraform.NewResourceConfigRaw(map[string]interface{}{argVaultUrl: "http://localhost:8200"})
	p := New("dev")()
	diags := p.Configure(ctx, rc)
	if diags.HasError() {
		t.Fatal(diags)
	}
}
func TestProvider_configure_url_env(t *testing.T) {
	ctx := context.TODO()
	resetEnv := func() {
		os.Unsetenv("VAULT_ADDR")
	}
	defer resetEnv()

	os.Setenv("VAULT_ADDR", "http://localhost:8200")

	rc := terraform.NewResourceConfigRaw(map[string]interface{}{})
	p := New("dev")()
	diags := p.Configure(ctx, rc)
	if diags.HasError() {
		t.Fatal(diags)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
