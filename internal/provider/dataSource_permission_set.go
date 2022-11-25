package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	psKey = []string{
		"id",
		"name",
	}
)

func dataSourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePermissionSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "Search set based on id.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: psKey,
			},
			"name": {
				Description:  "Search set based on name.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: psKey,
			},
			"permissions": {
				Description: "List of permissions.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePermissionSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Info(ctx, "Querying Looker Permission Set")
	var permissionSet = lookergo.PermissionSet{}
	if psId, exists := d.GetOk("id"); exists { // Query using ID
		ps, _, err := c.PermissionSets.Get(ctx, psId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if ps.Id != "" {
			permissionSet.Id = ps.Id
			permissionSet.Name = ps.Name
			permissionSet.Permissions = ps.Permissions
		}
	} else if psNameKey, exists := d.GetOk("name"); exists { // Query using Name
		psSet, _, err := c.PermissionSets.GetByName(ctx, psNameKey.(string), &lookergo.ListOptions{})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(psSet) > 0 {
			for _, ps := range psSet {
				if ps.Name == psNameKey.(string) {
					permissionSet.Id = ps.Id
					permissionSet.Name = ps.Name
					permissionSet.Permissions = ps.Permissions
					break
				}
			}
		} else {
			return diag.Errorf("Permission Set not found.")
		}

	} else {
		return diag.Errorf("Neither name, nor id provided.")
	}
	d.SetId(permissionSet.Id)
	d.Set("name", permissionSet.Name)
	d.Set("permissions", permissionSet.Permissions)
	return diags
}
