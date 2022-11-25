package lookergo

import (
	"context"
	"fmt"
	"net/url"
)

const permissionSetBasePath = "4.0/permission_sets"

type PermissionSetResource interface {
	List(context.Context, *ListOptions) ([]PermissionSet, *Response, error)
	Get(ctx context.Context, PermissionSetId string) (*PermissionSet, *Response, error)
	GetByName(ctx context.Context, PermissionSetName string, opt *ListOptions) ([]PermissionSet, *Response, error)
	Create(ctx context.Context, PermissionSet *PermissionSet) (*PermissionSet, *Response, error)
	Update(ctx context.Context, PermissionSetId string, PermissionSet *PermissionSet) (*PermissionSet, *Response, error)
	Delete(ctx context.Context, PermissionSetId string) (*Response, error)
}

type PermissionSetResourceOp struct {
	client *Client
}

var _ PermissionSetResource = &PermissionSetResourceOp{}

type PermissionSet struct {
	Can         map[string]bool `json:"can,omitempty"` // Operations the current user is able to perform on this object
	AllAccess   bool            `json:"all_access,omitempty"`
	BuiltIn     bool            `json:"built_in,omitempty"`
	Id          string          `json:"id,omitempty"`   // Unique Id
	Name        string          `json:"name,omitempty"` // Name of PermissionSet
	Permissions []string        `json:"permissions,omitempty"`
	Url         string          `json:"url,omitempty"` // Link to get this item
}

// List -
func (s *PermissionSetResourceOp) List(ctx context.Context, opt *ListOptions) ([]PermissionSet, *Response, error) {
	return doList(ctx, s.client, permissionSetBasePath, opt, new([]PermissionSet))
}

func (s *PermissionSetResourceOp) Get(ctx context.Context, PermissionSetId string) (*PermissionSet, *Response, error) {
	return doGetById(ctx, s.client, permissionSetBasePath, PermissionSetId, new(PermissionSet))
}

func (s *PermissionSetResourceOp) GetByName(ctx context.Context, PermissionSetName string, opt *ListOptions) ([]PermissionSet, *Response, error) {
	if PermissionSetName == "" {
		return nil, nil, NewArgError("name", "has to be non-empty")
	}
	qs := url.Values{}
	qs.Add("fields", "id,name,permissions")
	qs.Add("name", PermissionSetName)
	path := fmt.Sprintf("%s/search", permissionSetBasePath)
	return doListByX(ctx, s.client, path, opt, new([]PermissionSet), qs)
}

func (s *PermissionSetResourceOp) Create(ctx context.Context, permissionSet *PermissionSet) (*PermissionSet, *Response, error) {
	return doCreate(ctx, s.client, permissionSetBasePath, permissionSet, new(PermissionSet))
}

func (s *PermissionSetResourceOp) Update(ctx context.Context, PermissionSetId string, permissionSet *PermissionSet) (*PermissionSet, *Response, error) {
	return doUpdate(ctx, s.client, permissionSetBasePath, PermissionSetId, permissionSet, new(PermissionSet))
}

func (s *PermissionSetResourceOp) Delete(ctx context.Context, PermissionSetId string) (*Response, error) {
	return doDelete(ctx, s.client, permissionSetBasePath, PermissionSetId)
}
