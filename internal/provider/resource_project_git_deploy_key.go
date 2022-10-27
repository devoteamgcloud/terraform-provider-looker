package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGitDeployKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGitDeployKeyCreate,
		ReadContext:   resourceProjectGitDeployKeyRead,
		// UpdateContext: resourceProjectGitDeployKeyUpdate, All fields are ForceNew or Computed w/out Optional, Update is superfluous 
		DeleteContext: resourceProjectGitDeployKeyDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectGitDeployKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	projectName := d.Get("project_id").(string)

	pubKey, _, err := dc.Projects.GitDeployKeyCreate(ctx, projectName)
	if err != nil {
		pubKey, _, err = dc.Projects.GitDeployKeyGet(ctx, projectName)
		if err != nil {
			return diag.FromErr(err)
		} else if pubKey != nil {
			d.Set("public_key", pubKey)
		}
	} else {
		d.Set("public_key", pubKey)
	}

	d.SetId("-")

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectGitDeployKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	projectName := d.Get("project_id").(string)

	pubKey := new(string)
	pubKey, _, err = c.Projects.GitDeployKeyGet(ctx, projectName)
	if err != nil {
		pubKey, _, err = dc.Projects.GitDeployKeyGet(ctx, projectName)
		if err != nil {
			return logErrDiag(ctx, diags, "Could not read ssh public key", "err", err)
		}
	}

	d.Set("public_key", *pubKey)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectGitDeployKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// Looker API doesn't support the deletion of a git deploy key. Nothing happens in this function.
	_ = c

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
