package provider

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

const (
	envVaultAddr      = "VAULT_ADDR"
	provider          = "vaultoperator"
	resInit           = provider + "_init"
	argVaultUrl       = "vault_url"
	argVaultAddr      = "vault_addr"
	argRequestHeaders = "request_headers"
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
			DataSourcesMap: map[string]*schema.Resource{
				resInit: providerDatasource(),
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
			Optional:    true,
			Description: "Vault instance URL",
			Deprecated:  fmt.Sprintf("%q is deprecated, please use %q instead", argVaultUrl, argVaultAddr),
		},
		argVaultAddr: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Vault instance URL",
		},
		argRequestHeaders: {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		a := &apiClient{}

		if u := d.Get(argVaultAddr).(string); u != "" {
			a.url = u
		} else if u := d.Get(argVaultUrl).(string); u != "" {
			a.url = u
		} else {
			a.url = os.Getenv(envVaultAddr)
		}

		if a.url == "" {
			return nil, diag.Errorf("argument '%s' is required, or set VAULT_ADDR environment variable", argVaultUrl)
		}

		if c, err := api.NewClient(&api.Config{Address: a.url}); err != nil {
			logError("failed to create Vault API client: %v", err)
			return nil, diag.FromErr(err)
		} else {
			a.client = c
		}

		return a, nil
	}
}

func logError(fmt string, v ...interface{}) {
	log.Printf("[ERROR] "+fmt, v)
}

func logInfo(fmt string, v ...interface{}) {
	log.Printf("[INFO] "+fmt, v)
}

func logDebug(fmt string, v ...interface{}) {
	log.Printf("[DEBUG] "+fmt, v)
}
