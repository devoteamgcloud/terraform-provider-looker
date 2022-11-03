package provider

import (
	"context"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMemberCreate,
		ReadContext:   resourceGroupMemberRead,
		UpdateContext: resourceGroupMemberUpdate,
		DeleteContext: resourceGroupMemberDelete,
		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"group": {
				Type:     schema.TypeSet,
				Optional: true,
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

func parentGroup(ctx context.Context, d *schema.ResourceData, c *lookergo.Client) (*lookergo.Group, error) {
	tflog.Info(ctx, "Verifying parent group.")
	group, _, err := c.Groups.Get(ctx, idAsInt(d.Get("target_group_id").(string)))
	if err != nil {
		return nil, err
	} else {
		tflog.Info(ctx, "Found group with", map[string]interface{}{"id": group.Id, "name": group.Name})
		return group, nil
	}
}

func resourceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	pg, err := parentGroup(ctx, d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	userSet, ok := d.GetOk("user")
	if ok {
		userItems := make([]interface{}, userSet.(*schema.Set).Len())
		for i, raw := range userSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			val := obj["id"].(string)
			tflog.Info(ctx, "Add user", map[string]interface{}{"id": val})

			memberUser, _, err := c.Groups.AddMemberUser(ctx, idAsInt(pg.Id), idAsInt(val))
			if err != nil {
				return diag.FromErr(err)
			}
			userItems[i] = map[string]interface{}{"id": idAsString(memberUser.Id), "first_name": memberUser.FirstName, "last_name": memberUser.LastName}
		}
		d.Set("user", userItems)
	}

	groupSet, ok := d.GetOk("group")
	if ok {
		groupItems := make([]interface{}, groupSet.(*schema.Set).Len())
		for i, raw := range groupSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			val := obj["id"].(string)
			tflog.Info(ctx, "Add group", map[string]interface{}{"id": val})

			memberGroup, _, err := c.Groups.AddMemberGroup(ctx, idAsInt(pg.Id), idAsInt(val))
			if err != nil {
				return diag.FromErr(err)
			}
			groupItems[i] = map[string]interface{}{"id": idAsString(memberGroup.Id), "name": memberGroup.Name}
		}
		d.Set("group", groupItems)
	}
	d.SetId("-")
	return resourceGroupMemberRead(ctx, d, m)
}

func resourceGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	pg, err := parentGroup(ctx, d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "Read group members for", map[string]interface{}{"target_group_id": pg.Id})

	userSet, ok := d.GetOk("user")
	if ok {
		memberUsers, _, err := c.Groups.ListMemberUsers(ctx, pg.Id, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		var userItems []interface{}
		for _, raw := range userSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})

			for _, user := range memberUsers {
				if idAsString(user.Id) == obj["id"].(string) {
					userItems = append(userItems, map[string]interface{}{"id": idAsString(user.Id), "first_name": user.FirstName, "last_name": user.LastName})
				}
			}
		}
		d.Set("user", userItems)
	}

	groupSet, ok := d.GetOk("group")
	if ok {
		memberGroups, _, err := c.Groups.ListMemberGroups(ctx, pg.Id, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		var groupItems []interface{}
		for _, raw := range groupSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})

			for _, group := range memberGroups {
				if idAsString(group.Id) == obj["id"].(string) {
					groupItems = append(groupItems, map[string]interface{}{"id": idAsString(group.Id), "name": group.Name})
				}
			}
		}
		d.Set("group", groupItems)
	}

	return diags
}

func oldNewDiff(old, desired []string) (remove, create []string) {
	for _, item := range old { // remove
		if !contains(desired, item) {
			remove = append(remove, item)
		}
	}
	for _, item := range desired { // add
		if !contains(old, item) {
			create = append(create, item)
		}
	}
	return
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func idsFromMysteryInterface(m []interface{}) (ids []string) {
	for _, item := range m {
		ids = append(ids, item.(map[string]interface{})["id"].(string))
	}
	return
}

func resourceGroupMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	pg, err := parentGroup(ctx, d, c)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "Update group members for", map[string]interface{}{"target_group_id": pg.Id})

	if d.HasChange("user") {
		oldRaw, newRaw := d.GetChange("user")
		remove, create := oldNewDiff(
			idsFromMysteryInterface(oldRaw.(*schema.Set).List()),
			idsFromMysteryInterface(newRaw.(*schema.Set).List()))

		for _, item := range remove {
			_, err := c.Groups.RemoveMemberUser(ctx, idAsInt(pg.Id), idAsInt(item))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, item := range create {
			_, _, err := c.Groups.AddMemberUser(ctx, idAsInt(pg.Id), idAsInt(item))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("group") {
		oldRaw, newRaw := d.GetChange("group")
		remove, create := oldNewDiff(
			idsFromMysteryInterface(oldRaw.(*schema.Set).List()),
			idsFromMysteryInterface(newRaw.(*schema.Set).List()))

		for _, item := range remove {
			_, err := c.Groups.RemoveMemberGroup(ctx, idAsInt(pg.Id), idAsInt(item))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, item := range create {
			_, _, err := c.Groups.AddMemberGroup(ctx, idAsInt(pg.Id), idAsInt(item))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceGroupMemberRead(ctx, d, m)
}

func resourceGroupMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	pg, err := parentGroup(ctx, d, c)
	if err != nil {
		return diag.FromErr(err)
	}

	_, ok := d.GetOk("user")
	if ok {
		// override state with actual config
		userSet, _ := d.GetChange("user")
		for _, raw := range userSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			val := obj["id"].(string)
			tflog.Info(ctx, "Remove user from group", map[string]interface{}{"id": val})

			_, err := c.Groups.RemoveMemberUser(ctx, idAsInt(pg.Id), idAsInt(val))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	_, ok = d.GetOk("group")
	if ok {
		// override state with actual config
		groupSet, _ := d.GetChange("group")
		for _, raw := range groupSet.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			val := obj["id"].(string)
			tflog.Info(ctx, "Remove group from group", map[string]interface{}{"id": val})

			_, err := c.Groups.RemoveMemberGroup(ctx, idAsInt(pg.Id), idAsInt(val))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Finally mark as deleted
	d.SetId("")

	return diags
}
