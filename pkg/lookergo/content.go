package lookergo

import (
	"context"
	"net/url"
)

// ContentMetaGroupUserBasePath is the API call URI suffix for content metadata access
const ContentMetaGroupUserBasePath = "4.0/content_metadata_access"

// ContentMetaGroupUserResource is an interface that defines required methods
// for interacting with Looker content metadata access,e.g. a folder permission.
type ContentMetaGroupUserResource interface {
	// ListByID implements ContentMetaGroupUserResource. it is primarily meant to be used against Looker folders. It lists out all permissions a folder has.
	ListByID(context.Context, string, *ListOptions) ([]ContentMetaGroupUser, *Response, error)
}

// ContentMetaGroupUser is a struct that represents a Looker content metadata access,e.g. a folder permission. Reference : https://developers.looker.com/api/explorer/4.0/methods/Content/all_content_metadata_accesses?sdk=go
type ContentMetaGroupUser struct {
	Can               map[string]bool `json:"can,omitempty"`       // Operations the current user is able to perform on this object
	Id                string          `json:"id"`                  // Unique Id
	ContentMetadataId string          `json:"content_metadata_id"` // ID of associated content metadata
	PermissionType    string          `json:"permission_type"`     // Type of permission: "view" or "edit" Valid values are: "view", "edit"
	GroupId           string          `json:"group_id"`            // ID of associated group
	UserId            string          `json:"user_id"`             // ID of associated user
}

// Ensuring ContentMetaGroupUserResourceOp implements ContentMetaGroupUserResource
var _ ContentMetaGroupUserResource = &ContentMetaGroupUserResourceOp{}

// ContentMetaGroupUserResourceOp is a struct that implements ContentMetaGroupUserResource and has an internal http client for sending API calls
type ContentMetaGroupUserResourceOp struct {
	client *Client
}

// ListByID implements ContentMetaGroupUserResource. it is primarily meant to be used against Looker folders. It lists out all permissions a folder has.
func (s *ContentMetaGroupUserResourceOp) ListByID(ctx context.Context, contentMetadataId string, opt *ListOptions) ([]ContentMetaGroupUser, *Response, error) {
	if contentMetadataId == "" {
		return nil, nil, &ArgError{arg: "contentMetadataId", reason: "has to be non-empty"}
	}
	qs := url.Values{}
	qs.Add("content_metadata_id", contentMetadataId)
	return doListByX(ctx, s.client, ContentMetaGroupUserBasePath, opt, new([]ContentMetaGroupUser), qs)
}
