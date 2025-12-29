package common

const (
	// Programming languages

	Bash       = "bash"
	C          = "c"
	CPP        = "cpp"
	CSharp     = "csharp"
	Go         = "go"
	Java       = "java"
	JavaScript = "javascript"
	JSON       = "json"
	Kotlin     = "kotlin"
	PHP        = "php"
	Python     = "python"
	Ruby       = "ruby"
	Rust       = "rust"
	Scala      = "scala"
	Shell      = "shell"
	Swift      = "swift"
	Text       = "text"
	TypeScript = "typescript"
	Undefined  = "undefined"
	XML        = "xml"
	YAML       = "yaml"

	// File extensions

	BashExtension       = ".sh"
	CExtension          = ".c"
	CPPExtension        = ".cpp"
	CSharpExtension     = ".cs"
	GoExtension         = ".go"
	JavaExtension       = ".java"
	JavaScriptExtension = ".js"
	JSONExtension       = ".json"
	KotlinExtension     = ".kt"
	PHPExtension        = ".php"
	PythonExtension     = ".py"
	RubyExtension       = ".rb"
	RustExtension       = ".rs"
	ScalaExtension      = ".scala"
	ShellExtension      = ".sh"
	SwiftExtension      = ".swift"
	TextExtension       = ".txt"
	TypeScriptExtension = ".ts"
	UndefinedExtension  = ".txt"
	XMLExtension        = ".xml"
	YAMLExtension       = ".yaml"

	// Code example categories

	SyntaxExample              = "Syntax example"
	NonMongoCommand            = "Non-MongoDB command"
	ExampleReturnObject        = "Example return object"
	ExampleConfigurationObject = "Example configuration object"
	UsageExample               = "Usage example"

	/*
		The constants below for Products, SubProducts, and Directories are used in
		the `productInfoMap` in the GetProductInfo.go file. The GetProductInfo
		func relies on these constants to return the appropriate `product` and/or
		`sub_product` name for their respective fields on a per-docs-page basis
		in the code examples DB.
	*/
	// Products

	Atlas                        = "Atlas"
	AtlasArchitecture            = "Atlas Architecture Center"
	BIConnector                  = "BI Connector"
	CloudManager                 = "Cloud Manager"
	Compass                      = "Compass"
	DBTools                      = "Database Tools"
	Django                       = "Django Integration"
	Drivers                      = "Drivers"
	EnterpriseKubernetesOperator = "Enterprise Kubernetes Operator"
	EFCoreProvider               = "Entity Framework Core Provider"
	KafkaConnector               = "Kafka Connector"
	MCPServer                    = "MongoDB MCP Server"
	MDBCLI                       = "MongoDB CLI"
	Mongosh                      = "MongoDB Shell"
	Mongosync                    = "Mongosync"
	OpsManager                   = "Ops Manager"
	RelationalMigrator           = "Relational Migrator"
	Server                       = "Server"
	SparkConnector               = "Spark Connector"

	// SubProducts

	AtlasCLI         = "Atlas CLI"
	AtlasOperator    = "Kubernetes Operator"
	Charts           = "Charts"
	DataFederation   = "Data Federation"
	OnlineArchive    = "Online Archive"
	Search           = "Search"
	StreamProcessing = "Stream Processing"
	Terraform        = "Terraform"
	Triggers         = "Triggers"
	VectorSearch     = "Vector Search"

	// Directories that map to specific sub-products

	DataFederationDir   = "data-federation"
	OnlineArchiveDir    = "online-archive"
	SearchDir           = "atlas-search"
	StreamProcessingDir = "atlas-stream-processing"
	TriggersDir         = "triggers"
	VectorSearchDir     = "atlas-vector-search"
	AiIntegrationsDir   = "ai-integrations"
)

var CanonicalLanguages = []string{Bash, C, CPP,
	CSharp, Go, Java, JavaScript,
	JSON, Kotlin, PHP, Python,
	Ruby, Rust, Scala, Shell,
	Swift, Text, TypeScript, Undefined, XML, YAML,
}

var WesUrls = []string{
	"https://mongodb.com/docs/manual/administration/install-community",
	"https://mongodb.com/docs/search",
	"https://mongodb.com/docs/get-started",
	"https://mongodb.com/docs/manual/installation",
	"https://mongodb.com/docs",
	"https://mongodb.com/docs/manual",
	"https://mongodb.com/docs/compass/install",
	"https://mongodb.com/docs/mongodb-shell/install",
	"https://mongodb.com/docs/manual/reference/connection-string",
	"https://mongodb.com/docs/development",
	"https://mongodb.com/docs/database-tools/installation",
	"https://mongodb.com/docs/atlas/troubleshoot-connection",
	"https://mongodb.com/docs/drivers",
	"https://mongodb.com/docs/compass/query/filter",
	"https://mongodb.com/docs/atlas/security-add-mongodb-users",
	"https://mongodb.com/docs/manual/indexes",
	"https://mongodb.com/docs/mongodb-shell",
	"https://mongodb.com/docs/manual/crud",
	"https://mongodb.com/docs/drivers/node/current",
	"https://mongodb.com/docs/tools-and-connectors",
	"https://mongodb.com/docs/manual/tutorial/query-documents",
	"https://mongodb.com/docs/manual/core/document",
	"https://mongodb.com/docs/atlas",
	"https://mongodb.com/docs/manual/tutorial/install-mongodb-enterprise-on-windows",
	"https://mongodb.com/docs/manual/core/databases-and-collections",
	"https://mongodb.com/docs/atlas/security/ip-access-list",
	"https://mongodb.com/docs/atlas/atlas-vector-search/tutorials/vector-search-quick-start",
	"https://mongodb.com/docs/management",
	"https://mongodb.com/docs/mcp-server/get-started",
	"https://mongodb.com/docs/database-tools/mongodump",
	"https://mongodb.com/docs/compass",
	"https://mongodb.com/docs/404",
	"https://mongodb.com/docs/manual/aggregation",
	"https://mongodb.com/docs/manual/core/aggregation-pipeline",
	"https://mongodb.com/docs/drivers/compatibility",
	"https://mongodb.com/docs/manual/reference/method/db.collection.find",
	"https://mongodb.com/docs/compass/connect",
	"https://mongodb.com/docs/manual/reference/method/db.collection.updateMany",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/group",
	"https://mongodb.com/docs/manual/reference/limits",
	"https://mongodb.com/docs/atlas/atlas-search/tutorial",
	"https://mongodb.com/docs/v7.0/tutorial/install-mongodb-on-ubuntu",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/lookup",
	"https://mongodb.com/docs/atlas/atlas-vector-search/vector-search-overview",
	"https://mongodb.com/docs/manual/data-modeling",
	"https://mongodb.com/docs/compass/query/queries",
	"https://mongodb.com/docs/manual/tutorial/insert-documents",
	"https://mongodb.com/docs/manual/reference/configuration-options",
	"https://mongodb.com/docs/languages/javascript",
	"https://mongodb.com/docs/manual/changeStreams",
	"https://mongodb.com/docs/atlas/getting-started",
	"https://mongodb.com/docs/database-tools/mongorestore",
	"https://mongodb.com/docs/compass/query-with-natural-language/enable-natural-language-querying",
	"https://mongodb.com/docs/atlas/atlas-vector-search/create-embeddings",
	"https://mongodb.com/docs/manual/core/transactions",
	"https://mongodb.com/docs/atlas/atlas-search",
	"https://mongodb.com/docs/manual/reference/operator/query/regex",
	"https://mongodb.com/docs/mongodb-shell/connect",
	"https://mongodb.com/docs/mongodb-vscode/connect",
	"https://mongodb.com/docs/manual/core/index-ttl",
	"https://mongodb.com/docs/mongodb-shell/run-commands",
	"https://mongodb.com/docs/drivers/node/current/get-started",
	"https://mongodb.com/docs/languages/python/pymongo-driver/current",
	"https://mongodb.com/docs/atlas/compass-connection",
	"https://mongodb.com/docs/manual/sharding",
	"https://mongodb.com/docs/manual/core/indexes/index-types/index-compound",
	"https://mongodb.com/docs/manual/replication",
	"https://mongodb.com/docs/manual/reference/glossary",
	"https://mongodb.com/docs/database-tools",
	"https://mongodb.com/docs/atlas/atlas-search/manage-indexes",
	"https://mongodb.com/docs/manual/reference/bson-types",
	"https://mongodb.com/docs/manual/core/index-unique",
	"https://mongodb.com/docs/manual/tutorial/query-arrays",
	"https://mongodb.com/docs/languages/python",
	"https://mongodb.com/docs/database-tools/mongoexport",
	"https://mongodb.com/docs/drivers/node/current/connect/mongoclient",
	"https://mongodb.com/docs/manual/reference/built-in-roles",
	"https://mongodb.com/docs/atlas/cli/current/install-atlas-cli",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/match",
	"https://mongodb.com/docs/manual/reference/method/db.collection.updateOne",
	"https://mongodb.com/docs/manual/reference/connection-string-options",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/unwind",
	"https://mongodb.com/docs/manual/reference/database-users",
	"https://mongodb.com/docs/database-tools/mongoimport",
	"https://mongodb.com/docs/manual/release-notes",
	"https://mongodb.com/docs/drivers/node/current/connect",
	"https://mongodb.com/docs/manual/reference/operator/query/in",
	"https://mongodb.com/docs/manual/core/views",
	"https://mongodb.com/docs/languages/java",
	"https://mongodb.com/docs/legacy",
	"https://mongodb.com/docs/atlas/atlas-vector-search/vector-search-type",
	"https://mongodb.com/docs/languages/python/pymongo-driver/current/get-started",
	"https://mongodb.com/docs/manual/release-notes/8.0",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/project",
	"https://mongodb.com/docs/manual/core/timeseries-collections",
	"https://mongodb.com/docs/compass/import-export",
	"https://mongodb.com/docs/manual/administration/production-notes",
	"https://mongodb.com/docs/manual/reference/mql/expressions",
	"https://mongodb.com/docs/drivers/node/current/crud/insert",
	"https://mongodb.com/docs/drivers/csharp/current",
	"https://mongodb.com/docs/drivers/node/current/integrations/mongoose-get-started",
	"https://mongodb.com/docs/manual/tutorial/update-documents",
	"https://mongodb.com/docs/atlas/security-add-mongodb-roles",
	"https://mongodb.com/docs/manual/release-notes/8.2",
	"https://mongodb.com/docs/manual/reference/mql/aggregation-stages",
	"https://mongodb.com/docs/manual/reference/versioning",
	"https://mongodb.com/docs/manual/tutorial/backup-and-restore-tools",
	"https://mongodb.com/docs/manual/tutorial/install-mongodb-enterprise-on-ubuntu",
	"https://mongodb.com/docs/atlas/atlas-vector-search/vector-search-stage",
	"https://mongodb.com/docs/manual/reference/method/db.collection.update",
	"https://mongodb.com/docs/kubernetes/current",
	"https://mongodb.com/docs/manual/tutorial/project-fields-from-query-results",
	"https://mongodb.com/docs/atlas/cli/current",
	"https://mongodb.com/docs/manual/core/indexes/create-index",
	"https://mongodb.com/docs/manual/tutorial/query-embedded-documents",
	"https://mongodb.com/docs/atlas/monitor-cluster-metrics",
	"https://mongodb.com/docs/manual/reference/operator/query/exists",
	"https://mongodb.com/docs/atlas/architecture/current",
	"https://mongodb.com/docs/atlas/billing/atlas-flex-costs",
	"https://mongodb.com/docs/atlas/sample-data",
	"https://mongodb.com/docs/manual/reference/write-concern",
	"https://mongodb.com/docs/manual/reference/mql/query-predicates",
	"https://mongodb.com/docs/drivers/node/current/crud/query/retrieve",
	"https://mongodb.com/docs/atlas/cluster-autoscaling",
	"https://mongodb.com/docs/manual/reference/method/ObjectId",
	"https://mongodb.com/docs/v8.0/tutorial/install-mongodb-on-ubuntu",
	"https://mongodb.com/docs/manual/core/read-preference",
	"https://mongodb.com/docs/manual/tutorial/query-for-null-fields",
	"https://mongodb.com/docs/manual/tutorial/install-mongodb-enterprise-on-os-x",
	"https://mongodb.com/docs/manual/core/schema-validation",
	"https://mongodb.com/docs/manual/reference/command",
	"https://mongodb.com/docs/atlas/atlas-search/searching",
	"https://mongodb.com/docs/manual/release-notes/7.0",
	"https://mongodb.com/docs/manual/reference/method/Date",
	"https://mongodb.com/docs/manual/reference/operator/query/or",
	"https://mongodb.com/docs/manual/reference/program/mongod",
	"https://mongodb.com/docs/drivers/java/sync/current/get-started",
	"https://mongodb.com/docs/atlas/data-federation/query/connect-with-sql-composable",
	"https://mongodb.com/docs/atlas/reference/user-roles",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/sort",
	"https://mongodb.com/docs/manual/core/index-partial",
	"https://mongodb.com/docs/manual/core/gridfs",
	"https://mongodb.com/docs/manual/reference/operator/query/and",
	"https://mongodb.com/docs/manual/reference/operator/query/elemMatch",
	"https://mongodb.com/docs/atlas/security-private-endpoint",
	"https://mongodb.com/docs/manual/core/wiredtiger",
	"https://mongodb.com/docs/manual/reference/insert-methods",
	"https://mongodb.com/docs/manual/reference/parameters",
	"https://mongodb.com/docs/compass/upgrade",
	"https://mongodb.com/docs/manual/data-modeling/schema-design-process",
	"https://mongodb.com/docs/manual/administration/connection-pool-overview",
	"https://mongodb.com/docs/manual/tutorial/install-mongodb-enterprise-on-windows-unattended",
	"https://mongodb.com/docs/manual/reference/cheatsheet",
	"https://mongodb.com/docs/manual/reference/method/db.collection.createIndex",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/merge",
	"https://mongodb.com/docs/atlas/create-connect-deployments",
	"https://mongodb.com/docs/manual/tutorial/getting-started",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/addFields",
	"https://mongodb.com/docs/manual/reference/method/db.collection.distinct",
	"https://mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate",
	"https://mongodb.com/docs/manual/reference/operator/update/push",
	"https://mongodb.com/docs/drivers/java/sync/current",
	"https://mongodb.com/docs/mongodb-vscode/playgrounds",
	"https://mongodb.com/docs/atlas/atlas-search/tutorial/partial-match",
	"https://mongodb.com/docs/languages/php",
	"https://mongodb.com/docs/languages/python/pymongo-driver/current/connect",
	"https://mongodb.com/docs/manual/tutorial/remove-documents",
	"https://mongodb.com/docs/manual/tutorial/configure-ssl",
	"https://mongodb.com/docs/manual/core/indexes/index-types",
	"https://mongodb.com/docs/atlas/billing/cluster-configuration-costs",
	"https://mongodb.com/docs/atlas/setup-cluster-security",
	"https://mongodb.com/docs/atlas/configure-api-access",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/count",
	"https://mongodb.com/docs/manual/reference/sql-comparison",
	"https://mongodb.com/docs/atlas/atlas-vector-search/rag",
	"https://mongodb.com/docs/manual/tutorial/expire-data",
	"https://mongodb.com/docs/mongodb-shell/crud",
	"https://mongodb.com/docs/mongodb-shell/crud/read",
	"https://mongodb.com/docs/drivers/go/current",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/filter",
	"https://mongodb.com/docs/manual/core/views/create-view",
	"https://mongodb.com/docs/ops-manager/current",
	"https://mongodb.com/docs/manual/reference/method/db.collection.insertMany",
	"https://mongodb.com/docs/atlas/tutorial/deploy-free-tier-cluster",
	"https://mongodb.com/docs/atlas/import/c2c-pull-live-migration",
	"https://mongodb.com/docs/mongodb-shell/write-scripts",
	"https://mongodb.com/docs/manual/reference/operator/query/ne",
	"https://mongodb.com/docs/atlas/mongodb-users-roles-and-privileges",
	"https://mongodb.com/docs/mongodb-shell/reference/data-types",
	"https://mongodb.com/docs/manual/reference",
	"https://mongodb.com/docs/manual/reference/method/db.collection.findOne",
	"https://mongodb.com/docs/manual/reference/operator/update/set",
	"https://mongodb.com/docs/manual/geospatial-queries",
	"https://mongodb.com/docs/manual/reference/method/db.collection.insertOne",
	"https://mongodb.com/docs/manual/reference/mql/update",
	"https://mongodb.com/docs/manual/reference/log-messages",
	"https://mongodb.com/docs/manual/core/replica-set-oplog",
	"https://mongodb.com/docs/manual/tutorial/deploy-replica-set",
	"https://mongodb.com/docs/manual/core/indexes/index-types/index-text",
	"https://mongodb.com/docs/manual/tutorial/query-array-of-documents",
	"https://mongodb.com/docs/manual/reference/method/db.createCollection",
	"https://mongodb.com/docs/v7.0/administration/install-community",
	"https://mongodb.com/docs/manual/reference/operator/update/unset",
	"https://mongodb.com/docs/manual/core/security-scram",
	"https://mongodb.com/docs/manual/reference/method/db.collection.count",
	"https://mongodb.com/docs/mongocli/current",
	"https://mongodb.com/docs/atlas/pause-terminate-cluster",
	"https://mongodb.com/docs/atlas/atlas-vector-search/ai-agents",
	"https://mongodb.com/docs/atlas/security-vpc-peering",
	"https://mongodb.com/docs/atlas/backup-restore-cluster",
	"https://mongodb.com/docs/manual/core/index-sparse",
	"https://mongodb.com/docs/manual/reference/method/db.collection.bulkWrite",
	"https://mongodb.com/docs/atlas/billing",
	"https://mongodb.com/docs/languages/python/pymongo-driver/current/integrations/fastapi-integration",
	"https://mongodb.com/docs/manual/reference/mql",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/facet",
	"https://mongodb.com/docs/manual/reference/method",
	"https://mongodb.com/docs/manual/reference/operator/aggregation/cond",
	"https://mongodb.com/docs/atlas/atlas-search/operators-and-collectors",
	"https://mongodb.com/docs/manual/administration/install-enterprise",
	"https://mongodb.com/docs/manual/tutorial/equality-sort-range-guideline",
	"https://mongodb.com/docs/manual/core/views/join-collections-with-view",
	"https://mongodb.com/docs/atlas/atlas-search/index-definitions",
	"https://mongodb.com/docs/manual/reference/installation-ubuntu-community-troubleshooting",
	"https://mongodb.com/docs/manual/reference/method/db.collection.deleteMany",
	"https://mongodb.com/docs/atlas/tutorial/create-new-cluster",
	"https://mongodb.com/docs/manual/reference/method/db.collection.aggregate",
	"https://mongodb.com/docs/manual/tutorial/aggregation-complete-examples",
	"https://mongodb.com/docs/manual/reference/operator/update/pull",
	"https://mongodb.com/docs/atlas/access/manage-org-users",
	"https://mongodb.com/docs/manual/reference/method/cursor.sort",
	"https://mongodb.com/docs/atlas/import",
	"https://mongodb.com/docs/drivers/php/laravel-mongodb/current",
	"https://mongodb.com/docs/manual/tutorial/create-users",
	"https://mongodb.com/docs/mongosync/current",
	"https://mongodb.com/docs/languages/python/pymongo-driver/current/reference/migration",
	"https://mongodb.com/docs/manual/core/backups",
	"https://mongodb.com/docs/manual/reference/method/db.createUser",
	"https://mongodb.com/docs/drivers/node/current/databases-collections",
	"https://mongodb.com/docs/manual/core/aggregation-pipeline-optimization",
	"https://mongodb.com/docs/atlas/ai-integrations",
	"https://mongodb.com/docs/v7.0/installation",
	"https://mongodb.com/docs/manual/reference/method/db.collection.updatemany",
	"https://mongodb.com/docs/drivers/go/current/get-started",
	"https://mongodb.com/docs/languages/python/pymongo-driver/current/databases-collections",
	"https://mongodb.com/docs/manual/reference/operator/query/gt",
	"https://mongodb.com/docs/manual/reference/operator/query/size",
	"https://mongodb.com/docs/manual/tutorial/update-documents-with-aggregation-pipeline",
	"https://mongodb.com/docs/manual/reference/privilege-actions",
	"https://mongodb.com/docs/v8.0/administration/install-community",
}
