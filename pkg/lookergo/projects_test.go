package lookergo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjectsResourceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/4.0/projects/sandbox-with-sand", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
  "id": "sandbox-with-sand",
  "name": "sandbox-with-sand",
  "git_remote_url": "git@github.com:Sandbox/sandbox-with-sand.git",
  "git_username": "",
  "git_production_branch_name": "main",
  "use_git_cookie_auth": false,
  "git_username_user_attribute": null,
  "git_password_user_attribute": null,
  "git_service_name": "github",
  "git_application_server_http_port": null,
  "git_application_server_http_scheme": null,
  "pull_request_mode": "off",
  "validation_required": true,
  "git_release_mgmt_enabled": false,
  "allow_warnings": true,
  "uses_git": true,
  "is_example": false,
  "dependency_status": "install_none",
  "can": {
    "webhook_deploy": true,
    "show_manifest": true,
    "index": true,
    "show": true,
    "validate": true,
    "link_to_service": false,
    "update": true,
    "view_git_deploy_key": true,
    "show_branches": true,
    "deploy_ref_to_production": true
  }
}`)
	})

	result, resp, err := client.Projects.Get(ctx, "sandbox-with-sand")
	_ = resp
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	expected := &Project{
		Id:                      "sandbox-with-sand",
		Name:                    "sandbox-with-sand",
		UsesGit:                 boolPtr(true),
		GitRemoteUrl:            "git@github.com:Sandbox/sandbox-with-sand.git",
		GitProductionBranchName: "main",
		GitServiceName:          "github",
		PullRequestMode:         "off",
		ValidationRequired:      boolPtr(true),
		AllowWarnings:           boolPtr(true),
		DependencyStatus:        "install_none",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error(errGotWant("Projects.Get", result, expected))
	}
}
