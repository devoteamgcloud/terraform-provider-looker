package lookergo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const roleBasePath = "4.0/roles"

type Permission struct {
	Permission  string      `json:"permission"`
	Parent      interface{} `json:"parent"`
	Description string      `json:"description"`
}

type Role struct {
	Id              int           `json:"id,string,omitempty"`
	Name            string        `json:"name,omitempty"`
	PermissionSet   PermissionSet `json:"permission_set,omitempty"`
	PermissionSetID string        `json:"permission_set_id,omitempty"`
	ModelSet        ModelSet      `json:"model_set,omitempty"`
	ModelSetID      string        `json:"model_set_id,omitempty"`
}

type PermissionSetSlice []PermissionSet

// RolesResource is an interface for interfacing with the Role resource endpoints of the API.
// Ref: https://developers.looker.com/api/explorer/4.0/methods/Role
type RolesResource interface {
	List(context.Context, *ListOptions) ([]Role, *Response, error)
	ListByName(ctx context.Context, name string, opt *ListOptions) ([]Role, *Response, error)
	Get(context.Context, int) (*Role, *Response, error)
	Create(context.Context, *Role) (*Role, *Response, error)
	Update(context.Context, int, *Role) (*Role, *Response, error)
	Delete(context.Context, int) (*Response, error)
	RoleGroupsList(context.Context, int, *ListOptions) ([]Group, *Response, error)
	RoleGroupsSet(context.Context, int, []string) ([]Group, *Response, error)
	RoleUsersList(context.Context, int, *ListOptions) ([]User, *Response, error)
	RoleUsersSet(context.Context, int, []string) ([]User, *Response, error)
}

// RolesResourceOp handles operations between Role related methods of the API.
type RolesResourceOp struct {
	client *Client
}

var _ RolesResource = &RolesResourceOp{}

// List -
func (s *RolesResourceOp) List(ctx context.Context, opt *ListOptions) ([]Role, *Response, error) {
	return doList(ctx, s.client, roleBasePath, opt, new([]Role))
}

// ListByName -
func (s *RolesResourceOp) ListByName(ctx context.Context, name string, opt *ListOptions) ([]Role, *Response, error) {
	if name == "" {
		return nil, nil, NewArgError("name", "has to be non-empty")
	}
	qs := url.Values{}
	qs.Add("name", name)
	path := fmt.Sprintf("%s/search", roleBasePath)
	return doListByX(ctx, s.client, path, opt, new([]Role), qs)
}

// Get -
func (s *RolesResourceOp) Get(ctx context.Context, id int) (*Role, *Response, error) {
	return doGetById(ctx, s.client, roleBasePath, id, new(Role))
}

// Create -
func (s *RolesResourceOp) Create(ctx context.Context, createReq *Role) (*Role, *Response, error) {
	return doCreate(ctx, s.client, roleBasePath, createReq, new(Role))
}

// Update -
func (s *RolesResourceOp) Update(ctx context.Context, id int, updateReq *Role) (*Role, *Response, error) {
	return doUpdate(ctx, s.client, roleBasePath, id, updateReq, new(Role))
}

// Delete -
func (s *RolesResourceOp) Delete(ctx context.Context, id int) (*Response, error) {
	return doDelete(ctx, s.client, roleBasePath, id)
}

// RoleGroupsList -
func (s *RolesResourceOp) RoleGroupsList(ctx context.Context, id int, opt *ListOptions) ([]Group, *Response, error) {
	return doList(ctx, s.client, roleBasePath, opt, new([]Group), strconv.Itoa(id), "groups")
}

// RoleGroupsSet -
func (s *RolesResourceOp) RoleGroupsSet(ctx context.Context, id int, groupIds []string) ([]Group, *Response, error) {
	return doSet(ctx, s.client, roleBasePath, groupIds, new([]Group), strconv.Itoa(id), "groups")
}

// RoleUsersList -
func (s *RolesResourceOp) RoleUsersList(ctx context.Context, id int, opt *ListOptions) ([]User, *Response, error) {
	return doList(ctx, s.client, roleBasePath, opt, new([]User), strconv.Itoa(id), "users")
}

// RoleUsersSet -
func (s *RolesResourceOp) RoleUsersSet(ctx context.Context, id int, userIds []string) ([]User, *Response, error) {
	return doSet(ctx, s.client, roleBasePath, userIds, new([]User), strconv.Itoa(id), "users")
}
