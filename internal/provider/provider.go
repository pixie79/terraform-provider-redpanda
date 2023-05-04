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

type RedPandaProvider struct {
	version string
}

// RedPandaProviderModel describes the provider data model.
type RedPandaProviderModel struct {
	SchemaApiUrl     types.String `tfsdk:"schema_api_url"`
	BootstrapServers types.List   `tfsdk:"bootstrap_servers"`
}

func (p *RedPandaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "redpanda"
	resp.Version = p.version
}

func (p *RedPandaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schema_api_url": schema.StringAttribute{
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
	var data RedPandaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := data.SchemaApiUrl.ValueString()
	bootstrapServers := data.BootstrapServers.String()
	schemaClient := NewClientSchema(apiURL)
	topicClient := NewClientTopic(bootstrapServers)
	resp.DataSourceData = schemaClient
	resp.ResourceData = schemaClient
}

func (p *RedPandaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaResource,
		NewTopicResource,
	}
}

func (p *RedPandaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSchemaDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedPandaProvider{
			version: version,
		}
	}
}
