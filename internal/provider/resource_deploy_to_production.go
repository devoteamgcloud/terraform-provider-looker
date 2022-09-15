package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeployToProduction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployToProductionCreate,
		ReadContext:   resourceDeployToProductionRead,
		UpdateContext: resourceDeployToProductionUpdate,
		DeleteContext: resourceDeployToProductionDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
func resourceDeployToProductionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	projectName := d.Get("project_id").(string)
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	_, _, err = dc.Projects.DeployToProduction(ctx, projectName)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	d.Set("project_id", projectName)
	d.SetId("-")
	return resourceDeployToProductionRead(ctx, d, m)
}
func resourceDeployToProductionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	projectName := d.Get("project_id").(string)

	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	d.Set("project_id", projectName)
	return diags
}
func resourceDeployToProductionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	projectName := d.Get("project_id").(string)
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	_, _, err = dc.Projects.DeployToProduction(ctx, projectName)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("-")
	return resourceDeployToProductionRead(ctx, d, m)
}
func resourceDeployToProductionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	d.SetId("")
	return diags
}
