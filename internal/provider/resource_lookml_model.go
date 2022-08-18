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
				Description: "Name of the model. Also used as the unique identifier",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_name": {
				Description: "Name of project containing the model",
				Type:        schema.TypeString,
				Required:    true,
			},
			"allowed_db_connection_names": {
				Description: "Array of names of connections this model is allowed to use (looker_connection)",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"unlimited_db_connections": {
				Description: "Is this model allowed to use all current and future connections?",
				Type:        schema.TypeBool,
				Required:    false,
			},
			"label": {
				Description: "UI-friendly name for this model",
				Type:        schema.TypeString,
				Required:    false,
				Computed: true,
			},
			"has_content": {
				Description: "Does this model declaration have have lookml content?",
				Type:        schema.TypeBool,
				Required:    false,
				Computed: true,
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
