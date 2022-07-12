package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
				Description: "Name of the project. => lowercase letters, numbers, underscores, and dashes",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					value := v.(string)
					var diags diag.Diagnostics
					re := regexp.MustCompile(`(?m)^([a-z0-9]|_|-)*$`)
					match := re.FindString(value)
					if match == "" {
						diag1 := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "only lowercase letters, numbers, underscores, and dashes allowed",
						}
						diags = append(diags, diag1)
					}
					return diags
				},
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
	// c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, m.(*Config).Api, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	/*	if err := ensureDevClient(ctx, m); err != nil {
			return diagErrAppend(diags, err)
		}
		dc := m.(*Config).DevClient*/

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	project := &lookergo.Project{
		Name:      d.Get("name").(string),
		IsExample: boolPtr(d.Get("is_example").(bool)),
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: create project", currFuncName()))

	response, _, err := dc.Projects.Create(ctx, project)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: project created, Id: %v", currFuncName(), response.Id))
	d.SetId(project.Name)

	resourceProjectRead(ctx, d, m)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: project created, Id: %v", currFuncName(), project.Id))
	d.SetId(project.Name)

	resourceProjectRead(ctx, d, m)

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	if err := ensureDevClient(ctx, m); err != nil {
		return diagErrAppend(diags, err)
	}
	c := m.(*Config).Api // .(*lookergo.Client)
	dc := m.(*Config).DevClient
	// Refresh token for dev Api connection if not used before.
	err := dc.EnsureStaticToken(ctx, c, m.(*Config).ApiUserID)
	if err != nil {
		return diagErrAppend(diags, err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	projectId := d.Id()
	var project *lookergo.Project
	project, _, err = dc.Projects.Get(ctx, projectId)
	if err != nil {
		return diagErrAppend(diags, err)
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
	// c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))
	if err := ensureDevClient(ctx, m); err != nil {
		return diagErrAppend(diags, err)
	}
	dc := m.(*Config).DevClient

	projectId := d.Id()

	project, _, err := dc.Projects.Get(ctx, projectId)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	// API not available
	deletedProject := &lookergo.Project{Name: fmt.Sprintf("deleteme-%s-%s", project.Name, srand(4))}

	_, resp, err := dc.Projects.Update(ctx, projectId, deletedProject)
	if resp.StatusCode == http.StatusInternalServerError {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf("Tried to rename project '%s' to '%s', "+
				"but got internal server error", project.Name, deletedProject.Name),
			Detail: fmt.Sprintf("%v\nA server 500 error is to be expected. "+
				"The resource might actually have been renamed. "+
				"We will check if the previous resource still exists to verify it has been renamed.\n"+
				"Err:%v", `¯\_(ツ)_/¯`, err.Error()),
		})
		tflog.Debug(ctx, fmt.Sprintf("Action: tried renaming project, but got err 500"),
			map[string]interface{}{"orig_name": project.Name, "new_name": deletedProject.Name})
		time.Sleep(5 * time.Second)
	} else if 300 <= resp.StatusCode && err != nil {
		return diagErrAppend(diags, err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Action: renamed project, New name: %v", deletedProject.Name))

	_, resp, err = dc.Projects.Get(ctx, projectId)
	if resp.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, fmt.Sprintf("Action: project '%v' not found, so let's assume it has been renamed to '%s'.", project.Name, deletedProject.Name))
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Project '%v' not found, so let's assume it has been renamed to '%s'.", project.Name, deletedProject.Name),
			Detail:   fmt.Sprintf("%v\n Err:%v ", `¯\_(ツ)_/¯`, err.Error()),
		})
	} else if 300 <= resp.StatusCode && err != nil {
		return diagErrAppend(diags, err)
	}
	time.Sleep(5 * time.Second)

	renamedProject, _, err := dc.Projects.Get(ctx, deletedProject.Name)
	if renamedProject != nil && err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Renamed project '%v' was found, so it's safe to assume '%s' has been deleted", deletedProject.Name, project.Name),
			Detail:   `¯\_(ツ)_/¯`,
		})
	} else if err != nil {
		return diagErrAppend(diags, err)
	}

	// if it works, it works ...
	uaccEmail := os.Getenv("LOOKER_USERACC_EMAIL")
	uaccPass := os.Getenv("LOOKER_USERACC_PASS")
	if uaccPass != "" && uaccEmail != "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("will use private API with email '%v' and pass ****", uaccEmail),
		})
		err := dodgyProjectDelete(dc.BaseURL, uaccEmail, uaccPass, deletedProject.Name)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Could not delete using private API",
			})
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Project fully deleted!",
			})
		}
	}

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

func dodgyProjectDelete(u *url.URL, email string, pass string, projectName string) error {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"))

	var csrfToken string
	var csrfToken2 string
	var csrfToken3 string
	var devmode string

	c.OnHTML("head", func(e *colly.HTMLElement) {
		csrfToken = e.ChildAttr("meta[name=csrf-token]", "content")
	})
	err := c.Visit((&url.URL{Scheme: u.Scheme, Host: u.Host, Path: "login/email"}).String())
	if err != nil {
		return err
	}

	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println("response received", r.StatusCode)
	// })
	err = c.Post((&url.URL{Scheme: u.Scheme, Host: u.Host, Path: "login"}).String(), map[string]string{"csrf-token": csrfToken, "email": email, "password": pass})
	if err != nil {
		return err
	}

	c.OnHTML("head", func(e *colly.HTMLElement) {
		csrfToken2 = e.ChildAttr("meta[name=csrf-token]", "content")
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		devmode = e.ChildAttr("a[lk-track-name=dev_mode]", "lk-track-action")
	})
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println("response received", r.StatusCode)
	// })
	projectsUrl := (&url.URL{Scheme: u.Scheme, Host: u.Host, Path: "projects"}).String()
	err = c.Visit(projectsUrl)
	if err != nil {
		return err
	}

	// Set dev mode
	if devmode != "Exit Development Mode" {
		// c.OnResponse(func(r *colly.Response) {
		// 	log.Println("response received", r.StatusCode)
		// })
	}
	err = c.Post((&url.URL{Scheme: u.Scheme, Host: u.Host, Path: "account/developer-mode/enter"}).String(),
		map[string]string{"csrf-token": csrfToken2, "_method": "put"})
	if err != nil {
		// so be it…
		// return err
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", projectsUrl)
	})
	c.OnHTML("head", func(e *colly.HTMLElement) {
		csrfToken3 = e.ChildAttr("meta[name=csrf-token]", "content")
	})
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println("response received", r.StatusCode)
	// })
	confirmDeleteUrl := (&url.URL{Scheme: u.Scheme, Host: u.Host, Path: fmt.Sprintf("projects/%s/confirm_delete", projectName)}).String()
	err = c.Visit(confirmDeleteUrl)
	if err != nil {
		return err
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", confirmDeleteUrl)
	})
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println("response received", r.StatusCode)
	// })
	doDeleteUrl := (&url.URL{Scheme: u.Scheme, Host: u.Host, Path: fmt.Sprintf("projects/%s", projectName)}).String()
	err = c.Post(doDeleteUrl, map[string]string{"csrf-token": csrfToken3, "_method": "DELETE"})
	if err != nil {
		return err
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", projectsUrl)
	})
	err = c.Post((&url.URL{Scheme: u.Scheme, Host: u.Host, Path: "account/developer-mode/leave"}).String(),
		map[string]string{"csrf-token": csrfToken2, "_method": "put"})
	if err != nil {
		return err
	}

	return nil

}
