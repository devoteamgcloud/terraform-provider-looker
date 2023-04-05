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
				Required: true,
			},
			"dashboard_element_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comparison_type": {
				Type:     schema.TypeString,
				Required: true,
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
				Required: true,
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
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Optional: true,
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
			"owner_id": {
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
			"treshold": {
				Type:     schema.TypeFloat,
				Required: true,
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
		alert.Description = castToPtr(value.(string))
	}

	if value, ok := d.GetOk("cron"); ok {
		alert.Cron = value.(string)
	}

	if value, ok := d.GetOk("dashboard_element_id"); ok {
		alert.DashboardElementId = castToPtr(value.(string))
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
			switch dest_type {
			case string(lookergo.DestinationType_EMAIL):
				dest.DestinationType = lookergo.DestinationType_EMAIL
				val, ok := obj["email_address"]
				if ok {
					if val == "" {
						return diag.Errorf("email_address is required for destination_type EMAIL")
					}
					dest.EmailAddress = castToPtr(val.(string))
				}
			case string(lookergo.DestinationType_ACTION_HUB):
				dest.DestinationType = lookergo.DestinationType_ACTION_HUB
				if val, ok := obj["action_hub_integration_id"]; ok {
					dest.ActionHubIntegrationId = castToPtr(val.(string))
				} else {
					return diag.Errorf("action_hub_integration_id is required for destination_type ACTION_HUB")
				}
				if val, ok := obj["action_hub_form_params_json"]; ok {
					dest.ActionHubFormParamsJson = castToPtr(val.(string))
				}
			}
			dests = append(dests, dest)
		}
		alert.Destinations = dests
	}

	if value, ok := d.GetOk("comparison_type"); ok {
		s := value.(string)
		switch s {
		case string(lookergo.ComparisonType_EQUAL_TO):
			alert.ComparisonType = lookergo.ComparisonType_EQUAL_TO
		case string(lookergo.ComparisonType_GREATER_THAN):
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN
		case string(lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO):
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO
		case string(lookergo.ComparisonType_LESS_THAN):
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN
		case string(lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO):
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO
		case string(lookergo.ComparisonType_INCREASES_BY):
			alert.ComparisonType = lookergo.ComparisonType_INCREASES_BY
		case string(lookergo.ComparisonType_DECREASES_BY):
			alert.ComparisonType = lookergo.ComparisonType_DECREASES_BY
		case string(lookergo.ComparisonType_CHANGES_BY):
			alert.ComparisonType = lookergo.ComparisonType_CHANGES_BY
		default:
			return diag.Errorf("comparison type must respect the naming conventions")
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
						filter.FilterValue = castToPtr(val.(string))
					}
				}
			}
			alert.Field = field
		}
	}

	if value, ok := d.GetOk("custom_title"); ok {
		alert.CustomTitle = castToPtr(value.(string))
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
	if value, ok := d.GetOk("treshold"); ok {
		alert.Threshold = value.(float64)
	}
	if value, ok := d.GetOk("owner_id"); ok {
		alert.OwnerId = value.(string)
	}
	new_alert, _, err := c.Alerts.Create(ctx, alert)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*new_alert.Id)
	return resourceAlertsRead(ctx, d, m)
}

func resourceAlertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Id()
	alert, _, err := c.Alerts.Get(ctx, ID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("description", alert.Description)
	d.Set("cron", alert.Cron)
	d.Set("dashboard_element_id", alert.DashboardElementId)
	flattenAppliedDashboard := func(alertApp *[]lookergo.AlertAppliedDashboardFilter) []interface{} {
		appDashboards := make([]interface{}, len(*alertApp))
		for i, elem := range *alertApp {
			dashboard := make(map[string]interface{})
			dashboard["filter_title"] = elem.FilterTitle
			dashboard["field_name"] = elem.FieldName
			dashboard["filter_value"] = elem.FilterValue
			dashboard["filter_description"] = *elem.FilterDescription
			appDashboards[i] = dashboard
		}
		return appDashboards
	}
	appDashboards := flattenAppliedDashboard(alert.AppliedDashboardFilters)
	if err = d.Set("applied_dashboard_filters", appDashboards); err != nil {
		return diag.FromErr(err)
	}
	flattenDestinations := func(destinations *[]lookergo.AlertDestination) []interface{} {
		appDashboards := make([]interface{}, len(*destinations))
		for i, elem := range *destinations {
			dashboard := make(map[string]interface{})
			dashboard["destination_type"] = elem.DestinationType
			if elem.EmailAddress != nil {
				dashboard["email_address"] = elem.EmailAddress
			}
			if elem.ActionHubIntegrationId != nil {
				dashboard["action_hub_integration_id"] = elem.ActionHubIntegrationId
			}
			if elem.ActionHubFormParamsJson != nil {
				dashboard["action_hub_form_params_json"] = *elem.ActionHubFormParamsJson
			}
			appDashboards[i] = dashboard
		}
		return appDashboards
	}
	destinations := flattenDestinations(&alert.Destinations)
	if err = d.Set("destinations", destinations); err != nil {
		return diag.FromErr(err)
	}
	flattenField := func(field *lookergo.AlertField) interface{} {
		Dashboards := make([]interface{}, 1)
		dashboard := make(map[string]interface{})
		dashboard["title"] = field.Title
		dashboard["name"] = field.Name
		filters := make([]interface{}, len(*field.Filter))
		for i, elem := range *field.Filter {
			filter := make(map[string]interface{})
			filter["field_name"] = elem.FieldName
			filter["field_value"] = elem.FieldValue
			if elem.FilterValue != nil {
				filter["filter_value"] = elem.FilterValue
			}
			filters[i] = filter
		}
		dashboard["filter"] = filters
		Dashboards[0] = dashboard
		return Dashboards
	}

	field := flattenField(&alert.Field)
	if err = d.Set("field", field); err != nil {
		return diag.FromErr(err)
	}
	d.Set("comparison_type", alert.ComparisonType)
	d.Set("custom_title", alert.CustomTitle)
	d.Set("followable", alert.Followable)
	d.Set("is_disabled", alert.IsDisabled)
	d.Set("is_public", alert.IsPublic)
	return diags
}

func resourceAlertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	ID := d.Id()
	alert := &lookergo.Alert{}
	if value, ok := d.GetOk("description"); ok {
		alert.Description = castToPtr(value.(string))
	}

	if value, ok := d.GetOk("cron"); ok {
		alert.Cron = value.(string)
	}

	if value, ok := d.GetOk("dashboard_element_id"); ok {
		alert.DashboardElementId = castToPtr(value.(string))
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
					dest.EmailAddress = castToPtr(val.(string))
				} else {
					return diag.Errorf("email_address is required for destination_type EMAIL")
				}
			case lookergo.DestinationType_ACTION_HUB:
				dest.DestinationType = lookergo.DestinationType_ACTION_HUB
				if val, ok := obj["action_hub_integration_id"]; ok {
					dest.ActionHubIntegrationId = castToPtr(val.(string))
				} else {
					return diag.Errorf("action_hub_integration_id is required for destination_type ACTION_HUB")
				}
				if val, ok := obj["action_hub_form_params_json"]; ok {
					dest.ActionHubFormParamsJson = castToPtr(val.(string))
				}
			}
			dests = append(dests, dest)
		}
		alert.Destinations = dests
	}

	if value, ok := d.GetOk("comparison_type"); ok {
		s := value.(string)
		switch s {
		case string(lookergo.ComparisonType_EQUAL_TO):
			alert.ComparisonType = lookergo.ComparisonType_EQUAL_TO
		case string(lookergo.ComparisonType_GREATER_THAN):
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN
		case string(lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO):
			alert.ComparisonType = lookergo.ComparisonType_GREATER_THAN_OR_EQUAL_TO
		case string(lookergo.ComparisonType_LESS_THAN):
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN
		case string(lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO):
			alert.ComparisonType = lookergo.ComparisonType_LESS_THAN_OR_EQUAL_TO
		case string(lookergo.ComparisonType_INCREASES_BY):
			alert.ComparisonType = lookergo.ComparisonType_INCREASES_BY
		case string(lookergo.ComparisonType_DECREASES_BY):
			alert.ComparisonType = lookergo.ComparisonType_DECREASES_BY
		case string(lookergo.ComparisonType_CHANGES_BY):
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
						filter.FilterValue = castToPtr(val.(string))
					}
				}
			}
			alert.Field = field
		}
	}

	if value, ok := d.GetOk("custom_title"); ok {
		alert.CustomTitle = castToPtr(value.(string))
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
	if value, ok := d.GetOk("treshold"); ok {
		alert.Threshold = value.(float64)
	}
	if value, ok := d.GetOk("owner_id"); ok {
		alert.OwnerId = value.(string)
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
