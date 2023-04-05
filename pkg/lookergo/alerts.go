package lookergo

import (
	"context"
)

const AlertsBasePath = "4.0/alerts"

type AlertsResourceOp struct {
	client *Client
}

var _ AlertsResource = &AlertsResourceOp{}

type AlertsResource interface {
	//List(context.Context, *ListOptions) ([]Alert, *Response, error)
	Get(context.Context, string) (*Alert, *Response, error)
	//Get(context.Context,*ListOptions, string) ([]Alert, *Response, error)
	Create(context.Context, *Alert) (*Alert, *Response, error)
	Update(context.Context, string, *Alert) (*Alert, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

type DestinationType string

const (
	DestinationType_EMAIL      DestinationType = "EMAIL"
	DestinationType_ACTION_HUB DestinationType = "ACTION_HUB"
)

type ComparisonType string

const (
	ComparisonType_EQUAL_TO                 ComparisonType = "EQUAL_TO"
	ComparisonType_GREATER_THAN             ComparisonType = "GREATER_THAN"
	ComparisonType_GREATER_THAN_OR_EQUAL_TO ComparisonType = "GREATER_THAN_OR_EQUAL_TO"
	ComparisonType_LESS_THAN                ComparisonType = "LESS_THAN"
	ComparisonType_LESS_THAN_OR_EQUAL_TO    ComparisonType = "LESS_THAN_OR_EQUAL_TO"
	ComparisonType_INCREASES_BY             ComparisonType = "INCREASES_BY"
	ComparisonType_DECREASES_BY             ComparisonType = "DECREASES_BY"
	ComparisonType_CHANGES_BY               ComparisonType = "CHANGES_BY"
)

type AlertConditionState struct {
	PreviousTimeSeriesId *string `json:"previous_time_series_id,omitempty"` // (Write-Only) The second latest time string the alert has seen.
	LatestTimeSeriesId   *string `json:"latest_time_series_id,omitempty"`   // (Write-Only) Latest time string the alert has seen.
}

type AlertAppliedDashboardFilter struct {
	FilterTitle       string  `json:"filter_title"`                 // Field Title. Refer to `DashboardFilter.title` in [DashboardFilter](#!/types/DashboardFilter). Example `Name`
	FieldName         string  `json:"field_name"`                   // Field Name. Refer to `DashboardFilter.dimension` in [DashboardFilter](#!/types/DashboardFilter). Example `distribution_centers.name`
	FilterValue       string  `json:"filter_value"`                 // Field Value. [Filter Expressions](https://cloud.google.com/looker/docs/reference/filter-expressions). Example `Los Angeles CA`
	FilterDescription *string `json:"filter_description,omitempty"` // Human Readable Filter Description. This may be null or auto-generated. Example `is Los Angeles CA`
}

type AlertDestination struct {
	DestinationType         DestinationType `json:"destination_type"`                      // Type of destination that the alert will be sent to Valid values are: "EMAIL", "ACTION_HUB".
	EmailAddress            *string         `json:"email_address,omitempty"`               // Email address for the 'email' type
	ActionHubIntegrationId  *string         `json:"action_hub_integration_id,omitempty"`   // Action hub integration id for the 'action_hub' type. [Integration](#!/types/Integration)
	ActionHubFormParamsJson *string         `json:"action_hub_form_params_json,omitempty"` // Action hub form params json for the 'action_hub' type [IntegrationParam](#!/types/IntegrationParam)
}

type AlertFieldFilter struct {
	FieldName   string  `json:"field_name"`             // Field Name. Has format `<view>.<field>`
	FieldValue  string  `json:"field_value"`            // Field Value. Depends on the type of field - numeric or string. For [location](https://cloud.google.com/looker/docs/reference/field-reference/dimension-type-reference#location) type, it's a list of floats. Example `[1.0, 56.0]`
	FilterValue *string `json:"filter_value,omitempty"` // Filter Value. Usually null except for [location](https://cloud.google.com/looker/docs/reference/field-reference/dimension-type-reference#location) type. It'll be a string of lat,long ie `'1.0,56.0'`
}

type AlertField struct {
	Title  string              `json:"title"`            // Field's title. Usually auto-generated to reflect field name and its filters
	Name   string              `json:"name"`             // Field's name. Has the format `<view>.<field>` Refer to [docs](https://cloud.google.com/looker/docs/sharing-and-publishing/creating-alerts) for more details
	Filter *[]AlertFieldFilter `json:"filter,omitempty"` // (Optional / Advance Use) List of fields filter. This further restricts the alert to certain dashboard element's field values. This can be used on top of dashboard filters `applied_dashboard_filters`. To keep thing simple, it's suggested to just use dashboard filters. Example: `{ 'title': '12 Number on Hand', 'name': 'inventory_items.number_on_hand', 'filter': [{ 'field_name': 'inventory_items.id', 'field_value': 12, 'filter_value': null }] }`
}

type InvestigativeContentType string

const InvestigativeContentType_Dashboard InvestigativeContentType = "dashboard"

type Alert struct {
	AppliedDashboardFilters   *[]AlertAppliedDashboardFilter `json:"applied_dashboard_filters,omitempty"` // Filters coming from the dashboard that are applied. Example `[{ "filter_title": "Name", "field_name": "distribution_centers.name", "filter_value": "Los Angeles CA" }]`
	ComparisonType            ComparisonType                 `json:"comparison_type"`                     // This property informs the check what kind of comparison we are performing. Only certain condition types are valid for time series alerts. For details, refer to [Setting Alert Conditions](https://cloud.google.com/looker/docs/sharing-and-publishing/creating-alerts#setting_alert_conditions) Valid values are: "EQUAL_TO", "GREATER_THAN", "GREATER_THAN_OR_EQUAL_TO", "LESS_THAN", "LESS_THAN_OR_EQUAL_TO", "INCREASES_BY", "DECREASES_BY", "CHANGES_BY".
	Cron                      string                         `json:"cron"`                                // Vixie-Style crontab specification when to run. At minumum, it has to be longer than 15 minute intervals
	CustomUrlBase             *string                        `json:"custom_url_base,omitempty"`           // Domain for the custom url selected by the alert creator from the admin defined domain allowlist
	CustomUrlParams           *string                        `json:"custom_url_params,omitempty"`         // Parameters and path for the custom url defined by the alert creator
	CustomUrlLabel            *string                        `json:"custom_url_label,omitempty"`          // Label for the custom url defined by the alert creator
	ShowCustomUrl             *bool                          `json:"show_custom_url,omitempty"`           // Boolean to determine if the custom url should be used
	CustomTitle               *string                        `json:"custom_title,omitempty"`              // An optional, user-defined title for the alert
	DashboardElementId        *string                        `json:"dashboard_element_id,omitempty"`      // ID of the dashboard element associated with the alert. Refer to [dashboard_element()](#!/Dashboard/DashboardElement)
	Description               *string                        `json:"description,omitempty"`               // An optional description for the alert. This supplements the title
	Destinations              []AlertDestination             `json:"destinations"`                        // Array of destinations to send alerts to. Must be the same type of destination. Example `[{ "destination_type": "EMAIL", "email_address": "test@test.com" }]`
	Field                     AlertField                     `json:"field"`
	Followed                  *bool                          `json:"followed,omitempty"`                    // Whether or not the user follows this alert.
	Followable                *bool                          `json:"followable,omitempty"`                  // Whether or not the alert is followable
	Id                        *string                        `json:"id,omitempty"`                          // ID of the alert
	IsDisabled                *bool                          `json:"is_disabled,omitempty"`                 // Whether or not the alert is disabled
	DisabledReason            *string                        `json:"disabled_reason,omitempty"`             // Reason for disabling alert
	IsPublic                  *bool                          `json:"is_public,omitempty"`                   // Whether or not the alert is public
	InvestigativeContentType  *InvestigativeContentType      `json:"investigative_content_type,omitempty"`  // The type of the investigative content Valid values are: "dashboard".
	InvestigativeContentId    *string                        `json:"investigative_content_id,omitempty"`    // The ID of the investigative content. For dashboards, this will be the dashboard ID
	InvestigativeContentTitle *string                        `json:"investigative_content_title,omitempty"` // The title of the investigative content.
	LookmlDashboardId         *string                        `json:"lookml_dashboard_id,omitempty"`         // ID of the LookML dashboard associated with the alert
	LookmlLinkId              *string                        `json:"lookml_link_id,omitempty"`              // ID of the LookML dashboard element associated with the alert
	OwnerId                   string                         `json:"owner_id"`                              // User id of alert owner
	OwnerDisplayName          *string                        `json:"owner_display_name,omitempty"`          // Alert owner's display name
	Threshold                 float64                        `json:"threshold"`                             // Value of the alert threshold
	TimeSeriesConditionState  *AlertConditionState           `json:"time_series_condition_state,omitempty"`
}

func (s *AlertsResourceOp) Get(ctx context.Context, AlertId string) (*Alert, *Response, error) {
	return doGetById(ctx, s.client, AlertsBasePath, AlertId, new(Alert))
}

func (s *AlertsResourceOp) Create(ctx context.Context, requestAlert *Alert) (*Alert, *Response, error) {
	return doCreate(ctx, s.client, AlertsBasePath, requestAlert, new(Alert))
}

func (s *AlertsResourceOp) Update(ctx context.Context, AlertId string, requestAlert *Alert) (*Alert, *Response, error) {
	return doUpdate(ctx, s.client, AlertsBasePath, AlertId, requestAlert, new(Alert))
}

func (s *AlertsResourceOp) Delete(ctx context.Context, AlertId string) (*Response, error) {
	return doDelete(ctx, s.client, AlertsBasePath, AlertId)
}
