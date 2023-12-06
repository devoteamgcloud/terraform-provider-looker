package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
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
			// TODO
			// "begin_at": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "Timestamp for when this theme becomes active. Null=always",
			// },
			// "end_at": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "Timestamp for when this theme expires. Null=never",
			// },
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
						"tile_title_alignment": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The text alignment of tile titles (New Dashboards)",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringMatch(
									func() *regexp.Regexp {
										ret, _ := regexp.Compile("^(?i)(center|left|right)$")
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

func populateSettings(settingsMap map[string]interface{}) *lookergo.ThemeSettings {
	settings := &lookergo.ThemeSettings{
		BackgroundColor:          castToPtr(settingsMap["background_color"].(string)),
		BaseFontSize:             castToPtr(settingsMap["base_font_size"].(string)),
		ColorCollectionId:        castToPtr(settingsMap["color_collection_id"].(string)),
		FontColor:                castToPtr(settingsMap["font_color"].(string)),
		FontFamily:               castToPtr(settingsMap["font_family"].(string)),
		FontSource:               castToPtr(settingsMap["font_source"].(string)),
		PrimaryButtonColor:       castToPtr(settingsMap["primary_button_color"].(string)),
		ShowFiltersBar:           boolPtr(settingsMap["show_filters_bar"].(bool)),
		ShowTitle:                boolPtr(settingsMap["show_title"].(bool)),
		TextTileTextColor:        castToPtr(settingsMap["text_tile_text_color"].(string)),
		TileBackgroundColor:      castToPtr(settingsMap["tile_background_color"].(string)),
		TextTileBackgroundColor:  castToPtr(settingsMap["text_tile_background_color"].(string)),
		TileTextColor:            castToPtr(settingsMap["tile_text_color"].(string)),
		TitleColor:               castToPtr(settingsMap["title_color"].(string)),
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
	return settings
}

func resourceThemeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	theme := &lookergo.Theme{}
	// if value, ok := d.GetOk("begin_at"); ok {
	// 	parsedTime, err := time.Parse("2006-01-02 15:04:05",value.(string))
	// 	if err != nil {
	// 		fmt.Println("Error parsing time:", err)
	// 		return
	// 	}
	// 	theme.BeginAt = &parsedTime
	// }
	// if value, ok := d.GetOk("end_at"); ok {
	// 	theme.EndAt = value.(*time.Time)
	// }
	if value, ok := d.GetOk("name"); ok {
		theme.Name = castToPtr(value.(string))
	}
	if value, ok := d.GetOk("settings"); ok {
		for _, raw := range value.(*schema.Set).List() {
			settingsMap := raw.(map[string]interface{})
			theme.Settings = populateSettings(settingsMap)
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
	result := map[string]interface{}{
		"background_color":             settings.BackgroundColor,
		"base_font_size":               settings.BaseFontSize,
		"color_collection_id":          settings.ColorCollectionId,
		"font_color":                   settings.FontColor,
		"font_family":                  settings.FontFamily,
		"font_source":                  settings.FontSource,
		"primary_button_color":         settings.PrimaryButtonColor,
		"show_filters_bar":             settings.ShowFiltersBar,
		"show_title":                   settings.ShowTitle,
		"text_tile_text_color":         settings.TextTileTextColor,
		"tile_background_color":        settings.TileBackgroundColor,
		"text_tile_background_color":   settings.TextTileBackgroundColor,
		"tile_text_color":              settings.TileTextColor,
		"title_color":                  settings.TitleColor,
		"tile_title_alignment":         settings.TileTitleAlignment,
		"tile_shadow":                  settings.TileShadow,
		"show_last_updated_indicator":  settings.ShowLastUpdatedIndicator,
		"show_reload_data_icon":        settings.ShowReloadDataIcon,
		"show_dashboard_menu":          settings.ShowDashboardMenu,
		"show_filters_toggle":          settings.ShowFiltersToggle,
		"show_dashboard_header":        settings.ShowDashboardHeader,
		"center_dashboard_title":       settings.CenterDashboardTitle,
		"dashboard_title_font_size":    settings.DashboardTitleFontSize,
		"box_shadow":                   settings.BoxShadow,
		"page_margin_top":              settings.PageMarginTop,
		"page_margin_bottom":           settings.PageMarginBottom,
		"page_margin_sides":            settings.PageMarginSides,
		"show_explore_header":          settings.ShowExploreHeader,
		"show_explore_title":           settings.ShowExploreTitle,
		"show_explore_last_run":        settings.ShowExploreLastRun,
		"show_explore_timezone":        settings.ShowExploreTimezone,
		"show_explore_run_stop_button": settings.ShowExploreRunStopButton,
		"show_explore_actions_button":  settings.ShowExploreActionsButton,
		"show_look_header":             settings.ShowLookHeader,
		"show_look_title":              settings.ShowLookTitle,
		"show_look_last_run":           settings.ShowLookLastRun,
		"show_look_timezone":           settings.ShowLookTimezone,
		"show_look_run_stop_button":    settings.ShowLookRunStopButton,
		"show_look_actions_button":     settings.ShowLookActionsButton,
		"tile_title_font_size":         settings.TileTitleFontSize,
		"column_gap_size":              settings.ColumnGapSize,
		"row_gap_size":                 settings.RowGapSize,
		"border_radius":                settings.BorderRadius,
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
	// if theme.BeginAt != nil {
	// 	if err = d.Set("begin_at", theme.BeginAt.String()); err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// }
	// if theme.EndAt != nil {
	// 	if err = d.Set("end_at", theme.EndAt); err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// }

	var settingsItems []interface{}
	settingsItems = append(settingsItems, flattenThemeSettings(theme.Settings))
	d.Set("settings", settingsItems)
	return diags
}

func resourceThemeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	id := d.Id()
	theme, _, err := c.Themes.Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	// if d.HasChange("begin_at") {
	// 	if value, ok := d.GetOk("begin_at"); ok {
	// 		theme.BeginAt = value.(*time.Time)
	// 	}
	// }
	// if d.HasChange("end_at") {
	// 	if value, ok := d.GetOk("end_at"); ok {
	// 		theme.EndAt = value.(*time.Time)
	// 	}
	// }
	if value, ok := d.GetOk("name"); ok {
		theme.Name = castToPtr(value.(string))
	}
	if d.HasChange("settings") {
		if value, ok := d.GetOk("settings"); ok {
			for _, raw := range value.(*schema.Set).List() {
				settingsMap := raw.(map[string]interface{})
				theme.Settings = populateSettings(settingsMap)
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
