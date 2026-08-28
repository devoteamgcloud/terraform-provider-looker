package provider_test

// NOTE:  here were copied from github.com/hashicorp/terraform-provider-vault/vault/provider_test.go
import (
	"sync"
	"testing"

	"github.com/devoteamgcloud/terraform-provider-looker/internal/provider"
	"github.com/devoteamgcloud/terraform-provider-looker/internal/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const providerName = "looker"

var (
	devProvider     *schema.Provider
	providerLock    sync.Mutex
	providerVersion string = "0.0.0"
)

func testProvider() *schema.Provider {
	providerLock.Lock()
	defer providerLock.Unlock()
	if devProvider == nil {
		devProvider = provider.New(providerVersion)()

	}
	return devProvider
}

func testSteps(t *testing.T, steps []resource.TestStep) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
			testutils.SkipTestAcc(t)
		},
		Providers: map[string]*schema.Provider{
			"looker": testProvider(),
		},
		Steps: steps,
	})
}
