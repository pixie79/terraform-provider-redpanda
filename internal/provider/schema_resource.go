// Package provider.go
package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SchemaResource{}
var _ resource.ResourceWithImportState = &SchemaResource{}

func NewSchemaResource() resource.Resource {
	return &SchemaResource{}
}

// SchemaResource defines the resource implementation.
type SchemaResource struct {
	client *ClientSchema
}

func (r *SchemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *SchemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required:            true,
			},
			"schema_type": schema.StringAttribute{
				MarkdownDescription: "Schema Type, defaults to AVRO",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("AVRO"),
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

func (r *SchemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ClientSchema)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SchemaModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schemaModel := &SchemaModel{
		Subject:    data.Subject,
		Schema:     data.Schema,
		SchemaType: data.SchemaType,
	}

	err := r.client.CreateSchema(schemaModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	data.Version = schemaModel.Version
	data.Id = schemaModel.Id

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SchemaModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schemaModel := &SchemaModel{
		Subject:    data.Subject,
		Schema:     data.Schema,
		SchemaType: data.SchemaType,
	}

	err := r.client.UpdateSchema(schemaModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	data.Version = schemaModel.Version
	data.Id = schemaModel.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSchema(data.Subject.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

}

func (r *SchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("version"), req, resp)
}
