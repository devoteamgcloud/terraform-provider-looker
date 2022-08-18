package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

// -
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: `
`,
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Computed:     false,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"soft_delete": {
				Description: "Only delete terraform reference to resource, keep actual resource on remote.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_groups": {
				Type:     schema.TypeSet,
				Computed: true,
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

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	tflog.Info(ctx, "Creating Looker Group")

	newGroup, _, err := c.Groups.Create(ctx, &lookergo.Group{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idAsString(newGroup.Id))

	// Ensure resource created
	// resourceGroupRead(ctx, d, m)

	tflog.Info(ctx, "Created Looker Group", map[string]interface{}{"id": newGroup.Id, "name": newGroup.Name})

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	groupID := idAsInt(d.Id())

	group, _, err := c.Groups.Get(ctx, groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("id", idAsString(group.Id)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	listGroups, _, err := c.Groups.ListById(ctx, []int{idAsInt(group.Id)}, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	group = &listGroups[0]

	if len(group.RoleIds) > 0 {
		if err = d.Set("roles", group.RoleIds.ToSliceOfStrings()); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err = d.Set("roles", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(group.ParentGroupIds) > 0 {
		if err = d.Set("parent_groups", group.ParentGroupIds.ToSliceOfStrings()); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err = d.Set("parent_groups", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	groupID := idAsInt(d.Id())

	group, _, err := c.Groups.Get(ctx, groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges() {
		if d.HasChange("name") {
			group.Name = d.Get("name").(string)
		}

		if _, _, err = c.Groups.Update(ctx, groupID, group); err != nil {
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	groupID := idAsInt(d.Id())

	if !d.Get("soft_delete").(bool) {
		if _, err := c.Groups.Delete(ctx, groupID); err != nil {
			return diag.FromErr(err)
		}
	}

	// Finally mark as deleted
	d.SetId("")

	return diags
}
