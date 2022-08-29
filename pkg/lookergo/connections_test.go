package lookergo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestConnectionsResourceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/connections/testingpsql", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `  {
    "name": "testingpsql",
    "snippets": [
      {
        "name": "show_processes",
        "label": "Show Processes",
        "sql": "SELECT * FROM pg_stat_activity"
      }
    ],
    "host": "surus.db.elephantsql.com",
    "port": "5432",
    "database": "hegdgxme",
    "db_timezone": null,
    "query_timezone": null,
    "schema": null,
    "max_connections": 5,
    "max_billing_gigabytes": null,
    "ssl": true,
    "verify_ssl": false,
    "tmp_db_name": null,
    "jdbc_additional_params": "",
    "pool_timeout": 120,
    "created_at": "2022-06-14T12:26:50.000+00:00",
    "user_id": "123",
    "user_attribute_fields": [],
    "maintenance_cron": null,
    "last_regen_at": "1655254550",
    "last_reap_at": "",
    "sql_runner_precache_tables": true,
    "sql_writing_with_info_schema": false,
    "after_connect_statements": "",
    "pdt_concurrency": 1,
    "disable_context_comment": false,
    "oauth_application_id": null,
    "always_retry_failed_builds": null,
    "cost_estimate_enabled": false,
    "pdt_api_control_enabled": false,
    "dialect": {
      "supports_cost_estimate": true,
      "cost_estimate_style": "configurable",
      "automatically_run_sql_runner_snippets": true,
      "connection_tests": [
        "connect",
        "kill",
        "query",
        "database_timezone",
        "database_version",
        "tmp_db",
        "cdt",
        "tmp_db_views"
      ],
      "supports_inducer": false,
      "supports_multiple_databases": false,
      "supports_persistent_derived_tables": true,
      "has_ssl_support": true,
      "name": "postgres",
      "label": "PostgreSQL 9.5+",
      "supports_streaming": true,
      "persistent_table_indexes": "explicit",
      "persistent_table_sortkeys": "",
      "persistent_table_distkey": ""
    },
    "dialect_name": "postgres",
    "example": false,
    "managed": false,
    "pdts_enabled": false,
    "username": "hegdgxme",
    "uses_oauth": false,
    "pdt_context_override": null,
    "tunnel_id": "",
    "uses_application_default_credentials": false,
    "impersonated_service_account": null,
    "can": {
      "index": true,
      "index_limited": true,
      "show": true,
      "cost_estimate": true,
      "run_sql_queries": true,
      "access_data": true,
      "explore": true,
      "refresh_schemas": true,
      "destroy": true,
      "test": true,
      "create": true,
      "update": true
    }
}`)
	})

	result, resp, err := client.Connections.Get(ctx, "testingpsql")
	_ = resp
	if err != nil {
		t.Errorf("Connections.Get returned error: %v", err)
	}

	expected := &DBConnection{
		Name:        "testingpsql",
		Dialect:     &DBDialect{Name: "postgres"},
		DialectName: "postgres",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error(errGotWant("ModelSets.List", result, expected))
	}

}
