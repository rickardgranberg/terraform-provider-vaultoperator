package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	argInitialized = "initialized"
)

func providerDatasource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Resource for vault operator init",

		ReadContext: providerDatasourceRead,

		Schema: map[string]*schema.Schema{
			argInitialized: {
				Description: "The current initialization state of Vault.",
				Type: schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func providerDatasourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	d.SetId(client.url)

	res, err := client.client.Sys().InitStatus()
	if err != nil {
		logError("failed to read init status from Vault: %v", err)
		return diag.FromErr(err)
	}

	logDebug("response: %v", res)

	if err := d.Set(argInitialized, res); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
