package lookergo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const roleBasePath = "4.0/roles"
const permissionSetBasePath = "4.0/permission_sets"

type Permission struct {
	Permission  string      `json:"permission"`
	Parent      interface{} `json:"parent"`
	Description string      `json:"description"`
}

type Role struct {
	Id   int    `json:"id,string,omitempty"`
	Name string `json:"name,omitempty"`
	/*	PermissionSet struct {
			BuiltIn     bool     `json:"built_in"`
			Id          string   `json:"id"`
			AllAccess   bool     `json:"all_access"`
			Name        string   `json:"name"`
			Permissions []string `json:"permissions"`
			Url         string   `json:"url"`
		} `json:"permission_set,omitempty"`
		ModelSet struct {
			BuiltIn   bool     `json:"built_in"`
			Id        string   `json:"id"`
			AllAccess bool     `json:"all_access"`
			Models    []string `json:"models"`
			Name      string   `json:"name"`
			Url       string   `json:"url"`
		} `json:"model_set,omitempty"`
		Url      string `json:"url,omitempty"`
		UsersUrl string `json:"users_url,omitempty"` */
}

// RolesResource is an interface for interfacing with the Role resource endpoints of the API.
// Ref: https://developers.looker.com/api/explorer/4.0/methods/Role
type RolesResource interface {
	PermissionSetsList(context.Context, *ListOptions) ([]PermissionSet, *Response, error)
	PermissionSetListById(context.Context, []int, *ListOptions) ([]PermissionSet, *Response, error)
	PermissionSetListByName(context.Context, string, *ListOptions) ([]PermissionSet, *Response, error)
	PermissionSetGet(context.Context, int) (*PermissionSet, *Response, error)
	List(context.Context, *ListOptions) ([]Role, *Response, error)
	Create(context.Context, *Role) (*Role, *Response, error)
	Update(context.Context, int, *Role) (*Role, *Response, error)
	Delete(context.Context, int) (*Response, error)
	RoleGroupsList(context.Context, int, *ListOptions) ([]Role, *Response, error)
	RoleGroupsSet(context.Context, int, []int) ([]Group, *Response, error)
	RoleUsersList(context.Context, int, *ListOptions) ([]Role, *Response, error)
	RoleUsersSet(context.Context, int, []int) ([]User, *Response, error)
}

// RolesResourceOp handles operations between Role related methods of the API.
type RolesResourceOp struct {
	client *Client
}

var _ RolesResource = &RolesResourceOp{}

type PermissionSet struct {
	BuiltIn     bool     `json:"built_in"`
	Id          int      `json:"id,string"`
	AllAccess   bool     `json:"all_access"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// PermissionSetsList -
func (s *RolesResourceOp) PermissionSetsList(ctx context.Context, opt *ListOptions) ([]PermissionSet, *Response, error) {
	return doList(ctx, s.client, permissionSetBasePath, opt, new([]PermissionSet))
}

func (s *RolesResourceOp) PermissionSetListById(ctx context.Context, ids []int, opt *ListOptions) ([]PermissionSet, *Response, error) {
	if len(ids) == 0 {
		return nil, nil, NewArgError("id", "specify one or more id(s)")
	}

	var idsQString string
	for i, id := range ids {
		idsQString += strconv.Itoa(id)
		if i+1 < len(ids) {
			idsQString += ","
		}
	}

	qs := url.Values{}
	// qs.Add("fields", "id,name,user_count,role_ids")
	qs.Add("id", idsQString)

	path := fmt.Sprintf("%s/search", permissionSetBasePath)

	return doListByX(ctx, s.client, path, opt, new([]PermissionSet), qs)
}

func (s *RolesResourceOp) PermissionSetListByName(ctx context.Context, name string, opt *ListOptions) ([]PermissionSet, *Response, error) {
	if name == "" {
		return nil, nil, NewArgError("name", "has to be non-empty")
	}

	qs := url.Values{}
	// qs.Add("fields", "id,name,user_count,role_ids")
	qs.Add("name", name)

	path := fmt.Sprintf("%s/search", permissionSetBasePath)

	return doListByX(ctx, s.client, path, opt, new([]PermissionSet), qs)
}

func (s *RolesResourceOp) PermissionSetGet(ctx context.Context, id int) (*PermissionSet, *Response, error) {
	return doGetById(ctx, s.client, permissionSetBasePath, id, new(PermissionSet))
}

// List -
func (s *RolesResourceOp) List(ctx context.Context, opt *ListOptions) ([]Role, *Response, error) {
	return doList(ctx, s.client, roleBasePath, opt, new([]Role))
}

// Create -
func (s *RolesResourceOp) Create(ctx context.Context, createReq *Role) (*Role, *Response, error) {
	return doCreate(ctx, s.client, userBasePath, createReq, new(Role))
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
func (s *RolesResourceOp) RoleGroupsList(ctx context.Context, id int, opt *ListOptions) ([]Role, *Response, error) {
	return doList(ctx, s.client, roleBasePath, opt, new([]Role), strconv.Itoa(id), "groups")
}

// RoleGroupsSet -
func (s *RolesResourceOp) RoleGroupsSet(ctx context.Context, id int, groupIds []int) ([]Group, *Response, error) {
	return doSet(ctx, s.client, userBasePath, groupIds, new([]Group), strconv.Itoa(id), "groups")
}

// RoleUsersList -
func (s *RolesResourceOp) RoleUsersList(ctx context.Context, id int, opt *ListOptions) ([]Role, *Response, error) {
	return doList(ctx, s.client, roleBasePath, opt, new([]Role), strconv.Itoa(id), "users")
}

// RoleUsersSet -
func (s *RolesResourceOp) RoleUsersSet(ctx context.Context, id int, userIds []int) ([]User, *Response, error) {
	return doSet(ctx, s.client, userBasePath, userIds, new([]User), strconv.Itoa(id), "users")
}

/*
// Get All Permissions
// GET {{endpoint}}/4.0/permissions
// Delete Permission Set
// DELETE {{endpoint}}/4.0/permission_sets/{{permission_set_id}}
// Update Permission Set
// PATCH  {{endpoint}}/4.0/permission_sets/{{permission_set_id}}
// Create Permission Set
// POST {{endpoint}}/4.0/permission_sets
*/
