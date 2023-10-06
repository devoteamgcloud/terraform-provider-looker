package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		Description: `Manage permission sets.
`,
		CreateContext: resourcePermissionSetCreate,
		ReadContext:   resourcePermissionSetRead,
		UpdateContext: resourcePermissionSetUpdate,
		DeleteContext: resourcePermissionSetDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Computed:     false,
				Required:     true,
				Description:  "The name of the permission set.",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"permissions": {
				Description: "List of permissions.",
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

func resourcePermissionSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	tflog.Info(ctx, "Creating Permission Set")
	permset := &lookergo.PermissionSet{}
	if value, ok := d.GetOk("name"); ok {
		permset.Name = value.(string)
	}
	if value, ok := d.GetOk("permissions"); ok {
		var permission_list []string
		for _, raw := range value.(*schema.Set).List() {
			permission_list = append(permission_list, raw.(string))
		}
		permset.Permissions = permission_list
	}
	new_permset, _, err := c.PermissionSets.Create(ctx, permset)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(new_permset.Id)
	d.Set("name", new_permset.Name)
	d.Set("permissions", new_permset.Permissions)

	return resourcePermissionSetRead(ctx, d, m)
}

func resourcePermissionSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	permissionSetID := d.Id()
	permissionSet, response, err := c.PermissionSets.Get(ctx, permissionSetID)
	if response.StatusCode == 404 {
		d.SetId("") // Mark as deleted
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("id", permissionSet.Id); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", permissionSet.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("permissions", permissionSet.Permissions); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePermissionSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	permissionSetID := d.Id()

	permissionSet, _, err := c.PermissionSets.Get(ctx, permissionSetID)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		permissionSet.Name = d.Get("name").(string)
	}
	if d.HasChange("permissions") {
		var value = d.Get("permissions").(*schema.Set)
		var permission_list []string
		for _, raw := range value.List() {
			permission_list = append(permission_list, raw.(string))
		}
		permissionSet.Permissions = permission_list
	}
	if _, _, err = c.PermissionSets.Update(ctx, permissionSetID, permissionSet); err != nil {
		return diag.FromErr(err)
	}

	return resourcePermissionSetRead(ctx, d, m)
}

func resourcePermissionSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	PermissionSetId := d.Id()

	if _, err := c.PermissionSets.Delete(ctx, PermissionSetId); err != nil {
		return diag.FromErr(err)
	}
	// Finally mark as deleted
	d.SetId("")

	return diags
}
