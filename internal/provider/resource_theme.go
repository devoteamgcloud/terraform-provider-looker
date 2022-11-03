package provider

import (
	"context"
	"time"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTheme() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThemeCreate,
		ReadContext:   resourceThemeRead,
		UpdateContext: resourceThemeUpdate,
		DeleteContext: resourceThemeDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Theme id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Theme name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"begin_at": {
				Description: "Timestamp when theme becomes active",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"end_at": {
				Description: "Timestamp when theme expires",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"soft_delete": {
				Description: "Only delete terraform reference to resource, keep actual resource on remote.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceThemeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)

	tflog.Info(ctx, "Creating Looker theme")
	newTheme, _, err := c.Themes.Create(ctx, &lookergo.Theme{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idAsString(newTheme.Id))

	tflog.Info(ctx, "Created Looker Theme", map[string]interface{}{"id": newTheme.Id, "name": newTheme.Name})

	return resourceThemeRead(ctx, d, m)
}

func resourceThemeRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	themeID := idAsInt(d.Id())

	theme, _, err := c.Themes.Get(ctx, themeID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("id", theme.Id); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", theme.Name); err != nil {
		return diag.FromErr(err)
	}

	listThemes, _, err := c.Themes.List(ctx, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	theme = &listThemes[0]

	return diags
}

func resourceThemeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config).Api // .(*lookergo.Client)
	themeID := idAsInt(d.Id())

	theme, _, err := c.Themes.Get(ctx, themeID)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges() {
		if d.HasChange("name") {
			theme.Name = d.Get("name").(string)
		}

		if _, _, err = c.Themes.Update(ctx, themeID, theme); err != nil {
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceThemeRead(ctx, d, m)
}

func resourceThemeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	themeID := idAsInt(d.Id())

	if !d.Get("soft_delete").(bo[ol) {
		if _, err := c.Themes.Delete(ctx, themeID); err != nil {
			return diag.FromErr(err)
		}
	}

	// Finally mark as deleted
	d.SetId("")

	return diags
}
