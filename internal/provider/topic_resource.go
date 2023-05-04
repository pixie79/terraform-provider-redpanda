// Package provider.go
package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TopicResource{}
var _ resource.ResourceWithImportState = &TopicResource{}

func NewTopicResource() resource.Resource {
	return &TopicResource{}
}

// TopicResource defines the resource implementation.
type TopicResource struct {
	client *ClientTopic
}

func (r *TopicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_topic"
}

func (r *TopicResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Topic resource",

		Attributes: map[string]schema.Attribute{
			"topic": schema.StringAttribute{
				MarkdownDescription: "Topic name",
				Required:            true,
			},
			"replication_factor": schema.Int64Attribute{
				MarkdownDescription: "Topic Replication Factor",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"partitions": schema.Int64Attribute{
				MarkdownDescription: "Number of Partitions",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"cleanup_policy": schema.StringAttribute{
				MarkdownDescription: "Number of Partitions",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("compact"),
			},
		},
	}
}

func (r *TopicResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ClientTopic)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TopicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TopicModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	topicModel := &TopicModel{
		Topic:             data.Topic,
		Partitions:        data.Partitions,
		ReplicationFactor: data.ReplicationFactor,
	}

	err := r.client.CreateTopic(topicModel)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	data.Partitions = topicModel.Partitions
	data.ReplicationFactor = topicModel.ReplicationFactor

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TopicModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TopicModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//topicModel := &TopicModel{
	//	Topic:             data.Topic,
	//	ReplicationFactor: data.ReplicationFactor,
	//	Partitions:        data.Partitions,
	//}

	//err := r.client.UpdateTopic(topicModel)
	//if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
	//	return
	//}
	//
	//// Write logs using the tflog package
	//tflog.Trace(ctx, "created a resource")
	//
	//data.Partitions = topicModel.Partitions
	//data.ReplicationFactor = topicModel.ReplicationFactor

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TopicResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TopicModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//err := r.client.DeleteTopic(data.Topic.ValueString())
	//if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
	//	return
	//}
	//
	//// Write logs using the tflog package
	//tflog.Trace(ctx, "created a resource")

}

func (r *TopicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("version"), req, resp)
}
