package lookergo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersResourceOp_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/4.0/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
  {
    "avatar_url": "https://gravatar.lookercdn.com/avatar/longhash?s=156&d=blank",
    "avatar_url_without_sizing": "https://gravatar.lookercdn.com/avatar/longhash?d=blank",
    "credentials_api3": [],
    "credentials_email": {
      "created_at": "2021-11-09T16:11:29.000+00:00",
      "logged_in_at": "",
      "type": "email",
      "email": "demo1@example.com",
      "forced_password_reset_at_next_login": false,
      "is_disabled": true,
      "password_reset_url": "https://example.cloud.looker.com/password/reset/example",
      "url": "https://localhost:19999/api/4.0/users/10/credentials_email",
      "user_url": "https://localhost:19999/api/4.0/users/10",
      "can": {
        "show_password_reset_url": true
      }
    },
    "credentials_embed": [],
    "credentials_google": null,
    "credentials_ldap": null,
    "credentials_looker_openid": null,
    "credentials_oidc": null,
    "credentials_saml": null,
    "credentials_totp": null,
    "email": "demo1@example.com",
    "first_name": "DemoKermit",
    "id": "10",
    "last_name": "the Frog",
    "locale": "en",
    "looker_versions": [],
    "models_dir_validated": false,
    "ui_state": null,
    "embed_group_folder_id": null,
    "home_folder_id": "38",
    "personal_folder_id": "38",
    "presumed_looker_employee": false,
    "sessions": [],
    "verified_looker_employee": false,
    "roles_externally_managed": false,
    "allow_direct_roles": true,
    "allow_normal_group_membership": true,
    "allow_roles_from_normal_groups": true,
    "display_name": "",
    "group_ids": [],
    "is_disabled": true,
    "role_ids": [],
    "url": "https://localhost:19999/api/4.0/users/10",
    "can": {
      "show": true,
      "index": true,
      "show_details": true,
      "index_details": true,
      "sudo": false
    }
  },
  {
    "avatar_url": "https://gravatar.lookercdn.com/avatar/longhash?s=156&d=blank",
    "avatar_url_without_sizing": "https://gravatar.lookercdn.com/avatar/longhash?d=blank",
    "credentials_api3": [],
    "credentials_email": {
      "created_at": "2021-11-09T16:11:29.000+00:00",
      "logged_in_at": "",
      "type": "email",
      "email": "demo2@example.com",
      "forced_password_reset_at_next_login": false,
      "is_disabled": true,
      "password_reset_url": "https://example.cloud.looker.com/password/reset/example",
      "url": "https://localhost:19999/api/4.0/users/20/credentials_email",
      "user_url": "https://localhost:19999/api/4.0/users/20",
      "can": {
        "show_password_reset_url": true
      }
    },
    "credentials_embed": [],
    "credentials_google": null,
    "credentials_ldap": null,
    "credentials_looker_openid": null,
    "credentials_oidc": null,
    "credentials_saml": null,
    "credentials_totp": null,
    "email": "demo1@example.com",
    "first_name": "DemoPiggy",
    "id": "20",
    "last_name": "Star",
    "locale": "en",
    "looker_versions": [],
    "models_dir_validated": false,
    "ui_state": null,
    "embed_group_folder_id": null,
    "home_folder_id": "238",
    "personal_folder_id": "238",
    "presumed_looker_employee": false,
    "sessions": [],
    "verified_looker_employee": false,
    "roles_externally_managed": false,
    "allow_direct_roles": true,
    "allow_normal_group_membership": true,
    "allow_roles_from_normal_groups": true,
    "display_name": "",
    "group_ids": [],
    "is_disabled": true,
    "role_ids": [],
    "url": "https://localhost:19999/api/4.0/users/20",
    "can": {
      "show": true,
      "index": true,
      "show_details": true,
      "index_details": true,
      "sudo": false
    }
  }
]`)
	})

	returned, resp, err := client.Users.List(ctx, nil)
	_ = resp
	if err != nil {
		t.Errorf("users.List returned error: %v", err)
	}

	expected := []User{
		{Id: 10, FirstName: "DemoKermit", LastName: "the Frog", CredentialEmail: &CredentialEmail{
			Email:      "demo1@example.com",
			IsDisabled: true,
		}},
		{Id: 20, FirstName: "DemoPiggy", LastName: "Star", CredentialEmail: &CredentialEmail{
			Email:      "demo2@example.com",
			IsDisabled: true,
		}},
	}
	if !reflect.DeepEqual(returned, expected) {
		t.Error(errGotWant("users.List", returned, expected))
	}
}

func TestUsersResourceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/4.0/users/10", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
    "avatar_url": "https://gravatar.lookercdn.com/avatar/longhash?s=156&d=blank",
    "avatar_url_without_sizing": "https://gravatar.lookercdn.com/avatar/longhash?d=blank",
    "credentials_api3": [],
    "credentials_email": {
      "created_at": "2021-11-09T16:11:29.000+00:00",
      "logged_in_at": "",
      "type": "email",
      "email": "demo1@example.com",
      "forced_password_reset_at_next_login": false,
      "is_disabled": true,
      "password_reset_url": "https://example.cloud.looker.com/password/reset/example",
      "url": "https://localhost:19999/api/4.0/users/10/credentials_email",
      "user_url": "https://localhost:19999/api/4.0/users/10",
      "can": {
        "show_password_reset_url": true
      }
    },
    "credentials_embed": [],
    "credentials_google": null,
    "credentials_ldap": null,
    "credentials_looker_openid": null,
    "credentials_oidc": null,
    "credentials_saml": null,
    "credentials_totp": null,
    "email": "demo1@example.com",
    "first_name": "DemoKermit",
    "id": "10",
    "last_name": "the Frog",
    "locale": "en",
    "looker_versions": [],
    "models_dir_validated": false,
    "ui_state": null,
    "embed_group_folder_id": null,
    "home_folder_id": "38",
    "personal_folder_id": "38",
    "presumed_looker_employee": false,
    "sessions": [],
    "verified_looker_employee": false,
    "roles_externally_managed": false,
    "allow_direct_roles": true,
    "allow_normal_group_membership": true,
    "allow_roles_from_normal_groups": true,
    "display_name": "",
    "group_ids": [],
    "is_disabled": true,
    "role_ids": [],
    "url": "https://localhost:19999/api/4.0/users/10",
    "can": {
      "show": true,
      "index": true,
      "show_details": true,
      "index_details": true,
      "sudo": false
    }
  }`)
	})

	returned, resp, err := client.Users.Get(ctx, 10)
	_ = resp
	if err != nil {
		t.Errorf("users.List returned error: %v", err)
	}

	expected := &User{
		Id: 10, FirstName: "DemoKermit", LastName: "the Frog", Email: "demo1@example.com",
		CredentialEmail: &CredentialEmail{
			Email:      "demo1@example.com",
			IsDisabled: true,
		}}
	if !reflect.DeepEqual(returned, expected) {
		t.Error(errGotWant("users.List", returned, expected))
	}
}

func TestUsersResourceOp_CreateEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/4.0/users/146/credentials_email", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{
  "created_at": "2022-05-23T05:38:38.000+00:00",
  "logged_in_at": "",
  "type": "email",
  "email": "kermit009@ipv4.plus",
  "forced_password_reset_at_next_login": false,
  "is_disabled": true,
  "password_reset_url": "",
  "url": "https://localhost:19999/api/4.0/users/146/credentials_email",
  "user_url": "https://localhost:19999/api/4.0/users/146",
  "can": {
    "show_password_reset_url": true
  }
}`)
	})
	newEntry := &CredentialEmail{
		Email:       "kermit009@ipv4.plus",
		ForcedReset: false,
		IsDisabled:  true,
	}

	created, resp, err := client.Users.CreateEmail(ctx, 146, newEntry)
	_ = resp
	if err != nil {
		t.Errorf("groups.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(created, newEntry) {
		t.Error(errGotWant("groups.Create", created, newEntry))
	}

}

func TestUsersResourceOp_GetRoles(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/4.0/users/123/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
  {
    "id": "2",
    "name": "Admin",
    "permission_set": {
      "built_in": true,
      "id": "1",
      "all_access": true,
      "name": "Admin",
      "permissions": [
        "access_data",
        "see_lookml_dashboards",
        "see_looks",
        "see_user_dashboards",
        "explore",
        "create_table_calculations",
        "create_custom_fields",
        "save_content",
        "create_public_looks",
        "download_with_limit",
        "download_without_limit",
        "schedule_look_emails",
        "schedule_external_look_emails",
        "create_alerts",
        "follow_alerts",
        "send_to_s3",
        "send_to_sftp",
        "send_outgoing_webhook",
        "send_to_integration",
        "see_sql",
        "see_lookml",
        "develop",
        "deploy",
        "support_access_toggle",
        "use_sql_runner",
        "clear_cache_refresh",
        "can_copy_print",
        "see_drill_overlay",
        "manage_spaces",
        "manage_homepage",
        "manage_models",
        "manage_stereo",
        "create_prefetches",
        "login_special_email",
        "embed_browse_spaces",
        "embed_save_shared_space",
        "see_alerts",
        "see_queries",
        "see_logs",
        "see_users",
        "sudo",
        "see_schedules",
        "see_pdts",
        "see_datagroups",
        "update_datagroups",
        "see_system_activity",
        "administer",
        "mobile_app_access"
      ],
      "url": "https://localhost:19999/api/4.0/permission_sets/1",
      "can": {
        "show": true,
        "index": true,
        "update": true
      }
    },
    "model_set": {
      "built_in": true,
      "id": "1",
      "all_access": true,
      "models": [
        "usa-sharkninja",
        "reprise_th_test",
        "usa-leafguard",
        "usa-j3nutragena",
        "usa-levistrauss",
        "dach-amazon-wfs",
        "usa-aldius",
        "dach-cross-client",
        "test_us_ecomm_reviews",
        "nld-plussuper",
        "learning-customer-1",
        "cross_channel",
        "cross_client",
        "facebook_ads",
        "google_ads",
        "ga360",
        "ga4",
        "hub",
        "block_linkedin",
        "pinterest",
        "sa360",
        "snapchat",
        "tiktok",
        "twitter_ads",
        "block_youtube_dv360",
        "block_youtube_ga",
        "amazon_prime",
        "uk-tommee-t",
        "nld-hbm",
        "dach-dyson-ch",
        "dach-h-and-m",
        "canada-amz-xcm",
        "usa-amazon-gca",
        "sandbox-vincent1",
        "dach-dyson-de",
        "dach-boconcept",
        "us_ecommerce_reviews",
        "dach-teva-de",
        "gcp_billing",
        "lookml-diagram",
        "dach-ch"
      ],
      "name": "All",
      "url": "https://localhost:19999/api/4.0/model_sets/1",
      "can": {
        "show": true,
        "index": true,
        "update": true
      }
    },
    "url": "https://localhost:19999/api/4.0/roles/2",
    "users_url": "https://localhost:19999/api/4.0/roles/2/users",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  }
]`)
	})

	returned, resp, err := client.Users.GetRoles(ctx, 123)

	_ = resp
	if err != nil {
		t.Errorf("users.GetRoles returned error: %v", err)
	}

	expected := []Role{
		{Id: 2, Name: "Admin"},
	}
	if !reflect.DeepEqual(returned, expected) {
		t.Error(errGotWant("users.List", returned, expected))
	}
}

func TestUsersResourceOp_SetRoles(t *testing.T) {
	setup()
	defer teardown()

	putRequest := []int{88, 99}
	putResponse := []Role{
		{Id: 88, Name: "Admin"},
		{Id: 99, Name: "Admin"},
	}

	mux.HandleFunc("/4.0/users/123/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		fmt.Fprint(w, `[
  {
    "id": "88",
    "name": "Admin"
  },{
    "id": "99",
    "name": "Admin"
  }
]`)
	})

	created, resp, err := client.Users.SetRoles(ctx, 123, putRequest)
	_ = resp
	if err != nil {
		t.Errorf("groups.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(created, putResponse) {
		t.Error(errGotWant("groups.Create", created, putResponse))
	}
}
