package provider

import (
	"context"
	"fmt"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
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
				Optional:    true,
				Default:     false,
			},
			"label": {
				Description: "UI-friendly name for this model",
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
			},
		},
		Importer: nil,
	}
}

func resourceLookMlModelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	lmlMdlName := d.Get("name").(string)
	projectName := d.Get("project_name").(string)
	dbConnNames := schemaSetToStringSlice(d.Get("allowed_db_connection_names").(*schema.Set))
	unlimitedConn := d.Get("unlimited_db_connections").(bool)

	//logDebug(ctx, "Create MlModel", "lmlMdlName", lmlMdlName, "projectName", projectName, "dbConnNames", dbConnNames, "unlimitedConnections", unlimitedConn)
	var lookmlModelOptions = lookergo.LookMLModel{Name: lmlMdlName, ProjectName: projectName, AllowedDbConnectionNames: dbConnNames, UnlimitedDbConnections: unlimitedConn}

	createdModel, _, err := c.LookMLModel.Create(ctx, &lookmlModelOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	d.Set("name", createdModel.Name)
	d.SetId(createdModel.Name)
	return resourceLookMlModelRead(ctx, d, m)
}

func resourceLookMlModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))
	lmlMdlName := d.Get("name").(string)
	newModel, response, err := c.LookMLModel.Get(ctx, lmlMdlName)
	if err != nil {
		return diag.FromErr(err)
	}
	if response.StatusCode == 404 {
		d.SetId("") // Mark as deleted
		return diags
	}
	if newModel == nil {
		return diag.FromErr(new(lookergo.ArgError))
	}
	if err = d.Set("project_name", newModel.ProjectName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("allowed_db_connection_names", newModel.AllowedDbConnectionNames); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("unlimited_db_connections", newModel.UnlimitedDbConnections); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", newModel.Label); err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceLookMlModelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))
	lmlMdlName := d.Get("name").(string)
	projectName := d.Get("project_name").(string)
	dbConnNames := schemaSetToStringSlice(d.Get("allowed_db_connection_names").(*schema.Set))
	unlimitedConn := d.Get("unlimited_db_connections").(bool)
	oldLookMl := lookergo.LookMLModel{Name: lmlMdlName, ProjectName: projectName, AllowedDbConnectionNames: dbConnNames, UnlimitedDbConnections: unlimitedConn}
	lookerML, _, err := c.LookMLModel.Update(ctx, d.Get("name").(string), &oldLookMl)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("project_name", lookerML.ProjectName); err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceLookMlModelRead(ctx, d, m)
}

func resourceLookMlModelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	_, err := c.LookMLModel.Delete(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
