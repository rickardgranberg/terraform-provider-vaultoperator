package provider

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

const (
	provider    = "vaultoperator"
	resInit     = provider + "_init"
	argVaultUrl = "vault_url"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: providerSchema(),
			ResourcesMap: map[string]*schema.Resource{
				resInit: resourceInit(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
	client *api.Client
	url    string
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		argVaultUrl: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Vault instance URL",
		},
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		a := &apiClient{}
		a.url = os.Getenv("VAULT_ADDR")

		if u, ok := d.GetOk(argVaultUrl); ok {
			a.url = u.(string)
		}

		if a.url == "" {
			return nil, diag.Errorf("%s is required", argVaultUrl)
		}

		c, err := api.NewClient(&api.Config{
			Address: a.url,
		})

		if err != nil {
			logError("failed to create Vault API client: %v", err)
			return nil, diag.FromErr(err)
		}

		a.client = c

		return a, nil
	}
}

func logError(fmt string, v ...interface{}) {
	log.Printf("[ERROR] "+fmt, v)
}
