package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
	"sync"
	"time"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"rename_when_delete": {
				Description: "The looker API doesn't provide means to delete a project" +
					"however, we can rename the project from orig_name to deleteme-orig_name-X1Y2",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": {
				Description: "Name of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_example": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	project := &lookergo.Project{
		Name:      d.Get("name").(string),
		IsExample: boolPtr(d.Get("discriminator").(bool)),
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: create project", currFuncName()))

	// Create dev client
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: create dev client connection", currFuncName()))
	devClient, _, err := c.CreateDevConnection(ctx)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	_ = c
	_ = devClient
	// TODO

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: project created, Id: %v", currFuncName(), project.Id))
	d.SetId(project.Name)

	resourceProjectRead(ctx, d, m)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	projectId := d.Id()

	_, _ = c, projectId

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	// TODO
	_ = c

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	projectId := d.Id()

	// TODO
	_, _ = c, projectId

	// API not available
	// _, err = c.Projects.Delete(ctx, project.Id)

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

// rMutex ensures random seed
var rMutex sync.Mutex

// generates a random string of fixed size
func srand(size int) string {
	var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	rMutex.Lock()
	rand.Seed(time.Now().UnixNano())
	defer rMutex.Unlock()
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf)
}
