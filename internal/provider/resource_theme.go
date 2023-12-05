package provider

import (
	"context"
	"time"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTheme() *schema.Resource {
	return &schema.Resource{
		Description: `
`,
		CreateContext: resourceThemeCreate,
		ReadContext:   resourceThemeRead,
		UpdateContext: resourceThemeUpdate,
		DeleteContext: resourceThemeDelete,
		Schema: map[string]*schema.Schema{
			"begin_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Timestamp for when this theme becomes active. Null=always",
			},
			"end_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Timestamp for when this theme expires. Null=never",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Id",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of theme. Can only be alphanumeric and underscores.",
			},
			"settings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Settings for the theme",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"background_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default background color",
						},

						"base_font_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base font size for scaling fonts (only supported by legacy dashboards)",
						},
						"color_collection_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional. ID of color collection to use with the theme. Use an empty string for none.",
						},
						"font_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default font color",
						},
						"font_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Primary font family",
						},
						"font_source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source specification for font",
						},
						"info_button_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(DEPRECATED) Info button color",
						},
						"primary_button_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Primary button color",
						},
						"show_filters_bar": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show filters. Defaults to true.",
						},
						"show_title": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the title. Defaults to true.",
						},
						"text_tile_text_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text color for text tiles",
						},
						"tile_background_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Background color for tiles",
						},
						"text_tile_background_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Background color for text tiles",
						},
						"tile_text_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text color for tiles",
						},
						"title_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Color for titles",
						},
						"warn_button_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(DEPRECATED) Warning button color",
						},
						"tile_title_alignment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text alignment of tile titles (New Dashboards)",
						},
						"tile_shadow": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggles the tile shadow (not supported)",
						},
						"show_last_updated_indicator": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the dashboard last updated indicator. Defaults to true.",
						},
						"show_reload_data_icon": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show reload data icon/button. Defaults to true.",
						},
						"show_dashboard_menu": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the dashboard actions menu. Defaults to true.",
						},
						"show_filters_toggle": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the filters icon/toggle. Defaults to true.",
						},
						"show_dashboard_header": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the dashboard header. Defaults to true.",
						},
						"center_dashboard_title": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to center the dashboard title. Defaults to false.",
						},
						"dashboard_title_font_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dashboard title font size.",
						},
						"box_shadow": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default box shadow.",
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

func resourceThemeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	theme := &lookergo.Theme{}
	if value, ok := d.GetOk("begin_at"); ok {
		theme.BeginAt = value.(*time.Time)
	}
	if value, ok := d.GetOk("end_at"); ok {
		theme.EndAt = value.(*time.Time)
	}
	if value, ok := d.GetOk("name"); ok {
		theme.Name = castToPtr(value.(string))
	}
	if value, ok := d.GetOk("settings"); ok {
		for _, raw := range value.(*schema.Set).List() {
			settingsMap := raw.(map[string]interface{})
			theme.Settings = &lookergo.ThemeSettings{
				BackgroundColor:         castToPtr(settingsMap["background_color"].(string)),
				BaseFontSize:            castToPtr(settingsMap["base_font_size"].(string)),
				ColorCollectionId:       castToPtr(settingsMap["color_collection_id"].(string)),
				FontColor:               castToPtr(settingsMap["font_color"].(string)),
				FontFamily:              castToPtr(settingsMap["font_family"].(string)),
				FontSource:              castToPtr(settingsMap["font_source"].(string)),
				InfoButtonColor:         castToPtr(settingsMap["info_button_color"].(string)),
				PrimaryButtonColor:      castToPtr(settingsMap["primary_button_color"].(string)),
				ShowFiltersBar:          boolPtr(settingsMap["show_filters_bar"].(bool)),
				ShowTitle:               boolPtr(settingsMap["show_title"].(bool)),
				TextTileTextColor:       castToPtr(settingsMap["text_tile_text_color"].(string)),
				TileBackgroundColor:     castToPtr(settingsMap["tile_background_color"].(string)),
				TextTileBackgroundColor: castToPtr(settingsMap["text_tile_background_color"].(string)),
				TileTextColor:           castToPtr(settingsMap["tile_text_color"].(string)),
				TitleColor:              castToPtr(settingsMap["title_color"].(string)),
				WarnButtonColor:         castToPtr(settingsMap["warn_button_color"].(string)),
				TileTitleAlignment:      castToPtr(settingsMap["tile_title_alignment"].(string)),
			}
		}
	}
	new_theme, _, err := c.Themes.Create(ctx, theme)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*new_theme.Id)
	return resourceThemeRead(ctx, d, m)
}

func flattenThemeSettings(settings *schema.ResourceData) map[string]interface{} {
	if settings == nil {
		return nil
	}

	result := make(map[string]interface{})

	if value, ok := settings.GetOk("background_color"); ok {
		result["background_color"] = value.(string)
	}

	if value, ok := settings.GetOk("base_font_size"); ok {
		result["base_font_size"] = value.(string)
	}

	if value, ok := settings.GetOk("color_collection_id"); ok {
		result["color_collection_id"] = value.(string)
	}

	if value, ok := settings.GetOk("font_color"); ok {
		result["font_color"] = value.(string)
	}

	if value, ok := settings.GetOk("font_family"); ok {
		result["font_family"] = value.(string)
	}

	if value, ok := settings.GetOk("font_source"); ok {
		result["font_source"] = value.(string)
	}

	if value, ok := settings.GetOk("info_button_color"); ok {
		result["info_button_color"] = value.(string)
	}

	if value, ok := settings.GetOk("primary_button_color"); ok {
		result["primary_button_color"] = value.(string)
	}

	if value, ok := settings.GetOk("show_filters_bar"); ok {
		result["show_filters_bar"] = value.(bool)
	}

	if value, ok := settings.GetOk("show_title"); ok {
		result["show_title"] = value.(bool)
	}

	if value, ok := settings.GetOk("text_tile_text_color"); ok {
		result["text_tile_text_color"] = value.(string)
	}

	if value, ok := settings.GetOk("tile_background_color"); ok {
		result["tile_background_color"] = value.(string)
	}

	if value, ok := settings.GetOk("text_tile_background_color"); ok {
		result["text_tile_background_color"] = value.(string)
	}

	if value, ok := settings.GetOk("tile_text_color"); ok {
		result["tile_text_color"] = value.(string)
	}

	if value, ok := settings.GetOk("title_color"); ok {
		result["title_color"] = value.(string)
	}

	if value, ok := settings.GetOk("warn_button_color"); ok {
		result["warn_button_color"] = value.(string)
	}

	if value, ok := settings.GetOk("tile_title_alignment"); ok {
		result["tile_title_alignment"] = value.(string)
	}

	return result
}

func resourceThemeRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	id := d.Id()
	theme, _, err := c.Themes.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", *theme.Name); err != nil {
		return diag.FromErr(err)
	}
	if theme.BeginAt != nil {
		if err = d.Set("begin_at", theme.BeginAt); err != nil {
			return diag.FromErr(err)
		}
	}
	if theme.EndAt != nil {
		if err = d.Set("end_at", theme.EndAt); err != nil {
			return diag.FromErr(err)
		}
	}
	if theme.Settings != nil {
		if err = d.Set("settings", flattenThemeSettings(d)); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceThemeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	id := d.Id()
	theme, _, err := c.Themes.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if value, ok := d.GetOk("begin_at"); ok {
		theme.BeginAt = value.(*time.Time)
	}
	if value, ok := d.GetOk("end_at"); ok {
		theme.EndAt = value.(*time.Time)
	}
	if value, ok := d.GetOk("name"); ok {
		theme.Name = castToPtr(value.(string))
	}
	if value, ok := d.GetOk("settings"); ok {
		for _, raw := range value.(*schema.Set).List() {
			settingsMap := raw.(map[string]interface{})
			theme.Settings = &lookergo.ThemeSettings{
				BackgroundColor:         castToPtr(settingsMap["background_color"].(string)),
				BaseFontSize:            castToPtr(settingsMap["base_font_size"].(string)),
				ColorCollectionId:       castToPtr(settingsMap["color_collection_id"].(string)),
				FontColor:               castToPtr(settingsMap["font_color"].(string)),
				FontFamily:              castToPtr(settingsMap["font_family"].(string)),
				FontSource:              castToPtr(settingsMap["font_source"].(string)),
				InfoButtonColor:         castToPtr(settingsMap["info_button_color"].(string)),
				PrimaryButtonColor:      castToPtr(settingsMap["primary_button_color"].(string)),
				ShowFiltersBar:          boolPtr(settingsMap["show_filters_bar"].(bool)),
				ShowTitle:               boolPtr(settingsMap["show_title"].(bool)),
				TextTileTextColor:       castToPtr(settingsMap["text_tile_text_color"].(string)),
				TileBackgroundColor:     castToPtr(settingsMap["tile_background_color"].(string)),
				TextTileBackgroundColor: castToPtr(settingsMap["text_tile_background_color"].(string)),
				TileTextColor:           castToPtr(settingsMap["tile_text_color"].(string)),
				TitleColor:              castToPtr(settingsMap["title_color"].(string)),
				WarnButtonColor:         castToPtr(settingsMap["warn_button_color"].(string)),
				TileTitleAlignment:      castToPtr(settingsMap["tile_title_alignment"].(string)),
			}
		}
	}
	_, _, err = c.Themes.Update(ctx, id, theme)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceThemeRead(ctx, d, m)
}

func resourceThemeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	id := d.Id()
	_, err := c.Themes.Delete(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
