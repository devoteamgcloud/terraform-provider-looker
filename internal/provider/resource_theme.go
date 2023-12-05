package provider

import (
	"context"
	"time"
	"regexp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringMatch(
						func() *regexp.Regexp {
							ret, _ := regexp.Compile("^[a-zA-Z0-9_]+$")
							return ret
						}(),
						"'name' must be only alphanumeric characters and underscores",
					),
				),
			},
			"settings": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Settings for the theme",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"background_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Default background color",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
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
							Required:    true,
							Description: "Default font color",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"font_family": {
							Type:        schema.TypeString,
							Required:    true,
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
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"primary_button_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Primary button color",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
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
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"tile_background_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Background color for tiles",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"text_tile_background_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Background color for text tiles",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"tile_text_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Text color for tiles",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"title_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Color for titles",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"warn_button_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(DEPRECATED) Warning button color",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("#(?i)(?:[0-9A-F]{2}){2,4}")
										return ret
									}(),
									"color must be a valid color hex code",
								),
							),
						},
						"tile_title_alignment": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text alignment of tile titles (New Dashboards)",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("^(center|left|right)$")
										return ret
									}(),
									"Invalid value for alignment. Allowed values are 'center', 'left', or 'right'.",
								),
							),
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
						"page_margin_top": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dashboard page margin top.",
						},
						"page_margin_bottom": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dashboard page margin bottom.",
						},
						"page_margin_sides": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dashboard page margin left and right.",
						},
						"show_explore_header": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the explore page header. Defaults to true.",
						},
						"show_explore_title": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the explore page title. Defaults to true.",
						},
						"show_explore_last_run": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the explore page last run. Defaults to true.",
						},
						"show_explore_timezone": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the explore page timezone. Defaults to true.",
						},
						"show_explore_run_stop_button": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the explore page run button. Defaults to true.",
						},
						"show_explore_actions_button": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the explore page actions button. Defaults to true.",
						},
						"show_look_header": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the look page header. Defaults to true.",
						},
						"show_look_title": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the look page title. Defaults to true.",
						},
						"show_look_last_run": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the look page last run. Defaults to true.",
						},
						"show_look_timezone": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the look page timezone Defaults to true.",
						},
						"show_look_run_stop_button": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the look page run button. Defaults to true.",
						},
						"show_look_actions_button": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Toggle to show the look page actions button. Defaults to true.",
						},
						"tile_title_font_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Font size for tiles.",
						},
						"column_gap_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The vertical gap/gutter size between tiles.",
						},
						"row_gap_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The horizontal gap/gutter size between tiles.",
						},
						"border_radius": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The border radius for tiles.",
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
				BackgroundColor:          castToPtr(settingsMap["background_color"].(string)),
				BaseFontSize:             castToPtr(settingsMap["base_font_size"].(string)),
				ColorCollectionId:        castToPtr(settingsMap["color_collection_id"].(string)),
				FontColor:                castToPtr(settingsMap["font_color"].(string)),
				FontFamily:               castToPtr(settingsMap["font_family"].(string)),
				FontSource:               castToPtr(settingsMap["font_source"].(string)),
				InfoButtonColor:          castToPtr(settingsMap["info_button_color"].(string)),
				PrimaryButtonColor:       castToPtr(settingsMap["primary_button_color"].(string)),
				ShowFiltersBar:           boolPtr(settingsMap["show_filters_bar"].(bool)),
				ShowTitle:                boolPtr(settingsMap["show_title"].(bool)),
				TextTileTextColor:        castToPtr(settingsMap["text_tile_text_color"].(string)),
				TileBackgroundColor:      castToPtr(settingsMap["tile_background_color"].(string)),
				TextTileBackgroundColor:  castToPtr(settingsMap["text_tile_background_color"].(string)),
				TileTextColor:            castToPtr(settingsMap["tile_text_color"].(string)),
				TitleColor:               castToPtr(settingsMap["title_color"].(string)),
				WarnButtonColor:          castToPtr(settingsMap["warn_button_color"].(string)),
				TileTitleAlignment:       castToPtr(settingsMap["tile_title_alignment"].(string)),
				TileShadow:               boolPtr(settingsMap["tile_shadow"].(bool)),
				ShowLastUpdatedIndicator: boolPtr(settingsMap["show_last_updated_indicator"].(bool)),
				ShowReloadDataIcon:       boolPtr(settingsMap["show_reload_data_icon"].(bool)),
				ShowDashboardMenu:        boolPtr(settingsMap["show_dashboard_menu"].(bool)),
				ShowFiltersToggle:        boolPtr(settingsMap["show_filters_toggle"].(bool)),
				ShowDashboardHeader:      boolPtr(settingsMap["show_dashboard_header"].(bool)),
				CenterDashboardTitle:     boolPtr(settingsMap["center_dashboard_title"].(bool)),
				DashboardTitleFontSize:   castToPtr(settingsMap["dashboard_title_font_size"].(string)),
				BoxShadow:                castToPtr(settingsMap["box_shadow"].(string)),
				PageMarginTop:            castToPtr(settingsMap["page_margin_top"].(string)),
				PageMarginBottom:         castToPtr(settingsMap["page_margin_bottom"].(string)),
				PageMarginSides:          castToPtr(settingsMap["page_margin_sides"].(string)),
				ShowExploreHeader:        boolPtr(settingsMap["show_explore_header"].(bool)),
				ShowExploreTitle:         boolPtr(settingsMap["show_explore_title"].(bool)),
				ShowExploreLastRun:       boolPtr(settingsMap["show_explore_last_run"].(bool)),
				ShowExploreTimezone:      boolPtr(settingsMap["show_explore_timezone"].(bool)),
				ShowExploreRunStopButton: boolPtr(settingsMap["show_explore_run_stop_button"].(bool)),
				ShowExploreActionsButton: boolPtr(settingsMap["show_explore_actions_button"].(bool)),
				ShowLookHeader:           boolPtr(settingsMap["show_look_header"].(bool)),
				ShowLookTitle:            boolPtr(settingsMap["show_look_title"].(bool)),
				ShowLookLastRun:          boolPtr(settingsMap["show_look_last_run"].(bool)),
				ShowLookTimezone:         boolPtr(settingsMap["show_look_timezone"].(bool)),
				ShowLookRunStopButton:    boolPtr(settingsMap["show_look_run_stop_button"].(bool)),
				ShowLookActionsButton:    boolPtr(settingsMap["show_look_actions_button"].(bool)),
				TileTitleFontSize:        castToPtr(settingsMap["tile_title_font_size"].(string)),
				ColumnGapSize:            castToPtr(settingsMap["column_gap_size"].(string)),
				RowGapSize:               castToPtr(settingsMap["row_gap_size"].(string)),
				BorderRadius:             castToPtr(settingsMap["border_radius"].(string)),
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

func flattenThemeSettings(settings *lookergo.ThemeSettings) map[string]interface{} {
	if settings == nil {
		return nil
	}

	result := make(map[string]interface{})

	result["background_color"] = getValueOrDefault(settings.BackgroundColor)
	result["base_font_size"] = getValueOrDefault(settings.BaseFontSize)
	result["color_collection_id"] = getValueOrDefault(settings.ColorCollectionId)
	result["font_color"] = getValueOrDefault(settings.FontColor)
	result["font_family"] = getValueOrDefault(settings.FontFamily)
	result["font_source"] = getValueOrDefault(settings.FontSource)
	result["info_button_color"] = getValueOrDefault(settings.InfoButtonColor)
	result["primary_button_color"] = getValueOrDefault(settings.PrimaryButtonColor)
	result["show_filters_bar"] = getValueOrDefault(settings.ShowFiltersBar)
	result["show_title"] = getValueOrDefault(settings.ShowTitle)
	result["text_tile_text_color"] = getValueOrDefault(settings.TextTileTextColor)
	result["tile_background_color"] = getValueOrDefault(settings.TileBackgroundColor)
	result["text_tile_background_color"] = getValueOrDefault(settings.TextTileBackgroundColor)
	result["tile_text_color"] = getValueOrDefault(settings.TileTextColor)
	result["title_color"] = getValueOrDefault(settings.TitleColor)
	result["warn_button_color"] = getValueOrDefault(settings.WarnButtonColor)
	result["tile_title_alignment"] = getValueOrDefault(settings.TileTitleAlignment)
	result["tile_shadow"] = getValueOrDefault(settings.TileShadow)
	result["show_last_updated_indicator"] = getValueOrDefault(settings.ShowLastUpdatedIndicator)
	result["show_reload_data_icon"] = getValueOrDefault(settings.ShowReloadDataIcon)
	result["show_dashboard_menu"] = getValueOrDefault(settings.ShowDashboardMenu)
	result["show_filters_toggle"] = getValueOrDefault(settings.ShowFiltersToggle)
	result["show_dashboard_header"] = getValueOrDefault(settings.ShowDashboardHeader)
	result["center_dashboard_title"] = getValueOrDefault(settings.CenterDashboardTitle)
	result["dashboard_title_font_size"] = getValueOrDefault(settings.DashboardTitleFontSize)
	result["box_shadow"] = getValueOrDefault(settings.BoxShadow)
	result["page_margin_top"] = getValueOrDefault(settings.PageMarginTop)
	result["page_margin_bottom"] = getValueOrDefault(settings.PageMarginBottom)
	result["page_margin_sides"] = getValueOrDefault(settings.PageMarginSides)
	result["show_explore_header"] = getValueOrDefault(settings.ShowExploreHeader)
	result["show_explore_title"] = getValueOrDefault(settings.ShowExploreTitle)
	result["show_explore_last_run"] = getValueOrDefault(settings.ShowExploreLastRun)
	result["show_explore_timezone"] = getValueOrDefault(settings.ShowExploreTimezone)
	result["show_explore_run_stop_button"] = getValueOrDefault(settings.ShowExploreRunStopButton)
	result["show_explore_actions_button"] = getValueOrDefault(settings.ShowExploreActionsButton)
	result["show_look_header"] = getValueOrDefault(settings.ShowLookHeader)
	result["show_look_title"] = getValueOrDefault(settings.ShowLookTitle)
	result["show_look_last_run"] = getValueOrDefault(settings.ShowLookLastRun)
	result["show_look_timezone"] = getValueOrDefault(settings.ShowLookTimezone)
	result["show_look_run_stop_button"] = getValueOrDefault(settings.ShowLookRunStopButton)
	result["show_look_actions_button"] = getValueOrDefault(settings.ShowLookActionsButton)
	result["tile_title_font_size"] = getValueOrDefault(settings.TileTitleFontSize)
	result["column_gap_size"] = getValueOrDefault(settings.ColumnGapSize)
	result["row_gap_size"] = getValueOrDefault(settings.RowGapSize)
	result["border_radius"] = getValueOrDefault(settings.BorderRadius)

	return result
}

func getValueOrDefault(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return value
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
		if err = d.Set("begin_at", theme.BeginAt.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if theme.EndAt != nil {
		if err = d.Set("end_at", theme.EndAt); err != nil {
			return diag.FromErr(err)
		}
	}
	// if theme.Settings != nil {
	// 	if err = d.Set("settings", theme.Settings); err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// }
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
				BackgroundColor:          castToPtr(settingsMap["background_color"].(string)),
				BaseFontSize:             castToPtr(settingsMap["base_font_size"].(string)),
				ColorCollectionId:        castToPtr(settingsMap["color_collection_id"].(string)),
				FontColor:                castToPtr(settingsMap["font_color"].(string)),
				FontFamily:               castToPtr(settingsMap["font_family"].(string)),
				FontSource:               castToPtr(settingsMap["font_source"].(string)),
				InfoButtonColor:          castToPtr(settingsMap["info_button_color"].(string)),
				PrimaryButtonColor:       castToPtr(settingsMap["primary_button_color"].(string)),
				ShowFiltersBar:           boolPtr(settingsMap["show_filters_bar"].(bool)),
				ShowTitle:                boolPtr(settingsMap["show_title"].(bool)),
				TextTileTextColor:        castToPtr(settingsMap["text_tile_text_color"].(string)),
				TileBackgroundColor:      castToPtr(settingsMap["tile_background_color"].(string)),
				TextTileBackgroundColor:  castToPtr(settingsMap["text_tile_background_color"].(string)),
				TileTextColor:            castToPtr(settingsMap["tile_text_color"].(string)),
				TitleColor:               castToPtr(settingsMap["title_color"].(string)),
				WarnButtonColor:          castToPtr(settingsMap["warn_button_color"].(string)),
				TileTitleAlignment:       castToPtr(settingsMap["tile_title_alignment"].(string)),
				TileShadow:               boolPtr(settingsMap["tile_shadow"].(bool)),
				ShowLastUpdatedIndicator: boolPtr(settingsMap["show_last_updated_indicator"].(bool)),
				ShowReloadDataIcon:       boolPtr(settingsMap["show_reload_data_icon"].(bool)),
				ShowDashboardMenu:        boolPtr(settingsMap["show_dashboard_menu"].(bool)),
				ShowFiltersToggle:        boolPtr(settingsMap["show_filters_toggle"].(bool)),
				ShowDashboardHeader:      boolPtr(settingsMap["show_dashboard_header"].(bool)),
				CenterDashboardTitle:     boolPtr(settingsMap["center_dashboard_title"].(bool)),
				DashboardTitleFontSize:   castToPtr(settingsMap["dashboard_title_font_size"].(string)),
				BoxShadow:                castToPtr(settingsMap["box_shadow"].(string)),
				PageMarginTop:            castToPtr(settingsMap["page_margin_top"].(string)),
				PageMarginBottom:         castToPtr(settingsMap["page_margin_bottom"].(string)),
				PageMarginSides:          castToPtr(settingsMap["page_margin_sides"].(string)),
				ShowExploreHeader:        boolPtr(settingsMap["show_explore_header"].(bool)),
				ShowExploreTitle:         boolPtr(settingsMap["show_explore_title"].(bool)),
				ShowExploreLastRun:       boolPtr(settingsMap["show_explore_last_run"].(bool)),
				ShowExploreTimezone:      boolPtr(settingsMap["show_explore_timezone"].(bool)),
				ShowExploreRunStopButton: boolPtr(settingsMap["show_explore_run_stop_button"].(bool)),
				ShowExploreActionsButton: boolPtr(settingsMap["show_explore_actions_button"].(bool)),
				ShowLookHeader:           boolPtr(settingsMap["show_look_header"].(bool)),
				ShowLookTitle:            boolPtr(settingsMap["show_look_title"].(bool)),
				ShowLookLastRun:          boolPtr(settingsMap["show_look_last_run"].(bool)),
				ShowLookTimezone:         boolPtr(settingsMap["show_look_timezone"].(bool)),
				ShowLookRunStopButton:    boolPtr(settingsMap["show_look_run_stop_button"].(bool)),
				ShowLookActionsButton:    boolPtr(settingsMap["show_look_actions_button"].(bool)),
				TileTitleFontSize:        castToPtr(settingsMap["tile_title_font_size"].(string)),
				ColumnGapSize:            castToPtr(settingsMap["column_gap_size"].(string)),
				RowGapSize:               castToPtr(settingsMap["row_gap_size"].(string)),
				BorderRadius:             castToPtr(settingsMap["border_radius"].(string)),
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
