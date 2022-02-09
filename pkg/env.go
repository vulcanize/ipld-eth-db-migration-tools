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

	LOG_READ_GAPS_DIR  = "LOG_READ_GAPS_DIR"
	LOG_WRITE_GAPS_DIR = "LOG_WRITE_GAPS_DIR"

	MIGRATION_START                   = "MIGRATION_START"
	MIGRATION_STOP                    = "MIGRATION_STOP"
	MIGRATION_TABLE_NAMES             = "MIGRATION_TABLE_NAMES"
	MIGRATION_WORKERS_PER_TABLE       = "MIGRATION_WORKERS_PER_TABLE"
	MIGRATION_AUTO_RANGE              = "MIGRATION_AUTO_RANGE"
	MIGRATION_AUTO_RANGE_SEGMENT_SIZE = "MIGRATION_AUTO_RANGE_SEGMENT_SIZE"

	TRANSFER_TABLE_NAME   = "TRANSFER_TABLE_NAME"
	TRANSFER_SEGMENT_SIZE = "TRANSFER_SEGMENT_SIZE"
	LOG_TRANSFER_GAPS_DIR = "LOG_TRANSFER_GAPS_DIR"

	OLD_DATABASE_NAME                 = "OLD_DATABASE_NAME"
	OLD_DATABASE_HOSTNAME             = "OLD_DATABASE_HOSTNAME"
	OLD_DATABASE_PORT                 = "OLD_DATABASE_PORT"
	OLD_DATABASE_USER                 = "OLD_DATABASE_USER"
	OLD_DATABASE_PASSWORD             = "OLD_DATABASE_PASSWORD"
	OLD_DATABASE_MAX_IDLE_CONNECTIONS = "OLD_DATABASE_MAX_IDLE_CONNECTIONS"
	OLD_DATABASE_MAX_OPEN_CONNECTIONS = "OLD_DATABASE_MAX_OPEN_CONNECTIONS"
	OLD_DATABASE_MAX_CONN_LIFETIME    = "OLD_DATABASE_MAX_CONN_LIFETIME"

	NEW_DATABASE_NAME                 = "NEW_DATABASE_NAME"
	NEW_DATABASE_HOSTNAME             = "NEW_DATABASE_HOSTNAME"
	NEW_DATABASE_PORT                 = "NEW_DATABASE_PORT"
	NEW_DATABASE_USER                 = "NEW_DATABASE_USER"
	NEW_DATABASE_PASSWORD             = "NEW_DATABASE_PASSWORD"
	NEW_DATABASE_MAX_IDLE_CONNECTIONS = "NEW_DATABASE_MAX_IDLE_CONNECTIONS"
	NEW_DATABASE_MAX_OPEN_CONNECTIONS = "NEW_DATABASE_MAX_OPEN_CONNECTIONS"
	NEW_DATABASE_MAX_CONN_LIFETIME    = "NEW_DATABASE_MAX_CONN_LIFETIME"
)

// TOML mappings
const (
	TOML_LOGRUS_FILE  = "log.file"
	TOML_LOGRUS_LEVEL = "log.level"

	TOML_LOG_READ_GAPS_DIR  = "log.readGapsDir"
	TOML_LOG_WRITE_GAPS_DIR = "log.writeGapsDir"

	TOML_MIGRATION_RANGES                  = "migrator.ranges"
	TOML_MIGRATION_START                   = "migrator.start"
	TOML_MIGRATION_STOP                    = "migrator.stop"
	TOML_MIGRATION_TABLE_NAMES             = "migrator.migrationTableNames"
	TOML_MIGRATION_WORKERS_PER_TABLE       = "migrator.workersPerTable"
	TOML_MIGRATION_AUTO_RANGE              = "migrator.autoRange"
	TOML_MIGRATION_AUTO_RANGE_SEGMENT_SIZE = "migrator.segmentSize"

	TOML_TRANSFER_TABLE_NAME   = "migrator.transferTableName"
	TOML_TRANSFER_SEGMENT_SIZE = "migrator.pagesPerTx"
	TOML_LOG_TRANSFER_GAPS_DIR = "migrator.transferGapDir"

	TOML_OLD_DATABASE_NAME                 = "v2.databaseName"
	TOML_OLD_DATABASE_HOSTNAME             = "v2.databaseHostName"
	TOML_OLD_DATABASE_PORT                 = "v2.databasePort"
	TOML_OLD_DATABASE_USER                 = "v2.databaseUser"
	TOML_OLD_DATABASE_PASSWORD             = "v2.databasePassword"
	TOML_OLD_DATABASE_MAX_IDLE_CONNECTIONS = "v2.databaseMaxIdleConns"
	TOML_OLD_DATABASE_MAX_OPEN_CONNECTIONS = "v2.databaseMaxOpenConns"
	TOML_OLD_DATABASE_MAX_CONN_LIFETIME    = "v2.databaseMaxConnLifetime"

	TOML_NEW_DATABASE_NAME                 = "v3.databaseName"
	TOML_NEW_DATABASE_HOSTNAME             = "v3.databaseHostName"
	TOML_NEW_DATABASE_PORT                 = "v3.databasePort"
	TOML_NEW_DATABASE_USER                 = "v3.databaseUser"
	TOML_NEW_DATABASE_PASSWORD             = "v3.databasePassword"
	TOML_NEW_DATABASE_MAX_IDLE_CONNECTIONS = "v3.databaseMaxIdleConns"
	TOML_NEW_DATABASE_MAX_OPEN_CONNECTIONS = "v3.databaseMaxOpenConns"
	TOML_NEW_DATABASE_MAX_CONN_LIFETIME    = "v3.databaseMaxConnLifetime"
)

// CLI flags
const (
	CLI_LOGRUS_FILE  = "log-file"
	CLI_LOGRUS_LEVEL = "log-level"

	CLI_LOG_READ_GAPS_DIR  = "read-gaps-dir"
	CLI_LOG_WRITE_GAPS_DIR = "write-gaps-dir"

	CLI_MIGRATION_START                   = "start-height"
	CLI_MIGRATION_STOP                    = "stop-height"
	CLI_MIGRATION_TABLE_NAMES             = "migration-table-names"
	CLI_MIGRATION_WORKERS_PER_TABLE       = "workers-per-table"
	CLI_MIGRATION_AUTO_RANGE              = "auto-range"
	CLI_MIGRATION_AUTO_RANGE_SEGMENT_SIZE = "migration-segment-size"

	CLI_TRANSFER_TABLE_NAME   = "transfer-table-name"
	CLI_TRANSFER_SEGMENT_SIZE = "transfer-segment-size"
	CLI_LOG_TRANSFER_GAPS_DIR = "transfer-gap-dir"

	CLI_OLD_DATABASE_NAME                 = "v2-db-name"
	CLI_OLD_DATABASE_HOSTNAME             = "v2-db-hostname"
	CLI_OLD_DATABASE_PORT                 = "v2-db-port"
	CLI_OLD_DATABASE_USER                 = "v2-db-username"
	CLI_OLD_DATABASE_PASSWORD             = "v2-db-password"
	CLI_OLD_DATABASE_MAX_IDLE_CONNECTIONS = "v2-db-max-idle"
	CLI_OLD_DATABASE_MAX_OPEN_CONNECTIONS = "v2-db-max-open"
	CLI_OLD_DATABASE_MAX_CONN_LIFETIME    = "v2-db-max-lifetime"

	CLI_NEW_DATABASE_NAME                 = "v3-db-name"
	CLI_NEW_DATABASE_HOSTNAME             = "v3-db-hostname"
	CLI_NEW_DATABASE_PORT                 = "v3-db-port"
	CLI_NEW_DATABASE_USER                 = "v3-db-username"
	CLI_NEW_DATABASE_PASSWORD             = "v3-db-password"
	CLI_NEW_DATABASE_MAX_IDLE_CONNECTIONS = "v3-db-max-idle"
	CLI_NEW_DATABASE_MAX_OPEN_CONNECTIONS = "v3-db-max-open"
	CLI_NEW_DATABASE_MAX_CONN_LIFETIME    = "v3-db-max-lifetime"
)
