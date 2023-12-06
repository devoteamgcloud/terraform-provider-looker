package lookergo

import (
	"context"
	"time"
)

const ThemesBasePath = "4.0/themes"

type ThemesResource interface {
	Get(context.Context, string) (*Theme, *Response, error)
	Create(context.Context, *Theme) (*Theme, *Response, error)
	Update(context.Context, string, *Theme) (*Theme, *Response, error)
	Delete(context.Context, string) (*Response, error)
}
type ThemeSettings struct {
	BackgroundColor          *string `json:"background_color,omitempty"`             // Default background color
	BaseFontSize             *string `json:"base_font_size,omitempty"`               // Base font size for scaling fonts (only supported by legacy dashboards)
	ColorCollectionId        *string `json:"color_collection_id,omitempty"`          // Optional. ID of color collection to use with the theme. Use an empty string for none.
	FontColor                *string `json:"font_color,omitempty"`                   // Default font color
	FontFamily               *string `json:"font_family,omitempty"`                  // Primary font family
	FontSource               *string `json:"font_source,omitempty"`                  // Source specification for font
	InfoButtonColor          *string `json:"info_button_color,omitempty"`            // (DEPRECATED) Info button color
	PrimaryButtonColor       *string `json:"primary_button_color,omitempty"`         // Primary button color
	ShowFiltersBar           *bool   `json:"show_filters_bar,omitempty"`             // Toggle to show filters. Defaults to true.
	ShowTitle                *bool   `json:"show_title,omitempty"`                   // Toggle to show the title. Defaults to true.
	TextTileTextColor        *string `json:"text_tile_text_color,omitempty"`         // Text color for text tiles
	TileBackgroundColor      *string `json:"tile_background_color,omitempty"`        // Background color for tiles
	TextTileBackgroundColor  *string `json:"text_tile_background_color,omitempty"`   // Background color for text tiles
	TileTextColor            *string `json:"tile_text_color,omitempty"`              // Text color for tiles
	TitleColor               *string `json:"title_color,omitempty"`                  // Color for titles
	WarnButtonColor          *string `json:"warn_button_color,omitempty"`            // (DEPRECATED) Warning button color
	TileTitleAlignment       *string `json:"tile_title_alignment,omitempty"`         // The text alignment of tile titles (New Dashboards)
	TileShadow               *bool   `json:"tile_shadow,omitempty"`                  // Toggles the tile shadow (not supported)
	ShowLastUpdatedIndicator *bool   `json:"show_last_updated_indicator,omitempty"`  // Toggle to show the dashboard last updated indicator. Defaults to true.
	ShowReloadDataIcon       *bool   `json:"show_reload_data_icon,omitempty"`        // Toggle to show reload data icon/button. Defaults to true.
	ShowDashboardMenu        *bool   `json:"show_dashboard_menu,omitempty"`          // Toggle to show the dashboard actions menu. Defaults to true.
	ShowFiltersToggle        *bool   `json:"show_filters_toggle,omitempty"`          // Toggle to show the filters icon/toggle. Defaults to true.
	ShowDashboardHeader      *bool   `json:"show_dashboard_header,omitempty"`        // Toggle to show the dashboard header. Defaults to true.
	CenterDashboardTitle     *bool   `json:"center_dashboard_title,omitempty"`       // Toggle to center the dashboard title. Defaults to false.
	DashboardTitleFontSize   *string `json:"dashboard_title_font_size,omitempty"`    // Dashboard title font size.
	BoxShadow                *string `json:"box_shadow,omitempty"`                   // Default box shadow.
	PageMarginTop            *string `json:"page_margin_top,omitempty"`              // Dashboard page margin top.
	PageMarginBottom         *string `json:"page_margin_bottom,omitempty"`           // Dashboard page margin bottom.
	PageMarginSides          *string `json:"page_margin_sides,omitempty"`            // Dashboard page margin left and right.
	ShowExploreHeader        *bool   `json:"show_explore_header,omitempty"`          // Toggle to show the explore page header. Defaults to true.
	ShowExploreTitle         *bool   `json:"show_explore_title,omitempty"`           // Toggle to show the explore page title. Defaults to true.
	ShowExploreLastRun       *bool   `json:"show_explore_last_run,omitempty"`        // Toggle to show the explore page last run. Defaults to true.
	ShowExploreTimezone      *bool   `json:"show_explore_timezone,omitempty"`        // Toggle to show the explore page timezone. Defaults to true.
	ShowExploreRunStopButton *bool   `json:"show_explore_run_stop_button,omitempty"` // Toggle to show the explore page run button. Defaults to true.
	ShowExploreActionsButton *bool   `json:"show_explore_actions_button,omitempty"`  // Toggle to show the explore page actions button. Defaults to true.
	ShowLookHeader           *bool   `json:"show_look_header,omitempty"`             // Toggle to show the look page header. Defaults to true.
	ShowLookTitle            *bool   `json:"show_look_title,omitempty"`              // Toggle to show the look page title. Defaults to true.
	ShowLookLastRun          *bool   `json:"show_look_last_run,omitempty"`           // Toggle to show the look page last run. Defaults to true.
	ShowLookTimezone         *bool   `json:"show_look_timezone,omitempty"`           // Toggle to show the look page timezone Defaults to true.
	ShowLookRunStopButton    *bool   `json:"show_look_run_stop_button,omitempty"`    // Toggle to show the look page run button. Defaults to true.
	ShowLookActionsButton    *bool   `json:"show_look_actions_button,omitempty"`     // Toggle to show the look page actions button. Defaults to true.
	TileTitleFontSize        *string `json:"tile_title_font_size,omitempty"`         // Font size for tiles.
	ColumnGapSize            *string `json:"column_gap_size,omitempty"`              // The vertical gap/gutter size between tiles.
	RowGapSize               *string `json:"row_gap_size,omitempty"`                 // The horizontal gap/gutter size between tiles.
	BorderRadius             *string `json:"border_radius,omitempty"`                // The border radius for tiles.
}
type Theme struct {
	Can      *map[string]bool `json:"can,omitempty"`      // Operations the current user is able to perform on this object
	BeginAt  *time.Time       `json:"begin_at,omitempty"` // Timestamp for when this theme becomes active. Null=always
	EndAt    *time.Time       `json:"end_at,omitempty"`   // Timestamp for when this theme expires. Null=never
	Id       *string          `json:"id,omitempty"`       // Unique Id
	Name     *string          `json:"name,omitempty"`     // Name of theme. Can only be alphanumeric and underscores.
	Settings *ThemeSettings   `json:"settings,omitempty"`
}

type ThemesResourceOp struct {
	client *Client
}

var _ ThemesResource = &ThemesResourceOp{}

func (s *ThemesResourceOp) Get(ctx context.Context, ThemeId string) (*Theme, *Response, error) {
	return doGetById(ctx, s.client, ThemesBasePath, ThemeId, new(Theme))
}

func (s *ThemesResourceOp) Create(ctx context.Context, requestTheme *Theme) (*Theme, *Response, error) {
	return doCreate(ctx, s.client, ThemesBasePath, requestTheme, new(Theme))
}

func (s *ThemesResourceOp) Update(ctx context.Context, ThemeId string, requestTheme *Theme) (*Theme, *Response, error) {
	return doUpdate(ctx, s.client, ThemesBasePath, ThemeId, requestTheme, new(Theme))
}

func (s *ThemesResourceOp) Delete(ctx context.Context, ThemeId string) (*Response, error) {
	return doDelete(ctx, s.client, ThemesBasePath, ThemeId)
}
