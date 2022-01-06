// VulcanizeDB
// Copyright © 2022 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package migration_tools

// ENV bindings
const (
	LOGRUS_FILE  = "LOGRUS_FILE"
	LOGRUS_LEVEL = "LOGRUS_LEVEL"

	MIGRATION_START             = "MIGRATION_START"
	MIGRATION_STOP              = "MIGRATION_STOP"
	MIGRATION_TABLE_NAMES       = "MIGRATION_TABLE_NAMES"
	MIGRATION_WORKERS_PER_TABLE = "MIGRATION_WORKERS_PER_TABLE"

	V2_DATABASE_NAME                 = "V2_DATABASE_NAME"
	V2_DATABASE_HOSTNAME             = "V2_DATABASE_HOSTNAME"
	V2_DATABASE_PORT                 = "V2_DATABASE_PORT"
	V2_DATABASE_USER                 = "V2_DATABASE_USER"
	V2_DATABASE_PASSWORD             = "V2_DATABASE_PASSWORD"
	V2_DATABASE_MAX_IDLE_CONNECTIONS = "V2_DATABASE_MAX_IDLE_CONNECTIONS"
	V2_DATABASE_MAX_OPEN_CONNECTIONS = "V2_DATABASE_MAX_OPEN_CONNECTIONS"
	V2_DATABASE_MAX_CONN_LIFETIME    = "V2_DATABASE_MAX_CONN_LIFETIME"

	V3_DATABASE_NAME                 = "V3_DATABASE_NAME"
	V3_DATABASE_HOSTNAME             = "V3_DATABASE_HOSTNAME"
	V3_DATABASE_PORT                 = "V3_DATABASE_PORT"
	V3_DATABASE_USER                 = "V3_DATABASE_USER"
	V3_DATABASE_PASSWORD             = "V3_DATABASE_PASSWORD"
	V3_DATABASE_MAX_IDLE_CONNECTIONS = "V3_DATABASE_MAX_IDLE_CONNECTIONS"
	V3_DATABASE_MAX_OPEN_CONNECTIONS = "V3_DATABASE_MAX_OPEN_CONNECTIONS"
	V3_DATABASE_MAX_CONN_LIFETIME    = "V3_DATABASE_MAX_CONN_LIFETIME"
)

// TOML mappings
const (
	TOML_LOGRUS_FILE  = "log.file"
	TOML_LOGRUS_LEVEL = "log.level"

	TOML_MIGRATION_RANGES            = "migrator.ranges"
	TOML_MIGRATION_START             = "migrator.start"
	TOML_MIGRATION_STOP              = "migrator.stop"
	TOML_MIGRATION_TABLE_NAMES       = "migrator.tableNames"
	TOML_MIGRATION_WORKERS_PER_TABLE = "migrator_workerPerTable"

	TOML_V2_DATABASE_NAME                 = "v2.databaseName"
	TOML_V2_DATABASE_HOSTNAME             = "v2.databaseHostName"
	TOML_V2_DATABASE_PORT                 = "v2.databasePort"
	TOML_V2_DATABASE_USER                 = "v2.databaseUser"
	TOML_V2_DATABASE_PASSWORD             = "v2.databasePassword"
	TOML_V2_DATABASE_MAX_IDLE_CONNECTIONS = "v2.databaseMaxIdleConns"
	TOML_V2_DATABASE_MAX_OPEN_CONNECTIONS = "v2.databaseMaxOpenConns"
	TOML_V2_DATABASE_MAX_CONN_LIFETIME    = "v2.databaseMaxConnLifetime"

	TOML_V3_DATABASE_NAME                 = "v3.databaseName"
	TOML_V3_DATABASE_HOSTNAME             = "v3.databaseHostName"
	TOML_V3_DATABASE_PORT                 = "v3.databasePort"
	TOML_V3_DATABASE_USER                 = "v3.databaseUser"
	TOML_V3_DATABASE_PASSWORD             = "v3.databasePassword"
	TOML_V3_DATABASE_MAX_IDLE_CONNECTIONS = "v3.databaseMaxIdleConns"
	TOML_V3_DATABASE_MAX_OPEN_CONNECTIONS = "v3.databaseMaxOpenConns"
	TOML_V3_DATABASE_MAX_CONN_LIFETIME    = "v3.databaseMaxConnLifetime"
)

// CLI flags
const (
	CLI_LOGRUS_FILE  = "log-file"
	CLI_LOGRUS_LEVEL = "log-level"

	CLI_MIGRATION_START             = "start-height"
	CLI_MIGRATION_STOP              = "stop-height"
	CLI_MIGRATION_TABLE_NAMES       = "table-names"
	CLI_MIGRATION_WORKERS_PER_TABLE = "workers-per-table"

	CLI_V2_DATABASE_NAME                 = "v2-db-name"
	CLI_V2_DATABASE_HOSTNAME             = "v2-db-hostname"
	CLI_V2_DATABASE_PORT                 = "v2-db-port"
	CLI_V2_DATABASE_USER                 = "v2-db-username"
	CLI_V2_DATABASE_PASSWORD             = "v2-db-password"
	CLI_V2_DATABASE_MAX_IDLE_CONNECTIONS = "v2-db-max-idle"
	CLI_V2_DATABASE_MAX_OPEN_CONNECTIONS = "v2-db-max-open"
	CLI_V2_DATABASE_MAX_CONN_LIFETIME    = "v2-db-max-lifetime"

	CLI_V3_DATABASE_NAME                 = "v3-db-name"
	CLI_V3_DATABASE_HOSTNAME             = "v3-db-hostname"
	CLI_V3_DATABASE_PORT                 = "v3-db-port"
	CLI_V3_DATABASE_USER                 = "v3-db-username"
	CLI_V3_DATABASE_PASSWORD             = "v3-db-password"
	CLI_V3_DATABASE_MAX_IDLE_CONNECTIONS = "v3-db-max-idle"
	CLI_V3_DATABASE_MAX_OPEN_CONNECTIONS = "v3-db-max-open"
	CLI_V3_DATABASE_MAX_CONN_LIFETIME    = "v3-db-max-lifetime"
)
