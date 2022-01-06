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

import (
	"time"

	"github.com/spf13/viper"
)

// Config struct holds the configuration params for a Migrator
type Config struct {
	ReadDB          DBConfig
	WriteDB         DBConfig
	WorkersPerTable int
}

// DBConfig struct holds Postgres configuration params
type DBConfig struct {
	Username     string
	Password     string
	Hostname     string
	DatabaseName string
	Port         int

	MaxConns        int
	MaxIdle         int
	MaxConnLifetime time.Duration
}

// NewConfig returns a new Config
func NewConfig() *Config {
	viper.BindEnv(TOML_MIGRATION_START, MIGRATION_START)
	viper.BindEnv(TOML_MIGRATION_STOP, MIGRATION_STOP)
	viper.BindEnv(TOML_MIGRATION_TABLE_NAMES, MIGRATION_TABLE_NAMES)
	viper.BindEnv(TOML_MIGRATION_WORKERS_PER_TABLE, MIGRATION_WORKERS_PER_TABLE)

	viper.BindEnv(TOML_V2_DATABASE_NAME, V2_DATABASE_NAME)
	viper.BindEnv(TOML_V2_DATABASE_PASSWORD, V2_DATABASE_PASSWORD)
	viper.BindEnv(TOML_V2_DATABASE_PORT, V2_DATABASE_PORT)
	viper.BindEnv(TOML_V2_DATABASE_USER, V2_DATABASE_USER)
	viper.BindEnv(TOML_V2_DATABASE_HOSTNAME, V2_DATABASE_HOSTNAME)
	viper.BindEnv(TOML_V2_DATABASE_MAX_CONN_LIFETIME, V2_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_V2_DATABASE_MAX_OPEN_CONNECTIONS, V2_DATABASE_MAX_OPEN_CONNECTIONS)
	viper.BindEnv(TOML_V2_DATABASE_MAX_CONN_LIFETIME, V2_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_V2_DATABASE_MAX_IDLE_CONNECTIONS, V2_DATABASE_MAX_IDLE_CONNECTIONS)

	viper.BindEnv(TOML_V3_DATABASE_NAME, V3_DATABASE_NAME)
	viper.BindEnv(TOML_V3_DATABASE_PASSWORD, V3_DATABASE_PASSWORD)
	viper.BindEnv(TOML_V3_DATABASE_PORT, V3_DATABASE_PORT)
	viper.BindEnv(TOML_V3_DATABASE_USER, V3_DATABASE_USER)
	viper.BindEnv(TOML_V3_DATABASE_HOSTNAME, V3_DATABASE_HOSTNAME)
	viper.BindEnv(TOML_V3_DATABASE_MAX_CONN_LIFETIME, V3_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_V3_DATABASE_MAX_OPEN_CONNECTIONS, V3_DATABASE_MAX_OPEN_CONNECTIONS)
	viper.BindEnv(TOML_V3_DATABASE_MAX_CONN_LIFETIME, V3_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_V3_DATABASE_MAX_IDLE_CONNECTIONS, V3_DATABASE_MAX_IDLE_CONNECTIONS)

	return &Config{
		WorkersPerTable: viper.GetInt(TOML_MIGRATION_WORKERS_PER_TABLE),
		ReadDB: DBConfig{
			Username:        viper.GetString(TOML_V2_DATABASE_USER),
			Password:        viper.GetString(TOML_V2_DATABASE_PASSWORD),
			Hostname:        viper.GetString(TOML_V2_DATABASE_HOSTNAME),
			DatabaseName:    viper.GetString(TOML_V2_DATABASE_NAME),
			Port:            viper.GetInt(TOML_V2_DATABASE_PORT),
			MaxConns:        viper.GetInt(TOML_V2_DATABASE_MAX_OPEN_CONNECTIONS),
			MaxIdle:         viper.GetInt(TOML_V2_DATABASE_MAX_IDLE_CONNECTIONS),
			MaxConnLifetime: viper.GetDuration(TOML_V2_DATABASE_MAX_CONN_LIFETIME),
		},
		WriteDB: DBConfig{
			Username:        viper.GetString(TOML_V3_DATABASE_USER),
			Password:        viper.GetString(TOML_V3_DATABASE_PASSWORD),
			Hostname:        viper.GetString(TOML_V3_DATABASE_HOSTNAME),
			DatabaseName:    viper.GetString(TOML_V3_DATABASE_NAME),
			Port:            viper.GetInt(TOML_V3_DATABASE_PORT),
			MaxConns:        viper.GetInt(TOML_V3_DATABASE_MAX_OPEN_CONNECTIONS),
			MaxIdle:         viper.GetInt(TOML_V3_DATABASE_MAX_IDLE_CONNECTIONS),
			MaxConnLifetime: viper.GetDuration(TOML_V3_DATABASE_MAX_CONN_LIFETIME),
		},
	}
}
