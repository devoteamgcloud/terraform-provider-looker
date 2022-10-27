package lookergo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/beefsack/go-rate"
	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const (
	libraryVersion = "0.0.1-dev"
	defaultBaseURL = "https://api.example.com/"
	userAgent      = "API/" + libraryVersion
	mediaType      = "application/json"
)

var rl *rate.RateLimiter

func init() {
	rl = rate.New(1, 1000*time.Millisecond) // Once per second
}

// Rate contains the rate limit for the current client.
type Rate struct {
	// The number of request per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset Timestamp `json:"reset"`
}

// Client manages communication with the API.
type Client struct {
	// Pointer reference to a shared HTTP client for communicating with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for HTTP client
	UserAgent string

	// Rate contains the current rate limit for the client as determined by the most recent
	// API call. It is not thread-safe. Please consider using GetRate() instead.
	Rate    Rate
	ratemtx sync.Mutex

	mu sync.Mutex

	// Resources used for communicating with the API
	Groups      GroupsResource
	Users       UsersResource
	Roles       RolesResource
	Folders     FoldersResource
	Workspaces  WorkspacesResource
	Projects    ProjectsResource
	Sessions    SessionsResource
	ModelSets   ModelSetsResource
	Connections ConnectionsResource
	LookMLModel LookMlModelsResource

	// TODO: Expand

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string

	// Production or dev workspace
	Workspace string
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// ListOptions specifies the optional parameters to various List methods that
// support ~~pagination~~ limit/offset querystring.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Limit int `url:"limit,omitempty"`

	// For paginated result sets, the number of results to include per page.
	Offset int `url:"offset,omitempty"`
}

// Response is an API response. This wraps the standard http.Response returned from the API.
type Response struct {
	*http.Response

	Rate
}

type Error struct {
	Field            string `json:"field"`
	Code             string `json:"code"`
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
	// Error message
	Message string `json:"message"`
	// RequestID returned from the API, useful to contact support.
	DocsURL string `json:"docs"`
	// More specific error msg
	Errors []Error `json:"errors,omitempty"`
}

// NewClient -
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Groups = &GroupsResourceOp{client: c}
	c.Users = &UsersResourceOp{client: c}
	c.Roles = &RolesResourceOp{client: c}
	c.Folders = &FoldersResourceOp{client: c}
	c.Workspaces = &WorkspacesResourceOp{client: c}
	c.Projects = &ProjectsResourceOp{client: c}
	c.Sessions = &SessionsResourceOp{client: c}
	c.ModelSets = &ModelSetsResourceOp{client: c}
	c.Connections = &ConnectionsResourceOp{client: c}
	c.LookMLModel = &LookMlModelsResourceOp{client: c}

	c.headers = make(map[string]string)
	c.Workspace = "production"

	return c
}

// ApiConfig -
type ApiConfig struct {
	ClientId     string
	ClientSecret string
	BaseURL      string
	ClientCtx    context.Context
}

// AuthToken -
type AuthToken struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type,omitempty"`
	ExpiresIn    int     `json:"expires_in,omitempty"` // Differs from oauth2.Token{}: `json:"expiry,omitempty"
	RefreshToken *string `json:"refresh_token,omitempty"`
	raw          interface{}
}

// NewFromApiv3Creds -
func NewFromApiv3Creds(config ApiConfig) *Client {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	// clientcredentials.Config manages the token refreshing
	oauthConfig := clientcredentials.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		TokenURL:     fmt.Sprintf("%s4.0/login", config.BaseURL),
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	// ctx := context.Background()
	ts := oauthConfig.TokenSource(config.ClientCtx)

	return NewClient(oauth2.NewClient(config.ClientCtx, ts))
}

// NewFromStaticToken returns a new  API client using the given API token.
func NewFromStaticToken(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	return NewClient(oauth2.NewClient(ctx, ts))
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new API client instance.
// func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
// 	c := NewClient(httpClient)
// 	for _, opt := range opts {
// 		if err := opt(c); err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return c, nil
// }

// SetBaseURL is a client option for setting the base URL.
func (c *Client) SetBaseURL(bu string) error {
	u, err := url.Parse(bu)
	if err != nil {
		return err
	}

	c.BaseURL = u
	return nil
}

// SetUserAgent is a client option for setting the user agent.
func (c *Client) SetUserAgent(ua string) error {
	c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
	return nil
}

// SetRequestHeaders sets optional HTTP headers on the client that are
// sent on each HTTP request.
func (c *Client) SetRequestHeaders(headers map[string]string) error {
	for k, v := range headers {
		c.headers[k] = v
	}
	return nil
}

func (c *Client) SetOauthCredentials(ctx context.Context, clientId string, clientSecret string) error {
	var loginUrl url.URL
	if c.BaseURL != nil {
		loginUrl = *c.BaseURL
	} else {
		u, _ := url.Parse(defaultBaseURL)
		loginUrl = *u
	}
	loginUrl.Path = path.Join(loginUrl.Path, "4.0", "login")

	oauthConfig := clientcredentials.Config{
		ClientID:     strings.Trim(strings.TrimSpace(clientId), "'"),
		ClientSecret: strings.Trim(strings.TrimSpace(clientSecret), "'"),
		TokenURL:     loginUrl.String(),
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	tokenSource := oauthConfig.TokenSource(ctx)
	c.client = oauth2.NewClient(ctx, tokenSource)
	return nil
}

func (c *Client) SetOauthStaticToken(ctx context.Context, token *oauth2.Token) error {
	if token == nil {
		return fmt.Errorf("no token provided")
	}

	tokenSource := oauth2.StaticTokenSource(token)
	c.client = oauth2.NewClient(ctx, tokenSource)
	return nil
}

func (c *Client) EnsureStaticToken(ctx context.Context, parentClient *Client, apiUserID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client.Transport == nil && c.Workspace != "dev" {
		// Get duplicate API token for current user
		token, _, err := parentClient.Sessions.GetLoginUserToken(ctx, apiUserID)
		if err != nil {
			return err
		}

		err = c.SetOauthStaticToken(ctx, token)
		if err != nil {
			return err
		}

		// Set dev workspace for dup token
		session, _, err := c.Sessions.SetWorkspaceId(ctx, "dev")
		if err != nil {
			return err
		} else if session.WorkspaceId == "production" || session.WorkspaceId == "" {
			return fmt.Errorf("did not find dev workspace")
		} else {
			c.Workspace = session.WorkspaceId
		}
	}

	return nil
}

func (c *Client) CreateDevConnection(ctx context.Context, rc RequestCompletionCallback) (*Client, *Session, error) {

	// Get current user ID
	user, _, err := c.Sessions.GetCurrentUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get duplicate API token for current user
	token, _, err := c.Sessions.GetLoginUserToken(ctx, strconv.Itoa(user.Id))
	if err != nil {
		return nil, nil, err
	}

	// Create duplicate API client
	devClient := NewClient(nil)
	if err := devClient.SetBaseURL(c.BaseURL.String()); err != nil {
		return nil, nil, err
	}

	// Set token as static on dup client
	err = devClient.SetOauthStaticToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	devClient.OnRequestCompleted(rc)

	// Set dev workspace for dup token
	session, _, err := devClient.Sessions.SetWorkspaceId(ctx, "dev")
	if err != nil {
		return nil, nil, err
	}

	if session.WorkspaceId != "dev" {
		return nil, nil, fmt.Errorf("failed to set dev environment for" +
			"duplicate dev workspace connection")
	} else {
		c.Workspace = "dev"
	}

	return devClient, session, err
}

/*
TODO: Pagination.
It is provided in the headers
GET {{endpoint}}/4.0/users?limit=5
(H) X-Total-Count: 107
(H) Link: <https://x.cloud.looker.com:19999/api/4.0/users?limit=5&offset=0>; rel="first",<https://x.cloud.looker.com:19999/api/4.0/users?limit=5&offset=105>; rel="last",<https://reprise.cloud.looker.com:19999/api/4.0/users?limit=5&offset=5>; rel="next"

*/

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, _ = io.CopyN(ioutil.Discard, resp.Body, maxBodySlurpSize)
		}

		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	c.ratemtx.Lock()
	response := newResponse(resp)
	rl.Wait()

	c.Rate = response.Rate
	c.ratemtx.Unlock()

	err = CheckResponse(resp)

	if err != nil {
		return response, err
	}

	if v != nil {
		if resp.StatusCode == 204 {
			return response, nil
		} else if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
// If the API error response does not include the request ID in its body, the one from its header will be used.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

func (r Rate) String() string {
	return Stringify(r)
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

type service interface {
	Group | User | CredentialEmail | Role | PermissionSet | Session | Project | GitBranch | Folder
}

// addOptions -
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// doList is a generic list lookup
func doList[T any](ctx context.Context, client *Client, basePath string, opt *ListOptions, svc *[]T, pathSuffix ...string) ([]T, *Response, error) {
	path := fmt.Sprintf("%s%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"))
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, svc)
	if err != nil {
		return nil, resp, err
	}

	return *svc, resp, err
}

func doGet[T any](ctx context.Context, client *Client, basePath string, svc *T, pathSuffix ...string) (*T, *Response, error) {
	path := fmt.Sprintf("%s%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"))

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, svc)
	if err != nil {
		return nil, resp, err
	}

	return svc, resp, err
}

func doGetById[T any](ctx context.Context, client *Client, basePath string, id any, svc *T) (*T, *Response, error) {
	switch id.(type) {
	case int:
		if id.(int) < 1 {
			return nil, nil, NewArgError("id", "cannot be less than 1")
		}
	case string:
	default:
		panic("Invalid type for ID. Has to be either int or string")
	}

	path := fmt.Sprintf("%s/%v", basePath, id)

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, svc)
	if err != nil {
		return nil, resp, err
	}

	return svc, resp, err
}

func doListByX[T service](ctx context.Context, client *Client, basePath string, opt *ListOptions, svc *[]T, qs url.Values) ([]T, *Response, error) {
	path := fmt.Sprintf("%s?%s", basePath, qs.Encode())
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, svc)
	if err != nil {
		return nil, resp, err
	}

	return *svc, resp, err
}

func doCreate[T any, N any](ctx context.Context, client *Client, basePath string, svc *T, newSvc *N, pathSuffix ...string) (*N, *Response, error) {
	path := fmt.Sprintf("%s%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"))

	req, err := client.NewRequest(ctx, http.MethodPost, path, svc)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, newSvc)
	if err != nil {
		return nil, resp, err
	}

	return newSvc, resp, err
}

func doCreateX(ctx context.Context, client *Client, basePath string, newSvc *string, qs url.Values, pathSuffix ...string) (*string, *Response, error) {
	var path string
	if qs != nil {
		path = fmt.Sprintf("%s%s?%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"), qs.Encode())
	} else {
		path = fmt.Sprintf("%s%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"))
	}

	req, err := client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, newSvc)
	if err != nil {
		return nil, resp, err
	}

	return newSvc, resp, err
}

func doEmptyPost[N service](ctx context.Context, client *Client, basePath string, newSvc *N, pathSuffix ...string) (*N, *Response, error) {
	path := fmt.Sprintf("%s%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"))

	req, err := client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, newSvc)
	if err != nil {
		return nil, resp, err
	}

	return newSvc, resp, err
}

func doUpdate[T any, U any](ctx context.Context, client *Client, basePath string, id any, svc *T, uSvc *U, pathSuffix ...string) (*U, *Response, error) {
	var path string

	switch id.(type) {
	case int:
		if id.(int) < 1 {
			return nil, nil, NewArgError("id", "cannot be less than 1")
		}
		path = fmt.Sprintf("%s/%d%s", basePath, id.(int), strings.Join(append([]string{""}, pathSuffix...), "/"))
	case string:
		path = fmt.Sprintf("%s/%s%s", basePath, id.(string), strings.Join(append([]string{""}, pathSuffix...), "/"))
	default:
		panic("Invalid type for ID. Has to be either int or string")
	}

	req, err := client.NewRequest(ctx, http.MethodPatch, path, svc)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, uSvc)
	if err != nil {
		return nil, resp, err
	}

	return uSvc, resp, err
}

func doSet[T any](ctx context.Context, client *Client, basePath string, ids []string, svc *[]T, pathSuffix ...string) ([]T, *Response, error) {
	if len(ids) < 1 {
		return nil, nil, NewArgError("ids", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s%s", basePath, strings.Join(append([]string{""}, pathSuffix...), "/"))

	req, err := client.NewRequest(ctx, http.MethodPut, path, ids)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, svc)
	if err != nil {
		return nil, resp, err
	}

	return *svc, resp, err
}

func doDelete(ctx context.Context, client *Client, basePath string, id any, pathSuffix ...string) (*Response, error) {
	var path string

	switch id.(type) {
	case int:
		if id.(int) < 1 {
			return nil, NewArgError("id", "cannot be less than 1")
		}
		path = fmt.Sprintf("%s/%d%s", basePath, id.(int), strings.Join(append([]string{""}, pathSuffix...), "/"))
	case string:
		path = fmt.Sprintf("%s/%s%s", basePath, id.(string), strings.Join(append([]string{""}, pathSuffix...), "/"))
	default:
		panic("Invalid type for ID. Has to be either int or string")
	}

	req, err := client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func doAddMember[T service](ctx context.Context, client *Client, path string, svc *T, addNew interface{}) (*T, *Response, error) {

	req, err := client.NewRequest(ctx, http.MethodPost, path, addNew)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(ctx, req, svc)
	if err != nil {
		return nil, resp, err
	}

	return svc, resp, err
}

func boolPtr(b bool) *bool {
	return &b
}
