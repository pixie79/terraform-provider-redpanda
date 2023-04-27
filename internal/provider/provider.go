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
var _ provider.Provider = &SchemaProvider{}

type SchemaProvider struct {
	version string
}

// SchemaProviderModel describes the provider data model.
type SchemaProviderModel struct {
	SchemaApiUrl types.String `tfsdk:"schema_api_url"`
}

func (p *SchemaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "redpanda"
	resp.Version = p.version
}

func (p *SchemaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schema_api_url": schema.StringAttribute{
				MarkdownDescription: "Schema Registry URL",
				Required:            true,
			},
		},
	}
}

func (p *SchemaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SchemaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := data.SchemaApiUrl.ValueString()
	client := NewClientSchema(apiURL)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SchemaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaResource,
	}
}

func (p *SchemaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSchemaDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SchemaProvider{
			version: version,
		}
	}
}
