package provider

import (
	"context"
	"fmt"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	//"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlerts() *schema.Resource {
	return &schema.Resource{
		Description: `
`,
		CreateContext: resourceAlertsCreate,
		ReadContext:   resourceAlertsRead,
		UpdateContext: resourceAlertsUpdate,
		DeleteContext: resourceAlertsDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"cron": {
				Type: schema.TypeString,
				Optional: true,
			},
			"dashboard_element_id": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"comparison_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: func(val any, key string)(warns []string, errs []error) {
					v := val.(string)
					value := lookergo.ComparisonType(v)
					switch value {
						case lookergo.ComparisonType_EQUAL_TO, lookergo.ComparisonType_GREATER_THAN, lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO, lookergo.ComparisonType_LESS_THAN,
						lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO,lookergo.ComparisonType_INCREASES_BY,lookergo.ComparisonType_DECREASES_BY,lookergo.ComparisonType_CHANGES_BY:
							return nil
						}
					errs = append(errs, fmt.Errorf("comparison type must be a supported value, please refer to Looker documentation for more information"))
					return
				},
			},
			"applied_dashboard_filters": {
				Type: schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"field_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filter_description": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			
			"custom_title": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"followable": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"is_disabled": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"is_public": {
				Type: schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAlertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	//c := m.(*Config).Api // .(*lookergo.Client)
	d.SetId("OK")
	return resourceFolderRead(ctx, d, m)
}

func resourceAlertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	//c := m.(*Config).Api // .(*lookergo.Client)
	d.SetId("OK")

	return diags
}

func resourceAlertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	//c := m.(*Config).Api // .(*lookergo.Client)
	//FolderID := d.Id()
	return resourceFolderRead(ctx, d, m)
}

func resourceAlertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	//c := m.(*Config).Api // .(*lookergo.Client)

	// Finally mark as deleted
	d.SetId("")

	return diags
}
