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

const (
	MIGRATION_FIND_GAPS         = "MIGRATION_FIND_GAPS"
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

const (
	TOML_MIGRATION_FIND_GAPS         = "migrator.findGaps"
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
