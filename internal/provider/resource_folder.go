package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceFolder() *schema.Resource {
	return &schema.Resource{
		Description: `
`,
		CreateContext: resourceFolderCreate,
		ReadContext:   resourceFolderRead,
		UpdateContext: resourceFolderUpdate,
		DeleteContext: resourceFolderDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Computed:     false,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			}, "parent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	tflog.Info(ctx, "Creating Looker Folder")
	folder := &lookergo.Folder{}
	if value, ok := d.GetOk("name"); ok {
		folder.Name = value.(string)
	}
	if value, ok := d.GetOk("parent_id"); ok {
		folder.ParentId = value.(string)
	}
	newFolder, _, err := c.Folders.Create(ctx, folder)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newFolder.Id)
	d.Set("name", newFolder.Name)
	d.Set("parent_id", newFolder.ParentId)

	tflog.Info(ctx, "Created Looker Folder", map[string]interface{}{"id": newFolder.Id, "name": newFolder.Name})

	return resourceFolderRead(ctx, d, m)
}

func resourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	FolderID := d.Id()

	Folder, _, err := c.Folders.Get(ctx, FolderID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("id", Folder.Id); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", Folder.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("parent_id", Folder.ParentId); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	FolderID := d.Id()

	Folder, _, err := c.Folders.Get(ctx, FolderID)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges() {
		if d.HasChange("name") {
			Folder.Name = d.Get("name").(string)
		}
		if d.HasChange("parent_id") {
			Folder.ParentId = d.Get("parent_id").(string)
		}

		if _, _, err = c.Folders.Update(ctx, FolderID, Folder); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFolderRead(ctx, d, m)
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	FolderID := d.Id()

	if _, err := c.Folders.Delete(ctx, FolderID); err != nil {
		return diag.FromErr(err)
	}
	// Finally mark as deleted
	d.SetId("")

	return diags
}
