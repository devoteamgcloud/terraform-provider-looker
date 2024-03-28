package lookergo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const groupBasePath = "4.0/groups"

// GroupsResource is an interface for interfacing with the Group resource endpoints of the API.
// Ref: https://developers.looker.com/api/explorer/4.0/methods/Group
//
//	https://blob.b-cdn.net/looker_api4.0_ref-1652781627.html#operation/group
type GroupsResource interface {
	List(context.Context, *ListOptions) ([]Group, *Response, error)
	ListByName(context.Context, string, *ListOptions) ([]Group, *Response, error)
	ListById(context.Context, []int, *ListOptions) ([]Group, *Response, error)
	Get(context.Context, int) (*Group, *Response, error)
	Create(context.Context, *Group) (*Group, *Response, error)
	Update(context.Context, int, *Group) (*Group, *Response, error)
	Delete(context.Context, int) (*Response, error)
	ListMemberGroups(context.Context, int, *ListOptions) ([]Group, *Response, error)
	AddMemberGroup(context.Context, int, int) (*Group, *Response, error)
	RemoveMemberGroup(context.Context, int, int) (*Response, error)
	ListMemberUsers(context.Context, int, *ListOptions) ([]User, *Response, error)
	AddMemberUser(context.Context, int, int) (*User, *Response, error)
	RemoveMemberUser(context.Context, int, int) (*Response, error)
}

// GroupsResourceOp handles operations between Group related methods of the API.
type GroupsResourceOp struct {
	client *Client
}

var _ GroupsResource = &GroupsResourceOp{}

// Group -
// Ref: https://developers.looker.com/api/explorer/4.0/types/Group/Group
type Group struct {
	CanAddToContentMetadata bool            `json:"can_add_to_content_metadata,omitempty"`
	ExternalGroupId         bool            `json:"external_group_id,omitempty"`
	Id                      int             `json:"id,string,omitempty"`
	Name                    string          `json:"name"`
	UserCount               int             `json:"user_count,omitempty"` // ! not stringified
	ExternallyManaged       bool            `json:"externally_managed,omitempty"`
	IncludeByDefault        bool            `json:"include_by_default,omitempty"`
	ContainsCurrentUser     bool            `json:"contains_current_user,omitempty"`
	ParentGroupIds          SliceStringInts `json:"parent_group_ids,omitempty"`
	RoleIds                 SliceStringInts `json:"role_ids,omitempty"`
}

// List all groups
func (s *GroupsResourceOp) List(ctx context.Context, opt *ListOptions) ([]Group, *Response, error) {
	return doList(ctx, s.client, groupBasePath, opt, new([]Group))
}

func (s *GroupsResourceOp) ListByName(ctx context.Context, name string, opt *ListOptions) ([]Group, *Response, error) {
	if name == "" {
		return nil, nil, NewArgError("name", "has to be non-empty")
	}

	qs := url.Values{}
	qs.Add("fields", "id,name,user_count,role_ids")
	qs.Add("name", name)

	path := fmt.Sprintf("%s/search/with_hierarchy", groupBasePath)

	return doListByX(ctx, s.client, path, opt, new([]Group), qs)
}

func (s *GroupsResourceOp) ListById(ctx context.Context, ids []int, opt *ListOptions) ([]Group, *Response, error) {
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
	qs.Add("fields", "id,name,user_count,role_ids,parent_groups")
	qs.Add("id", idsQString)

	path := fmt.Sprintf("%s/search/with_hierarchy", groupBasePath)

	return doListByX(ctx, s.client, path, opt, new([]Group), qs)
}

// Get a group by ID.
func (s *GroupsResourceOp) Get(ctx context.Context, id int) (*Group, *Response, error) {
	return doGetById(ctx, s.client, groupBasePath, id, new(Group))
}

// Create a group by ID.
func (s *GroupsResourceOp) Create(ctx context.Context, createReq *Group) (*Group, *Response, error) {
	return doCreate(ctx, s.client, groupBasePath, createReq, new(Group))
}

// Update a group by ID.
func (s *GroupsResourceOp) Update(ctx context.Context, id int, updateReq *Group) (*Group, *Response, error) {
	return doUpdate(ctx, s.client, groupBasePath, id, updateReq, new(Group))
}

// Delete a group by ID.
func (s *GroupsResourceOp) Delete(ctx context.Context, id int) (*Response, error) {
	return doDelete(ctx, s.client, groupBasePath, id)
}

// ListMemberGroups gets all member groups inside a group.
func (s *GroupsResourceOp) ListMemberGroups(ctx context.Context, id int, opt *ListOptions) ([]Group, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/groups", groupBasePath, id)

	return doList(ctx, s.client, path, opt, new([]Group))
}

type NewGroupMemberGroup struct {
	GroupID int `json:"group_id"`
}

// AddMemberGroup -
func (s *GroupsResourceOp) AddMemberGroup(ctx context.Context, parentID int, memberID int) (*Group, *Response, error) {
	if parentID < 1 || memberID < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/groups", groupBasePath, parentID)

	return doAddMember(ctx, s.client, path, new(Group), NewGroupMemberGroup{GroupID: memberID})
}

// RemoveMemberGroup -
func (s *GroupsResourceOp) RemoveMemberGroup(ctx context.Context, parentID int, memberID int) (*Response, error) {
	if parentID < 1 || memberID < 1 {
		return nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/groups", groupBasePath, parentID)

	return doDelete(ctx, s.client, path, memberID)
}

// ListMemberUsers gets all member groups inside a group.
func (s *GroupsResourceOp) ListMemberUsers(ctx context.Context, id int, opt *ListOptions) ([]User, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/users", groupBasePath, id)

	return doList(ctx, s.client, path, opt, new([]User))
}

type NewGroupMemberUser struct {
	UserID int `json:"user_id"`
}

// AddMemberUser -
func (s *GroupsResourceOp) AddMemberUser(ctx context.Context, parentID int, memberID int) (*User, *Response, error) {
	if parentID < 1 || memberID < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/users", groupBasePath, parentID)

	return doAddMember(ctx, s.client, path, new(User), NewGroupMemberUser{UserID: memberID})
}

// RemoveMemberUser -
func (s *GroupsResourceOp) RemoveMemberUser(ctx context.Context, parentID int, memberID int) (*Response, error) {
	if parentID < 1 || memberID < 1 {
		return nil, NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/users", groupBasePath, parentID)

	return doDelete(ctx, s.client, path, memberID)
}
