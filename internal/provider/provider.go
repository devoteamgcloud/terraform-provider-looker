package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultTimeout     = 5 * time.Minute
	minimumRefreshWait = 3 * time.Second
	checkDelay         = 10 * time.Second
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
	//   /**/-
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"base_url": {
					Description: "For base_url, provide the URL including /api/ ! " +
						"Normally, a REST API should not have api in it's path, " +
						"therefore we don't add the /api/ inside the provider. ",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("LOOKER_BASE_URL", nil),
				},
				"client_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_CLIENT_ID", nil),
				},
				"client_secret": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_CLIENT_SECRET", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"looker_user":  dataSourceUser(),
				"looker_group": dataSourceGroup(),
				"looker_project": dataSourceProject(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"looker_user":                   resourceUser(),
				"looker_group":                  resourceGroup(),
				"looker_group_member":           resourceGroupMember(),
				"looker_role":                   resourceRole(),
				"looker_role_groups":            resourceRoleMember(),
				"looker_connection":             resourceConnection(),
				"looker_project":                resourceProject(),
				"looker_project_git_deploy_key": resourceProjectGitDeployKey(),
				"looker_project_git_repo":       resourceProjectGitRepo(),
				"looker_lookml_model":           resourceLookMlModel(),
				"looker_model_set":              resourceModelSet(),
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return providerConfigure(ctx, d, p, version)
		}

		return p
	}
}

type Workspace int

const (
	WorkspaceProduction Workspace = iota
	WorkspaceDev
)

type Config struct {
	Api                       *lookergo.Client
	ApiUserID                 string
	DevClient                 *lookergo.Client
	Workspace                 Workspace
	RequestCompletionCallback lookergo.RequestCompletionCallback
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider, version string) (interface{}, diag.Diagnostics) {
	tflog.Debug(ctx, "Configure provider", map[string]interface{}{"conninfo": d.ConnInfo(), "schema": p.Schema})
	tflog.Debug(ctx, "Provider config", map[string]interface{}{"client_id": d.Get("client_id").(string)})
	
	userAgent := p.UserAgent("terraform-provider-looker", version)
	var diags diag.Diagnostics

	client := lookergo.NewClient(nil)
	devClient := lookergo.NewClient(nil)

	old_url := d.Get("base_url").(string)
	var newURL string
	if len(old_url) > 5{
		switch old_url[len(old_url)-4:] {
		case "api/":
			newURL = old_url
		case ".com":
			newURL = old_url + "/api/"
		case "com/":
			newURL = old_url + "api/"
		case "/api":
			newURL = old_url + "/"
		}
	}else {
		newURL = old_url
	}

	if err := client.SetBaseURL(newURL); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set looker API endpoint",
			Detail:   "Err: " + err.Error(),
		})
		return nil, diags
	}
	devClient.SetBaseURL(newURL)

	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	if err := client.SetOauthCredentials(ctx, clientId, clientSecret); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set looker API client_id/client_secret",
			Detail:   "Err: " + err.Error(),
		})
		return nil, diags
	}
	if err := client.SetUserAgent(userAgent); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set looker API user-agent",
			Detail:   "Err: " + err.(*lookergo.ErrorResponse).Message,
		})
		return nil, diags
	}
	devClient.SetUserAgent(userAgent)

	var config Config

	config.RequestCompletionCallback = func(req *http.Request, resp *http.Response) {
		if code := resp.StatusCode; code >= 200 && code <= 299 {
			tflog.Debug(ctx, "HTTP Request", map[string]interface{}{"req_url": req.URL.String(), "req_method": req.Method, "resp_status": resp.Status})
		} else {
			tflog.Debug(ctx, "HTTP Error", map[string]interface{}{"req_url": req.URL.String(), "req_method": req.Method, "resp_status": resp.Status, "resp_length": resp.ContentLength, "resp_headers": resp.Header})
		}
	}

	client.OnRequestCompleted(config.RequestCompletionCallback)
	devClient.OnRequestCompleted(config.RequestCompletionCallback)

	session, _, err := client.Sessions.Get(ctx)
	if err != nil {
		errMsg := err.Error()
		var errBodyMd string
		// Otherwise the error is the full html doc. :O
		if reflect.TypeOf(err).Elem() == reflect.TypeOf((*url.Error)(nil)).Elem() {
			if reflect.TypeOf(err.(*url.Error).Err).Elem() == reflect.TypeOf((*oauth2.RetrieveError)(nil)).Elem() {
				errMsg = fmt.Sprintf("oauth2: cannot fetch token: %v", err.(*url.Error).Err.(*oauth2.RetrieveError).Response.Status)
				converter := md.NewConverter("", true, nil)
				mdBytes, _ := converter.ConvertBytes(err.(*url.Error).Err.(*oauth2.RetrieveError).Body)
				errBodyMd = string(mdBytes)
				// errBodyMd = string(err.(*url.Error).Err.(*oauth2.RetrieveError).Body)
			}
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Looker client: " + errMsg,
			Detail:   "Unable to authenticate user for authenticated Looker client\n" + errBodyMd,
		})
		return nil, diags
	}

	// Get current user ID
	user, _, err := client.Sessions.GetCurrentUser(ctx)
	if err != nil {
		return nil, diagErrAppend(diags, err)
	}

	switch session.WorkspaceId {
	case "production":
		config = Config{Api: client, ApiUserID: strconv.Itoa(user.Id), DevClient: devClient, Workspace: WorkspaceProduction}
	case "dev":
		config = Config{Api: client, ApiUserID: strconv.Itoa(user.Id), DevClient: devClient, Workspace: WorkspaceDev}
	default:
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Session Workspace ID is nome of production/dev",
			Detail:   "Unable to find workspace ID in /session call",
		})
		return nil, diags
	}

	return &config, nil
}

func diagErrAppend(diags diag.Diagnostics, err error) diag.Diagnostics {
	switch err.(type) {
	case *lookergo.ErrorResponse:
		if len(err.(*lookergo.ErrorResponse).Errors) >= 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.(*lookergo.ErrorResponse).Message,
			})
			for _, errRespErr := range err.(*lookergo.ErrorResponse).Errors {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("field: %v, code: %v", errRespErr.Field, errRespErr.Code),
					Detail:   errRespErr.Message,
					AttributePath: cty.Path{
						cty.GetAttrStep{Name: errRespErr.Field},
					},
				})
			}
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
	default:
		return diag.FromErr(err)
	}
	return diags
}

func ensureDevClient(ctx context.Context, m interface{}) error {
	if m.(*Config).DevClient == nil {
		tflog.Debug(ctx, fmt.Sprintf("Fn: %v, Action: create dev client connection", currFuncName()))
		devClient, _, err := m.(*Config).Api.CreateDevConnection(ctx, func(req *http.Request, resp *http.Response) {
			if code := resp.StatusCode; code >= 200 && code <= 299 {
				tflog.Debug(ctx, "DevClient: HTTP Request", map[string]interface{}{"req_url": req.URL.String(), "req_method": req.Method, "resp_status": resp.Status})
			} else {
				tflog.Debug(ctx, "DevClient: HTTP Error", map[string]interface{}{"req_url": req.URL.String(), "req_method": req.Method, "resp_status": resp.Status, "resp_length": resp.ContentLength, "resp_headers": resp.Header})
			}
		})
		if err != nil {
			return err
		} else {
			m.(*Config).DevClient = devClient
		}
	}
	return nil
}
