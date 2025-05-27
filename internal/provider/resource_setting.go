package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingUpdate,
		ReadContext:   resourceSettingRead,
		UpdateContext: resourceSettingUpdate,
		DeleteContext: resourceSettingDelete, // Setting is not deletable, but DeleteContext is required
		Schema: map[string]*schema.Schema{
			"instance_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Externally available instance configuration information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature_flags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Feature flags for the instance",
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
						"license_features": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "License features enabled on the instance",
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
					},
				},
			},
			"extension_framework_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Toggle extension framework on or off",
			},
			"extension_load_url_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "(DEPRECATED) Toggle extension load url on or off. Do not use. This is temporary setting that will eventually become a noop and subsequently deleted",
			},
			"marketplace_auto_install_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "(DEPRECATED) Toggle marketplace auto install on or off. Deprecated - do not use. Auto install can now be enabled via marketplace automation settings",
			},
			"marketplace_automation": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Marketplace automation settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether marketplace auto installation is enabled",
						},
						"update_looker_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether marketplace auto update is enabled for looker extensions",
						},
						"update_third_party_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether marketplace auto update is enabled for third party extensions",
						},
					},
				},
			},
			"marketplace_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Toggle marketplace on or off",
			},
			"marketplace_site": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of Looker marketplace CDN",
			},
			"marketplace_terms_accepted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Accept marketplace terms by setting this value to true, or get the current status. Marketplace terms CANNOT be declined once accepted. Accepting marketplace terms automatically enables the marketplace. The marketplace can still be disabled after it has been enabled",
			},
			"privatelabel_configuration": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Optional:    true,
				Description: "Private label configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logo_file": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Customer logo image. Expected base64 encoded data (write-only)",
						},
						"logo_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Logo image url (read-only)",
						},
						"favicon_file": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Customer favicon image. Expected base64 encoded data (write-only)",
						},
						"favicon_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Favicon image url (read-only)",
						},
						"default_title": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Default page title",
						},
						"show_help_menu": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Boolean to toggle showing help menus",
						},
						"show_docs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Boolean to toggle showing docs",
						},
						"show_email_sub_options": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Boolean to toggle showing email subscription options",
						},
						"allow_looker_mentions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Boolean to toggle mentions of Looker in emails",
						},
						"allow_looker_links": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Boolean to toggle links to Looker in emails",
						},
						"custom_welcome_email_advanced": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Allow subject line and email heading customization in customized emails",
						},
						"setup_mentions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Remove the word Looker from appearing in the account setup page",
						},
						"alerts_logo": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Remove Looker logo from Alerts",
						},
						"alerts_links": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Remove Looker links from Alerts",
						},
						"folders_mentions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Remove Looker mentions in home folder page when you don't have any items saved",
						},
					},
				},
			},
			"custom_welcome_email": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Custom welcome email configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If true, custom email content will replace the default body of welcome emails",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Requires custom_welcome_email to be enabled. The HTML to use as custom content for welcome emails. Script elements and other potentially dangerous markup will be removed",
						},
						"subject": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Requires custom_welcome_email and privatelabel_configuration.custom_welcome_email_advanced to be enabled. The text to appear in the email subject line. Only available with a whitelabel license and whitelabel_configuration.advanced_custom_welcome_email enabled",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Requires custom_welcome_email and privatelabel_configuration.custom_welcome_email_advanced to be enabled. The text to appear in the header line of the email body. Only available with a whitelabel license and whitelabel_configuration.advanced_custom_welcome_email enabled",
						},
					},
				},
			},
			"onboarding_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Toggle onboarding on or off",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Change instance-wide default timezone",
			},
			"allow_user_timezones": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Toggle user-specific timezones on or off",
			},
			"data_connector_default_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Toggle default future connectors on or off",
			},
			"host_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Change the base portion of your Looker instance URL setting",
			},
			"override_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "(Write-Only) If warnings are preventing a host URL change, this parameter allows for overriding warnings to force update the setting. Does not directly change any Looker settings",
			},
			"email_domain_allowlist": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of email domains that are allowed to be used for user creation",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"embed_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Embed configuration. Requires embedding to be enabled https://cloud.google.com/looker/docs/r/looker-core-feature-embed",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_allowlist": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"alert_url_allowlist": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"alert_url_param_owner": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Owner of who defines the alert/schedule params on the base url",
						},
						"alert_url_label": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Label for the alert/schedule url",
						},
						"sso_auth_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Is SSO embedding enabled for this Looker",
						},
						"embed_cookieless_v2": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Is Cookieless embedding enabled for this Looker",
						},
						"embed_content_navigation": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Is embed content navigation enabled for this looker",
						},
						"embed_content_management": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Is embed content management enabled for this Looker",
						},
						"strict_sameorigin_for_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "When true, prohibits the use of Looker login pages in non-Looker iframes. When false, Looker login pages may be used in non-Looker hosted iframes",
						},
						"look_filters": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "When true, filters are enabled on embedded Looks",
						},
						"hide_look_navigation": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "When true, removes navigation to Looks from embedded dashboards and explores",
						},
						"embed_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "True if embedding is licensed for this Looker instance",
						},
					},
				},
			},
			"login_notification_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Toggle login notification on or off",
			},
			"login_notification_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text to display in the login notification banner",
			},
			"dashboard_autorefresh_restriction": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Toggle dashboard auto refresh restriction",
			},
			"dashboard_auto_refresh_minimum_interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Minimum time interval for dashboard element automatic refresh. Examples: (30 seconds, 1 minute)",
			},
			"managed_certificate_uri": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m any) (diags diag.Diagnostics) {
	c := m.(*Config).Api

	d.SetId("looker_settings")

	setting, _, err := c.Setting.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("extension_framework_enabled", setting.ExtensionFrameworkEnabled)
	d.Set("extension_load_url_enabled", setting.ExtensionLoadUrlEnabled)
	d.Set("marketplace_auto_install_enabled", setting.MarketplaceAutoInstallEnabled)
	d.Set("marketplace_enabled", setting.MarketplaceEnabled)
	d.Set("marketplace_site", setting.MarketplaceSite)
	d.Set("marketplace_terms_accepted", setting.MarketplaceTermsAccepted)
	d.Set("onboarding_enabled", setting.OnboardingEnabled)
	d.Set("timezone", setting.Timezone)
	d.Set("allow_user_timezones", setting.AllowUserTimezones)
	d.Set("data_connector_default_enabled", setting.DataConnectorDefaultEnabled)
	d.Set("host_url", setting.HostUrl)
	d.Set("login_notification_enabled", setting.LoginNotificationEnabled)
	d.Set("login_notification_text", setting.LoginNotificationText)
	d.Set("dashboard_autorefresh_restriction", setting.DashboardAutorefreshRestriction)
	d.Set("dashboard_auto_refresh_minimum_interval", setting.DashboardAutoRefreshMinimumInterval)

	d.Set("email_domain_allowlist", setting.EmailDomainAllowlist)
	d.Set("managed_certificate_uri", setting.ManagedCertificateUri)

	readInstanceConfig(d, setting)
	readMarketplaceAutomation(d, setting)
	readPrivatelabelConfiguration(d, setting)
	readCustomWelcomeEmail(d, setting)
	readEmbedConfig(d, setting)

	return diags
}

func readInstanceConfig(d *schema.ResourceData, setting *lookergo.Setting) {
	if setting.InstanceConfig != nil {
		instanceConfigMap := map[string]any{
			"feature_flags":    make(map[string]bool),
			"license_features": make(map[string]bool),
		}

		if setting.InstanceConfig.FeatureFlags != nil {
			instanceConfigMap["feature_flags"] = *setting.InstanceConfig.FeatureFlags
		}

		if setting.InstanceConfig.LicenseFeatures != nil {
			instanceConfigMap["license_features"] = *setting.InstanceConfig.LicenseFeatures
		}

		d.Set("instance_config", []map[string]any{instanceConfigMap})
	}
}

func readMarketplaceAutomation(d *schema.ResourceData, setting *lookergo.Setting) {
	if setting.MarketplaceAutomation != nil {
		marketplaceAutomation := setting.MarketplaceAutomation
		marketplaceAutomationMap := map[string]any{
			"install_enabled":            marketplaceAutomation.InstallEnabled,
			"update_looker_enabled":      marketplaceAutomation.UpdateLookerEnabled,
			"update_third_party_enabled": marketplaceAutomation.UpdateThirdPartyEnabled,
		}

		d.Set("marketplace_automation", []map[string]any{marketplaceAutomationMap})
	}
}

func readPrivatelabelConfiguration(d *schema.ResourceData, setting *lookergo.Setting) {
	if setting.PrivatelabelConfiguration != nil {
		privatelabelConfiguration := setting.PrivatelabelConfiguration
		privatelabelConfigurationMap := map[string]any{
			"logo_url":                      privatelabelConfiguration.LogoUrl,
			"favicon_url":                   privatelabelConfiguration.FaviconUrl,
			"default_title":                 privatelabelConfiguration.DefaultTitle,
			"show_help_menu":                privatelabelConfiguration.ShowHelpMenu,
			"show_docs":                     privatelabelConfiguration.ShowDocs,
			"show_email_sub_options":        privatelabelConfiguration.ShowEmailSubOptions,
			"allow_looker_mentions":         privatelabelConfiguration.AllowLookerMentions,
			"allow_looker_links":            privatelabelConfiguration.AllowLookerLinks,
			"custom_welcome_email_advanced": privatelabelConfiguration.CustomWelcomeEmailAdvanced,
			"setup_mentions":                privatelabelConfiguration.SetupMentions,
			"alerts_logo":                   privatelabelConfiguration.AlertsLogo,
			"alerts_links":                  privatelabelConfiguration.AlertsLinks,
			"folders_mentions":              privatelabelConfiguration.FoldersMentions,
		}

		d.Set("privatelabel_configuration", []map[string]any{privatelabelConfigurationMap})
	}
}

func readCustomWelcomeEmail(d *schema.ResourceData, setting *lookergo.Setting) {
	if setting.CustomWelcomeEmail != nil {
		customWelcomeEmail := setting.CustomWelcomeEmail
		customWelcomeEmailMap := map[string]any{
			"enabled": customWelcomeEmail.Enabled,
			"content": customWelcomeEmail.Content,
			"subject": customWelcomeEmail.Subject,
			"header":  customWelcomeEmail.Header,
		}

		d.Set("custom_welcome_email", []map[string]any{customWelcomeEmailMap})
	}
}

func readEmbedConfig(d *schema.ResourceData, setting *lookergo.Setting) {
	if setting.EmbedConfig != nil {
		embedConfig := setting.EmbedConfig
		embedConfigMap := map[string]any{
			"domain_allowlist":            embedConfig.DomainAllowlist,
			"alert_url_allowlist":         embedConfig.AlertUrlAllowlist,
			"alert_url_param_owner":       embedConfig.AlertUrlParamOwner,
			"alert_url_label":             embedConfig.AlertUrlLabel,
			"sso_auth_enabled":            embedConfig.SsoAuthEnabled,
			"embed_cookieless_v2":         embedConfig.EmbedCookielessV2,
			"embed_content_navigation":    embedConfig.EmbedContentNavigation,
			"embed_content_management":    embedConfig.EmbedContentManagement,
			"strict_sameorigin_for_login": embedConfig.StrictSameoriginForLogin,
			"look_filters":                embedConfig.LookFilters,
			"hide_look_navigation":        embedConfig.HideLookNavigation,
			"embed_enabled":               embedConfig.EmbedEnabled,
		}

		d.Set("embed_config", []map[string]any{embedConfigMap})
	}
}

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m any) (diags diag.Diagnostics) {
	c := m.(*Config).Api
	setting, _, err := c.Setting.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var x []byte
	x, _ = json.MarshalIndent(setting, "", "  ")
	tflog.Debug(ctx, fmt.Sprintf("Current settings: %s", string(x)))

	updateSetting("extension_framework_enabled", setting.ExtensionFrameworkEnabled, d)
	updateSetting("extension_load_url_enabled", setting.ExtensionLoadUrlEnabled, d)
	updateSetting("marketplace_auto_install_enabled", setting.MarketplaceAutoInstallEnabled, d)
	updateSetting("marketplace_enabled", setting.MarketplaceEnabled, d)
	updateSetting("marketplace_terms_accepted", setting.MarketplaceTermsAccepted, d)
	updateSetting("onboarding_enabled", setting.OnboardingEnabled, d)
	updateSetting("allow_user_timezones", setting.AllowUserTimezones, d)
	updateSetting("data_connector_default_enabled", setting.DataConnectorDefaultEnabled, d)
	updateSetting("override_warnings", setting.OverrideWarnings, d)
	updateSetting("dashboard_autorefresh_restriction", setting.DashboardAutorefreshRestriction, d)

	updateSetting("timezone", setting.Timezone, d)
	updateSetting("host_url", setting.HostUrl, d)
	updateSetting("dashboard_auto_refresh_minimum_interval", setting.DashboardAutoRefreshMinimumInterval, d)

	updateStringList("email_domain_allowlist", setting.EmailDomainAllowlist, d)
	updateStringList("managed_certificate_uri", setting.ManagedCertificateUri, d)

	updateMarketplaceAutomation(d, setting)
	updatePrivatelabelConfiguration(d, setting)

	if err = updateEmbedConfig(d, setting); err != nil {
		return diag.FromErr(err)
	}

	if err = updateCustomWelcomeEmail(d, setting); err != nil {
		return diag.FromErr(err)
	}

	x, _ = json.MarshalIndent(setting, "", "  ")
	tflog.Debug(ctx, fmt.Sprintf("Updated settings: %s", string(x)))

	if err = handleDependencies(ctx, d, setting, c); err != nil {
		return diag.FromErr(err)
	}

	setting.CleanFromReadOnly()

	_, _, err = c.Setting.Update(ctx, setting)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Settings updated successfully")

	cleanStateFromWriteOnly(d)

	return resourceSettingRead(ctx, d, m)
}

func updateSetting[T any](key string, target *T, d *schema.ResourceData) {
	if d.HasChange(key) {
		val := d.Get(key).(T)
		tflog.Debug(context.Background(), fmt.Sprintf("Updating %s: %v -> %v", key, *target, val))
		*target = val
	}
}

func updateStringList(key string, target *[]string, d *schema.ResourceData) {
	if d.HasChange(key) {
		values := d.Get(key).([]any)
		result := make([]string, len(values))
		for i, v := range values {
			result[i] = v.(string)
		}
		*target = result
	}
}

func updateMarketplaceAutomation(d *schema.ResourceData, setting *lookergo.Setting) {
	updateSetting("marketplace_automation.0.install_enabled", setting.MarketplaceAutomation.InstallEnabled, d)
	updateSetting("marketplace_automation.0.update_looker_enabled", setting.MarketplaceAutomation.UpdateLookerEnabled, d)
	updateSetting("marketplace_automation.0.update_third_party_enabled", setting.MarketplaceAutomation.UpdateThirdPartyEnabled, d)
}

func updatePrivatelabelConfiguration(d *schema.ResourceData, setting *lookergo.Setting) {
	updateSetting("privatelabel_configuration.0.logo_file", setting.PrivatelabelConfiguration.LogoFile, d)
	updateSetting("privatelabel_configuration.0.favicon_file", setting.PrivatelabelConfiguration.FaviconFile, d)
	updateSetting("privatelabel_configuration.0.default_title", setting.PrivatelabelConfiguration.DefaultTitle, d)
	updateSetting("privatelabel_configuration.0.show_help_menu", setting.PrivatelabelConfiguration.ShowHelpMenu, d)
	updateSetting("privatelabel_configuration.0.show_docs", setting.PrivatelabelConfiguration.ShowDocs, d)
	updateSetting("privatelabel_configuration.0.show_email_sub_options", setting.PrivatelabelConfiguration.ShowEmailSubOptions, d)
	updateSetting("privatelabel_configuration.0.allow_looker_mentions", setting.PrivatelabelConfiguration.AllowLookerMentions, d)
	updateSetting("privatelabel_configuration.0.allow_looker_links", setting.PrivatelabelConfiguration.AllowLookerLinks, d)
	updateSetting("privatelabel_configuration.0.custom_welcome_email_advanced", setting.PrivatelabelConfiguration.CustomWelcomeEmailAdvanced, d)
	updateSetting("privatelabel_configuration.0.setup_mentions", setting.PrivatelabelConfiguration.SetupMentions, d)
	updateSetting("privatelabel_configuration.0.alerts_logo", setting.PrivatelabelConfiguration.AlertsLogo, d)
	updateSetting("privatelabel_configuration.0.alerts_links", setting.PrivatelabelConfiguration.AlertsLinks, d)
	updateSetting("privatelabel_configuration.0.folders_mentions", setting.PrivatelabelConfiguration.FoldersMentions, d)
}

func updateCustomWelcomeEmail(d *schema.ResourceData, setting *lookergo.Setting) error {
	updateSetting("custom_welcome_email.0.enabled", setting.CustomWelcomeEmail.Enabled, d)

	// content, subject, and header should only be updated if custom_welcome_email is enabled
	if setting.CustomWelcomeEmail.Enabled != nil {
		if *setting.CustomWelcomeEmail.Enabled {
			updateSetting("custom_welcome_email.0.content", setting.CustomWelcomeEmail.Content, d)
			updateSetting("custom_welcome_email.0.subject", setting.CustomWelcomeEmail.Subject, d)
			updateSetting("custom_welcome_email.0.header", setting.CustomWelcomeEmail.Header, d)
		} else {
			// Check if any of these fields are explicitly set in the configuration
			rawConfig := d.GetRawConfig()
			customWelcomeEmailConfig := rawConfig.GetAttr("custom_welcome_email")

			if !customWelcomeEmailConfig.IsNull() && customWelcomeEmailConfig.LengthInt() > 0 {
				emailBlock := customWelcomeEmailConfig.Index(cty.NumberIntVal(0))

				fields := []string{"content", "subject", "header"}
				for _, field := range fields {
					fieldValue := emailBlock.GetAttr(field)
					if !fieldValue.IsNull() && fieldValue.IsKnown() && fieldValue.AsString() != "" {
						return fmt.Errorf("custom_welcome_email.%s cannot be set when custom_welcome_email.enabled is false (was %q)", field, fieldValue.AsString())
					}
				}
			}

			setting.CustomWelcomeEmail.Content = nil
			setting.CustomWelcomeEmail.Subject = nil
			setting.CustomWelcomeEmail.Header = nil
		}
	}

	return nil
}

func updateEmbedConfig(d *schema.ResourceData, setting *lookergo.Setting) error {
	updateStringList("embed_config.0.domain_allowlist", setting.EmbedConfig.DomainAllowlist, d)
	updateStringList("embed_config.0.alert_url_allowlist", setting.EmbedConfig.AlertUrlAllowlist, d)
	updateSetting("embed_config.0.alert_url_param_owner", setting.EmbedConfig.AlertUrlParamOwner, d)
	updateSetting("embed_config.0.alert_url_label", setting.EmbedConfig.AlertUrlLabel, d)
	updateSetting("embed_config.0.sso_auth_enabled", setting.EmbedConfig.SsoAuthEnabled, d)
	updateSetting("embed_config.0.embed_cookieless_v2", setting.EmbedConfig.EmbedCookielessV2, d)
	updateSetting("embed_config.0.embed_content_navigation", setting.EmbedConfig.EmbedContentNavigation, d)
	updateSetting("embed_config.0.embed_content_management", setting.EmbedConfig.EmbedContentManagement, d)
	updateSetting("embed_config.0.strict_sameorigin_for_login", setting.EmbedConfig.StrictSameoriginForLogin, d)
	updateSetting("embed_config.0.look_filters", setting.EmbedConfig.LookFilters, d)
	updateSetting("embed_config.0.hide_look_navigation", setting.EmbedConfig.HideLookNavigation, d)
	updateSetting("embed_config.0.embed_enabled", setting.EmbedConfig.EmbedEnabled, d)

	if setting.EmbedConfig != nil &&
		setting.EmbedConfig.EmbedCookielessV2 != nil && setting.EmbedConfig.EmbedEnabled != nil &&
		*setting.EmbedConfig.EmbedCookielessV2 && !*setting.EmbedConfig.EmbedEnabled {
		return fmt.Errorf("embed_config.embed_cookieless_v2 cannot be set to true when embed_config.embed_enabled is false")
	}

	return nil
}

func handleDependencies(ctx context.Context, d *schema.ResourceData, setting *lookergo.Setting, c *lookergo.Client) error {
	var errors []error

	if d.HasChange("privatelabel_configuration.0.custom_welcome_email_advanced") {
		v := d.Get("privatelabel_configuration.0.custom_welcome_email_advanced")

		tflog.Debug(
			ctx,
			"privatelabel_configuration.custom_welcome_email_advanced has changed, updating it first, before the rest",
		)

		singularSetting := &lookergo.Setting{
			PrivatelabelConfiguration: &lookergo.PrivatelabelConfiguration{
				CustomWelcomeEmailAdvanced: castToPtr(v.(bool)),
			},
		}

		_, _, err := c.Setting.Update(ctx, singularSetting)

		if err != nil {
			errors = append(errors, err)
		}
	}

	if setting.PrivatelabelConfiguration != nil &&
		setting.PrivatelabelConfiguration.CustomWelcomeEmailAdvanced != nil &&
		!*setting.PrivatelabelConfiguration.CustomWelcomeEmailAdvanced {
		// custom_welcome_email.subject, and custom_welcome_email.header
		// should only be included if custom_welcome_email_advanced is enabled
		// BUT if they are explicitly configured, it's an error
		rawConfig := d.GetRawConfig()
		customWelcomeEmailConfig := rawConfig.GetAttr("custom_welcome_email")

		if !customWelcomeEmailConfig.IsNull() && customWelcomeEmailConfig.LengthInt() > 0 {
			emailBlock := customWelcomeEmailConfig.Index(cty.NumberIntVal(0))

			fields := []string{"subject", "header"}
			for _, field := range fields {
				fieldValue := emailBlock.GetAttr(field)
				if !fieldValue.IsNull() && fieldValue.IsKnown() && fieldValue.AsString() != "" {
					return fmt.Errorf("custom_welcome_email.%s cannot be set when privatelabel_configuration.custom_welcome_email_advanced is false (was %q)", field, fieldValue.AsString())
				}
			}
		}

		setting.CustomWelcomeEmail.Subject = nil
		setting.CustomWelcomeEmail.Header = nil
	}

	if len(errors) > 0 {
		return fmt.Errorf("one or more errors occurred: %v", errors)
	}

	return nil
}

func cleanStateFromWriteOnly(d *schema.ResourceData) {
	if d.HasChange("override_warnings") {
		d.Set("override_warnings", nil)
	}

	if d.HasChange("privatelabel_configuration.0.logo_file") {
		d.Set("privatelabel_configuration.0.logo_file", nil)
	}

	if d.HasChange("privatelabel_configuration.0.favicon_file") {
		d.Set("privatelabel_configuration.0.favicon_file", nil)
	}
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m any) (diags diag.Diagnostics) {
	d.SetId("")
	return diags
}
