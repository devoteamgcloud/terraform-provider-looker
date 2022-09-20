package lookergo

import (
	"context"
)

const FoldersBasePath = "4.0/folders"

type FoldersResource interface {
	List(context.Context, *ListOptions) ([]Folder, *Response, error)
	Get(context.Context, string) (*Folder, *Response, error)
	//Get(context.Context,*ListOptions, string) ([]Folder, *Response, error)
	Create(context.Context, *Folder) (*Folder, *Response, error)
	Update(context.Context, string, *Folder) (*Folder, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

type FoldersResourceOp struct {
	client *Client
}

var _ FoldersResource = &FoldersResourceOp{}

type LookModel struct {
	Id    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type Dashboard struct {
	ContentFavoriteId  string     `json:"content_favorite_id,omitempty"`
	ContentMetadataId  string     `json:"content_metadata_id,omitempty"`
	Description        string     `json:"description,omitempty"`
	Hidden             bool       `json:"hidden,omitempty"`
	Id                 string     `json:"id,omitempty"`
	Model              *LookModel `json:"model,omitempty"`
	QueryTimezone      string     `json:"query_timezone,omitempty"`
	Readonly           bool       `json:"readonly,omitempty"`
	RefreshInterval    string     `json:"refresh_interval,omitempty"`
	RefreshIntervalToI int64      `json:"refresh_interval_to_i,omitempty"`
	Folder             *Folder    `json:"folder,omitempty"`
	Title              string     `json:"title,omitempty"`
	UserId             string     `json:"user_id,omitempty"`
	Slug               string     `json:"slug,omitempty"`
	PreferredViewer    string     `json:"preferred_viewer,omitempty"`
}

type Look struct {
	ContentMetadataId        string     `json:"content_metadata_id,omitempty"`
	Id                       string     `json:"id,omitempty"`
	Title                    string     `json:"title,omitempty"`
	UserId                   string     `json:"user_id,omitempty"`
	ContentFavoriteId        string     `json:"content_favorite_id,omitempty"`
	CreatedAt                string     `json:"created_at,omitempty"`
	Deleted                  bool       `json:"deleted,omitempty"`
	DeletedAt                string     `json:"deleted_at,omitempty"`
	Description              string     `json:"description,omitempty"`
	EmbededUrl               string     `json:"embeded_url,omitempty"`
	ExcelFileUrl             string     `json:"excel_file_url,omitempty"`
	FavoriteCount            int64      `json:"favorite_count,omitempty"`
	GoogleSpreadsheetFormula string     `json:"google_spreadsheet_formula,omitempty"`
	ImageEmbedUrl            string     `json:"image_embed_url,omitempty"`
	IsRunOnLoad              bool       `json:"is_run_on_load,omitempty"`
	LastAccessedAt           string     `json:"last_accessed_at,omitempty"`
	LastUpdaterId            string     `json:"last_updater_id,omitempty"`
	LastViewedAt             string     `json:"last_viewed_at,omitempty"`
	Model                    *LookModel `json:"model,omitempty"`
	Public                   bool       `json:"public,omitempty"`
	PublicSlug               string     `json:"public_slug,omitempty"`
	PublicUrl                string     `json:"public_url,omitempty"`
	QueryId                  string     `json:"query_id,omitempty"`
	ShortUrl                 string     `json:"short_url,omitempty"`
	Folder                   *Folder    `json:"folder,omitempty"`
	FolderId                 string     `json:"folder_id,omitempty"`
	UpdatedAt                string     `json:"updated_at,omitempty"`
	ViewCount                int64      `json:"view_count,omitempty"`
	Dashboards               *Dashboard `json:"dashboards,omitempty"`
}

type Folder struct {
	Name                 string `json:"name,omitempty"`
	ParentId             string `json:"parent_id,omitempty"`
	Id                   string `json:"id,omitempty"`
	ContentMetadataId    string `json:"content_metadata_id,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	CreatorId            string `json:"creator_id,omitempty"`
	ChildCount           int64  `json:"child_count,omitempty"`
	ExternalId           string `json:"external_id,omitempty"`
	IsEmbed              bool   `json:"is_embed,omitempty"`
	IsEmbedSharedRoot    bool   `json:"is_embed_shared_root,omitempty"`
	IsEmbedUsersRoot     bool   `json:"is_embed_users_root,omitempty"`
	IsPersonal           bool   `json:"is_personal,omitempty"`
	IsPersonalDescendant bool   `json:"is_personal_descendant,omitempty"`
	IsSharedRoot         bool   `json:"is_shared_root,omitempty"`
	IsUsersRoot          bool   `json:"is_users_root,omitempty"`
	//Dashboards           *Dashboard `json:"dashboards,omitempty"`
}

/*
List
ListById
Get -> recurse children.
Create
Update
Delete
*/

func (s *FoldersResourceOp) List(ctx context.Context, opt *ListOptions) ([]Folder, *Response, error) {
	return doList(ctx, s.client, FoldersBasePath, opt, new([]Folder))
}

func (s *FoldersResourceOp) Get(ctx context.Context, FolderId string) (*Folder, *Response, error) {
	return doGetById(ctx, s.client, FoldersBasePath, FolderId, new(Folder))
}

// func (s *FoldersResourceOp) Get(ctx context.Context, opt *ListOptions, FolderId string) ([]Folder, *Response, error) {
// 	if FolderId == "" {
// 		return nil, nil, NewArgError("name", "has to be non-empty")
// 	}
// 	qs := url.Values{}

// 	path := fmt.Sprintf("%s/%s/children", FoldersBasePath, FolderId)
// 	return doListByX(ctx, s.client, path, opt, new([]Folder), qs)
// }

func (s *FoldersResourceOp) Create(ctx context.Context, requestFolder *Folder) (*Folder, *Response, error) {
	return doCreate(ctx, s.client, FoldersBasePath, requestFolder, new(Folder))
}

func (s *FoldersResourceOp) Update(ctx context.Context, FolderId string, requestFolder *Folder) (*Folder, *Response, error) {
	return doUpdate(ctx, s.client, FoldersBasePath, FolderId, requestFolder, new(Folder))
}

func (s *FoldersResourceOp) Delete(ctx context.Context, FolderId string) (*Response, error) {
	return doDelete(ctx, s.client, FoldersBasePath, FolderId)
}
