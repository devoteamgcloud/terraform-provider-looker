package lookergo

import (
	"context"
)

const themeBasePath = "4.0/themes"

// ThemesResource is an interface for interfacing with the Theme resource endpoints of the API.
// Ref: https://developers.looker.com/api/explorer/4.0/methods/Theme
type ThemesResource interface {
	List(context.Context, *ListOptions) ([]Theme, *Response, error)
	Get(context.Context, int) (*Theme, *Response, error)
	Create(context.Context, *Theme) (*Theme, *Response, error)
	Update(context.Context, int, *Theme) (*Theme, *Response, error)
	Delete(context.Context, int) (*Response, error)
}

// ThemesResourceOp handles operations between Theme related methods of the API.
type ThemesResourceOp struct {
	client *Client
}

var _ ThemesResource = &ThemesResourceOp{}

type Settings struct {
	BackgroundColor     string `json:"background_color"`
	BaseFontSize        string `json:"base_font_size"`
	ColorCollectionId   string `json:"color_collection_id"`
	FontColor           string `json:"font_color"`
	FontFamily          string `json:"font_family"`
	FontSource          string `json:"font_source"`
	InfoButtonColor     string `json:"info_button_color"`
	PrimaryBtnColor     string `json:"primary_button_color"`
	ShowFiltersBar      bool   `json:"show_filters_bar"`
	ShowTitle           bool   `json:"show_title"`
	TextTileTextColor   string `json:"text_tile_text_color"`
	TileBackgroundColor string `json:"tile_background_color"`
	TextColor           string `json:"tile_text_color"`
	TitleColor          string `json:"title_color"`
	WarnBtnColor        string `json:"warn_button_color"`
	TileTitleAlignment  string `json:"tile_title_alignment"`
	TileShadow          bool   `json:"tile_shadow"`
}

type Theme struct {
	Can      *map[string]bool `json:"can,omitempty"`
	Id       int              `json:"id,string,omitempty"`
	Name     string           `json:"name,omitempty"`
	BeginAt  string           `json:"begin_at,omitempty"`
	EndAt    string           `json:"end_at,omitempty"`
	Settings Settings         `json:"settings,omitempty"`
}

// methods

// List all themes
func (s *ThemesResourceOp) List(ctx context.Context, opt *ListOptions) ([]Theme, *Response, error) {
	return doList(ctx, s.client, themeBasePath, opt, new([]Theme))
}

// Get a theme by ID.
func (s *ThemesResourceOp) Get(ctx context.Context, id int) (*Theme, *Response, error) {
	return doGetById(ctx, s.client, themeBasePath, id, new(Theme))
}

// Create a theme by ID.
func (s *ThemesResourceOp) Create(ctx context.Context, createReq *Theme) (*Theme, *Response, error) {
	return doCreate(ctx, s.client, themeBasePath, createReq, new(Theme))
}

// Update a theme by ID.
func (s *ThemesResourceOp) Update(ctx context.Context, id int, updateReq *Theme) (*Theme, *Response, error) {
	return doUpdate(ctx, s.client, themeBasePath, id, updateReq, new(Theme))
}

// Delete a theme by ID.
func (s *ThemesResourceOp) Delete(ctx context.Context, id int) (*Response, error) {
	return doDelete(ctx, s.client, themeBasePath, id)
}
