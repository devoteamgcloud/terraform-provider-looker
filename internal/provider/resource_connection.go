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
				Type:        schema.TypeString,
				Required:    true,
				Description: dialect_names_tab,
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
	if err = d.Set("name", connection.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("host", connection.Host); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("user_attribute_fields", connection.UserAttributeFields); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("port", connection.Port); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("username", connection.Username); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("password", connection.Password); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("certificate", connection.Certificate); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("file_type", connection.FileType); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("database", connection.Database); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("db_timezone", connection.DbTimezone); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("query_timezone", connection.QueryTimezone); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("schema", connection.Schema); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("max_billing_gigabytes", connection.MaxBillingGigabytes); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("tmp_db_name", connection.TmpDbName); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("jdbc_additional_params", connection.JdbcAdditionalParams); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("dialect_name", connection.DialectName); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("maintenance_cron", connection.MaintenanceCron); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("after_connect_statements", connection.AfterConnectStatements); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("tunnel_id", connection.TunnelId); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("oauth_application_id", connection.OauthApplicationId); err != nil {
		return diag.FromErr(err)
	}

	if connection.Ssl != nil {
		if err = d.Set("ssl", *connection.Ssl); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.VerifySsl != nil {
		if err = d.Set("verify_ssl", *connection.VerifySsl); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.UserDbCredentials != nil {
		if err = d.Set("user_db_credentials", *connection.UserDbCredentials); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.SqlRunnerPrecacheTables != nil {
		if err = d.Set("sql_runner_precache_tables", *connection.SqlRunnerPrecacheTables); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.SqlWritingWithInfoSchema != nil {
		if err = d.Set("sql_writing_with_info_schema", *connection.SqlWritingWithInfoSchema); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.DisableContextComment != nil {
		if err = d.Set("disable_context_comment", *connection.DisableContextComment); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.AlwaysRetryFailedBuilds != nil {
		if err = d.Set("always_retry_failed_builds", *connection.AlwaysRetryFailedBuilds); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.CostEstimateEnabled != nil {
		if err = d.Set("cost_estimate_enabled", *connection.CostEstimateEnabled); err != nil {
			return diag.FromErr(err)
		}
	}

	if connection.PdtApiControlEnabled != nil {
		if err = d.Set("pdt_api_control_enabled", *connection.PdtApiControlEnabled); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = d.Set("max_connections", connection.MaxConnections); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("pool_timeout", connection.PoolTimeout); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("pdt_concurrency", connection.PdtConcurrency); err != nil {
		return diag.FromErr(err)
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

	if d.HasChange("host") {
		connection.Host = d.Get("host").(string)
	}
	if d.HasChange("user_attribute_fields") {
		connection.UserAttributeFields = interfaceListToStringList(d.Get("user_attribute_fields").([]interface{}))
	}
	if d.HasChange("port") {
		connection.Port = d.Get("port").(string)
	}
	if d.HasChange("username") {
		connection.Username = d.Get("username").(string)
	}
	if d.HasChange("password") {
		connection.Password = d.Get("password").(string)
	}
	if d.HasChange("certificate") {
		connection.Certificate = d.Get("certificate").(string)
	}
	if d.HasChange("file_type") {
		connection.FileType = d.Get("file_type").(string)
	}
	if d.HasChange("database") {
		connection.Database = d.Get("database").(string)
	}
	if d.HasChange("db_timezone") {
		connection.DbTimezone = d.Get("db_timezone").(string)
	}
	if d.HasChange("query_timezone") {
		connection.QueryTimezone = d.Get("query_timezone").(string)
	}
	if d.HasChange("schema") {
		connection.Schema = d.Get("schema").(string)
	}
	if d.HasChange("max_billing_gigabytes") {
		connection.MaxBillingGigabytes = d.Get("max_billing_gigabytes").(string)
	}
	if d.HasChange("tmp_db_name") {
		connection.TmpDbName = d.Get("tmp_db_name").(string)
	}
	if d.HasChange("jdbc_additional_params") {
		connection.JdbcAdditionalParams = d.Get("jdbc_additional_params").(string)
	}
	if d.HasChange("dialect_name") {
		connection.DialectName = d.Get("dialect_name").(string)
	}
	if d.HasChange("maintenance_cron") {
		connection.MaintenanceCron = d.Get("maintenance_cron").(string)
	}
	if d.HasChange("after_connect_statements") {
		connection.AfterConnectStatements = d.Get("after_connect_statements").(string)
	}
	if d.HasChange("tunnel_id") {
		connection.TunnelId = d.Get("tunnel_id").(string)
	}
	if d.HasChange("oauth_application_id") {
		connection.OauthApplicationId = d.Get("oauth_application_id").(string)
	}
	if d.HasChange("ssl") {
		connection.Ssl = boolPtr(d.Get("ssl").(bool))
	}
	if d.HasChange("verify_ssl") {
		connection.VerifySsl = boolPtr(d.Get("verify_ssl").(bool))
	}
	if d.HasChange("user_db_credentials") {
		connection.UserDbCredentials = boolPtr(d.Get("user_db_credentials").(bool))
	}
	if d.HasChange("sql_runner_precache_tables") {
		connection.SqlRunnerPrecacheTables = boolPtr(d.Get("sql_runner_precache_tables").(bool))
	}
	if d.HasChange("sql_writing_with_info_schema") {
		connection.SqlWritingWithInfoSchema = boolPtr(d.Get("sql_writing_with_info_schema").(bool))
	}
	if d.HasChange("disable_context_comment") {
		connection.DisableContextComment = boolPtr(d.Get("disable_context_comment").(bool))
	}
	if d.HasChange("always_retry_failed_builds") {
		connection.AlwaysRetryFailedBuilds = boolPtr(d.Get("always_retry_failed_builds").(bool))
	}
	if d.HasChange("cost_estimate_enabled") {
		connection.CostEstimateEnabled = boolPtr(d.Get("cost_estimate_enabled").(bool))
	}
	if d.HasChange("pdt_api_control_enabled") {
		connection.PdtApiControlEnabled = boolPtr(d.Get("pdt_api_control_enabled").(bool))
	}
	if d.HasChange("max_connections") {
		connection.MaxConnections = int64(d.Get("max_connections").(int))
	}
	if d.HasChange("pool_timeout") {
		connection.PoolTimeout = int64(d.Get("pool_timeout").(int))
	}
	if d.HasChange("pdt_concurrency") {
		connection.PdtConcurrency = int64(d.Get("pdt_concurrency").(int))
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
