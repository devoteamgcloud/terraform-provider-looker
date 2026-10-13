package testutils

// NOTE: the functions here were copied from github.com/hashicorp/terraform-provider-vault/testutil/testutil.go
import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	EnvVarSkipVaultNext = "SKIP_LOOKER_NEXT_TESTS"
	EnvVarTfAccEnt      = "TF_ACC_ENTERPRISE"
)

func TestAccPreCheck(t *testing.T) {
	t.Helper()
	GetTestLookerCreds(t)
}

func SkipTestAcc(t *testing.T) {
	t.Helper()
	SkipTestEnvUnset(t, resource.EnvTfAcc)
}

// SkipTestEnvSet skips the test if any of the provided environment variables
// have a non-empty value.
func SkipTestEnvSet(t *testing.T, envVars ...string) []string {
	t.Helper()
	return handleTestEnvSetF(t, t.Skipf, envVars...)
}

// SkipTestEnvUnset skips the test if any of the provided environment variables
// are empty/unset.
func SkipTestEnvUnset(t *testing.T, envVars ...string) []string {
	t.Helper()
	return handleTestEnvUnsetF(t, t.Skipf, envVars...)
}

// FatalTestEnvUnset fails the test if any of the provided environment variables
// have non-empty values.
func FatalTestEnvUnset(t *testing.T, envVars ...string) []string {
	t.Helper()
	return handleTestEnvUnsetF(t, t.Fatalf, envVars...)
}

func handleTestEnvUnsetF(t *testing.T, f func(f string, args ...interface{}), envVars ...string) []string {
	t.Helper()
	return handleTestEnv(t, func(k, v string) {
		t.Helper()
		if v == "" {
			f("%q must be set", k)
		}
	}, envVars...)
}

func handleTestEnvSetF(t *testing.T, f func(f string, args ...interface{}), envVars ...string) []string {
	t.Helper()
	return handleTestEnv(t, func(k, v string) {
		t.Helper()
		if v != "" {
			f("%q is set", k)
		}
	}, envVars...)
}

func handleTestEnv(t *testing.T, f func(k string, v string), envVars ...string) []string {
	t.Helper()
	var result []string
	for _, k := range envVars {
		v := os.Getenv(k)
		f(k, v)
		result = append(result, v)
	}
	return result
}

type LookerTestConf struct {
	URL, ClientID, ClientSecret string
}

func GetTestLookerCreds(t *testing.T) *LookerTestConf {
	v := SkipTestEnvUnset(t,
		"LOOKER_BASE_URL",
		"LOOKER_API_CLIENT_ID",
		"LOOKER_API_CLIENT_SECRET")

	return &LookerTestConf{
		URL:          v[0],
		ClientID:     v[1],
		ClientSecret: v[2],
	}
}

func TestCheckResourceAttrJSON(name, key, expectedValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %q", name)
		}
		instanceState := resourceState.Primary
		if instanceState == nil {
			return fmt.Errorf("%q has no primary instance state", name)
		}
		v, ok := instanceState.Attributes[key]
		if !ok {
			return fmt.Errorf("%s: attribute not found %q", name, key)
		}
		if expectedValue == "" && v == expectedValue {
			return nil
		}
		if v == "" {
			return fmt.Errorf("%s: attribute %q expected %#v, got %#v", name, key, expectedValue, v)
		}

		var stateJSON, expectedJSON interface{}
		err := json.Unmarshal([]byte(v), &stateJSON)
		if err != nil {
			return fmt.Errorf("%s: attribute %q not JSON: %s", name, key, err)
		}
		err = json.Unmarshal([]byte(expectedValue), &expectedJSON)
		if err != nil {
			return fmt.Errorf("expected value %q not JSON: %s", expectedValue, err)
		}
		if !reflect.DeepEqual(stateJSON, expectedJSON) {
			return fmt.Errorf("%s: attribute %q expected %#v, got %#v", name, key, expectedJSON, stateJSON)
		}
		return nil
	}
}

func GetResourceFromRootModule(s *terraform.State, resourceName string) (*terraform.ResourceState, error) {
	if rs, ok := s.RootModule().Resources[resourceName]; ok {
		return rs, nil
	}

	return nil, fmt.Errorf("expected resource %q, not found in state", resourceName)
}

// CheckJSONData from an expected string for a given resource attribute.
func CheckJSONData(resourceName, attr, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, err := GetResourceFromRootModule(s, resourceName)
		if err != nil {
			return err
		}

		actual, ok := rs.Primary.Attributes[attr]
		if !ok {
			return fmt.Errorf("resource %q has no attribute %q", resourceName, attr)
		}

		var e map[string]interface{}
		if err := json.Unmarshal([]byte(expected), &e); err != nil {
			return nil
		}

		var a map[string]interface{}
		if err := json.Unmarshal([]byte(actual), &a); err != nil {
			return nil
		}

		if !reflect.DeepEqual(e, a) {
			return fmt.Errorf("expected %#v, got %#v for resource attr %s.%s", e, a, resourceName, attr)
		}

		return nil
	}
}

// GetImportTestStep for resource name. If a custom ImportStateCheck function is not desired, pass
// a nil value. Optionally include field names that should be ignored during the import
// verification, typically ignore fields should only be provided for values that are not returned
// from the provisioning API.
func GetImportTestStep(resourceName string, skipVerify bool, check resource.ImportStateCheckFunc, ignoreFields ...string) resource.TestStep {
	ts := resource.TestStep{
		ResourceName:            resourceName,
		ImportState:             true,
		ImportStateVerify:       !skipVerify,
		ImportStateVerifyIgnore: ignoreFields,
	}

	if check != nil {
		ts.ImportStateCheck = check
	}

	return ts
}
