// resource_redpanda_topic.go
package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRedPandaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedPandaTopicCreate,
		ReadContext:   resourceRedPandaTopicRead,
		UpdateContext: resourceRedPandaTopicUpdate,
		DeleteContext: resourceRedPandaTopicDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"replication_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceRedPandaTopicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	topic := &Topic{
		Name:             d.Get("name").(string),
		Partitions:       d.Get("partitions").(int),
		ReplicationFactor: d.Get("replication_factor").(int),
	}

	err := client.CreateTopic(topic)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(topic.Name)

	return resourceRedPandaTopicRead(ctx, d, m)
}

func resourceRedPandaTopicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	topic, err := client.GetTopic(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", topic.Name)
	d.Set("partitions", topic.Partitions)
	d.Set("replication_factor", topic.ReplicationFactor)

	return nil
}

func resourceRedPandaTopicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	topic := &Topic{
		Name:             d.Get("name").(string),
		Partitions:       d.Get("partitions").(int),
		ReplicationFactor: d.Get("replication_factor").(int),
	}

	err := client.UpdateTopic(topic)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRedPandaTopicRead(ctx, d, m)
}

func resourceRedPandaTopicDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	err := client.DeleteTopic(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
