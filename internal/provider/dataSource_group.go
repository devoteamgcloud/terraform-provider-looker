package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	queryKey = []string{
		"id",
		"name",
	}
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "Search group based on id",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: queryKey,
			},
			"name": {
				Description:  "Search group based on name",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: queryKey,
			},
			"user_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Info(ctx, "Querying Looker Group")
	var group = lookergo.Group{}
	if groupIDKey, exists := d.GetOk("id"); exists { // Query using ID
		groups, _, err := c.Groups.ListById(ctx, []int{idAsInt(groupIDKey)}, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(groups) == 1 {
			group = groups[0]
		} else {
			return diag.Errorf("No results found for: %v", groupIDKey)
		}
	} else if groupNameKey, exists := d.GetOk("name"); exists { // Query using Name
		groups, _, err := c.Groups.ListByName(ctx, groupNameKey.(string), nil)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(groups) == 1 {
			group = groups[0]
		} else {
			return diag.Errorf("No results found for: %v", groupNameKey)
		}
	} else {
		return diag.Errorf("Neither name, nor id provided.")
	}

	d.SetId(idAsString(group.Id))
	d.Set("name", group.Name)
	d.Set("user_count", group.UserCount)
	d.Set("parent_groups", group.ParentGroupIds.ToSliceOfStrings())
	d.Set("roles", group.RoleIds.ToSliceOfStrings())

	return diags
}
