package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	folderKey = []string{
		"id",
		"name",
	}
)

func dataSourceFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFolderRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "Search folder based on id",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: folderKey,
			},
			"name": {
				Description:  "Search folder based on name",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: folderKey,
			},
			"parent_id": {
				Description:  "Id of the parent folder",
				Type:         schema.TypeString,
				Computed:     true,
			},
		},
	}
}

func dataSourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Info(ctx, "Querying Looker Group")
	var folder = lookergo.Folder{}
	if folderId, exists := d.GetOk("id"); exists { // Query using ID
		newfolder, _, err := c.Folders.Get(ctx, folderId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if newfolder != nil{
			folder.Id = newfolder.Id
			folder.Name = newfolder.Name
		}else {
			return diag.Errorf("Folder not found.")
		}
		
	} else if folderNameKey, exists := d.GetOk("name"); exists { // Query using Name
		folders, _, err := c.Folders.ListByName(ctx, folderNameKey.(string), &lookergo.ListOptions{})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(folders) > 0 {
			for _, newfolder := range folders {
				if newfolder.Name == folderNameKey.(string) {
					folder.Id = newfolder.Id
					folder.Name = newfolder.Name
				}
			}
		} else {
			return diag.Errorf("Folder not found.")
		}
	} else {
		return diag.Errorf("Neither name, nor id provided.")
	}

	d.SetId(folder.Id)
	d.Set("name", folder.Name)
	return diags
}
