package provider

import (
	"context"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFolderPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFolderPermissionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Search folder based access metadata based on folder id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permissions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of permissions for the folder.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"can": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Operations the current user is able to perform on this object",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique Id",
						},
						"permission_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of permission: \"view\" or \"edit\" Valid values are: \"view\", \"edit\"",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of associated group",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of associated user",
						},
					},
				},
			},
		},
	}
}
func dataSourceFolderPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "Querying Looker Folder Permission")
	folderId := d.Get("id")
	if folderId == nil {
		return diag.Errorf("Folder ID is missing.")
	}
	d.SetId(folderId.(string))
	c := m.(*Config).Api
	permissionsHCL, _, err := c.ContentMetaGroupUser.ListByID(ctx, folderId.(string), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	permissions := flattenFolderPermissions(permissionsHCL)
	err = d.Set("permissions", permissions)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return nil
}
func flattenFolderPermissions(permissions []lookergo.ContentMetaGroupUser) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, permission := range permissions {
		result = append(result, map[string]interface{}{
			"can":             permission.Can,
			"id":              permission.Id,
			"permission_type": permission.PermissionType,
			"group_id":        permission.GroupId,
			"user_id":         permission.UserId,
		})
	}
	return result
}
