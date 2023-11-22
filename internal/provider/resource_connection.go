package provider

import (
	"context"
	"fmt"
	"github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
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
				ForceNew:    true,
			},
			"host": {
				Description: "Host name/address of server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"port": {
				Description: "Port number on server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"username": {
				Description: "Username for server authentication",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "(Write-Only) Password for server authentication",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"certificate": {
				Description: "(Write-Only) Base64 encoded Certificate body for server authentication (when appropriate for dialect).",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"file_type": {
				Description: "(Write-Only) Certificate keyfile type - .json or .p12",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"database": {
				Description: "Database name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"db_timezone": {
				Description: "Time zone of database",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"query_timezone": {
				Description: "Timezone to use in queries",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"schema": {
				Description: "Scheme name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"max_billing_gigabytes": {
				Description: "Maximum size of query in GBs (BigQuery only, can be a user_attribute name)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tmp_db_name": {
				Description: "Name of temporary database (if used)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"jdbc_additional_params": {
				Description: "Additional params to add to JDBC connection string",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dialect_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: dialect_names_tab,
			},
			"maintenance_cron": {
				Description: "Cron string specifying when maintenance such as PDT trigger checks and drops should be performed",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tunnel_id": {
				Description: "The Id of the ssh tunnel this connection uses",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"oauth_application_id": {
				Description: "An External OAuth Application to use for authenticating to the database",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ssl": {
				Description: "Use SSL/TLS when connecting to server",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"verify_ssl": {
				Description: "Verify the SSL",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"user_db_credentials": {
				Description: "(Limited access feature) Are per user db credentials enabled. Enabling will remove previously set username and password",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"sql_runner_precache_tables": {
				Description: "Precache tables in the SQL Runner",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"sql_writing_with_info_schema": {
				Description: "Fetch Information Schema For SQL Writing",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"disable_context_comment": {
				Description: "When disable_context_comment is true comment will not be added to SQL",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"always_retry_failed_builds": {
				Description: "When true, error PDTs will be retried every regenerator cycle",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"cost_estimate_enabled": {
				Description: "When true, query cost estimate will be displayed in explore",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"pdt_api_control_enabled": {
				Description: "PDT builds on this connection can be kicked off and cancelled via API",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"max_connections": {
				Description: "Maximum number of concurrent connection to use",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"pool_timeout": {
				Description: "Connection Pool Timeout, in seconds",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"pdt_concurrency": {
				Description: "Maximum number of threads to use to build PDTs in parallel",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"pdt_context_override": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"context": {
							Description: "Context in which to override (`pdt` is the only allowed value)",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "pdt",
						},
						"host": {
							Description: "Host name/address of server",
							Type:        schema.TypeString,
							Required:    true,
						},
						"username": {
							Description: "Username for server authentication",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"port": {
							Description: "Port number on server",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"password": {
							Description: "(Write-Only) Password for server authentication",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"has_password": {
							Description: "Whether or not the password is overridden in this context",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     false,
						},
						"certificate": {
							Description: "(Write-Only) Base64 encoded Certificate body for server authentication (when appropriate for dialect)",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
						},
						"file_type": {
							Description: "(Write-Only) Certificate keyfile type - .json or .p12",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "json",
						},
						"jdbc_additional_params": {
							Description: "Additional params to add to JDBC connection string",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
						},
						"after_connect_statements": {
							Description: "SQL statements (semicolon separated) to issue after connecting to the database. Requires `custom_after_connect_statements` license feature",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
						},
					},
				},
			},
			"user_attribute_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: nil,
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
	if val, ok := d.GetOk("user_attribute_fields"); ok {
		nc.UserAttributeFields = interfaceListToStringList(val.([]interface{}))
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
	if val, ok := d.GetOk("pdt_context_override"); ok {
		obj := val.(*schema.Set).List()[0].(map[string]interface{})
		context := obj["context"].(string)
		username := obj["username"].(string)
		host := obj["host"].(string)
		port := obj["port"].(string)
		password := obj["password"].(string)
		certificate := obj["certificate"].(string)
		file_type := obj["file_type"].(string)
		jdbc_additional_params := obj["jdbc_additional_params"].(string)
		after_connect_statements := obj["after_connect_statements"].(string)
		pdt_connection := lookergo.DBConnectionOverride{
			Context: &context, Username: &username, Password: &password, Certificate: &certificate, FileType: &file_type, JdbcAdditionalParams: &jdbc_additional_params,
			AfterConnectStatements: &after_connect_statements, Host: &host, Port: &port,
		}
		nc.PdtContextOverride = &pdt_connection
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: test conf", currFuncName()))
	validateConfig, _, err := c.Connections.ValidateConfig(ctx, nc)
	if err != nil {
		return diagErrAppend(diags, err)
	}

	for i, dbcv := range validateConfig {
		tflog.Debug(ctx, fmt.Sprintf("Config [%v] validation message: %v", i, dbcv.Message))
		tflog.Debug(ctx, fmt.Sprintf("Config [%v] validation status: %v", i, dbcv.Status))
	}

	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: create", currFuncName()))
	connection, _, err := c.Connections.Create(ctx, nc)
	if err != nil {
		return diagErrAppend(diags, err)
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
	if d.Set("name", connection.Name) != nil {
		return diag.FromErr(err)
	}
	if d.Set("host", connection.Host) != nil {
		return diag.FromErr(err)
	}
	if d.Set("port", connection.Port) != nil {
		return diag.FromErr(err)
	}
	if d.Set("file_type", connection.FileType) != nil {
		return diag.FromErr(err)
	}
	if d.Set("database", connection.Database) != nil {
		return diag.FromErr(err)
	}
	if d.Set("db_timezone", connection.DbTimezone) != nil {
		return diag.FromErr(err)
	}
	if d.Set("query_timezone", connection.QueryTimezone) != nil {
		return diag.FromErr(err)
	}
	if d.Set("schema", connection.Schema) != nil {
		return diag.FromErr(err)
	}
	if d.Set("max_billing_gigabytes", connection.MaxBillingGigabytes) != nil {
		return diag.FromErr(err)
	}
	if d.Set("tmp_db_name", connection.TmpDbName) != nil {
		return diag.FromErr(err)
	}
	if d.Set("jdbc_additional_params", connection.JdbcAdditionalParams) != nil {
		return diag.FromErr(err)
	}
	if d.Set("dialect_name", connection.DialectName) != nil {
		return diag.FromErr(err)
	}
	if d.Set("maintenance_cron", connection.MaintenanceCron) != nil {
		return diag.FromErr(err)
	}
	if d.Set("after_connect_statements", connection.AfterConnectStatements) != nil {
		return diag.FromErr(err)
	}
	if d.Set("tunnel_id", connection.TunnelId) != nil {
		return diag.FromErr(err)
	}
	if d.Set("oauth_application_id", connection.OauthApplicationId) != nil {
		return diag.FromErr(err)
	}
	if d.Set("ssl", connection.Ssl) != nil {
		return diag.FromErr(err)
	}
	if d.Set("verify_ssl", connection.VerifySsl) != nil {
		return diag.FromErr(err)
	}
	if d.Set("user_db_credentials", connection.UserDbCredentials) != nil {
		return diag.FromErr(err)
	}
	if d.Set("sql_runner_precache_tables", connection.SqlRunnerPrecacheTables) != nil {
		return diag.FromErr(err)
	}
	if d.Set("sql_writing_with_info_schema", connection.SqlWritingWithInfoSchema) != nil {
		return diag.FromErr(err)
	}
	if d.Set("disable_context_comment", connection.DisableContextComment) != nil {
		return diag.FromErr(err)
	}
	if d.Set("always_retry_failed_builds", connection.AlwaysRetryFailedBuilds) != nil {
		return diag.FromErr(err)
	}
	if d.Set("cost_estimate_enabled", connection.CostEstimateEnabled) != nil {
		return diag.FromErr(err)
	}
	if d.Set("pdt_api_control_enabled", connection.PdtApiControlEnabled) != nil {
		return diag.FromErr(err)
	}
	if d.Set("max_connections", connection.MaxConnections) != nil {
		return diag.FromErr(err)
	}
	if d.Set("pool_timeout", connection.PoolTimeout) != nil {
		return diag.FromErr(err)
	}
	if d.Set("pdt_concurrency", connection.PdtConcurrency) != nil {
		return diag.FromErr(err)
	}
	if connection.PdtContextOverride != nil {
		pdt := connection.PdtContextOverride
		var items []interface{}
		items = append(items, map[string]interface{}{
			"context":                  pdt.Context,
			"username":                 pdt.Username,
			"host":                     pdt.Host,
			"port":                     pdt.Port,
			"has_password":             pdt.HasPassword,
			"jdbc_additional_params":   pdt.JdbcAdditionalParams,
			"after_connect_statements": pdt.AfterConnectStatements,
		})
		d.Set("pdt_context_override", items)
	}
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: end", currFuncName()))
	return diags
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	c := m.(*Config).Api // .(*lookergo.Client)
	tflog.Trace(ctx, fmt.Sprintf("Fn: %v, Action: start", currFuncName()))

	connectionName := d.Get("name").(string)
	connection, _, err := c.Connections.Get(ctx, connectionName)
	if err != nil {
		return diag.FromErr(err)
	}

	if val, ok := d.GetOk("host"); ok {
		connection.Host = val.(string)
	}
	if val, ok := d.GetOk("port"); ok {
		connection.Port = val.(string)
	}
	if val, ok := d.GetOk("username"); ok {
		connection.Username = val.(string)
	}
	if val, ok := d.GetOk("password"); ok {
		connection.Password = val.(string)
	}
	if val, ok := d.GetOk("certificate"); ok {
		connection.Certificate = val.(string)
	}
	if val, ok := d.GetOk("file_type"); ok {
		connection.FileType = val.(string)
	}
	if val, ok := d.GetOk("database"); ok {
		connection.Database = val.(string)
	}
	if val, ok := d.GetOk("db_timezone"); ok {
		connection.DbTimezone = val.(string)
	}
	if val, ok := d.GetOk("query_timezone"); ok {
		connection.QueryTimezone = val.(string)
	}
	if val, ok := d.GetOk("schema"); ok {
		connection.Schema = val.(string)
	}
	if val, ok := d.GetOk("max_billing_gigabytes"); ok {
		connection.MaxBillingGigabytes = val.(string)
	}
	if val, ok := d.GetOk("tmp_db_name"); ok {
		connection.TmpDbName = val.(string)
	}
	if val, ok := d.GetOk("jdbc_additional_params"); ok {
		connection.JdbcAdditionalParams = val.(string)
	}
	if val, ok := d.GetOk("dialect_name"); ok {
		connection.DialectName = val.(string)
	}
	if val, ok := d.GetOk("maintenance_cron"); ok {
		connection.MaintenanceCron = val.(string)
	}
	if val, ok := d.GetOk("after_connect_statements"); ok {
		connection.AfterConnectStatements = val.(string)
	}
	if val, ok := d.GetOk("tunnel_id"); ok {
		connection.TunnelId = val.(string)
	}
	if val, ok := d.GetOk("oauth_application_id"); ok {
		connection.OauthApplicationId = val.(string)
	}
	if val, ok := d.GetOk("ssl"); ok {
		connection.Ssl = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("verify_ssl"); ok {
		connection.VerifySsl = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("user_db_credentials"); ok {
		connection.UserDbCredentials = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("sql_runner_precache_tables"); ok {
		connection.SqlRunnerPrecacheTables = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("sql_writing_with_info_schema"); ok {
		connection.SqlWritingWithInfoSchema = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("disable_context_comment"); ok {
		connection.DisableContextComment = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("always_retry_failed_builds"); ok {
		connection.AlwaysRetryFailedBuilds = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("cost_estimate_enabled"); ok {
		connection.CostEstimateEnabled = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("pdt_api_control_enabled"); ok {
		connection.PdtApiControlEnabled = boolPtr(val.(bool))
	}
	if val, ok := d.GetOk("max_connections"); ok {
		connection.MaxConnections = int64(val.(int))
	}
	if val, ok := d.GetOk("pool_timeout"); ok {
		connection.PoolTimeout = int64(val.(int))
	}
	if val, ok := d.GetOk("pdt_concurrency"); ok {
		connection.PdtConcurrency = int64(val.(int))
	}
	if val, ok := d.GetOk("pdt_context_override"); ok {
		obj := val.(*schema.Set).List()[0].(map[string]interface{})
		context := obj["context"].(string)
		username := obj["username"].(string)
		host := obj["host"].(string)
		port := obj["port"].(string)
		password := obj["password"].(string)
		certificate := obj["certificate"].(string)
		file_type := obj["file_type"].(string)
		jdbc_additional_params := obj["jdbc_additional_params"].(string)
		after_connect_statements := obj["after_connect_statements"].(string)
		pdt_connection := lookergo.DBConnectionOverride{
			Context: &context, Username: &username, Password: &password, Certificate: &certificate, FileType: &file_type, JdbcAdditionalParams: &jdbc_additional_params,
			AfterConnectStatements: &after_connect_statements, Host: &host, Port: &port,
		}
		connection.PdtContextOverride = &pdt_connection
	}
	_, _, err = c.Connections.Update(ctx, connectionName, connection)
	if err != nil {
		return diag.FromErr(err)
	}
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

var dialect_names_tab = `<table>
<thead>
<tr>
<th>name</th>
<th>dialect_name</th>
</tr>
</thead>
<tbody>
<tr>
<td>Actian Avalanche</td>
<td>actian_avalanche</td>
</tr>
<tr>
<td>Amazon Athena</td>
<td>athena</td>
</tr>
<tr>
<td>Amazon Aurora MySQL</td>
<td>amazonaurora</td>
</tr>
<tr>
<td>Amazon Redshift</td>
<td>redshift</td>
</tr>
<tr>
<td>Apache Druid</td>
<td>druid</td>
</tr>
<tr>
<td>Apache Druid 0.13+</td>
<td>druid_13</td>
</tr>
<tr>
<td>Apache Druid 0.18+</td>
<td>druid_18</td>
</tr>
<tr>
<td>Apache Hive 2</td>
<td>hive2</td>
</tr>
<tr>
<td>Apache Hive 2.3+</td>
<td>hive2_3</td>
</tr>
<tr>
<td>Apache Hive 3.1.2+</td>
<td>hive3</td>
</tr>
<tr>
<td>Apache Spark 1.5+</td>
<td>spark1_5</td>
</tr>
<tr>
<td>Apache Spark 2.0</td>
<td>spark2_0</td>
</tr>
<tr>
<td>Apache Spark 3+</td>
<td>spark_3</td>
</tr>
<tr>
<td>ClickHouse</td>
<td>clickhouse</td>
</tr>
<tr>
<td>Cloudera Impala</td>
<td>impala</td>
</tr>
<tr>
<td>Cloudera Impala 3.1+</td>
<td>impala_3_1</td>
</tr>
<tr>
<td>Cloudera Impala 3.1+ with Native Driver</td>
<td>impala_native_3_1</td>
</tr>
<tr>
<td>Cloudera Impala with Native Driver</td>
<td>impala_native</td>
</tr>
<tr>
<td>Clustrix</td>
<td>clustrix</td>
</tr>
<tr>
<td>DataVirtuality</td>
<td>datavirtuality</td>
</tr>
<tr>
<td>Databricks</td>
<td>databricks</td>
</tr>
<tr>
<td>Denodo</td>
<td>denodo</td>
</tr>
<tr>
<td>Denodo 7</td>
<td>denodo7</td>
</tr>
<tr>
<td>Denodo 8</td>
<td>denodo8</td>
</tr>
<tr>
<td>Dremio</td>
<td>dremio</td>
</tr>
<tr>
<td>Dremio 11+</td>
<td>dremio_11</td>
</tr>
<tr>
<td>Exasol</td>
<td>exasol</td>
</tr>
<tr>
<td>Firebolt</td>
<td>firebolt</td>
</tr>
<tr>
<td>Google BigQuery Legacy SQL</td>
<td>bigquery</td>
</tr>
<tr>
<td>Google BigQuery Standard SQL</td>
<td>bigquery_standard_sql</td>
</tr>
<tr>
<td>Google Cloud PostgreSQL</td>
<td>google_cloud_postgres</td>
</tr>
<tr>
<td>Google Cloud SQL</td>
<td>googlecloudsql</td>
</tr>
<tr>
<td>Google Cloud Spanner</td>
<td>spanner</td>
</tr>
<tr>
<td>Greenplum</td>
<td>greenplum</td>
</tr>
<tr>
<td>IBM DB2</td>
<td>db2</td>
</tr>
<tr>
<td>IBM DB2 for AS400 and System i</td>
<td>as400</td>
</tr>
<tr>
<td>IBM Netezza</td>
<td>netezza</td>
</tr>
<tr>
<td>MariaDB</td>
<td>mariadb</td>
</tr>
<tr>
<td>Microsoft Azure PostgreSQL</td>
<td>azure_postgres</td>
</tr>
<tr>
<td>Microsoft Azure SQL Database</td>
<td>msazuresql</td>
</tr>
<tr>
<td>Microsoft Azure Synapse Analytics</td>
<td>mssqldw</td>
</tr>
<tr>
<td>Microsoft SQL Server 2008+</td>
<td>mssql_2008</td>
</tr>
<tr>
<td>Microsoft SQL Server 2012+</td>
<td>mssql_2012</td>
</tr>
<tr>
<td>Microsoft SQL Server 2016</td>
<td>mssql_2016</td>
</tr>
<tr>
<td>Microsoft SQL Server 2017+</td>
<td>mssql_2017</td>
</tr>
<tr>
<td>MongoBI</td>
<td>mongobi</td>
</tr>
<tr>
<td>MySQL</td>
<td>mysql</td>
</tr>
<tr>
<td>MySQL 8.0.12+</td>
<td>mysql_8</td>
</tr>
<tr>
<td>Oracle</td>
<td>oracle</td>
</tr>
<tr>
<td>Oracle ADWC</td>
<td>oracle_dwcs</td>
</tr>
<tr>
<td>PostgreSQL 9.5+</td>
<td>postgres</td>
</tr>
<tr>
<td>PostgreSQL pre-9.5</td>
<td>postgres9_1</td>
</tr>
<tr>
<td>PrestoDB</td>
<td>presto</td>
</tr>
<tr>
<td>PrestoSQL</td>
<td>prestosql</td>
</tr>
<tr>
<td>Qubole Presto</td>
<td>qubole_presto</td>
</tr>
<tr>
<td>SAP HANA</td>
<td>hana</td>
</tr>
<tr>
<td>SAP HANA 2+</td>
<td>hana_2</td>
</tr>
<tr>
<td>SingleStore</td>
<td>memsql</td>
</tr>
<tr>
<td>SingleStore 7+</td>
<td>memsql_7</td>
</tr>
<tr>
<td>Snowflake</td>
<td>snowflake</td>
</tr>
<tr>
<td>Teradata</td>
<td>teradata</td>
</tr>
<tr>
<td>Trino</td>
<td>trino</td>
</tr>
<tr>
<td>Vector</td>
<td>vector</td>
</tr>
<tr>
<td>Vertica</td>
<td>vertica</td>
</tr>
</tbody>
</table>
`
