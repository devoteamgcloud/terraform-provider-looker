package provider

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)

	userId, _ := strconv.Atoi(d.Get("id").(string))

	user, r, err := c.Users.Get(ctx, userId)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error fetching user",
			Detail:   "Unable to find user " + d.Get("user_id").(string) + " ; HTTP Status: " + r.Response.Status,
		})

		return diag.FromErr(err)
	}

	if err := d.Set("first_name", &user.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}

	if user.RoleIds != nil {
		err = d.Set("roles", user.RoleIds.ToSliceOfStrings())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if user.CredentialEmail != nil {
		if err := d.Set("email", user.CredentialEmail.Email); err != nil {
			return diag.FromErr(err)
		}
	} else if err := d.Set("email", ""); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(user.Id))

	return diags
}
