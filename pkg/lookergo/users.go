package lookergo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const userBasePath = "4.0/users"

// CredentialsEmail -
type CredentialsEmail struct {
	Can                            map[string]bool `json:"can,omitempty"`                                 // Operations the current user is able to perform on this object
	CreatedAt                      string          `json:"created_at,omitempty"`                          // Timestamp for the creation of this credential
	Email                          string          `json:"email"`                                         // EMail address used for user login
	ForcedPasswordResetAtNextLogin bool            `json:"forced_password_reset_at_next_login,omitempty"` // Force the user to change their password upon their next login
	IsDisabled                     bool            `json:"is_disabled,omitempty"`                         // Has this credential been disabled?
	LoggedInAt                     string          `json:"logged_in_at,omitempty"`                        // Timestamp for most recent login using credential
	PasswordResetUrl               string          `json:"password_reset_url,omitempty"`                  // Url with one-time use secret token that the user can use to reset password
	Type                           string          `json:"type,omitempty"`                                // Short name for the type of this kind of credential
	Url                            string          `json:"url,omitempty"`                                 // Link to get this item
	UserUrl                        string          `json:"user_url,omitempty"`                            // Link to get this user
}

type CredentialsEmbed struct {
	Can             map[string]bool `json:"can,omitempty"`               // Operations the current user is able to perform on this object
	CreatedAt       string          `json:"created_at,omitempty"`        // Timestamp for the creation of this credential
	ExternalGroupId string          `json:"external_group_id,omitempty"` // Embedder's id for a group to which this user was added during the most recent login
	ExternalUserId  string          `json:"external_user_id,omitempty"`  // Embedder's unique id for the user
	Id              string          `json:"id,omitempty"`                // Unique Id
	IsDisabled      bool            `json:"is_disabled,omitempty"`       // Has this credential been disabled?
	LoggedInAt      string          `json:"logged_in_at,omitempty"`      // Timestamp for most recent login using credential
	Type            string          `json:"type,omitempty"`              // Short name for the type of this kind of credential
	Url             string          `json:"url,omitempty"`               // Link to get this item
}

type CredentialsSaml struct {
	Can        map[string]bool `json:"can,omitempty"`          // Operations the current user is able to perform on this object
	CreatedAt  string          `json:"created_at,omitempty"`   // Timestamp for the creation of this credential
	Email      string          `json:"email,omitempty"`        // EMail address
	IsDisabled bool            `json:"is_disabled,omitempty"`  // Has this credential been disabled?
	LoggedInAt string          `json:"logged_in_at,omitempty"` // Timestamp for most recent login using credential
	SamlUserId string          `json:"saml_user_id,omitempty"` // Saml IdP's Unique ID for this user
	Type       string          `json:"type,omitempty"`         // Short name for the type of this kind of credential
	Url        string          `json:"url,omitempty"`          // Link to get this item
}

// SliceStringInts is an intermediate/shim type to make the json result return stringified integers,
// while being able to use correct integer types internally.
type SliceStringInts []int

// User defines a user in the database
// Ref: https://github.com/looker-open-source/sdk-codegen/blob/main/go/sdk/v4/models.go#L3508
type User struct {
	Can                    *map[string]bool    `json:"can,omitempty"`                       // Operations the current user is able to perform on this object
	AvatarUrl              string              `json:"avatar_url,omitempty"`                // URL for the avatar image (may be generic)
	AvatarUrlWithoutSizing string              `json:"avatar_url_without_sizing,omitempty"` // URL for the avatar image (may be generic), does not specify size
	CredentialsEmail       *CredentialsEmail   `json:"credentials_email,omitempty"`
	CredentialsEmbed       *[]CredentialsEmbed `json:"credentials_embed,omitempty"` // Embed credentials
	//CredentialsApi3            *[]CredentialsApi3       `json:"credentials_api3,omitempty"`                // API 3 credentials
	//CredentialsGoogle          *CredentialsGoogle       `json:"credentials_google,omitempty"`
	//CredentialsLdap            *CredentialsLDAP         `json:"credentials_ldap,omitempty"`
	//CredentialsLookerOpenid    *CredentialsLookerOpenid `json:"credentials_looker_openid,omitempty"`
	//CredentialsOidc            *CredentialsOIDC         `json:"credentials_oidc,omitempty"`
	//CredentialsTotp            *CredentialsTotp         `json:"credentials_totp,omitempty"`
	CredentialsSaml            *CredentialsSaml       `json:"credentials_saml,omitempty"`
	DisplayName                string                 `json:"display_name,omitempty"`                   // Full name for display (available only if both first_name and last_name are set)
	Email                      string                 `json:"email,omitempty"`                          // EMail address
	EmbedGroupSpaceId          string                 `json:"embed_group_space_id,omitempty"`           // (DEPRECATED) (Embed only) ID of user's group space based on the external_group_id optionally specified during embed user login
	FirstName                  string                 `json:"first_name,omitempty"`                     // First name
	GroupIds                   []string               `json:"group_ids,omitempty"`                      // Array of ids of the groups for this user
	HomeFolderId               string                 `json:"home_folder_id,omitempty"`                 // ID string for user's home folder
	Id                         string                 `json:"id,omitempty"`                             // Unique Id
	IsDisabled                 bool                   `json:"is_disabled,omitempty"`                    // Account has been disabled
	LastName                   string                 `json:"last_name,omitempty"`                      // Last name
	Locale                     string                 `json:"locale,omitempty"`                         // User's preferred locale. User locale takes precedence over Looker's system-wide default locale. Locale determines language of display strings and date and numeric formatting in API responses. Locale string must be a 2 letter language code or a combination of language code and region code: 'en' or 'en-US', for example.
	LookerVersions             []string               `json:"looker_versions,omitempty"`                // Array of strings representing the Looker versions that this user has used (this only goes back as far as '3.54.0')
	ModelsDirValidated         bool                   `json:"models_dir_validated,omitempty"`           // User's dev workspace has been checked for presence of applicable production projects
	PersonalFolderId           string                 `json:"personal_folder_id,omitempty"`             // ID of user's personal folder
	PresumedLookerEmployee     bool                   `json:"presumed_looker_employee,omitempty"`       // User is identified as an employee of Looker
	RoleIds                    SliceStringInts        `json:"role_ids,omitempty"`                       // Array of ids of the roles for this user
	UiState                    map[string]interface{} `json:"ui_state,omitempty"`                       // Per user dictionary of undocumented state information owned by the Looker UI.
	VerifiedLookerEmployee     bool                   `json:"verified_looker_employee,omitempty"`       // User is identified as an employee of Looker who has been verified via Looker corporate authentication
	RolesExternallyManaged     bool                   `json:"roles_externally_managed,omitempty"`       // User's roles are managed by an external directory like SAML or LDAP and can not be changed directly.
	AllowDirectRoles           bool                   `json:"allow_direct_roles,omitempty"`             // User can be directly assigned a role.
	AllowNormalGroupMembership bool                   `json:"allow_normal_group_membership,omitempty"`  // User can be a direct member of a normal Looker group.
	AllowRolesFromNormalGroups bool                   `json:"allow_roles_from_normal_groups,omitempty"` // User can inherit roles from a normal Looker group.
	EmbedGroupFolderId         string                 `json:"embed_group_folder_id,omitempty"`          // (Embed only) ID of user's group folder based on the external_group_id optionally specified during embed user login
	Url                        string                 `json:"url,omitempty"`                            // Link to get this item
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
	ListById(context.Context, []string, *ListOptions) ([]User, *Response, error)
	ListByEmail(context.Context, string, *ListOptions) ([]User, *Response, error)
	Get(context.Context, string) (*User, *Response, error)
	Create(context.Context, *User) (*User, *Response, error)
	Update(context.Context, string, *User) (*User, *Response, error)
	Delete(context.Context, string) (*Response, error)
	CreateEmail(context.Context, string, *CredentialsEmail) (*CredentialsEmail, *Response, error)
	GetEmail(context.Context, string) (*CredentialsEmail, *Response, error)
	UpdateEmail(context.Context, string, *CredentialsEmail) (*CredentialsEmail, *Response, error)
	DeleteEmail(context.Context, string) (*Response, error)
	CreatePasswordReset(context.Context, string) (*CredentialsEmail, *Response, error)
	SendPasswordReset(context.Context, string) (*CredentialsEmail, *Response, error)
	GetRoles(context.Context, string) ([]Role, *Response, error)
	SetRoles(context.Context, string, []string) ([]Role, *Response, error)
}

// UsersResourceOp handles operations between User related methods of the API.
type UsersResourceOp struct {
	client *Client
}

// List all users
func (s *UsersResourceOp) List(ctx context.Context, opt *ListOptions) ([]User, *Response, error) {
	return doList(ctx, s.client, userBasePath, opt, new([]User))
}

func (s *UsersResourceOp) ListById(ctx context.Context, ids []string, opt *ListOptions) ([]User, *Response, error) {
	if len(ids) == 0 {
		return nil, nil, NewArgError("id", "specify one or more id(s)")
	}

	var idsQString string
	for i, id := range ids {
		idsQString += id
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
		return nil, nil, NewArgError("email", "has to be non-empty")
	}

	qs := url.Values{}
	// qs.Add("fields", "id,name,user_count,role_ids")
	qs.Add("email", email)

	path := fmt.Sprintf("%s/search", userBasePath)

	return doListByX(ctx, s.client, path, opt, new([]User), qs)
}

// Get -
func (s *UsersResourceOp) Get(ctx context.Context, id string) (*User, *Response, error) {
	return doGetById(ctx, s.client, userBasePath, id, new(User))
}

// Create -
func (s *UsersResourceOp) Create(ctx context.Context, createReq *User) (*User, *Response, error) {
	return doCreate(ctx, s.client, userBasePath, createReq, new(User))
}

// Update -
func (s *UsersResourceOp) Update(ctx context.Context, id string, updateReq *User) (*User, *Response, error) {
	return doUpdate(ctx, s.client, userBasePath, id, updateReq, new(User))
}

// Delete -
func (s *UsersResourceOp) Delete(ctx context.Context, id string) (*Response, error) {
	return doDelete(ctx, s.client, userBasePath, id)
}

// CreateEmail -
func (s *UsersResourceOp) CreateEmail(ctx context.Context, id string, createReq *CredentialsEmail) (*CredentialsEmail, *Response, error) {
	return doCreate(ctx, s.client, userBasePath, createReq, new(CredentialsEmail), id, "credentials_email")
}

// GetEmail -
func (s *UsersResourceOp) GetEmail(ctx context.Context, id string) (*CredentialsEmail, *Response, error) {
	return doGet(ctx, s.client, userBasePath, new(CredentialsEmail), id, "credentials_email")
}

// UpdateEmail -
func (s *UsersResourceOp) UpdateEmail(ctx context.Context, id string, updateReq *CredentialsEmail) (*CredentialsEmail, *Response, error) {
	return doUpdate(ctx, s.client, userBasePath, id, updateReq, new(CredentialsEmail), "credentials_email")
}

// DeleteEmail -
func (s *UsersResourceOp) DeleteEmail(ctx context.Context, id string) (*Response, error) {
	return doDelete(ctx, s.client, userBasePath, id, "credentials_email")
}

// CreatePasswordReset -
func (s *UsersResourceOp) CreatePasswordReset(ctx context.Context, id string) (*CredentialsEmail, *Response, error) {
	return doEmptyPost(ctx, s.client, userBasePath, new(CredentialsEmail),
		id, "credentials_email", "password_reset")
}

// SendPasswordReset -
func (s *UsersResourceOp) SendPasswordReset(ctx context.Context, id string) (*CredentialsEmail, *Response, error) {
	return doEmptyPost(ctx, s.client, userBasePath, new(CredentialsEmail),
		id, "credentials_email", "send_password_reset")
}

// GetRoles -
func (s *UsersResourceOp) GetRoles(ctx context.Context, id string) ([]Role, *Response, error) {
	return doList(ctx, s.client, userBasePath, nil, new([]Role), id, "roles")
}

// SetRoles -
func (s *UsersResourceOp) SetRoles(ctx context.Context, id string, roleIds []string) ([]Role, *Response, error) {
	return doSet(ctx, s.client, userBasePath, roleIds, new([]Role), id, "roles")
}
