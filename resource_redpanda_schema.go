// resource_redpanda_schema.go
package main

import (
	"context"
	"fmt"
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
	if err := d.Set("version", schema.Version); err != nil {
		return diag.FromErr(err)
	}

	return resourceRedPandaSchemaRead(ctx, d, m)
}

func resourceRedPandaSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	subject := d.Id()

	var version string
	if d.Get("version") == 0 {
		version = "latest"
	} else {
		version = fmt.Sprintf("%d", d.Get("version"))
	}

	schema, err := client.GetSchema(subject, version)
	if err != nil {
		return diag.FromErr(err)
	}

	if schema == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("subject", schema.Subject); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", schema.Schema); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema_type", schema.SchemaType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", schema.Version); err != nil {
		return diag.FromErr(err)
	}

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
