package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGitRepo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGitRepoCreate,
		ReadContext:   resourceProjectGitRepoRead,
		UpdateContext: resourceProjectGitRepoCreate,
		DeleteContext: resourceProjectGitRepoDelete,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_remote_url": {Type: schema.TypeString, Required: true},
			"git_username":   {Type: schema.TypeString, Required: true},
			"git_production_branch_name": {Type: schema.TypeString, Optional: true,
				Default: "main",
				Description: "Git production branch name. Defaults to ~~master~~ main. " +
					"Supported only in Looker 21.0 and higher."},
			"use_git_cookie_auth": {Type: schema.TypeBool, Required: true,
				Description: "If true, the project uses a git cookie for authentication."},
			"git_service_name": {Type: schema.TypeString, Required: true,
				Description: "Name of the git service provider"},
			"pull_request_mode": {Type: schema.TypeString, Optional: true,
				Description: "The git pull request policy for this project. " +
					"Valid values are: `off`, `links`, `recommended`, `required`."},
			"validation_required": {Type: schema.TypeString, Optional: true, Default: false},
			"allow_warnings":      {Type: schema.TypeBool, Optional: true, Default: true},
			"is_example":          {Type: schema.TypeBool, Optional: true},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// https://developers.looker.com/api/explorer/4.0/methods/Project/update_project
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

func resourceProjectGitRepoRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	// c := m.(*Config).Api // .(*lookergo.Client)
	if err := ensureDevClient(ctx, m); err != nil {
		return diagErrAppend(diags, err)
	}
	dc := m.(*Config).DevClient
	projectName := d.Get("project_name").(string)

	project, _, err := dc.Projects.Get(ctx, projectName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("git_remote_url", project.GitRemoteUrl)
	d.Set("git_username", project.GitUsername)
	d.Set("git_production_branch_name", project.GitProductionBranchName)
	d.Set("use_git_cookie_auth", project.UseGitCookieAuth)
	d.Set("git_service_name", project.GitServiceName)
	d.Set("pull_request_mode", project.PullRequestMode)
	d.Set("validation_required", project.ValidationRequired)
	d.Set("allow_warnings", project.AllowWarnings)
	d.Set("is_example", project.IsExample)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

/*
	If someone writes lookml, push to production doesn't work if there are

*/

func resourceProjectGitRepoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	// c := m.(*Config).Api // .(*lookergo.Client)
	if err := ensureDevClient(ctx, m); err != nil {
		return diagErrAppend(diags, err)
	}
	dc := m.(*Config).DevClient
	projectName := d.Get("project_name").(string)

	stageOneUpdate := lookergo.Project{}
	if value, ok := d.GetOk("allow_warnings"); ok {
		stageOneUpdate.AllowWarnings = boolPtr(value.(bool))
	}

	_, _, err := dc.Projects.Update(ctx, projectName, &stageOneUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	stageTwoUpdate := lookergo.Project{}
	if value, ok := d.GetOk("git_remote_url"); ok {
		stageTwoUpdate.GitRemoteUrl = value.(string)
	}
	if value, ok := d.GetOk("git_username"); ok {
		stageTwoUpdate.GitUsername = value.(string)
	}
	if value, ok := d.GetOk("git_production_branch_name"); ok {
		stageTwoUpdate.GitProductionBranchName = value.(string)
	}
	if value, ok := d.GetOk("use_git_cookie_auth"); ok {
		stageTwoUpdate.UseGitCookieAuth = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_service_name"); ok {
		stageTwoUpdate.GitServiceName = value.(string)
	}
	if value, ok := d.GetOk("pull_request_mode"); ok {
		stageTwoUpdate.PullRequestMode = value.(string)
	}
	if value, ok := d.GetOk("validation_required"); ok {
		stageTwoUpdate.ValidationRequired = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("is_example"); ok {
		stageTwoUpdate.IsExample = boolPtr(value.(bool))
	}

	_, _, err = dc.Projects.Update(ctx, projectName, &stageTwoUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(projectName)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectGitRepoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// TODO
	_ = c

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceProjectGitRepoRead(ctx, d, m)
}

func resourceProjectGitRepoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	if err := ensureDevClient(ctx, m); err != nil {
		return diagErrAppend(diags, err)
	}
	dc := m.(*Config).DevClient
	projectName := d.Get("project_name").(string)

	_, err := dc.Projects.DeleteGitRepo(ctx, projectName)
	if err != nil {
		return diag.FromErr(err)
	}

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
