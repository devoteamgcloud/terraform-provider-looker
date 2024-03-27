package provider_test

import (
	"fmt"
	"testing"

	"github.com/devoteamgcloud/terraform-provider-looker/internal/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ─── ENV VARS ────────────────────────────────────────────────────────────────
// export TF_ACC="true";
// export LOOKER_BASE_URL="${LOOKER_BASE_URL:-$(pass looker/url)}"
// export LOOKER_API_CLIENT_SECRET="${LOOKER_API_CLIENT_SECRET:-$(pass looker/client_secret)}"
// export LOOKER_API_CLIENT_ID="${LOOKER_API_CLIENT_ID:-$(pass looker/client_id)}"
// export LOOKER_ALLOW_UNVERIFIED_SSL="true";
// ─────────────────────────────────────────────────────────────────────────────
// go test -v -timeout 30s -run ^TestDataSourceFolderPermissions$ github.com/devoteamgcloud/terraform-provider-looker/internal/provider
func TestDataSourceFolderPermissions(t *testing.T) {
	datasourceName := "data.looker_folder_permissions.this"
	var id uint = 1
	testSteps(t, []resource.TestStep{{
		Config: testDataSourceIdentityGroup_configName(id),
		Check: resource.ComposeTestCheckFunc(
			testutils.Resource(datasourceName, testutils.Attribute("id", testutils.IsNotEmpty())),
			testutils.Resource(datasourceName,
				testutils.Nested("permissions",
					testutils.CheckInList(
						// NOTE: "can" attribute is not tested as it may not be present
						// ─────────────────────────────────────────────────────────────────────────────
						// NOTE: user_id or group_id may be empty, that is why these are
						// commented out for now
						// testutils.Attribute("group_id", testutils.IsNotEmpty()),
						// testutils.Attribute("user_id", testutils.IsNotEmpty()),
						testutils.Attribute("id", testutils.IsNotEmpty()),
						testutils.Attribute("permission_type", testutils.IsNotEmpty()),
					),
				),
			),
		),
	}})
}
func testDataSourceIdentityGroup_configName(id uint) string {
	return fmt.Sprintf(`

data "looker_folder_permissions" "this" {
  id = %d
}
`, id)
}
