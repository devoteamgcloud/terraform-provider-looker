package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func dataSourcePublicEgressIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicEgressIpsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Egress IP addresses.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePublicEgressIpsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	egressPublicIps, _, err := c.EgressIpAddresses.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("public_ips", egressPublicIps.EgressIpAddresses); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("-")
	return diags
}
