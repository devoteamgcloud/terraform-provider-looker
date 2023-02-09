package provider

import (
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

var (
	dataUserKey = []string{
		"id",
		"email",
	}
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: dataUserKey,
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: dataUserKey,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
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
	user := lookergo.User{}

	if userId, exists := d.GetOk("id"); exists {
		localUser, _, err := c.Users.Get(ctx, userId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		user = lookergo.User(*localUser)
	}

	if email, exists := d.GetOk("email"); exists {
		localUser, _, err := c.Users.ListByEmail(ctx, email.(string), &lookergo.ListOptions{})
		if err != nil {
			return diag.FromErr(err)
		}
		user = lookergo.User(localUser[0])
	}

	if err := d.Set("first_name", user.FirstName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}

	if user.RoleIds != nil {
		err := d.Set("roles", user.RoleIds.ToSliceOfStrings())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if user.CredentialsEmail != nil {
		if err := d.Set("email", user.CredentialsEmail.Email); err != nil {
			return diag.FromErr(err)
		}
	} else if user.CredentialsSaml != nil {
		if err := d.Set("email", user.CredentialsSaml.Email); err != nil {
			return diag.FromErr(err)
		}
	} else if err := d.Set("email", ""); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.Id)

	return diags
}
