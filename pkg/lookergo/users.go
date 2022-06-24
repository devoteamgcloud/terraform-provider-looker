package lookergo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const userBasePath = "4.0/users"

// CredentialEmail -
type CredentialEmail struct {
	Email       string `json:"email"`
	ForcedReset bool   `json:"forced_password_reset_at_next_login"`
	IsDisabled  bool   `json:"is_disabled,omitempty"`
}

// SliceStringInts is an intermediate/shim type to make the json result return stringified integers,
// while being able to use correct integer types internally.
type SliceStringInts []int

// User defines a user in the database
// Ref: https://github.com/looker-open-source/sdk-codegen/blob/main/go/sdk/v4/models.go#L3508
type User struct {
	Id              int              `json:"id,string,omitempty"`
	FirstName       string           `json:"first_name,omitempty"`
	LastName        string           `json:"last_name,omitempty"`
	CredentialEmail *CredentialEmail `json:"credentials_email,omitempty" `
	DisplayName     string           `json:"display_name,omitempty"`
	RoleIds         SliceStringInts  `json:"role_ids,omitempty"`
	Email           string           `json:"email,omitempty"`
}

// JSON parsing

func (c *SliceStringInts) MarshalJSON() ([]byte, error) {
	var r []string
	for _, item := range *c {
		r = append(r, fmt.Sprintf("%v", item))
	}
	return json.Marshal(r)
}

func (c *SliceStringInts) UnmarshalJSON(b []byte) error {
	var n SliceStringInts

	var nums []interface{}

	err := json.Unmarshal(b, &nums)
	if err != nil {
		return err
	}

	for _, item := range nums {
		switch value := item.(type) {
		case int:
			n = append(n, item.(int))
		case float64:
			n = append(n, int(value))
		case string:
			num, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			n = append(n, num)
		}
	}

	*c = n
	return nil
}

func (c *SliceStringInts) ToSliceOfStrings() []string {
	var ret []string

	for _, item := range *c {
		ret = append(ret, strconv.Itoa(item))
	}

	return ret
}

func FromSliceOfStrings(s []string) SliceStringInts {
	var ret SliceStringInts

	for _, item := range s {
		itemInt, _ := strconv.Atoi(item)
		ret = append(ret, itemInt)
	}

	return ret
}

// UsersResource is an interface for interfacing with the User resource endpoints of the API.
// Ref: https://developers.looker.com/api/explorer/4.0/methods/User
type UsersResource interface {
	List(context.Context, *ListOptions) ([]User, *Response, error)
	ListById(context.Context, []int, *ListOptions) ([]User, *Response, error)
	ListByEmail(context.Context, string, *ListOptions) ([]User, *Response, error)
	Get(context.Context, int) (*User, *Response, error)
	Create(context.Context, *User) (*User, *Response, error)
	Update(context.Context, int, *User) (*User, *Response, error)
	Delete(context.Context, int) (*Response, error)
	CreateEmail(context.Context, int, *CredentialEmail) (*CredentialEmail, *Response, error)
	GetEmail(context.Context, int) (*CredentialEmail, *Response, error)
	UpdateEmail(context.Context, int, *CredentialEmail) (*CredentialEmail, *Response, error)
	DeleteEmail(context.Context, int) (*Response, error)
	CreatePasswordReset(context.Context, int) (*CredentialEmail, *Response, error)
	SendPasswordReset(context.Context, int) (*CredentialEmail, *Response, error)
	GetRoles(context.Context, int) ([]Role, *Response, error)
	SetRoles(context.Context, int, []string) ([]Role, *Response, error)
}

// UsersResourceOp handles operations between User related methods of the API.
type UsersResourceOp struct {
	client *Client
}

// List all users
func (s *UsersResourceOp) List(ctx context.Context, opt *ListOptions) ([]User, *Response, error) {
	return doList(ctx, s.client, userBasePath, opt, new([]User))
}

func (s *UsersResourceOp) ListById(ctx context.Context, ids []int, opt *ListOptions) ([]User, *Response, error) {
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

	path := fmt.Sprintf("%s/search", userBasePath)

	return doListByX(ctx, s.client, path, opt, new([]User), qs)
}
func (s *UsersResourceOp) ListByEmail(ctx context.Context, email string, opt *ListOptions) ([]User, *Response, error) {
	if email == "" {
		return nil, nil, NewArgError("name", "has to be non-empty")
	}

	qs := url.Values{}
	// qs.Add("fields", "id,name,user_count,role_ids")
	qs.Add("email", email)

	path := fmt.Sprintf("%s/search", userBasePath)

	return doListByX(ctx, s.client, path, opt, new([]User), qs)
}

// Get -
func (s *UsersResourceOp) Get(ctx context.Context, id int) (*User, *Response, error) {
	return doGetById(ctx, s.client, userBasePath, id, new(User))
}

// Create -
func (s *UsersResourceOp) Create(ctx context.Context, createReq *User) (*User, *Response, error) {
	return doCreate(ctx, s.client, userBasePath, createReq, new(User))
}

// Update -
func (s *UsersResourceOp) Update(ctx context.Context, id int, updateReq *User) (*User, *Response, error) {
	return doUpdate(ctx, s.client, userBasePath, id, updateReq, new(User))
}

// Delete -
func (s *UsersResourceOp) Delete(ctx context.Context, id int) (*Response, error) {
	return doDelete(ctx, s.client, userBasePath, id)
}

// CreateEmail -
func (s *UsersResourceOp) CreateEmail(ctx context.Context, id int, createReq *CredentialEmail) (*CredentialEmail, *Response, error) {
	return doCreate(ctx, s.client, userBasePath, createReq, new(CredentialEmail), strconv.Itoa(id), "credentials_email")
}

// GetEmail -
func (s *UsersResourceOp) GetEmail(ctx context.Context, id int) (*CredentialEmail, *Response, error) {
	return doGet(ctx, s.client, userBasePath, new(CredentialEmail), strconv.Itoa(id), "credentials_email")
}

// UpdateEmail -
func (s *UsersResourceOp) UpdateEmail(ctx context.Context, id int, updateReq *CredentialEmail) (*CredentialEmail, *Response, error) {
	return doUpdate(ctx, s.client, userBasePath, id, updateReq, new(CredentialEmail), "credentials_email")
}

// DeleteEmail -
func (s *UsersResourceOp) DeleteEmail(ctx context.Context, id int) (*Response, error) {
	return doDelete(ctx, s.client, userBasePath, id, "credentials_email")
}

// CreatePasswordReset -
func (s *UsersResourceOp) CreatePasswordReset(ctx context.Context, id int) (*CredentialEmail, *Response, error) {
	return doEmptyPost(ctx, s.client, userBasePath, new(CredentialEmail),
		strconv.Itoa(id), "credentials_email", "password_reset")
}

// SendPasswordReset -
func (s *UsersResourceOp) SendPasswordReset(ctx context.Context, id int) (*CredentialEmail, *Response, error) {
	return doEmptyPost(ctx, s.client, userBasePath, new(CredentialEmail),
		strconv.Itoa(id), "credentials_email", "send_password_reset")
}

// GetRoles -
func (s *UsersResourceOp) GetRoles(ctx context.Context, id int) ([]Role, *Response, error) {
	return doList(ctx, s.client, userBasePath, nil, new([]Role), strconv.Itoa(id), "roles")
}

// SetRoles -
func (s *UsersResourceOp) SetRoles(ctx context.Context, id int, roleIds []string) ([]Role, *Response, error) {
	return doSet(ctx, s.client, userBasePath, roleIds, new([]Role), strconv.Itoa(id), "roles")
}
