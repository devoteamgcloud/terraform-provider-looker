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

	_ = c
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

/*

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	var projectName string

	if d.Get("discriminator").(bool) {
		projectName = fmt.Sprintf("%s-%s", d.Get("name").(string), srand(4))
	} else {
		projectName = fmt.Sprintf("%s", d.Get("name").(string))
	}

	np := new(lookergo.Project)

	np.Name = projectName

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: create project", currFuncName()))

	session, _, err := c.Sessions.SetWorkspaceId(ctx, "dev")
	if err != nil {
		return diag.FromErr(err)
	} else if session.WorkspaceId != "dev" {
		return diag.Errorf("Could not set workspace to 'dev'")
	}

	project, _, err := c.Projects.Create(ctx, np)
	if err != nil {
		switch err.(type) {
		case *lookergo.ErrorResponse:
			if len(err.(*lookergo.ErrorResponse).Errors) >= 1 {
				for _, errRespErr := range err.(*lookergo.ErrorResponse).Errors {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  fmt.Sprintf("field: %v, code: %v, msg: %v ", errRespErr.Field, errRespErr.Code, errRespErr.Message),
						AttributePath: cty.Path{
							cty.GetAttrStep{Name: errRespErr.Field},
						},
					})
				}
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.(*lookergo.ErrorResponse).Message,
			}, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		default:
			return diag.FromErr(err)
		}
	}

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
	var project *lookergo.Project
	project, _, err := c.Projects.Get(ctx, projectId)
	if err != nil {
		switch err.(type) {
		case *lookergo.ErrorResponse:
			if err.(*lookergo.ErrorResponse).Response.StatusCode == http.StatusNotFound {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Could not find project with id '%v', check if it only exists in dev.", projectId),
				})
				session, _, err := c.Sessions.SetWorkspaceId(ctx, "dev")
				if err != nil {
					return diag.FromErr(err)
				} else if session.WorkspaceId != "dev" {
					return diag.Errorf("Could not set workspace to 'dev'")
				}

				project, _, err = c.Projects.Get(ctx, projectId)
				if err != nil {
					return diag.FromErr(err)
				}
			} else {
				return diag.FromErr(err)
			}
		default:
			return diag.FromErr(err)
		}
	}

	if project.Name != "" {
		d.Set("name", project.Name)
	}
	if project.IsExample != nil {
		d.Set("is_example", project.IsExample)
	}

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
	var project *lookergo.Project
	project, _, err := c.Projects.Get(ctx, projectId)
	if err != nil {
		switch err.(type) {
		case *lookergo.ErrorResponse:
			if err.(*lookergo.ErrorResponse).Response.StatusCode == http.StatusNotFound {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Could not find project with id '%v', check if it only exists in dev.", projectId),
				})
				session, _, err := c.Sessions.SetWorkspaceId(ctx, "dev")
				if err != nil {
					return append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  err.Error(),
					})
				} else if session.WorkspaceId != "dev" {
					return diag.Errorf("Could not set workspace to 'dev'")
				}

				project, _, err = c.Projects.Get(ctx, projectId)
				if err != nil {
					return append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  err.Error(),
					})
				}
			} else {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  err.Error(),
				})
			}
		default:
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
	}

	if d.Get("discriminator").(bool) {

		patchProject := lookergo.Project{Name: fmt.Sprintf("deleteme-%s", project.Name)}
		projectUpdated, _, err := c.Projects.Update(ctx, projectId, &patchProject)
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf("Project deletion not supported by the API, "+
				"Project name has been updated to '%s' instead.", projectUpdated.Name),
			AttributePath: cty.Path{
				cty.GetAttrStep{Name: "name"},
			},
		})
	}

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
*/
