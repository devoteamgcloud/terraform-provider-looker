package lookergo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/crypto/ssh"
)

const api_version = "4.0"
const projectsBasePath = api_version + "/projects"

// Ref: https://developers.looker.com/api/explorer/4.0/types/Project

type ProjectsResource interface {
	Get(ctx context.Context, projectName string) (*Project, *Response, error)
	// !!! NOT rest compliant !!!
	// name is required. git_remote_url is not allowed.
	// To configure Git for the newly created project, follow the instructions in update_project.
	Create(ctx context.Context, proj *Project) (*Project, *Response, error)
	Update(ctx context.Context, projectName string, proj *Project) (*Project, *Response, error)
	Delete(ctx context.Context, projectName string) (*Response, error)
	AllowWarnings(ctx context.Context, projectName string, value bool) (*Response, error)
	DeleteGitRepo(ctx context.Context, projectName string) (*Response, error)
	GitBranchesList(ctx context.Context, projectName string, opt *ListOptions) ([]GitBranch, *Response, error)
	GitBranchActiveGet(ctx context.Context, projectName string) (*GitBranch, *Response, error)
	GitBranchCheckout(ctx context.Context, projectName string, gbr *GitBranchRef) (*GitBranch, *Response, error)
	GitBranchUpdate(ctx context.Context, projectName string, gbr *GitBranchRef) (*GitBranch, *Response, error)
	GitBranchListByName(ctx context.Context, projectName string, branchName string) (*GitBranch, *Response, error)
	GitBranchDelete(ctx context.Context, projectName string, branchName string) (*Response, error)
	GitBranchDeployToProduction(ctx context.Context, projectName string, branch string) (*string, *Response, error)
	GitRefDeployToProduction(ctx context.Context, projectName string, ref string) (*string, *Response, error)
	DeployToProduction(ctx context.Context, projectName string) (*string, *Response, error)
	GitDeployKeyGet(ctx context.Context, projectName string) (*string, *Response, error)
	GitDeployKeyCreate(ctx context.Context, projectName string) (*string, *Response, error)
	// GitDeployKeyDelete(ctx context.Context, projectName string) (*Response, error) // Doesn't exist
}

type ProjectsResourceOp struct {
	client *Client
}

var _ ProjectsResource = &ProjectsResourceOp{}

type PullRequestMode string

const PullRequestMode_Off PullRequestMode = "off"
const PullRequestMode_Links PullRequestMode = "links"
const PullRequestMode_Recommended PullRequestMode = "recommended"
const PullRequestMode_Required PullRequestMode = "required"

// Dynamic writeable type for Project removes:
// can, id, uses_git, is_example
type WorkProject struct {
	Name                           string `json:"name,omitempty"`                               // Project display name
	GitRemoteUrl                   string `json:"git_remote_url,omitempty"`                     // Git remote repository url
	GitUsername                    string `json:"git_username,omitempty"`                       // Git username for HTTPS authentication. (For production only, if using user attributes.)
	GitPassword                    string `json:"git_password,omitempty"`                       // (Write-Only) Git password for HTTPS authentication. (For production only, if using user attributes.)
	GitProductionBranchName        string `json:"git_production_branch_name,omitempty"`         // Git production branch name. Defaults to master. Supported only in Looker 21.0 and higher.
	UseGitCookieAuth               *bool  `json:"use_git_cookie_auth,omitempty"`                // If true, the project uses a git cookie for authentication.
	GitUsernameUserAttribute       string `json:"git_username_user_attribute,omitempty"`        // User attribute name for username in per-user HTTPS authentication.
	GitPasswordUserAttribute       string `json:"git_password_user_attribute,omitempty"`        // User attribute name for password in per-user HTTPS authentication.
	GitServiceName                 string `json:"git_service_name,omitempty"`                   // Name of the git service provider
	GitApplicationServerHttpPort   int64  `json:"git_application_server_http_port,omitempty"`   // Port that HTTP(S) application server is running on (for PRs, file browsing, etc.)
	GitApplicationServerHttpScheme string `json:"git_application_server_http_scheme,omitempty"` // Scheme that is running on application server (for PRs, file browsing, etc.)
	DeploySecret                   string `json:"deploy_secret,omitempty"`                      // (Write-Only) Optional secret token with which to authenticate requests to the webhook deploy endpoint. If not set, endpoint is unauthenticated.
	UnsetDeploySecret              *bool  `json:"unset_deploy_secret,omitempty"`                // (Write-Only) When true, unsets the deploy secret to allow unauthenticated access to the webhook deploy endpoint.
	PullRequestMode                string `json:"pull_request_mode,omitempty"`                  // The git pull request policy for this project. Valid values are: "off", "links", "recommended", "required".
	ValidationRequired             *bool  `json:"validation_required,omitempty"`                // Validation policy: If true, the project must pass validation checks before project changes can be committed to the git repository
	GitReleaseMgmtEnabled          *bool  `json:"git_release_mgmt_enabled,omitempty"`           // If true, advanced git release management is enabled for this project
	AllowWarnings                  *bool  `json:"allow_warnings,omitempty"`                     // Validation policy: If true, the project can be committed with warnings when `validation_required` is true. (`allow_warnings` does nothing if `validation_required` is false).
	DependencyStatus               string `json:"dependency_status,omitempty"`                  // Status of dependencies in your manifest & lockfile
}

// Project struct for Project
type Project struct {
	// Project Id
	Id string `json:"id,omitempty"`
	// Project display name
	Name string `json:"name,omitempty"`
	// If true the project is configured with a git repository
	UsesGit *bool `json:"uses_git,omitempty"`
	// Git remote repository url
	GitRemoteUrl string `json:"git_remote_url,omitempty"`
	// Git username for HTTPS authentication. (For production only, if using user attributes.)
	GitUsername string `json:"git_username,omitempty"`
	// (Write-Only) Git password for HTTPS authentication. (For production only, if using user attributes.)
	GitPassword string `json:"git_password,omitempty"`
	// Git production branch name. Defaults to master. Supported only in Looker 21.0 and higher.
	GitProductionBranchName string `json:"git_production_branch_name,omitempty"`
	// If true, the project uses a git cookie for authentication.
	UseGitCookieAuth *bool `json:"use_git_cookie_auth,omitempty"`
	// User attribute name for username in per-user HTTPS authentication.
	GitUsernameUserAttribute string `json:"git_username_user_attribute,omitempty"`
	// User attribute name for password in per-user HTTPS authentication.
	GitPasswordUserAttribute string `json:"git_password_user_attribute,omitempty"`
	// Name of the git service provider
	GitServiceName string `json:"git_service_name,omitempty"`
	// Port that HTTP(S) application server is running on (for PRs, file browsing, etc.)
	GitApplicationServerHttpPort int64 `json:"git_application_server_http_port,omitempty"`
	// Scheme that is running on application server (for PRs, file browsing, etc.)
	GitApplicationServerHttpScheme string `json:"git_application_server_http_scheme,omitempty"`
	// (Write-Only) Optional secret token with which to authenticate requests to the webhook deploy endpoint. If not set, endpoint is unauthenticated.
	DeploySecret string `json:"deploy_secret,omitempty"`
	// (Write-Only) When true, unsets the deploy secret to allow unauthenticated access to the webhook deploy endpoint.
	UnsetDeploySecret *bool `json:"unset_deploy_secret,omitempty"`
	// The git pull request policy for this project. Valid values are: \"off\", \"links\", \"recommended\", \"required\".
	PullRequestMode string `json:"pull_request_mode,omitempty"`
	// Validation policy: If true, the project must pass validation checks before project changes can be committed to the git repository
	ValidationRequired *bool `json:"validation_required,omitempty"`
	// If true, advanced git release management is enabled for this project
	GitReleaseMgmtEnabled *bool `json:"git_release_mgmt_enabled,omitempty"`
	// Validation policy: If true, the project can be committed with warnings when `validation_required` is true. (`allow_warnings` does nothing if `validation_required` is false).
	AllowWarnings *bool `json:"allow_warnings,omitempty"`
	// If true the project is an example project and cannot be modified
	IsExample *bool `json:"is_example,omitempty"`
	// Status of dependencies in your manifest & lockfile
	DependencyStatus string `json:"dependency_status,omitempty"`
}

// GitBranch struct for GitBranch
type GitBranch struct {
	// The short name on the local. Updating `name` results in `git checkout <new_name>`
	Name string `json:"name,omitempty"`
	// The name of the remote
	Remote string `json:"remote,omitempty"`
	// The short name on the remote
	RemoteName string `json:"remote_name,omitempty"`
	// Name of error
	Error string `json:"error,omitempty"`
	// Message describing an error if present
	Message string `json:"message,omitempty"`
	// Name of the owner of a personal branch
	OwnerName string `json:"owner_name,omitempty"`
	// Whether or not this branch is readonly
	Readonly *bool `json:"readonly,omitempty"`
	// Whether or not this branch is a personal branch - readonly for all developers except the owner
	Personal *bool `json:"personal,omitempty"`
	// Whether or not a local ref exists for the branch
	IsLocal *bool `json:"is_local,omitempty"`
	// Whether or not a remote ref exists for the branch
	IsRemote *bool `json:"is_remote,omitempty"`
	// Whether or not this is the production branch
	IsProduction *bool `json:"is_production,omitempty"`
	// Number of commits the local branch is ahead of the remote
	AheadCount int64 `json:"ahead_count,omitempty"`
	// Number of commits the local branch is behind the remote
	BehindCount int64 `json:"behind_count,omitempty"`
	// UNIX timestamp at which this branch was last committed.
	CommitAt int64 `json:"commit_at,omitempty"`
	// The resolved ref of this branch. Updating `ref` results in `git reset --hard <new_ref>``.
	Ref string `json:"ref,omitempty"`
	// The resolved ref of this branch remote.
	RemoteRef string `json:"remote_ref,omitempty"`
}

// GitBranchRef -
type GitBranchRef struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
}

func (s *ProjectsResourceOp) Get(ctx context.Context, projectName string) (*Project, *Response, error) {
	return doGet(ctx, s.client, projectsBasePath, new(Project), projectName)
}

/*
	Create Project

	Create A Project

	dev mode required.

	Call update_session to select the 'dev' workspace.
	name is required. git_remote_url is not allowed. To configure Git for the newly created project, follow the instructions in update_project.
*/
func (s *ProjectsResourceOp) Create(ctx context.Context, proj *Project) (*Project, *Response, error) {
	return doCreate(ctx, s.client, projectsBasePath, proj, new(Project))
}

/*
	Update Project

	Update Project Configuration

	Apply changes to a project's configuration.

	Configuring Git for a Project

	To set up a Looker project with a remote git repository, follow these steps:

	Call update_session to select the 'dev' workspace.
	Call create_git_deploy_key to create a new deploy key for the project
	Copy the deploy key text into the remote git repository's ssh key configuration
	Call update_project to set project's git_remote_url ()and git_service_name, if necessary).
	When you modify a project's git_remote_url, Looker connects to the remote repository to fetch metadata.
	The remote git repository MUST be configured with the Looker-generated deploy key for this project prior to setting the project's git_remote_url.

	To set up a Looker project with a git repository residing on the Looker server (a 'bare' git repo):

	Call update_session to select the 'dev' workspace.
	Call update_project setting git_remote_url to null and git_service_name to "bare".
*/
func (s *ProjectsResourceOp) Update(ctx context.Context, projectName string, proj *Project) (*Project, *Response, error) {
	return doUpdate(ctx, s.client, projectsBasePath, projectName, proj, new(Project))
}

func (s *ProjectsResourceOp) Delete(ctx context.Context, projectName string) (*Response, error) {
	return doDelete(ctx, s.client, projectsBasePath, projectName)
}

func (s *ProjectsResourceOp) GitBranchesList(ctx context.Context, projectName string, opt *ListOptions) ([]GitBranch, *Response, error) {
	return doList(ctx, s.client, projectsBasePath, opt, new([]GitBranch), projectName, "git_branches")
}

func (s *ProjectsResourceOp) GitBranchActiveGet(ctx context.Context, projectName string) (*GitBranch, *Response, error) {
	return doGet(ctx, s.client, projectsBasePath, new(GitBranch), projectName, "git_branch")
}

func (s *ProjectsResourceOp) GitBranchCheckout(ctx context.Context, projectName string, gbr *GitBranchRef) (*GitBranch, *Response, error) {
	return doCreate(ctx, s.client, projectsBasePath, new(GitBranchRef), new(GitBranch), projectName, "git_branch")
}

func (s *ProjectsResourceOp) GitBranchUpdate(ctx context.Context, projectName string, gbr *GitBranchRef) (*GitBranch, *Response, error) {
	return doUpdate(ctx, s.client, projectsBasePath, projectName, gbr, new(GitBranch), projectName, "git_branch")
}

func (s *ProjectsResourceOp) GitBranchListByName(ctx context.Context, projectName string, branchName string) (*GitBranch, *Response, error) {
	panic("Not implemented")
}

func (s *ProjectsResourceOp) GitBranchDelete(ctx context.Context, projectName string, branchName string) (*Response, error) {
	return doDelete(ctx, s.client, projectsBasePath, projectName, "git_branch", branchName)
}

func (s *ProjectsResourceOp) GitBranchDeployToProduction(ctx context.Context, projectName string, branch string) (*string, *Response, error) {
	qs := url.Values{}
	qs.Add("branch", branch)

	return doCreateX(ctx, s.client, projectsBasePath, new(string), qs, projectName, "deploy_ref_to_production")
}

func (s *ProjectsResourceOp) GitRefDeployToProduction(ctx context.Context, projectName string, ref string) (*string, *Response, error) {
	qs := url.Values{}
	qs.Add("ref", ref)

	return doCreateX(ctx, s.client, projectsBasePath, new(string), qs, projectName, "deploy_ref_to_production")
}

func (s *ProjectsResourceOp) DeployToProduction(ctx context.Context, projectName string) (*string, *Response, error) {
	return doCreateX(ctx, s.client, projectsBasePath, new(string), nil, projectName, "deploy_to_production")
}

func (s *ProjectsResourceOp) GitDeployKeyGet(ctx context.Context, projectName string) (*string, *Response, error) {
	path := fmt.Sprintf("%s/%s/%s", projectsBasePath, projectName, "git/deploy_key")
	var gitPubKey string

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var buf bytes.Buffer
	resp, err := s.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(&buf)
		if err != nil {
			return nil, resp, err
		}
		publicKey, _, _, _, err := ssh.ParseAuthorizedKey(bodyBytes)
		if err != nil {
			panic(err)
		}

		gitPubKey = string(ssh.MarshalAuthorizedKey(publicKey))
	}

	return &gitPubKey, resp, err
}

func (s *ProjectsResourceOp) GitDeployKeyCreate(ctx context.Context, projectName string) (*string, *Response, error) {
	path := fmt.Sprintf("%s/%s/%s", projectsBasePath, projectName, "git/deploy_key")
	var gitPubKey string

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return &gitPubKey, nil, err
	}

	var buf bytes.Buffer
	resp, err := s.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(&buf)
		if err != nil {
			return nil, resp, err
		}
		publicKey, _, _, _, err := ssh.ParseAuthorizedKey(bodyBytes)
		if err != nil {
			panic(err)
		}

		gitPubKey = string(ssh.MarshalAuthorizedKey(publicKey))
	}

	return &gitPubKey, resp, err

}

func (s *ProjectsResourceOp) DeleteGitRepo(ctx context.Context, projectName string) (*Response, error) {
	path := fmt.Sprintf("%s/%s/", projectsBasePath, projectName)
	log.Printf("[DEBUG] Trying to delete git repo for project %s to path %s", projectName, path)
	gsn := "bare"
	b := map[string]*string{
		"git_remote_url":   nil,
		"git_service_name": &gsn,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(b)
	log.Printf("[DEBUG] Adding body %s", buf)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, b)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return nil, err
}

func (s *ProjectsResourceOp) AllowWarnings(ctx context.Context, projectName string, value bool) (*Response, error) {
	path := fmt.Sprintf("%s/%s", projectsBasePath, projectName)

	b := map[string]bool{
		"allow_warnings": value,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, b)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	resp, err := s.client.Do(ctx, req, &buf)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(&buf)
		if err != nil {
			return resp, err
		}
	}

	return resp, err
}

// func (s *ProjectsResourceOp) GitDeployKeyDelete(ctx context.Context, projectName string) (*Response, error) {
//
// }
