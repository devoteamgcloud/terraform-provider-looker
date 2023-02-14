package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceProjectGitRepo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGitRepoCreate,
		ReadContext:   resourceProjectGitRepoRead,
		UpdateContext: resourceProjectGitRepoUpdate,
		DeleteContext: resourceProjectGitRepoDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"git_remote_url": {
				Type: schema.TypeString, Required: true,
			},
			"git_username": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Git username for HTTPS authentication. For SSH authentication " +
					"skip this option and create project_git_deploy_key resource.",
			},
			"git_password": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Git password for HTTPS authentication. For SSH authentication " +
					"skip this option and create project_git_deploy_key resource.",
				Sensitive: true,
			},
			"use_git_cookie_auth": {
				Type: schema.TypeBool, Optional: true,
				Description: "If true, the project uses a git cookie for authentication.",
			},
			"git_service_name": {
				Type: schema.TypeString, Required: true,
				Description: "Name of the git service provider",
			},
			"git_production_branch_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "main",
				Description: "Git production branch name. Defaults to ~~master~~ main. " +
					"Supported only in Looker 21.0 and higher.",
			},
			"pull_request_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "off",
				Description: "The git pull request policy for this project. " +
					"Valid values are: `off`, `links`, `recommended`, `required`.",
			},
			"validation_required": {
				Type: schema.TypeBool, Optional: true, Default: true,
			},
			"allow_warnings": {
				Type: schema.TypeBool, Optional: true, Default: true,
			},
			"is_example": {
				Type: schema.TypeBool, Optional: true,
			},
			"git_release_mgmt_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Advanced Deploy Mode - Required for Webhook",
			},
			"deploy_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secret Value for Authentication Webhook",
			},
			"deploy_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Branch which will be deployed to Production after " +
					"creation of Project Resource. Required: Advanced Deploy Mode.",
			},
		},
		Importer: nil,
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
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	projectName := d.Get("project_id").(string)

	project, _, err := dc.Projects.Get(ctx, projectName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("git_remote_url", project.GitRemoteUrl)
	d.Set("git_username", project.GitUsername)
	d.Set("use_git_cookie_auth", project.UseGitCookieAuth)
	d.Set("git_service_name", project.GitServiceName)
	d.Set("pull_request_mode", project.PullRequestMode)
	d.Set("validation_required", project.ValidationRequired)
	d.Set("git_production_branch_name", project.GitProductionBranchName)
	d.Set("allow_warnings", project.AllowWarnings)
	d.Set("is_example", project.IsExample)
	d.Set("git_release_mgmt_enabled", project.GitReleaseMgmtEnabled)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectGitRepoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	projectName := d.Get("project_id").(string)

	projectGitRepoUpdate := lookergo.Project{}
	if value, ok := d.GetOk("allow_warnings"); ok {
		projectGitRepoUpdate.AllowWarnings = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_remote_url"); ok {
		projectGitRepoUpdate.GitRemoteUrl = value.(string)
	}
	if value, ok := d.GetOk("git_username"); ok {
		projectGitRepoUpdate.GitUsername = value.(string)
	}
	if value, ok := d.GetOk("use_git_cookie_auth"); ok {
		projectGitRepoUpdate.UseGitCookieAuth = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_service_name"); ok {
		projectGitRepoUpdate.GitServiceName = value.(string)
	}
	if value, ok := d.GetOk("pull_request_mode"); ok {
		projectGitRepoUpdate.PullRequestMode = value.(string)
	}
	if value, ok := d.GetOk("git_production_branch_name"); ok {
		projectGitRepoUpdate.GitProductionBranchName = value.(string)
	}
	if value, ok := d.GetOk("validation_required"); ok {
		projectGitRepoUpdate.ValidationRequired = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("is_example"); ok {
		projectGitRepoUpdate.IsExample = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_release_mgmt_enabled"); ok {
		projectGitRepoUpdate.GitReleaseMgmtEnabled = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("deploy_secret"); ok {
		projectGitRepoUpdate.DeploySecret = value.(string)
	}

	if value, ok := d.GetOk("git_password"); ok {
		payload := lookergo.Project{}
		payload.GitRemoteUrl = projectGitRepoUpdate.GitRemoteUrl
		projectGitRepoUpdate.GitPassword = value.(string)
		payload.GitUsername = projectGitRepoUpdate.GitUsername
		payload.GitPassword = projectGitRepoUpdate.GitPassword
		payload.GitServiceName = projectGitRepoUpdate.GitServiceName
		if !strings.HasPrefix(projectGitRepoUpdate.GitRemoteUrl, "https://") {
			return diag.Errorf("HTTPS Authentication requires URL starts with http://..")
		}
		_, _, err = dc.Projects.Update(ctx, projectName, &payload)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		payload := lookergo.Project{}
		payload.GitRemoteUrl = projectGitRepoUpdate.GitRemoteUrl
		if !strings.HasPrefix(projectGitRepoUpdate.GitRemoteUrl, "git@") && !strings.HasPrefix(payload.GitRemoteUrl, "ssh://") {
			return diag.Errorf("SSH Authentication requires URL starts with git@.. or ssh://..")
		}
		_, _, err = dc.Projects.Update(ctx, projectName, &payload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, _, err = dc.Projects.Update(ctx, projectName, &projectGitRepoUpdate)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	if value, ok := d.GetOk("deploy_branch"); ok {
		branchName := value.(string)
		_, _, err = dc.Projects.GitBranchDeployToProduction(ctx, projectName, branchName)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		_, _, err = dc.Projects.DeployToProduction(ctx, projectName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(projectName)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectGitRepoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	projectName := d.Get("project_id").(string)

	projectGitRepoUpdate := lookergo.Project{}
	if value, ok := d.GetOk("allow_warnings"); ok {
		projectGitRepoUpdate.AllowWarnings = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_remote_url"); ok {
		projectGitRepoUpdate.GitRemoteUrl = value.(string)
	}
	if value, ok := d.GetOk("git_username"); ok {
		projectGitRepoUpdate.GitUsername = value.(string)
	}
	if value, ok := d.GetOk("use_git_cookie_auth"); ok {
		projectGitRepoUpdate.UseGitCookieAuth = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_service_name"); ok {
		projectGitRepoUpdate.GitServiceName = value.(string)
	}
	if value, ok := d.GetOk("pull_request_mode"); ok {
		projectGitRepoUpdate.PullRequestMode = value.(string)
	}
	if value, ok := d.GetOk("git_production_branch_name"); ok {
		projectGitRepoUpdate.GitProductionBranchName = value.(string)
	}
	if value, ok := d.GetOk("validation_required"); ok {
		projectGitRepoUpdate.ValidationRequired = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("is_example"); ok {
		projectGitRepoUpdate.IsExample = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("git_release_mgmt_enabled"); ok {
		projectGitRepoUpdate.GitReleaseMgmtEnabled = boolPtr(value.(bool))
	}
	if value, ok := d.GetOk("deploy_secret"); ok {
		projectGitRepoUpdate.DeploySecret = value.(string)
	}
	if value, ok := d.GetOk("git_password"); ok {
		projectGitRepoUpdate.GitPassword = value.(string)
		if !strings.HasPrefix(projectGitRepoUpdate.GitRemoteUrl, "https://") {
			return diag.Errorf("HTTPS Authentication requires URL starts with http://..")
		}
	} else {
		if !strings.HasPrefix(projectGitRepoUpdate.GitRemoteUrl, "git@") && !strings.HasPrefix(projectGitRepoUpdate.GitRemoteUrl, "ssh://") {
			return diag.Errorf("SSH Authentication requires URL starts with git@.. or ssh://..")
		}
	}

	_, _, err = dc.Projects.Update(ctx, projectName, &projectGitRepoUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	if value, ok := d.GetOk("deploy_branch"); ok {
		branchName := value.(string)
		_, _, err = dc.Projects.GitBranchDeployToProduction(ctx, projectName, branchName)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		_, _, err = dc.Projects.DeployToProduction(ctx, projectName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceProjectGitRepoRead(ctx, d, m)
}

// 	Call update_project setting git_remote_url to null and git_service_name to "bare".

// resourceProjectGitRepoDelete
func resourceProjectGitRepoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	projectName := d.Get("project_id").(string)

	_, err = dc.Projects.DeleteGitRepo(ctx, projectName)
	if err != nil {
		return diag.FromErr(err)
	}

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
