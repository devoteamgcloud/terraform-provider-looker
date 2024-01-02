package lookergo

import (
	"context"
)

const ContentMetaGroupUserBasePath = "4.0/content_metadata_access"

type ContentMetaGroupUserResource interface {
	ListByID(context.Context, string, *ListOptions) ([]ContentMetaGroupUser, *Response, error)
}
type ContentMetaGroupUser struct {
	Can               map[string]bool `json:"can,omitempty"`
	Id                string          `json:"id"`
	ContentMetadataId string          `json:"content_metadata_id"`
	PermissionType    string          `json:"permission_type"`
	GroupId           string          `json:"group_id"`
	UserId            string          `json:"user_id"`
}
