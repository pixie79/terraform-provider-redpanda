// provider.go
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REDPANDA_API_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"redpanda_topic":  resourceRedPandaTopic(),
			"redpanda_schema": resourceRedPandaSchema(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiURL := d.Get("api_url").(string)
	return NewClient(apiURL), nil
}
