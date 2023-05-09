// Package provider

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
var _ resource.Resource = &SchemaRegistryResource{}
var _ resource.ResourceWithImportState = &SchemaRegistryResource{}

func NewSchemaRegistryResource() resource.Resource {
	return &SchemaRegistryResource{}
}

// SchemaRegistryResource defines the resource implementation.
type SchemaRegistryResource struct {
	client *RedPandaProvider
}

func (r *SchemaRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *SchemaRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *SchemaRegistryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	apiURL := req.ProviderData.(*RedPandaProvider).SchemaRegistryApiUrl

	client := NewClientSchemaRegistry(apiURL)

	//if !ok {
	//	resp.Diagnostics.AddError(
	//		"Unexpected Resource Configure Type",
	//		fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	//	)
	//
	//	return
	//}

	r.client = client
}

func (r *SchemaRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SchemaRegistryModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	SchemaRegistryModel := &SchemaRegistryModel{
		Subject:    data.Subject,
		Schema:     data.Schema,
		SchemaType: data.SchemaType,
	}

	err := r.client.CreateSchema(SchemaRegistryModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	data.Version = SchemaRegistryModel.Version
	data.Id = SchemaRegistryModel.Id

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaRegistryModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SchemaRegistryModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	SchemaRegistryModel := &SchemaRegistryModel{
		Subject:    data.Subject,
		Schema:     data.Schema,
		SchemaType: data.SchemaType,
	}

	err := r.client.UpdateSchema(SchemaRegistryModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	data.Version = SchemaRegistryModel.Version
	data.Id = SchemaRegistryModel.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaRegistryModel

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

func (r *SchemaRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("version"), req, resp)
}
