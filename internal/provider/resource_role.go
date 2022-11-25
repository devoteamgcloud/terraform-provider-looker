package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
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
				Type:         schema.TypeInt,
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
				Description: "Modelset ID",
				Type:        schema.TypeInt,
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
	role, response, err := c.Roles.Get(ctx, idAsInt(d.Id()))
	if response.StatusCode == 404 {
		d.SetId("") // Mark as deleted
		return diags
	}
	switch e := err.(type) {
	case *lookergo.ErrorResponse:
		if errResp := e.Response; errResp.StatusCode == http.StatusNotFound {
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
	d.Set("permission_set_id", role.PermissionSet.Id)
	d.Set("permission_set_name", role.PermissionSet.Name)
	d.Set("model_set_id", role.ModelSet.Id)

	return diags
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))
	var permissionSet lookergo.PermissionSet
	if psId, ok := d.GetOk("permission_set_id"); ok {
		perm, _, err := c.PermissionSets.Get(ctx, strconv.Itoa(psId.(int)))
		if err != nil {
			return logErrDiag(ctx, diags, "PermissionSet not found", "permission_set_id", psId)
		}
		permissionSet.Id = perm.Id
		permissionSet.Name = perm.Name
	} else if psName, ok := d.GetOk("permission_set_name"); ok {
		permissions, _, err := c.PermissionSets.GetByName(ctx, psName.(string), nil)
		if err != nil {
			return logErrDiag(ctx, diags, "Failed to query permission sets", "err", err)
		}
		for _, permission := range permissions {
			if permission.Name == psName.(string) {
				permissionSet.Id = permission.Id
				permissionSet.Name = permission.Name
			}
		}

	}
	model_set_id := strconv.Itoa(d.Get("model_set_id").(int))
	modelSet, _, err := c.ModelSets.Get(ctx, model_set_id)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to find ModelSet", "model_set_id", err)
	}

	roleName := d.Get("name").(string)
	if permissionSet.Id == "" {
		return logErrDiag(ctx, diags, "Failed to find permission set", "permission_set_id", err)
	}

	role := lookergo.Role{
		Name:            roleName,
		PermissionSetID: permissionSet.Id,
		PermissionSet:   permissionSet,
		ModelSetID:      modelSet.Id,
	}

	newRole, _, err := c.Roles.Create(ctx, &role)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to create Role", "err", err)
	}
	d.Set("name", newRole.Name)
	d.Set("permission_set_id", newRole.PermissionSet.Id)
	d.Set("permission_set_name", newRole.PermissionSet.Name)
	d.Set("model_set_id", newRole.ModelSet.Id)

	d.SetId(idAsString(newRole.Id))

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	var permissionSet lookergo.PermissionSet
	if psId, ok := d.GetOk("permission_set_id"); ok {
		perm, _, err := c.PermissionSets.Get(ctx, strconv.Itoa(psId.(int)))
		if err != nil {
			return logErrDiag(ctx, diags, "PermissionSet not found", "permission_set_id", psId)
		}
		permissionSet.Id = perm.Id
	} else if psName, ok := d.GetOk("permission_set_name"); ok {
		permissions, _, err := c.PermissionSets.GetByName(ctx, psName.(string), nil)
		if err != nil {
			return logErrDiag(ctx, diags, "Failed to query permission sets", "err", err)
		}
		for _, permission := range permissions {
			if permission.Name == psName.(string) {
				permissionSet.Id = permission.Id
				permissionSet.Name = permission.Name
			}
		}
	}

	role := lookergo.Role{
		Name:            d.Get("name").(string),
		PermissionSetID: permissionSet.Id,
		ModelSetID:      strconv.Itoa(d.Get("model_set_id").(int)),
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
