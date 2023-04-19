// resource_redpanda_schema.go
package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRedPandaSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedPandaSchemaCreate,
		ReadContext:   resourceRedPandaSchemaRead,
		UpdateContext: resourceRedPandaSchemaUpdate,
		DeleteContext: resourceRedPandaSchemaDelete,
		Schema: map[string]*schema.Schema{
			"subject": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AVRO",
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRedPandaSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	schema := &Schema{
		Subject:    d.Get("subject").(string),
		Schema:     d.Get("schema").(string),
		SchemaType: d.Get("schema_type").(string),
	}

	err := client.CreateSchema(schema)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(schema.Subject)
	d.Set("version", schema.Version)

	return resourceRedPandaSchemaRead(ctx, d, m)
}

func resourceRedPandaSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	schema, err := client.GetSchema(d.Id(), "latest")
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("subject", schema.Subject)
	d.Set("schema", schema.Schema)
	d.Set("schema_type", schema.SchemaType)
	d.Set("version", schema.Version)

	return nil
}

func resourceRedPandaSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	schema := &Schema{
		Subject:    d.Get("subject").(string),
		Schema:     d.Get("schema").(string),
		SchemaType: d.Get("schema_type").(string),
	}

	err := client.UpdateSchema(schema)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRedPandaSchemaRead(ctx, d, m)
}

func resourceRedPandaSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	err := client.DeleteSchema(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
