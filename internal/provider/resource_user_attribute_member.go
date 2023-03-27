package provider

import (
	"context"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributeMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeMemberCreate,
		ReadContext:   resourceUserAttributeMemberRead,
		UpdateContext: resourceUserAttributeMemberUpdate,
		DeleteContext: resourceUserAttributeMemberDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_attribute_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserAttributeMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	var diags diag.Diagnostics
	ID := d.Get("user_attribute_id").(string)
	groupSet, ok := d.GetOk("group")
	userAttrs := []lookergo.UserAttributeGroupValue{}
	if ok {
		for _, raw := range groupSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			val := obj["id"].(string)
			att := lookergo.UserAttributeGroupValue{}
			att.GroupId = val
			att.UserAttributeId = d.Get("user_attribute_id").(string)
			att.Value = obj["value"].(string)
			userAttrs = append(userAttrs, att)
		}
	}
	_, _, err := c.UserAttributes.SetUserAttributeValue(ctx, userAttrs, ID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("-")
	resourceUserAttributeMemberRead(ctx, d, m)
	return diags
}

func resourceUserAttributeMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	var diags diag.Diagnostics
	ID := d.Get("user_attribute_id").(string)
	attrs, _, err := c.UserAttributes.GetUserAttributeValue(ctx,ID)
	if err != nil {
		return diag.FromErr(err)
	}
	var attrItems []interface{}
	for _, attr := range *attrs {
		group,_,err := c.Groups.Get(ctx, idAsInt(attr.GroupId))
		if err != nil {
			return diag.FromErr(err)
		} 
		attrItems = append(attrItems, map[string]interface{}{"id": idAsString(group.Id), "name": group.Name, "value":attr.Value})
	}
	d.Set("group", attrItems)
	return diags
}

func resourceUserAttributeMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Get("user_attribute_id").(string)
	groupSet, ok := d.GetOk("group")
	userAttrs := []lookergo.UserAttributeGroupValue{}
	if ok {
		for _, raw := range groupSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			val := obj["id"].(string)
			att := lookergo.UserAttributeGroupValue{}
			att.GroupId = val
			att.UserAttributeId = d.Get("user_attribute_id").(string)
			att.Value = obj["value"].(string)
			userAttrs = append(userAttrs, att)
		}
	}
	_, _, err := c.UserAttributes.SetUserAttributeValue(ctx, userAttrs, ID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("-")
	return resourceUserAttributeMemberRead(ctx, d, m)
}

func resourceUserAttributeMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	var diags diag.Diagnostics
	att := []lookergo.UserAttributeGroupValue{}
	ID := d.Get("user_attribute_id").(string)
	c.UserAttributes.SetUserAttributeValue(ctx, att, ID)
	d.SetId("")
	return diags
}
