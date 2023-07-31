package provider

import (
	"context"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeCreate,
		ReadContext:   resourceUserAttributeRead,
		UpdateContext: resourceUserAttributeUpdate,
		DeleteContext: resourceUserAttributeDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "Name of user attribute.",
			},
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "Human-friendly label for user attribute",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Description:  "Type of user attribute ('string', 'number', 'datetime', 'yesno', 'zipcode')",
			},
			"value_is_hidden": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "If true, users will not be able to view values of this attribute.",
			},
			"user_can_view": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Non-admin users can see the values of their attributes and use them in filters.",
			},
			"user_can_edit": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Users can change the value of this attribute for themselves.",
			},
			"default_value": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Description:  "Default user attribute value",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserAttributeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	var diags diag.Diagnostics
	userAttr := &lookergo.UserAttribute{}
	if value, ok := d.GetOk("name"); ok {
		userAttr.Name = value.(string)
	}
	if value, ok := d.GetOk("label"); ok {
		userAttr.Label = value.(string)
	}
	if value, ok := d.GetOk("type"); ok {
		userAttr.Type = value.(string)
	}
	if value, ok := d.GetOk("value_is_hidden"); ok {
		userAttr.ValueIsHidden = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("user_can_view"); ok {
		userAttr.UserCanView = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("user_can_edit"); ok {
		userAttr.UserCanEdit = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("default_value"); ok {
		userAttr.DefaultValue = castToPtr(value.(string))
	}
	newAtt, _, err := c.UserAttributes.Create(ctx, userAttr)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*newAtt.Id)
	resourceUserAttributeRead(ctx, d, m)
	return diags
}

func resourceUserAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	var diags diag.Diagnostics
	UserAttr, _, err := c.UserAttributes.Get(ctx, idAsInt(d.Get("id")))
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("id", UserAttr.Id); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", UserAttr.Label); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", UserAttr.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", UserAttr.Type); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value_is_hidden", UserAttr.ValueIsHidden); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_can_view", UserAttr.UserCanView); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_can_edit", UserAttr.UserCanEdit); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("default_value", UserAttr.DefaultValue); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceUserAttributeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Id()
	UserAttr, _, err := c.UserAttributes.Get(ctx, idAsInt(d.Get("id")))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChanges() {
		if d.HasChange("name") {
			UserAttr.Name = d.Get("name").(string)
		}
		if d.HasChange("label") {
			UserAttr.Label = d.Get("label").(string)
		}
		if d.HasChange("type") {
			UserAttr.Type = d.Get("type").(string)
		}
		if d.HasChange("value_is_hidden") {
			UserAttr.ValueIsHidden = boolPtr(d.Get("value_is_hidden").(bool))
		}
		if d.HasChange("user_can_view") {
			UserAttr.UserCanView = boolPtr(d.Get("user_can_view").(bool))
		}
		if d.HasChange("user_can_edit") {
			UserAttr.UserCanEdit = boolPtr(d.Get("user_can_edit").(bool))
		}
		if d.HasChange("default_value") {
			UserAttr.DefaultValue = castToPtr(d.Get("default_value").(string))
		}

		if _, _, err = c.UserAttributes.Update(ctx, ID, UserAttr); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserAttributeRead(ctx, d, m)
}

func resourceUserAttributeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	var diags diag.Diagnostics
	ID := d.Id()
	_, err := c.UserAttributes.Delete(ctx, ID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
