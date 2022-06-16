package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the connection. Also used as the unique identifier",
				Type:        schema.TypeString,
				Required:    true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"file_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_billing_gigabytes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tmp_db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdbc_additional_params": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dialect_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"maintenance_cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"after_connect_statements": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth_application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_db_credentials": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sql_runner_precache_tables": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sql_writing_with_info_schema": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_context_comment": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"always_retry_failed_builds": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cost_estimate_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pdt_api_control_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pool_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pdt_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	connectionName := d.Get("name").(string)

	nc := new(lookergo.DBConnection)
	nc.Name = connectionName

	if val, ok := d.GetOk("host"); ok {
		nc.Host = val.(string)
	}
	if val, ok := d.GetOk("port"); ok {
		nc.Port = val.(string)
	}
	if val, ok := d.GetOk("username"); ok {
		nc.Username = val.(string)
	}
	if val, ok := d.GetOk("password"); ok {
		nc.Password = val.(string)
	}
	if val, ok := d.GetOk("certificate"); ok {
		nc.Certificate = val.(string)
	}
	if val, ok := d.GetOk("file_type"); ok {
		nc.FileType = val.(string)
	}
	if val, ok := d.GetOk("database"); ok {
		nc.Database = val.(string)
	}
	if val, ok := d.GetOk("db_timezone"); ok {
		nc.DbTimezone = val.(string)
	}
	if val, ok := d.GetOk("query_timezone"); ok {
		nc.QueryTimezone = val.(string)
	}
	if val, ok := d.GetOk("schema"); ok {
		nc.Schema = val.(string)
	}
	if val, ok := d.GetOk("max_billing_gigabytes"); ok {
		nc.MaxBillingGigabytes = val.(string)
	}
	if val, ok := d.GetOk("tmp_db_name"); ok {
		nc.TmpDbName = val.(string)
	}
	if val, ok := d.GetOk("jdbc_additional_params"); ok {
		nc.JdbcAdditionalParams = val.(string)
	}
	if val, ok := d.GetOk("dialect_name"); ok {
		nc.DialectName = val.(string)
	}
	if val, ok := d.GetOk("maintenance_cron"); ok {
		nc.MaintenanceCron = val.(string)
	}
	if val, ok := d.GetOk("after_connect_statements"); ok {
		nc.AfterConnectStatements = val.(string)
	}
	if val, ok := d.GetOk("tunnel_id"); ok {
		nc.TunnelId = val.(string)
	}
	if val, ok := d.GetOk("oauth_application_id"); ok {
		nc.OauthApplicationId = val.(string)
	}
	if val, ok := d.GetOk("ssl"); ok {
		nc.Ssl = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("verify_ssl"); ok {
		nc.VerifySsl = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("user_db_credentials"); ok {
		nc.UserDbCredentials = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("sql_runner_precache_tables"); ok {
		nc.SqlRunnerPrecacheTables = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("sql_writing_with_info_schema"); ok {
		nc.SqlWritingWithInfoSchema = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("disable_context_comment"); ok {
		nc.DisableContextComment = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("always_retry_failed_builds"); ok {
		nc.AlwaysRetryFailedBuilds = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("cost_estimate_enabled"); ok {
		nc.CostEstimateEnabled = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("pdt_api_control_enabled"); ok {
		nc.PdtApiControlEnabled = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("max_connections"); ok {
		nc.MaxConnections = int64(val.(int))
	}
	if val, ok := d.GetOk("pool_timeout"); ok {
		nc.PoolTimeout = int64(val.(int))
	}
	if val, ok := d.GetOk("pdt_concurrency"); ok {
		nc.PdtConcurrency = int64(val.(int))
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: test conf", currFuncName()))
	validateConfig, _, err := c.Connections.ValidateConfig(ctx, nc)
	if err != nil {
		switch err.(type) {
		case *json.InvalidUnmarshalError:
			return diag.FromErr(err)
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

		return diags
	}

	for i, dbcv := range validateConfig {
		tflog.Debug(ctx, fmt.Sprintf("Config [%v] validation message: %v", i, dbcv.Message))
		tflog.Debug(ctx, fmt.Sprintf("Config [%v] validation status: %v", i, dbcv.Status))
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: create", currFuncName()))
	connection, _, err := c.Connections.Create(ctx, nc)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: created, Name: %v", currFuncName(), connection.Name))

	d.SetId(connection.Name)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	connectionName := d.Get("name").(string)
	connection, _, err := c.Connections.Get(ctx, connectionName)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("name", "host", "port", "username") {
		if d.Set("name", connection.Name) != nil {
			return diag.FromErr(err)
		}
		if d.Set("host", connection.Host) != nil {
			return diag.FromErr(err)
		}
		if d.Set("port", connection.Port) != nil {
			return diag.FromErr(err)
		}
		if d.Set("username", connection.Username) != nil {
			return diag.FromErr(err)
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	connectionName := d.Get("name").(string)
	_, _, err := c.Connections.Get(ctx, connectionName)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "WARNING connection update is not implemented",
		AttributePath: cty.Path{
			cty.GetAttrStep{Name: "name"},
		},
	})

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	connectionName := d.Id()
	_, err := c.Connections.Delete(ctx, connectionName)
	if err != nil {
		return diag.FromErr(err)
	}

	// Finally mark as deleted
	d.SetId("")
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}
