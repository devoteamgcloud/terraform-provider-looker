package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func dataSourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserAttributeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of user attribute.",
			},
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-friendly label for user attribute",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of user attribute ('string', 'number', 'datetime', 'yesno', 'zipcode', 'relative_url', 'advanced_filter_string', 'advanced_filter_datetime', 'advanced_filter_number')",
			},
			"value_is_hidden": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, users will not be able to view values of this attribute.",
			},
			"user_can_view": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Non-admin users can see the values of their attributes and use them in filters.",
			},
			"user_can_edit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Users can change the value of this attribute for themselves.",
			},
			"default_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default value for when no value is set on the user.",
			},
		},
	}
}

func dataSourceUserAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	userAttr, _, err := c.UserAttributes.Get(ctx, idAsInt(d.Get("id")))
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("id", userAttr.Id); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", userAttr.Label); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", userAttr.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("default_value", userAttr.DefaultValue); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", userAttr.Type); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value_is_hidden", userAttr.ValueIsHidden); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_can_view", userAttr.UserCanView); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_can_edit", userAttr.UserCanEdit); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(userAttr.Id)
	return diags
}
