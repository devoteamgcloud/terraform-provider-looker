package provider

import (
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

var (
	dataRoleKey = []string{
		"id",
		"name",
	}
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: dataRoleKey,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: dataRoleKey,
			},
			"permission_set_id": {
				Description: "PermissionSet ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"model_set_id": {
				Description: "Modelset ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	role := lookergo.Role{}

	if roleId, exists := d.GetOk("id"); exists {
		localRole, _, err := c.Roles.Get(ctx, roleId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
		role = lookergo.Role(*localRole)
	}

	if roleName, exists := d.GetOk("name"); exists {
		localRole, _, err := c.Roles.ListByName(ctx, roleName.(string), &lookergo.ListOptions{})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(localRole) > 0 {
			for _, irole := range localRole {
				if strings.EqualFold(irole.Name, roleName.(string)) {
					role = irole
				}
			}
			if role.Id == 0 {
				return diag.Errorf("Role not found.")
			}
		} else {
			return diag.Errorf("Role not found.")
		}
	}

	if err := d.Set("id", role.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", role.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("permission_set_id", role.PermissionSet.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("model_set_id", role.ModelSet.Id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(role.Id))
	return diags
}
