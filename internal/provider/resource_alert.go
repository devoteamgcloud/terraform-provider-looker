package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dashboard_element_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comparison_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					value := lookergo.ComparisonType(v)
					switch value {
					case lookergo.ComparisonType_EQUAL_TO, lookergo.ComparisonType_GREATER_THAN, lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO, lookergo.ComparisonType_LESS_THAN,
						lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO, lookergo.ComparisonType_INCREASES_BY, lookergo.ComparisonType_DECREASES_BY, lookergo.ComparisonType_CHANGES_BY:
					default:
						errs = append(errs, fmt.Errorf("comparison type must be a supported value, please refer to Looker documentation for more information"))
					}
					return warns, errs
				},
			},
			"applied_dashboard_filters": {
				Type:     schema.TypeSet,
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
							Optional: true,
						},
					},
				},
			},
			"destinations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								v := val.(string)
								value := lookergo.DestinationType(v)
								switch value {
								case lookergo.DestinationType_EMAIL, lookergo.DestinationType_ACTION_HUB:
									return
								default:
									errs = append(errs, fmt.Errorf("destination type must be a supported value, please refer to Looker documentation for more information"))
								}
								return warns, errs
							},
						},
						"email_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_hub_integration_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action_hub_form_params_json": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								v := val.(string)
								parts := strings.Split(v, ".")
								if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
									return
								}
								errs = append(errs, fmt.Errorf("alert field name must be in following format <view>.<field>, please refer to Looker documentation for more information"))
								return nil, errs
							},
						},
						"filter": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"field_value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"filter_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"custom_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"followable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAlertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	alert := &lookergo.Alert{}

	if value, ok := d.GetOk("description"); ok {
		alert.Description = value.(*string)
	}

	if value, ok := d.GetOk("cron"); ok {
		alert.Cron = value.(string)
	}

	if value, ok := d.GetOk("dashboard_element_id"); ok {
		alert.Id = value.(*string)
	}

	if value, ok := d.GetOk("applied_dashboard_filters"); ok {
		filters := []lookergo.AlertAppliedDashboardFilter{}
		for _, raw := range value.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			filter := lookergo.AlertAppliedDashboardFilter{}
			filter.FilterTitle = obj["filter_title"].(string)
			filter.FilterDescription = obj["filter_description"].(*string)
			filter.FilterValue = obj["filter_value"].(string)
			filter.FieldName = obj["field_name"].(string)
			filters = append(filters, filter)
		}
		alert.AppliedDashboardFilters = &filters
	}

	if value, ok := d.GetOk("destinations"); ok {
		dests := []lookergo.AlertDestination{}
		for _, raw := range value.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			dest := lookergo.AlertDestination{}
			dest_type := obj["destination_type"].(string)
			value := lookergo.DestinationType(dest_type)
			switch value {
			case lookergo.DestinationType_EMAIL:
				dest.DestinationType = lookergo.DestinationType_EMAIL
				val, ok := obj["email_address"]
				if ok {
					dest.EmailAddress = val.(*string)
				} else {
					return diag.Errorf("email_address is required for destination_type EMAIL")
				}
			case lookergo.DestinationType_ACTION_HUB:
				dest.DestinationType = lookergo.DestinationType_ACTION_HUB
				if val, ok := obj["action_hub_integration_id"]; ok {
					dest.ActionHubIntegrationId = val.(*string)
				} else {
					return diag.Errorf("action_hub_integration_id is required for destination_type ACTION_HUB")
				}
				if val, ok := obj["action_hub_form_params_json"]; ok {
					dest.ActionHubFormParamsJson = val.(*string)
				}
			}
			dests = append(dests, dest)
		}
		alert.Destinations = dests
	}

	if value, ok := d.GetOk("comparison_type"); ok {
		s := value.(string)
		switch lookergo.ComparisonType(s) {
		case lookergo.ComparisonType_EQUAL_TO:
			alert.ComparisonType = lookergo.ComparisonType_EQUAL_TO
		case lookergo.ComparisonType_GREATER_THAN:
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN
		case lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO:
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO
		case lookergo.ComparisonType_LESS_THAN:
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN
		case lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO:
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO
		case lookergo.ComparisonType_INCREASES_BY:
			alert.ComparisonType = lookergo.ComparisonType_INCREASES_BY
		case lookergo.ComparisonType_DECREASES_BY:
			alert.ComparisonType = lookergo.ComparisonType_DECREASES_BY
		case lookergo.ComparisonType_CHANGES_BY:
			alert.ComparisonType = lookergo.ComparisonType_CHANGES_BY
		}
	}

	if value, ok := d.GetOk("field"); ok {
		if len(value.(*schema.Set).List()) > 1 {
			return diag.Errorf("maximum one field object can be defined")
		}
		for _, raw := range value.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			field := lookergo.AlertField{}
			field.Title = obj["title"].(string)
			field.Name = obj["name"].(string)
			if value, ok := obj["filter"]; ok {
				if len(value.(*schema.Set).List()) > 1 {
					return diag.Errorf("maximum one filter object can be defined")
				}
				for _, raw := range value.(*schema.Set).List() {
					obj := raw.(map[string]interface{})
					filter := lookergo.AlertFieldFilter{}
					filter.FieldName = obj["field_name"].(string)
					filter.FieldValue = obj["field_value"].(string)
					if val, ok := obj["filter_value"]; ok {
						filter.FilterValue = val.(*string)
					}
				}
			}
		}
	}

	if value, ok := d.GetOk("custom_title"); ok {
		alert.CustomTitle = value.(*string)
	}

	if value, ok := d.GetOk("followable"); ok {
		alert.Followable = boolPtr(value.(bool))
	}

	if value, ok := d.GetOk("is_disabled"); ok {
		alert.IsDisabled = boolPtr(value.(bool))
	}

	if value, ok := d.GetOk("is_public"); ok {
		alert.IsPublic = boolPtr(value.(bool))
	}
	new_alert, _, err := c.Alerts.Create(ctx, alert)
	if err != nil {
		diag.FromErr(err)
	}
	d.SetId(*new_alert.Id)
	return resourceFolderRead(ctx, d, m)
}

func resourceAlertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Id()
	_, _, err := c.Alerts.Get(ctx, ID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("OK")

	return diags
}

func resourceAlertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Id()
	alert := &lookergo.Alert{}
	if value, ok := d.GetOk("description"); ok {
		alert.Description = value.(*string)
	}

	if value, ok := d.GetOk("cron"); ok {
		alert.Cron = value.(string)
	}

	if value, ok := d.GetOk("dashboard_element_id"); ok {
		alert.Id = value.(*string)
	}

	if value, ok := d.GetOk("applied_dashboard_filters"); ok {
		filters := []lookergo.AlertAppliedDashboardFilter{}
		for _, raw := range value.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			filter := lookergo.AlertAppliedDashboardFilter{}
			filter.FilterTitle = obj["filter_title"].(string)
			filter.FilterDescription = obj["filter_description"].(*string)
			filter.FilterValue = obj["filter_value"].(string)
			filter.FieldName = obj["field_name"].(string)
			filters = append(filters, filter)
		}
		alert.AppliedDashboardFilters = &filters
	}

	if value, ok := d.GetOk("destinations"); ok {
		dests := []lookergo.AlertDestination{}
		for _, raw := range value.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			dest := lookergo.AlertDestination{}
			dest_type := obj["destination_type"].(string)
			value := lookergo.DestinationType(dest_type)
			switch value {
			case lookergo.DestinationType_EMAIL:
				dest.DestinationType = lookergo.DestinationType_EMAIL
				val, ok := obj["email_address"]
				if ok {
					dest.EmailAddress = val.(*string)
				} else {
					return diag.Errorf("email_address is required for destination_type EMAIL")
				}
			case lookergo.DestinationType_ACTION_HUB:
				dest.DestinationType = lookergo.DestinationType_ACTION_HUB
				if val, ok := obj["action_hub_integration_id"]; ok {
					dest.ActionHubIntegrationId = val.(*string)
				} else {
					return diag.Errorf("action_hub_integration_id is required for destination_type ACTION_HUB")
				}
				if val, ok := obj["action_hub_form_params_json"]; ok {
					dest.ActionHubFormParamsJson = val.(*string)
				}
			}
			dests = append(dests, dest)
		}
		alert.Destinations = dests
	}

	if value, ok := d.GetOk("comparison_type"); ok {
		s := value.(string)
		switch lookergo.ComparisonType(s) {
		case lookergo.ComparisonType_EQUAL_TO:
			alert.ComparisonType = lookergo.ComparisonType_EQUAL_TO
		case lookergo.ComparisonType_GREATER_THAN:
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN
		case lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO:
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO
		case lookergo.ComparisonType_LESS_THAN:
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN
		case lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO:
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO
		case lookergo.ComparisonType_INCREASES_BY:
			alert.ComparisonType = lookergo.ComparisonType_INCREASES_BY
		case lookergo.ComparisonType_DECREASES_BY:
			alert.ComparisonType = lookergo.ComparisonType_DECREASES_BY
		case lookergo.ComparisonType_CHANGES_BY:
			alert.ComparisonType = lookergo.ComparisonType_CHANGES_BY
		}
	}

	if value, ok := d.GetOk("field"); ok {
		if len(value.(*schema.Set).List()) > 1 {
			return diag.Errorf("maximum one field object can be defined")
		}
		for _, raw := range value.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			field := lookergo.AlertField{}
			field.Title = obj["title"].(string)
			field.Name = obj["name"].(string)
			if value, ok := obj["filter"]; ok {
				if len(value.(*schema.Set).List()) > 1 {
					return diag.Errorf("maximum one filter object can be defined")
				}
				for _, raw := range value.(*schema.Set).List() {
					obj := raw.(map[string]interface{})
					filter := lookergo.AlertFieldFilter{}
					filter.FieldName = obj["field_name"].(string)
					filter.FieldValue = obj["field_value"].(string)
					if val, ok := obj["filter_value"]; ok {
						filter.FilterValue = val.(*string)
					}
				}
			}
		}
	}

	if value, ok := d.GetOk("custom_title"); ok {
		alert.CustomTitle = value.(*string)
	}

	if value, ok := d.GetOk("followable"); ok {
		alert.Followable = boolPtr(value.(bool))
	}

	if value, ok := d.GetOk("is_disabled"); ok {
		alert.IsDisabled = boolPtr(value.(bool))
	}

	if value, ok := d.GetOk("is_public"); ok {
		alert.IsPublic = boolPtr(value.(bool))
	}
	c.Alerts.Update(ctx, ID, alert)
	return resourceAlertsRead(ctx, d, m)
}

func resourceAlertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Id()
	_, err := c.Alerts.Delete(ctx, ID)
	if err != nil {
		return diag.FromErr(err)
	}
	// Finally mark as deleted
	d.SetId("")
	return diags
}
