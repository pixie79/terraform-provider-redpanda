package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSourceRedPandaSchema{}

func NewDataSourceRedPandaSchema() datasource.DataSource {
	return &DataSourceRedPandaSchema{}
}

// DataSourceRedPandaSchema defines the data source implementation.
type DataSourceRedPandaSchema struct {
	client *ClientSchema
}

// DataSourceRedPandaSchemaModel describes the data source data model.
type DataSourceRedPandaSchemaModel struct {
	Subject    types.String `tfsdk:"subject"`
	Schema     types.String `tfsdk:"schema"`
	SchemaType types.String `tfsdk:"schemaType,omitempty"`
	Version    types.Int64  `tfsdk:"version,omitempty"`
	Id         types.Int64  `tfsdk:"id,omitempty"`
}

func (d *DataSourceRedPandaSchema) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redpanda"
}

func (d *DataSourceRedPandaSchema) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema resource",

		Attributes: map[string]schema.Attribute{
			"subject": schema.StringAttribute{
				MarkdownDescription: "Schema Subject",
				Optional:            false,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version of Schema",
				Optional:            false,
			},
		},
	}
}

func (d *DataSourceRedPandaSchema) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ClientSchema)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DataSourceRedPandaSchema) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceRedPandaSchemaModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	//data.Version = types.Int64Value(data.version)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
