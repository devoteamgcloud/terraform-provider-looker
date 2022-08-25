package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	_ "golang.org/x/exp/slices"
	"strconv"
	"net/http"
)

func resourceRoleMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleMemberCreate,
		ReadContext:   resourceRoleMemberRead,
		UpdateContext: resourceRoleMemberUpdate,
		DeleteContext: resourceRoleMemberDelete,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRoleMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	roleMemberGroupsSet, ok := d.GetOk("group")
	if ok {
		// Fetch and verify existence
		roleMemberGroups, _, err := c.Roles.RoleGroupsList(ctx, idAsInt(d.Get("role_id")), nil)
		switch err.(type) {
		case *lookergo.ErrorResponse:
			if errResp := err.(*lookergo.ErrorResponse).Response; errResp.StatusCode == http.StatusNotFound {
				logTrace(ctx, "role not found", "role_id", d.Get("role_id").(string))
				d.SetId("")
				return // Resource was not found.
			}
		case error:
			return logErrDiag(ctx, diags, "unable to query role", "role_id", d.Get("role_id").(string)) // Connection error.
		default:
			logTrace(ctx, "role group members", "roleMemberGroups", roleMemberGroups)
		}

		// Flatten
		var groupItems []interface{}
		for _, raw := range roleMemberGroupsSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})

			for _, group := range roleMemberGroups {
				if idAsString(group.Id) == obj["id"].(string) {
					groupItems = append(groupItems, map[string]interface{}{"id": idAsString(group.Id), "name": group.Name})
				}
			}
		}
		d.Set("group", groupItems)
		d.SetId("-")
	} else {
		d.SetId("")
	}

	return diags
}

func resourceRoleMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	managedGroupIds := getSetIds(d, "group")

	var unmanagedGroupIds []string
	role_id, err := strconv.Atoi(d.Get("role_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	roleMemberGroups, _, err := c.Roles.RoleGroupsList(ctx, role_id, nil)
	if err == nil {
		for _, group := range roleMemberGroups {
			unmanagedGroupIds = append(unmanagedGroupIds, idAsString(group.Id))
		}
	}

	groupIds := append(managedGroupIds, unmanagedGroupIds...)
	slices.Sort(groupIds)
	groupIds = slices.Compact(groupIds)
	if _, _, err = c.Roles.RoleGroupsSet(ctx, role_id, groupIds); err != nil {
		return logErrDiag(ctx, diags, "Failed to update Role member Groups", "err", err)
	}

	d.SetId("-")
	return resourceRoleMemberRead(ctx, d, m)
}

func resourceRoleMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	role_id, err := strconv.Atoi(d.Get("role_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	var currentGroupIds []string
	roleMemberGroups, _, err := c.Roles.RoleGroupsList(ctx, role_id, nil)
	if err == nil {
		for _, group := range roleMemberGroups {
			currentGroupIds = append(currentGroupIds, idAsString(group.Id))
		}
	}

	_, _, oldIds, newIds := getSetChangeIdsDiff(d, "group")
	finalTemp := newIds
	for _, id := range currentGroupIds {
		if !slices.Contains(oldIds, id) {
			finalTemp = append(finalTemp, id)
		}
	}
	slices.Sort(finalTemp)
	finalIds := slices.Compact(finalTemp)
	
	_, _, err = c.Roles.RoleGroupsSet(ctx, role_id, finalIds)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to update Role member Groups", "err", err)
	}

	d.SetId("-")

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceRoleMemberRead(ctx, d, m)
}

func resourceRoleMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	role_id, err := strconv.Atoi(d.Get("role_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	var currentGroupIds []string
	roleMemberGroups, _, err := c.Roles.RoleGroupsList(ctx, role_id, nil)
	if err == nil {
		for _, group := range roleMemberGroups {
			currentGroupIds = append(currentGroupIds, idAsString(group.Id))
		}
	}

	_, _, oldIds, _ := getSetChangeIdsDiff(d, "group")
	var finalIds []string
	for _, id := range currentGroupIds {
		if !slices.Contains(oldIds, id) {
			finalIds = append(finalIds, id)
		}
	}

	_, _, err = c.Roles.RoleGroupsSet(ctx, role_id, finalIds)
	if err != nil {
		return logErrDiag(ctx, diags, "Failed to update Role member Groups", "err", err)
	}

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
