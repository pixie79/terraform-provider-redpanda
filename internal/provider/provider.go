// Package provider provider.go
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure RedPandaProvider satisfies various provider interfaces.
var _ provider.Provider = &RedPandaProvider{}

// RedPandaProvider describes the provider data model.
type RedPandaProvider struct {
	version              string
	SchemaRegistryApiUrl types.String `tfsdk:"schema_registry_api_url"`
	BootstrapServers     types.String `tfsdk:"bootstrap_servers"`
}

func (p *RedPandaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "redpanda"
	resp.Version = p.version
}

func (p *RedPandaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schema_registry_api_url": schema.StringAttribute{
				MarkdownDescription: "Schema Registry URL",
				Optional:            true,
			},
			"bootstrap_servers": schema.StringAttribute{
				MarkdownDescription: "List of Bootstrap servers [\"localhost:9092\"]",
				Optional:            true,
			},
		},
	}
}

func (p *RedPandaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RedPandaProvider

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = p
	resp.ResourceData = p
}

func (p *RedPandaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaRegistryResource,
		//kafka.NewTopicResource,
	}
}

func (p *RedPandaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSchemaRegistryDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedPandaProvider{
			version: version,
		}
	}
}
