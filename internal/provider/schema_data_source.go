// Package provider.go
package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SchemaDataSource{}

func NewSchemaDataSource() datasource.DataSource {
	return &SchemaDataSource{}
}

// SchemaDataSource defines the data source implementation.
type SchemaDataSource struct {
	client *ClientSchema
}

func (d *SchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (d *SchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema resource",

		Attributes: map[string]schema.Attribute{
			"subject": schema.StringAttribute{
				MarkdownDescription: "Schema Subject",
				Required:            true,
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "Schema - string encoded",
				Computed:            true,
			},
			"schema_type": schema.StringAttribute{
				MarkdownDescription: "Schema Type, defaults to AVRO",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Version of Schema",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Id of Schema",
			},
		},
	}
}

func (d *SchemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SchemaModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	version, err := d.client.GetLatestVersion(data.Subject.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	schemaModel, err := d.client.GetSchema(data.Subject.ValueString(), version)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}
	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	data.Version = types.Int64Value(version)
	data.Id = schemaModel.Id
	data.Subject = schemaModel.Subject
	data.Schema = schemaModel.Schema
	data.SchemaType = schemaModel.SchemaType

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
