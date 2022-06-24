package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLookMlModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLookMlModelCreate,
		ReadContext:   resourceLookMlModelRead,
		UpdateContext: resourceLookMlModelUpdate,
		DeleteContext: resourceLookMlModelDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "LookML Model name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_name": {
				Description: "Project name LookML Model belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"allowed_db_connection_names": {
				Description: "List of allowed db connections (looker_connection)",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLookMlModelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	lmlMdlName := d.Get("name").(string)
	projectName := d.Get("project_name").(string)
	dbConnNames := schemaSetToStringSlice(d.Get("allowed_db_connection_names").(*schema.Set))

	logDebug(ctx, "Create MlModel", "lmlMdlName", lmlMdlName, "projectName", projectName, "dbConnNames", dbConnNames)

	_ = c

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceLookMlModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	lmlMdlName := d.Get("name").(string)
	projectName := d.Get("project_name").(string)
	dbConnNames := schemaSetToStringSlice(d.Get("allowed_db_connection_names").(*schema.Set))

	logDebug(ctx, "Create MlModel", "lmlMdlName", lmlMdlName, "projectName", projectName, "dbConnNames", dbConnNames)

	_ = c

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceLookMlModelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// TODO
	_ = c

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceLookMlModelRead(ctx, d, m)
}

func resourceLookMlModelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// TODO
	_ = c

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
