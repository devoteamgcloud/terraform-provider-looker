package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Role name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permission_set_id": {
				Description:  "PermissionSet ID",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"permission_set_id", "permission_set_name"},
			},
			"permission_set_name": {
				Description:  "PermissionSet Name",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"permission_set_id", "permission_set_name"},
			},
			"model_set_id": {
				Description: "Modelset name",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	logTrace(ctx, "query role", "role_id", d.Id())

	role, _, err := c.Roles.Get(ctx, idAsInt(d.Id()))
	switch err.(type) {
	case *lookergo.ErrorResponse:
		if errResp := err.(*lookergo.ErrorResponse).Response; errResp.StatusCode == http.StatusNotFound {
			logTrace(ctx, "role not found", "role_id", d.Id())
			d.SetId("")
			return // Resource was not found.
		}
	case error:
		return logErrDiag(ctx, diags, "unable to query role", "role_id", d.Id()) // Connection error.
	default:
		logTrace(ctx, "role found", "role", role)
	}

	d.Set("name", role.Name)
	d.Set("permission_set_id", idAsString(role.PermissionSet.Id))
	d.Set("permission_set_name", role.PermissionSet.Name)
	d.Set("model_set_id", role.ModelSet.Id)

	return diags
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	permissionSets, _, err := c.Roles.PermissionSetsList(ctx, nil)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to query permission sets", "err", err)
	}

	var permissionSet lookergo.PermissionSet
	if psId, ok := d.GetOk("permission_set_id"); ok {
		if has, ps := permissionSets.HasById(idAsInt(psId)); has {
			permissionSet = *ps
		} else {
			return logErrDiag(ctx, diags, "PermissionSet not found", "permission_set_id", psId)
		}
	} else if psName, ok := d.GetOk("permission_set_name"); ok {
		if has, ps := permissionSets.HasByName(psName.(string)); has {
			permissionSet = *ps
		} else {
			return logErrDiag(ctx, diags, "PermissionSet not found", "permission_set_name", psId)
		}
	}

	modelSet, _, err := c.ModelSets.Get(ctx, idAsString(d.Get("model_set_id")))
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to find ModelSet", "model_set_id", err)
	}

	roleName := d.Get("name").(string)

	role := lookergo.Role{
		Name:            roleName,
		PermissionSetID: idAsString(permissionSet.Id),
		ModelSetID:      modelSet.Id,
	}

	newRole, _, err := c.Roles.Create(ctx, &role)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to create Role", "err", err)
	}
	d.Set("name", newRole.Name)
	d.Set("permission_set_id", idAsString(newRole.PermissionSet.Id))
	d.Set("permission_set_name", newRole.PermissionSet.Name)
	d.Set("model_set_id", newRole.ModelSet.Id)

	d.SetId(idAsString(newRole.Id))

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	permissionSets, _, err := c.Roles.PermissionSetsList(ctx, nil)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to query permission sets", "err", err)
	}

	var permissionSet lookergo.PermissionSet
	if psId, ok := d.GetOk("permission_set_id"); ok {
		if has, ps := permissionSets.HasById(idAsInt(psId)); has {
			permissionSet = *ps
		} else {
			return logErrDiag(ctx, diags, "PermissionSet not found", "permission_set_id", psId)
		}
	} else if psName, ok := d.GetOk("permission_set_name"); ok {
		if has, ps := permissionSets.HasByName(idAsString(psName)); has {
			permissionSet = *ps
		} else {
			return logErrDiag(ctx, diags, "PermissionSet not found", "permission_set_name", psId)
		}
	}

	role := lookergo.Role{
		Name:            d.Get("name").(string),
		PermissionSetID: idAsString(permissionSet.Id),
		ModelSetID:      d.Get("model_set_id").(string),
	}

	newRole, _, err := c.Roles.Update(ctx, idAsInt(d.Id()), &role)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to create Role", "err", err)
	}
	d.Set("permission_set_id", newRole.PermissionSet.Id)
	d.Set("permission_set_name", newRole.PermissionSet.Name)

	logTrace(ctx, "updated role", "new_role", newRole)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	_, err := c.Roles.Delete(ctx, idAsInt(d.Id()))
	if err != nil {
		return logErrDiag(ctx, diags, "failed to delete role", "err", err)
	}

	d.SetId("") // Finally mark as deleted
	return diags
}
