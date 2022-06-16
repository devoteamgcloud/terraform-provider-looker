package lookergo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRolesResourceOp_PermissionSetsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/4.0/permission_sets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
[
  {
    "built_in": true,
    "id": "1",
    "all_access": true,
    "name": "Admin",
    "permissions": ["access_data", "see_lookml_dashboards", "see_looks", "see_user_dashboards", "explore", "create_table_calculations", "create_custom_fields", "save_content", "create_public_looks", "download_with_limit", "download_without_limit", "schedule_look_emails", "schedule_external_look_emails", "create_alerts", "follow_alerts", "send_to_s3", "send_to_sftp", "send_outgoing_webhook", "send_to_integration", "see_sql", "see_lookml", "develop", "deploy", "support_access_toggle", "use_sql_runner", "clear_cache_refresh", "can_copy_print", "see_drill_overlay", "manage_spaces", "manage_homepage", "manage_models", "manage_stereo", "create_prefetches", "login_special_email", "embed_browse_spaces", "embed_save_shared_space", "see_alerts", "see_queries", "see_logs", "see_users", "sudo", "see_schedules", "see_pdts", "see_datagroups", "update_datagroups", "see_system_activity", "administer", "mobile_app_access"],
    "url": "https://localhost:19999/api/4.0/permission_sets/1",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "2",
    "all_access": false,
    "name": "Developer",
    "permissions": ["access_data", "clear_cache_refresh", "create_alerts", "create_custom_fields", "create_public_looks", "create_table_calculations", "deploy", "develop", "download_without_limit", "explore", "follow_alerts", "login_special_email", "manage_spaces", "mobile_app_access", "save_content", "schedule_look_emails", "see_drill_overlay", "see_lookml", "see_lookml_dashboards", "see_looks", "see_sql", "see_user_dashboards", "send_to_integration", "use_sql_runner"],
    "url": "https://localhost:19999/api/4.0/permission_sets/2",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "3",
    "all_access": false,
    "name": "User",
    "permissions": ["access_data", "clear_cache_refresh", "create_alerts", "create_custom_fields", "create_table_calculations", "download_without_limit", "explore", "login_special_email", "manage_spaces", "mobile_app_access", "save_content", "schedule_look_emails", "see_drill_overlay", "see_lookml", "see_lookml_dashboards", "see_looks", "see_sql", "see_user_dashboards", "send_to_integration"],
    "url": "https://localhost:19999/api/4.0/permission_sets/3",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "4",
    "all_access": false,
    "name": "LookML Dashboard User",
     "permissions": ["access_data", "clear_cache_refresh", "mobile_app_access", "see_lookml_dashboards"],
    "url": "https://localhost:19999/api/4.0/permission_sets/4",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "5",
    "all_access": false,
    "name": "User who can't view LookML",
    "permissions": ["access_data", "clear_cache_refresh", "create_custom_fields", "create_table_calculations", "download_without_limit", "explore", "manage_spaces", "mobile_app_access", "save_content", "schedule_look_emails", "see_lookml_dashboards", "see_looks", "see_user_dashboards", "send_to_integration"],
    "url": "https://localhost:19999/api/4.0/permission_sets/5",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "6",
    "all_access": false,
    "name": "Viewer",
    "permissions": ["access_data", "clear_cache_refresh", "create_alerts", "download_without_limit", "follow_alerts", "login_special_email", "mobile_app_access", "schedule_look_emails", "see_drill_overlay", "see_lookml_dashboards", "see_looks", "see_user_dashboards", "send_to_integration"],
    "url": "https://localhost:19999/api/4.0/permission_sets/6",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "7",
    "all_access": false,
    "name": "Permission Set",
    "permissions": ["access_data", "download_with_limit", "schedule_external_look_emails", "schedule_look_emails", "see_user_dashboards"],
    "url": "https://localhost:19999/api/4.0/permission_sets/7",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  },
  {
    "built_in": false,
    "id": "8",
    "all_access": false,
    "name": "Client Admin",
    "permissions": ["access_data", "clear_cache_refresh", "create_alerts", "create_custom_fields", "create_prefetches", "create_public_looks", "create_table_calculations", "deploy", "develop", "download_with_limit", "download_without_limit", "explore", "follow_alerts", "login_special_email", "manage_spaces", "save_content", "schedule_external_look_emails", "schedule_look_emails", "see_drill_overlay", "see_lookml", "see_lookml_dashboards", "see_looks", "see_sql", "see_user_dashboards", "support_access_toggle", "use_sql_runner"],
    "url": "https://localhost:19999/api/4.0/permission_sets/8",
    "can": {
      "show": true,
      "index": true,
      "update": true
    }
  }
]`)
	})

	permSets, resp, err := client.Roles.PermissionSetsList(ctx, nil)
	_ = resp
	if err != nil {
		t.Errorf("Roles.PermissionSetsList returned error: %v", err)
	}

	expectedPermSets := []PermissionSet{
		{Id: 8, Name: "Client Admin"},
		{Id: 7, Name: "Permission set"}}
	if !reflect.DeepEqual(permSets, expectedPermSets) {
		t.Error(errGotWant("groups.List", permSets, expectedPermSets))
	}
}
