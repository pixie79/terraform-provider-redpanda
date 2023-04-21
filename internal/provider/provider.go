package main

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
	Endpoint     types.String `tfsdk:"endpoint"`
	SchemaApiUrl types.String `tfsdk:"schema_api_url"`
}

func (p *RedPandaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "RedPanda"
	resp.Version = p.version
}

func (p *RedPandaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"schema_api_url": schema.StringAttribute{
				MarkdownDescription: "Schema Registry URL",
				Optional:            true,
			},
		},
	}
}

// func (p *RedPandaProvider) Schema(_ context.Context, _ *provider.SchemaRequest) (*provider.SchemaResponse, error) {
// 	return &provider.SchemaResponse{
// 		Schema: &schema.Schema{
//             ResourcesMap: map[string]*schema.Resource{
//                 "RedPanda_schema": RedPandaSchemaResourceType(),
//             },
// 		},
// 	}, nil
// }

func (p *RedPandaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RedPandaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	apiURL := data.SchemaApiUrl.ValueString()
	client := NewClientSchema(apiURL)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RedPandaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceRedPandaSchema,
	}
}

func (p *RedPandaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceRedPandaSchema,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedPandaProvider{
			version: version,
		}
	}
}
