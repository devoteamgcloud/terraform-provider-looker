package acceptance_tests

import (
	"context"
	"strings"
	"testing"

	lookergo "github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
)

// go test -v -timeout 30s -run ^TestFolder$ github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo/acceptance_tests
func TestFolder(t *testing.T) {
	// NOTE: the expectation is that the client token has full admin permissions
	client := lookergo.NewClient(nil)
	old_url := ac.BaseURL
	client.DisableTLSVerification()
	newURL := strings.TrimSuffix(old_url, "/")
	if !strings.HasSuffix(newURL, "/api") {
		newURL += "/api/"
	} else {
		newURL += "/"
	}
	err := client.SetBaseURL(newURL)
	assert.NoError(t, err)

	err = client.SetOauthCredentials(ctx, ac.ClientId, ac.ClientSecret)
	assert.NoError(t, err)

	user, resp, err := client.Sessions.GetCurrentUser(ctx)
	assert.NotNil(t, user)
	assert.NotNil(t, resp)
	assert.NoError(t, err)

	t.Run("get", func(t *testing.T) {
		folder, resp, err := client.Folders.Get(context.Background(), "1")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, folder)
	})
	t.Run("create", func(t *testing.T) {
		folderName := acctest.RandomWithPrefix("test-" + t.Name())
		folder := &lookergo.Folder{
			Name:     folderName,
			ParentId: "1",
		}
		folder, resp, err = client.Folders.Create(context.Background(), folder)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, folder)
		resp, err = client.Folders.Delete(context.Background(), folder.Id)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
	t.Run("get_permissions", func(t *testing.T) {
		// NOTE: this simple tests assumes that the 'Shared' folder has permissions
		// set on it which should be the default
		metadata, resp, err := client.ContentMetaGroupUser.ListByID(ctx, "10", nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, metadata)
		pp.Print(metadata)
	})
	t.Run("get-permissions-empty", func(t *testing.T) {
		// NOTE: this simple tests assumes that the 'Shared' folder has permissions
		// set on it which should be the default
		folderName := acctest.RandomWithPrefix("test-" + t.Name())
		folder := &lookergo.Folder{
			Name:     folderName,
			ParentId: "1",
		}
		folder, resp, err = client.Folders.Create(context.Background(), folder)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, folder)
		defer func() {
			resp, err = client.Folders.Delete(context.Background(), folder.Id)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		}()
		metadata, resp, err := client.ContentMetaGroupUser.ListByID(ctx, folder.ContentMetadataId, nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, metadata)
	})
}
