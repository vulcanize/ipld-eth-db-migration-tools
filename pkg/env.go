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
	LOGRUS_FILE           = "LOGRUS_FILE"
	LOGRUS_LEVEL          = "LOGRUS_LEVEL"
	LOG_READ_GAPS_DIR     = "LOG_READ_GAPS_DIR"
	LOG_WRITE_GAPS_DIR    = "LOG_WRITE_GAPS_DIR"
	LOG_TRANSFER_GAPS_DIR = "LOG_TRANSFER_GAPS_DIR"

	MIGRATION_START                   = "MIGRATION_START"
	MIGRATION_STOP                    = "MIGRATION_STOP"
	MIGRATION_TABLE_NAMES             = "MIGRATION_TABLE_NAMES"
	MIGRATION_WORKERS_PER_TABLE       = "MIGRATION_WORKERS_PER_TABLE"
	MIGRATION_AUTO_RANGE              = "MIGRATION_AUTO_RANGE"
	MIGRATION_AUTO_RANGE_SEGMENT_SIZE = "MIGRATION_AUTO_RANGE_SEGMENT_SIZE"

	TRANSFER_TABLE_NAME   = "TRANSFER_TABLE_NAME"
	TRANSFER_SEGMENT_SIZE = "TRANSFER_SEGMENT_SIZE"

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
	TOML_LOGRUS_FILE           = "log.file"
	TOML_LOGRUS_LEVEL          = "log.level"
	TOML_LOG_READ_GAPS_DIR     = "log.readGapsDir"
	TOML_LOG_WRITE_GAPS_DIR    = "log.writeGapsDir"
	TOML_LOG_TRANSFER_GAPS_DIR = "log.transferGapDir"

	TOML_MIGRATION_RANGES                  = "migrator.ranges"
	TOML_MIGRATION_START                   = "migrator.start"
	TOML_MIGRATION_STOP                    = "migrator.stop"
	TOML_MIGRATION_TABLE_NAMES             = "migrator.migrationTableNames"
	TOML_MIGRATION_WORKERS_PER_TABLE       = "migrator.workersPerTable"
	TOML_MIGRATION_AUTO_RANGE              = "migrator.autoRange"
	TOML_MIGRATION_AUTO_RANGE_SEGMENT_SIZE = "migrator.segmentSize"
	TOML_TRANSFER_TABLE_NAME               = "migrator.transferTableName"
	TOML_TRANSFER_SEGMENT_SIZE             = "migrator.pagesPerTx"

	TOML_OLD_DATABASE_NAME                 = "old.databaseName"
	TOML_OLD_DATABASE_HOSTNAME             = "old.databaseHostName"
	TOML_OLD_DATABASE_PORT                 = "old.databasePort"
	TOML_OLD_DATABASE_USER                 = "old.databaseUser"
	TOML_OLD_DATABASE_PASSWORD             = "old.databasePassword"
	TOML_OLD_DATABASE_MAX_IDLE_CONNECTIONS = "old.databaseMaxIdleConns"
	TOML_OLD_DATABASE_MAX_OPEN_CONNECTIONS = "old.databaseMaxOpenConns"
	TOML_OLD_DATABASE_MAX_CONN_LIFETIME    = "old.databaseMaxConnLifetime"

	TOML_NEW_DATABASE_NAME                 = "new.databaseName"
	TOML_NEW_DATABASE_HOSTNAME             = "new.databaseHostName"
	TOML_NEW_DATABASE_PORT                 = "new.databasePort"
	TOML_NEW_DATABASE_USER                 = "new.databaseUser"
	TOML_NEW_DATABASE_PASSWORD             = "new.databasePassword"
	TOML_NEW_DATABASE_MAX_IDLE_CONNECTIONS = "new.databaseMaxIdleConns"
	TOML_NEW_DATABASE_MAX_OPEN_CONNECTIONS = "new.databaseMaxOpenConns"
	TOML_NEW_DATABASE_MAX_CONN_LIFETIME    = "new.databaseMaxConnLifetime"
)

// CLI flags
const (
	CLI_LOGRUS_FILE           = "log-file"
	CLI_LOGRUS_LEVEL          = "log-level"
	CLI_LOG_READ_GAPS_DIR     = "read-gaps-dir"
	CLI_LOG_WRITE_GAPS_DIR    = "write-gaps-dir"
	CLI_LOG_TRANSFER_GAPS_DIR = "transfer-gap-dir"

	CLI_MIGRATION_START                   = "start-height"
	CLI_MIGRATION_STOP                    = "stop-height"
	CLI_MIGRATION_TABLE_NAMES             = "migration-table-names"
	CLI_MIGRATION_WORKERS_PER_TABLE       = "workers-per-table"
	CLI_MIGRATION_AUTO_RANGE              = "auto-range"
	CLI_MIGRATION_AUTO_RANGE_SEGMENT_SIZE = "migration-segment-size"

	CLI_TRANSFER_TABLE_NAME   = "transfer-table-name"
	CLI_TRANSFER_SEGMENT_SIZE = "transfer-segment-size"

	CLI_OLD_DATABASE_NAME                 = "old-db-name"
	CLI_OLD_DATABASE_HOSTNAME             = "old-db-hostname"
	CLI_OLD_DATABASE_PORT                 = "old-db-port"
	CLI_OLD_DATABASE_USER                 = "old-db-username"
	CLI_OLD_DATABASE_PASSWORD             = "old-db-password"
	CLI_OLD_DATABASE_MAX_IDLE_CONNECTIONS = "old-db-max-idle"
	CLI_OLD_DATABASE_MAX_OPEN_CONNECTIONS = "old-db-max-open"
	CLI_OLD_DATABASE_MAX_CONN_LIFETIME    = "old-db-max-lifetime"

	CLI_NEW_DATABASE_NAME                 = "new-db-name"
	CLI_NEW_DATABASE_HOSTNAME             = "new-db-hostname"
	CLI_NEW_DATABASE_PORT                 = "new-db-port"
	CLI_NEW_DATABASE_USER                 = "new-db-username"
	CLI_NEW_DATABASE_PASSWORD             = "new-db-password"
	CLI_NEW_DATABASE_MAX_IDLE_CONNECTIONS = "new-db-max-idle"
	CLI_NEW_DATABASE_MAX_OPEN_CONNECTIONS = "new-db-max-open"
	CLI_NEW_DATABASE_MAX_CONN_LIFETIME    = "new-db-max-lifetime"
)
