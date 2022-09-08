package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Porject name",
				Type:         schema.TypeString,
				Required:     true,
			},
			"uses_git": {
				Description: "If true, the project uses a git repository.",
				Type:	schema.TypeBool,
				Computed: true,
			},
			"git_remote_url": {
				Description: "Git remote repository url.",
				Type:	schema.TypeString,
				Computed: true,
			},
			"git_username": {
				Description: "Git username for HTTPS authentication.",
				Type:	schema.TypeString,
				Computed: true,
			},
			"validation_required": {
				Description: "Validation policy: If true, the project must pass validation checks before project changes can be committed to the git repository",
				Type:	schema.TypeBool,
				Computed: true,
			},
			"allow_warnings": {
				Description: "Validation policy: If true, the project can be committed with warnings when ",
				Type:	schema.TypeBool,
				Computed: true,
			},
			"is_example": {
				Description: "If true, the project is an example project and cannot be modified.",
				Type:	schema.TypeBool,
				Computed: true,
			},
			"git_service_name": {
				Description: "Name of the git service provider.",
				Type:	schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	if err := ensureDevClient(ctx, m); err != nil {
		return diagErrAppend(diags, err)
	}
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	projectId := d.Get("name").(string)
	var project *lookergo.Project
	project, _, err = dc.Projects.Get(ctx, projectId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(project.Id)
	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("uses_git", project.UsesGit); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("git_remote_url", project.GitRemoteUrl); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("git_username", project.GitUsername); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("validation_required", project.ValidationRequired); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("allow_warnings", project.AllowWarnings); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("is_example", project.IsExample); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("git_service_name", project.GitServiceName); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
