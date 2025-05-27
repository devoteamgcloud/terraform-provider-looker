package lookergo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSettingResourceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/4.0/setting", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
    "instance_config": {
        "feature_flags": {
            "g3_lookml_assistant": true,
            "content_validator_scoping": true,
            "unification_content_lifecycle": true,
            "manage_schedules_permission_orig": true,
            "new_lexp_interpreter": true
        },
        "license_features": {
            "allow_private_embed": true,
            "ldap": true,
            "custom_home_page": true,
            "eula_for_all": true,
            "saml": true
        }
    },
    "extension_framework_enabled": true,
    "extension_load_url_enabled": true,
    "marketplace_auto_install_enabled": false,
    "marketplace_automation": {
        "install_enabled": false,
        "update_looker_enabled": false,
        "update_third_party_enabled": false
    },
    "marketplace_enabled": true,
    "marketplace_site": "https://static-a.cdn.looker.app/marketplace/",
    "marketplace_terms_accepted": true,
    "privatelabel_configuration": {
        "logo_file": null,
        "logo_url": "/some/logo.png",
        "favicon_file": null,
        "favicon_url": "/some/favicon.png",
        "default_title": "",
        "show_help_menu": true,
        "show_docs": true,
        "show_email_sub_options": false,
        "allow_looker_mentions": true,
        "allow_looker_links": true,
        "custom_welcome_email_advanced": false,
        "setup_mentions": false,
        "alerts_logo": false,
        "alerts_links": false,
        "folders_mentions": false
    },
    "custom_welcome_email": {
        "enabled": false,
        "content": "Example content",
        "subject": "Welcome to Looker",
        "header": "You&#39;ve been invited to join Looker!"
    },
    "onboarding_enabled": true,
    "timezone": "America/Los_Angeles",
    "allow_user_timezones": true,
    "data_connector_default_enabled": false,
    "host_url": "https://acme.cloud.looker.com",
    "email_domain_allowlist": [],
    "embed_cookieless_v2": true,
    "embed_enabled": true,
    "embed_config": {
        "domain_allowlist": [
            "https://*.acme.com",
            "https://*.acme.com/"
        ],
        "alert_url_allowlist": [],
        "alert_url_param_owner": "",
        "alert_url_label": "",
        "sso_auth_enabled": true,
        "embed_cookieless_v2": true,
        "embed_content_navigation": false,
        "embed_content_management": false,
        "strict_sameorigin_for_login": false,
        "look_filters": false,
        "hide_look_navigation": false,
        "embed_enabled": true
    },
    "dashboard_auto_refresh_restriction": false,
    "dashboard_auto_refresh_minimum_interval": null,
    "managed_certificate_uri": null
}`)
	})

	result, resp, err := client.Setting.Get(ctx)
	_ = resp
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	expected := &Setting{
		InstanceConfig: &InstanceConfig{
			FeatureFlags: &map[string]bool{
				"g3_lookml_assistant":              true,
				"content_validator_scoping":        true,
				"unification_content_lifecycle":    true,
				"manage_schedules_permission_orig": true,
				"new_lexp_interpreter":             true,
			},
			LicenseFeatures: &map[string]bool{
				"allow_private_embed": true,
				"ldap":                true,
				"custom_home_page":    true,
				"eula_for_all":        true,
				"saml":                true,
			},
		},
		ExtensionFrameworkEnabled:     castToPtr(true),
		ExtensionLoadUrlEnabled:       castToPtr(true),
		MarketplaceAutoInstallEnabled: castToPtr(false),
		MarketplaceAutomation: &MarketplaceAutomation{
			InstallEnabled:          castToPtr(false),
			UpdateLookerEnabled:     castToPtr(false),
			UpdateThirdPartyEnabled: castToPtr(false),
		},
		MarketplaceEnabled:       castToPtr(true),
		MarketplaceSite:          castToPtr("https://static-a.cdn.looker.app/marketplace/"),
		MarketplaceTermsAccepted: castToPtr(true),
		PrivatelabelConfiguration: &PrivatelabelConfiguration{
			LogoFile:                   nil,
			LogoUrl:                    castToPtr("/some/logo.png"),
			FaviconFile:                nil,
			FaviconUrl:                 castToPtr("/some/favicon.png"),
			DefaultTitle:               castToPtr(""),
			ShowHelpMenu:               castToPtr(true),
			ShowDocs:                   castToPtr(true),
			ShowEmailSubOptions:        castToPtr(false),
			AllowLookerMentions:        castToPtr(true),
			AllowLookerLinks:           castToPtr(true),
			CustomWelcomeEmailAdvanced: castToPtr(false),
			SetupMentions:              castToPtr(false),
			AlertsLogo:                 castToPtr(false),
			AlertsLinks:                castToPtr(false),
			FoldersMentions:            castToPtr(false),
		},
		CustomWelcomeEmail: &CustomWelcomeEmail{
			Enabled: castToPtr(false),
			Content: castToPtr("Example content"),
			Subject: castToPtr("Welcome to Looker"),
			Header:  castToPtr("You&#39;ve been invited to join Looker!"),
		},
		OnboardingEnabled:           castToPtr(true),
		Timezone:                    castToPtr("America/Los_Angeles"),
		AllowUserTimezones:          castToPtr(true),
		DataConnectorDefaultEnabled: castToPtr(false),
		HostUrl:                     castToPtr("https://acme.cloud.looker.com"),
		EmailDomainAllowlist:        castToPtr([]string{}),
		EmbedCookielessV2:           castToPtr(true),
		EmbedEnabled:                castToPtr(true),
		EmbedConfig: &EmbedConfig{
			DomainAllowlist:          castToPtr([]string{"https://*.acme.com", "https://*.acme.com/"}),
			AlertUrlAllowlist:        castToPtr([]string{}),
			AlertUrlParamOwner:       castToPtr(""),
			AlertUrlLabel:            castToPtr(""),
			SsoAuthEnabled:           castToPtr(true),
			EmbedCookielessV2:        castToPtr(true),
			EmbedContentNavigation:   castToPtr(false),
			EmbedContentManagement:   castToPtr(false),
			StrictSameoriginForLogin: castToPtr(false),
			LookFilters:              castToPtr(false),
			HideLookNavigation:       castToPtr(false),
			EmbedEnabled:             castToPtr(true),
		},
		DashboardAutorefreshRestriction:     castToPtr(false),
		DashboardAutoRefreshMinimumInterval: nil,
		ManagedCertificateUri:               nil,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Error(errGotWant("Projects.Get", result, expected))
	}
}
