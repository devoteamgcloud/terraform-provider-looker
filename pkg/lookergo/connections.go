package lookergo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const connectionsBasePath = "4.0/connections"

type ConnectionsResource interface {
	Get(ctx context.Context, connectionName string) (*DBConnection, *Response, error)
	Create(ctx context.Context, connection *DBConnection) (*DBConnection, *Response, error)
	Update(ctx context.Context, connectionName string, connection *DBConnection) (*DBConnection, *Response, error)
	Delete(ctx context.Context, connectionName string) (*Response, error)
	ValidateConfig(ctx context.Context, connection *DBConnection) ([]DBConnectionValidation, *Response, error)
	ValidateConnection(ctx context.Context, connectionName string, tests []string) ([]DBConnectionValidation, *Response, error)
}

type ConnectionsResourceOp struct {
	client *Client
}

var _ ConnectionsResource = &ConnectionsResourceOp{}

// <editor-fold desc="type_structs">
// DBConnection struct for DBConnection
type DBConnection struct {
	// Name of the connection. Also used as the unique identifier
	Name    string     `json:"name,omitempty"`
	Dialect *DBDialect `json:"dialect,omitempty"`
	// SQL Runner snippets for this connection
	Snippets []Snippet `json:"snippets,omitempty"`
	// True if PDTs are enabled on this connection
	PdtsEnabled *bool `json:"pdts_enabled,omitempty"`
	// Host name/address of server
	Host string `json:"host,omitempty"`
	// Port number on server
	Port string `json:"port,omitempty"`
	// Username for server authentication
	Username string `json:"username,omitempty"`
	// (Write-Only) Password for server authentication
	Password string `json:"password,omitempty"`
	// Whether the connection uses OAuth for authentication.
	UsesOauth *bool `json:"uses_oauth,omitempty"`
	// (Write-Only) Base64 encoded Certificate body for server authentication (when appropriate for dialect).
	Certificate string `json:"certificate,omitempty"`
	// (Write-Only) Certificate keyfile type - .json or .p12
	FileType string `json:"file_type,omitempty"`
	// Database name
	Database string `json:"database,omitempty"`
	// Time zone of database
	DbTimezone string `json:"db_timezone,omitempty"`
	// Timezone to use in queries
	QueryTimezone string `json:"query_timezone,omitempty"`
	// Scheme name
	Schema string `json:"schema,omitempty"`
	// Maximum number of concurrent connection to use
	MaxConnections int64 `json:"max_connections,omitempty"`
	// Maximum size of query in GBs (BigQuery only, can be a user_attribute name)
	MaxBillingGigabytes string `json:"max_billing_gigabytes,omitempty"`
	// Use SSL/TLS when connecting to server
	Ssl *bool `json:"ssl,omitempty"`
	// Verify the SSL
	VerifySsl *bool `json:"verify_ssl,omitempty"`
	// Name of temporary database (if used)
	TmpDbName string `json:"tmp_db_name,omitempty"`
	// Additional params to add to JDBC connection string
	JdbcAdditionalParams string `json:"jdbc_additional_params,omitempty"`
	// Connection Pool Timeout, in seconds
	PoolTimeout int64 `json:"pool_timeout,omitempty"`
	// (Read/Write) SQL Dialect name
	DialectName string `json:"dialect_name,omitempty"`
	// Creation date for this connection
	CreatedAt string `json:"created_at,omitempty"`
	// Id of user who last modified this connection configuration
	UserId string `json:"user_id,omitempty"`
	// Is this an example connection?
	Example *bool `json:"example,omitempty"`
	// (Limited access feature) Are per user db credentials enabled. Enabling will remove previously set username and password
	UserDbCredentials *bool `json:"user_db_credentials,omitempty"`
	// Fields whose values map to user attribute names
	UserAttributeFields []string `json:"user_attribute_fields,omitempty"`
	// Cron string specifying when maintenance such as PDT trigger checks and drops should be performed
	MaintenanceCron string `json:"maintenance_cron,omitempty"`
	// Unix timestamp at start of last completed PDT trigger check process
	LastRegenAt string `json:"last_regen_at,omitempty"`
	// Unix timestamp at start of last completed PDT reap process
	LastReapAt string `json:"last_reap_at,omitempty"`
	// Precache tables in the SQL Runner
	SqlRunnerPrecacheTables *bool `json:"sql_runner_precache_tables,omitempty"`
	// Fetch Information Schema For SQL Writing
	SqlWritingWithInfoSchema *bool `json:"sql_writing_with_info_schema,omitempty"`
	// SQL statements (semicolon separated) to issue after connecting to the database. Requires `custom_after_connect_statements` license feature
	AfterConnectStatements string `json:"after_connect_statements,omitempty"`
	// Not implemented
	// PdtContextOverride DBConnectionOverride `json:"pdt_context_override,omitempty"`
	// Is this connection created and managed by Looker
	Managed *bool `json:"managed,omitempty"`
	// The Id of the ssh tunnel this connection uses
	TunnelId string `json:"tunnel_id,omitempty"`
	// Maximum number of threads to use to build PDTs in parallel
	PdtConcurrency int64 `json:"pdt_concurrency,omitempty"`
	// When disable_context_comment is true comment will not be added to SQL
	DisableContextComment *bool `json:"disable_context_comment,omitempty"`
	// An External OAuth Application to use for authenticating to the database
	OauthApplicationId string `json:"oauth_application_id,omitempty"`
	// When true, error PDTs will be retried every regenerator cycle
	AlwaysRetryFailedBuilds *bool `json:"always_retry_failed_builds,omitempty"`
	// When true, query cost estimate will be displayed in explore.
	CostEstimateEnabled *bool `json:"cost_estimate_enabled,omitempty"`
	// PDT builds on this connection can be kicked off and cancelled via API.
	PdtApiControlEnabled *bool `json:"pdt_api_control_enabled,omitempty"`
}

type DBDialect struct {
	// The name of the dialect
	Name string `json:"name,omitempty"`
	// The human-readable label of the connection
	Label string `json:"label,omitempty"`
	// Whether the dialect supports query cost estimates
	SupportsCostEstimate *bool `json:"supports_cost_estimate,omitempty"`
	// How the dialect handles cost estimation
	CostEstimateStyle string `json:"cost_estimate_style,omitempty"`
	// PDT index columns
	PersistentTableIndexes string `json:"persistent_table_indexes,omitempty"`
	// PDT sortkey columns
	PersistentTableSortkeys string `json:"persistent_table_sortkeys,omitempty"`
	// PDT distkey column
	PersistentTableDistkey string `json:"persistent_table_distkey,omitempty"`
	// Suports streaming results
	SupportsStreaming *bool `json:"supports_streaming,omitempty"`
	// Should SQL Runner snippets automatically be run
	AutomaticallyRunSqlRunnerSnippets *bool `json:"automatically_run_sql_runner_snippets,omitempty"`
	// Array of names of the tests that can be run on a connection using this dialect
	ConnectionTests []string `json:"connection_tests,omitempty"`
	// Is supported with the inducer (i.e. generate from sql)
	SupportsInducer *bool `json:"supports_inducer,omitempty"`
	// Can multiple databases be accessed from a connection using this dialect
	SupportsMultipleDatabases *bool `json:"supports_multiple_databases,omitempty"`
	// Whether the dialect supports allowing Looker to build persistent derived tables
	SupportsPersistentDerivedTables *bool `json:"supports_persistent_derived_tables,omitempty"`
	// Does the database have client SSL support settable through the JDBC string explicitly?
	HasSslSupport *bool `json:"has_ssl_support,omitempty"`
}

type Snippet struct {
	// Name of the snippet
	Name string `json:"name,omitempty"`
	// Label of the snippet
	Label string `json:"label,omitempty"`
	// SQL text of the snippet
	Sql string `json:"sql,omitempty"`
}

/*
{
  "message": "Validation Failed",
  "errors": [
    {
      "field": "dialect",
      "code": "missing_field",
      "message": "This field is required.",
      "documentation_url": "https://docs.looker.com/"
    }
  ],
  "documentation_url": "https://docs.looker.com/"
}
*/

type DBConnectionValidation struct {
	ConnectionString string `json:"connection_string"`
	Message          string `json:"message"`
	Name             string `json:"name"`
	Status           string `json:"status"`
}

// </editor-fold>

func (s ConnectionsResourceOp) Get(ctx context.Context, connectionName string) (*DBConnection, *Response, error) {
	return doGet(ctx, s.client, connectionsBasePath, new(DBConnection), url.QueryEscape(connectionName))
}

func (s ConnectionsResourceOp) Create(ctx context.Context, connection *DBConnection) (*DBConnection, *Response, error) {
	return doCreate(ctx, s.client, connectionsBasePath, connection, new(DBConnection))
}

func (s ConnectionsResourceOp) Update(ctx context.Context, connectionName string, connection *DBConnection) (*DBConnection, *Response, error) {
	return doUpdate(ctx, s.client, connectionsBasePath, url.QueryEscape(connectionName), connection, new(DBConnection))
}

func (s ConnectionsResourceOp) Delete(ctx context.Context, connectionName string) (*Response, error) {

	return doDelete(ctx, s.client, connectionsBasePath, url.QueryEscape(connectionName))
}

func (s ConnectionsResourceOp) ValidateConfig(ctx context.Context, connection *DBConnection) (dbcv []DBConnectionValidation, resp *Response, err error) {
	path := fmt.Sprintf("%s%s", connectionsBasePath, strings.Join(append([]string{""}, "test"), "/"))

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, connection)
	if err != nil {
		return nil, nil, err
	}

	resp, err = s.client.Do(ctx, req, &dbcv)
	if err != nil {
		return nil, resp, err
	}

	return
}

// ValidateConnection -
// Possible for most db's: first do connect test; next do kill,query test
func (s ConnectionsResourceOp) ValidateConnection(ctx context.Context, connectionName string, tests []string) (dbcv []DBConnectionValidation, resp *Response, err error) {
	if len(tests) == 0 {
		tests = []string{"connect"}
	}
	path := fmt.Sprintf("%s/%s?tests=%v",
		connectionsBasePath,
		connectionName, url.QueryEscape(strings.Join(tests, ",")))

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err = s.client.Do(ctx, req, &dbcv)
	if err != nil {
		return nil, resp, err
	}

	return
}
