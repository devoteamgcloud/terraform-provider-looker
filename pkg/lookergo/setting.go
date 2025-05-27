package lookergo

import (
	"context"
)

const SettingBasePath = "4.0/setting"

type SettingResourceOp struct {
	client *Client
}

type SettingResource interface {
	Get(context.Context) (*Setting, *Response, error)
	Update(context.Context, *Setting) (*Setting, *Response, error)
}

var _ SettingResource = &SettingResourceOp{}

type Setting struct {
	InstanceConfig                      *InstanceConfig            `json:"instance_config,omitempty"`                         // Externally available instance configuration information (read-only)
	ExtensionFrameworkEnabled           *bool                      `json:"extension_framework_enabled,omitempty"`             // Toggle extension framework on or off
	ExtensionLoadUrlEnabled             *bool                      `json:"extension_load_url_enabled,omitempty"`              // (DEPRECATED) Toggle extension load url on or off. Do not use. This is temporary setting that will eventually become a noop and subsequently deleted
	MarketplaceAutoInstallEnabled       *bool                      `json:"marketplace_auto_install_enabled,omitempty"`        // (DEPRECATED) Toggle marketplace auto install on or off. Deprecated - do not use. Auto install can now be enabled via marketplace automation settings
	MarketplaceAutomation               *MarketplaceAutomation     `json:"marketplace_automation,omitempty"`                  // Marketplace automation settings
	MarketplaceEnabled                  *bool                      `json:"marketplace_enabled,omitempty"`                     // Toggle marketplace on or off
	MarketplaceSite                     *string                    `json:"marketplace_site,omitempty"`                        // Location of Looker marketplace CDN (read-only)
	MarketplaceTermsAccepted            *bool                      `json:"marketplace_terms_accepted,omitempty"`              // Accept marketplace terms by setting this value to true, or get the current status. Marketplace terms CANNOT be declined once accepted. Accepting marketplace terms automatically enables the marketplace. The marketplace can still be disabled after it has been enabled
	PrivatelabelConfiguration           *PrivatelabelConfiguration `json:"privatelabel_configuration,omitempty"`              // Private label configuration
	CustomWelcomeEmail                  *CustomWelcomeEmail        `json:"custom_welcome_email,omitempty"`                    // Custom welcome email configuration
	OnboardingEnabled                   *bool                      `json:"onboarding_enabled,omitempty"`                      // Toggle onboarding on or off
	Timezone                            *string                    `json:"timezone,omitempty"`                                // Change instance-wide default timezone
	AllowUserTimezones                  *bool                      `json:"allow_user_timezones,omitempty"`                    // Toggle user-specific timezones on or off
	DataConnectorDefaultEnabled         *bool                      `json:"data_connector_default_enabled,omitempty"`          // Toggle default future connectors on or off
	HostUrl                             *string                    `json:"host_url,omitempty"`                                // Change the base portion of your Looker instance URL setting
	OverrideWarnings                    *bool                      `json:"override_warnings,omitempty"`                       // (Write-Only) If warnings are preventing a host URL change, this parameter allows for overriding warnings to force update the setting. Does not directly change any Looker settings
	EmailDomainAllowlist                *[]string                  `json:"email_domain_allowlist,omitempty"`                  //
	EmbedCookielessV2                   *bool                      `json:"embed_cookieless_v2,omitempty"`                     // (DEPRECATED) Use embed_config.embed_cookieless_v2 instead. If embed_config.embed_cookieless_v2 is specified, it overrides this value
	EmbedEnabled                        *bool                      `json:"embed_enabled,omitempty"`                           // True if embedding is enabled https://cloud.google.com/looker/docs/r/looker-core-feature-embed, false otherwise (read-only)
	EmbedConfig                         *EmbedConfig               `json:"embed_config,omitempty"`                            // Embed configuration. Requires embedding to be enabled https://cloud.google.com/looker/docs/r/looker-core-feature-embed (read-only)
	LoginNotificationEnabled            *bool                      `json:"login_notification_enabled,omitempty"`              // Toggle login notification on or off (read-only)
	LoginNotificationText               *string                    `json:"login_notification_text,omitempty"`                 // Text to display in the login notification banner (read-only)
	DashboardAutorefreshRestriction     *bool                      `json:"dashboard_auto_refresh_restriction,omitempty"`      // Toggle Dashboard Auto Refresh restriction
	DashboardAutoRefreshMinimumInterval *string                    `json:"dashboard_auto_refresh_minimum_interval,omitempty"` // Minimum time interval for dashboard element automatic refresh. Examples: (30 seconds, 1 minute)
	ManagedCertificateUri               *[]string                  `json:"managed_certificate_uri,omitempty"`                 //
}

type InstanceConfig struct {
	FeatureFlags    *map[string]bool `json:"feature_flags,omitempty"`    // Feature flags for the instance (read-only)
	LicenseFeatures *map[string]bool `json:"license_features,omitempty"` // License features enabled on the instance (read-only)
}

type MarketplaceAutomation struct {
	InstallEnabled          *bool `json:"install_enabled,omitempty"`            // Whether marketplace auto installation is enabled
	UpdateLookerEnabled     *bool `json:"update_looker_enabled,omitempty"`      // Whether marketplace auto update is enabled for looker extensions
	UpdateThirdPartyEnabled *bool `json:"update_third_party_enabled,omitempty"` // Whether marketplace auto update is enabled for third party extensions
}

type PrivatelabelConfiguration struct {
	LogoFile                   *string `json:"logo_file,omitempty"`                     // Customer logo image. Expected base64 encoded data (write-only)
	LogoUrl                    *string `json:"logo_url,omitempty"`                      // Logo image url (read-only)
	FaviconFile                *string `json:"favicon_file,omitempty"`                  // Custom favicon image. Expected base64 encoded data (write-only)
	FaviconUrl                 *string `json:"favicon_url,omitempty"`                   // Favicon image url (read-only)
	DefaultTitle               *string `json:"default_title,omitempty"`                 // Default page title
	ShowHelpMenu               *bool   `json:"show_help_menu,omitempty"`                // Boolean to toggle showing help menus
	ShowDocs                   *bool   `json:"show_docs,omitempty"`                     // Boolean to toggle showing docs
	ShowEmailSubOptions        *bool   `json:"show_email_sub_options,omitempty"`        // Boolean to toggle showing email subscription options
	AllowLookerMentions        *bool   `json:"allow_looker_mentions,omitempty"`         // Boolean to toggle mentions of Looker in emails
	AllowLookerLinks           *bool   `json:"allow_looker_links,omitempty"`            // Boolean to toggle links to Looker in emails
	CustomWelcomeEmailAdvanced *bool   `json:"custom_welcome_email_advanced,omitempty"` // Allow subject line and email heading customization in customized emails
	SetupMentions              *bool   `json:"setup_mentions,omitempty"`                // Remove the word Looker from appearing in the account setup page
	AlertsLogo                 *bool   `json:"alerts_logo,omitempty"`                   // Remove Looker logo from Alerts
	AlertsLinks                *bool   `json:"alerts_links,omitempty"`                  // Remove Looker links from Alerts
	FoldersMentions            *bool   `json:"folders_mentions,omitempty"`              // Remove Looker mentions in home folder page when you don’t have any items saved
}

type CustomWelcomeEmail struct {
	Enabled *bool   `json:"enabled"`           // If true, custom email content will replace the default body of welcome emails
	Content *string `json:"content,omitempty"` // The HTML to use as custom content for welcome emails. Script elements and other potentially dangerous markup will be removed
	Subject *string `json:"subject,omitempty"` // The text to appear in the email subject line. Only available with a whitelabel license and whitelabel_configuration.advanced_custom_welcome_email enabled
	Header  *string `json:"header,omitempty"`  // The text to appear in the header line of the email body. Only available with a whitelabel license and whitelabel_configuration.advanced_custom_welcome_email enabled
}

type EmbedConfig struct {
	DomainAllowlist          *[]string `json:"domain_allowlist,omitempty"`
	AlertUrlAllowlist        *[]string `json:"alert_url_allowlist,omitempty"`         // List of allowed alert/schedule urls
	AlertUrlParamOwner       *string   `json:"alert_url_param_owner,omitempty"`       // Owner of who defines the alert/schedule params on the base url
	AlertUrlLabel            *string   `json:"alert_url_label,omitempty"`             // Label for the alert/schedule url
	SsoAuthEnabled           *bool     `json:"sso_auth_enabled,omitempty"`            // Is SSO embedding enabled for this Looker
	EmbedCookielessV2        *bool     `json:"embed_cookieless_v2,omitempty"`         // Is Cookieless embedding enabled for this Looker
	EmbedContentNavigation   *bool     `json:"embed_content_navigation,omitempty"`    // Is embed content navigation enabled for this looker
	EmbedContentManagement   *bool     `json:"embed_content_management,omitempty"`    // Is embed content management enabled for this Looker
	StrictSameoriginForLogin *bool     `json:"strict_sameorigin_for_login,omitempty"` // When true, prohibits the use of Looker login pages in non-Looker iframes. When false, Looker login pages may be used in non-Looker hosted iframes
	LookFilters              *bool     `json:"look_filters,omitempty"`                // When true, filters are enabled on embedded Looks
	HideLookNavigation       *bool     `json:"hide_look_navigation,omitempty"`        // When true, removes navigation to Looks from embedded dashboards and explores
	EmbedEnabled             *bool     `json:"embed_enabled,omitempty"`               // True if embedding is licensed for this Looker instance (read-only)
}

func (s *SettingResourceOp) Get(ctx context.Context) (*Setting, *Response, error) {
	return doGet(ctx, s.client, SettingBasePath, new(Setting))
}

func (s *SettingResourceOp) Update(ctx context.Context, requestSetting *Setting) (*Setting, *Response, error) {
	return doUpdate(ctx, s.client, SettingBasePath, "", requestSetting, new(Setting))
}

func (s *Setting) CleanFromReadOnly() {
	s.InstanceConfig = nil
	s.MarketplaceSite = nil
	s.PrivatelabelConfiguration.LogoUrl = nil
	s.PrivatelabelConfiguration.FaviconUrl = nil
	s.EmbedEnabled = nil
	s.LoginNotificationEnabled = nil
	s.LoginNotificationText = nil
}
