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
	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql/postgres"
	"github.com/spf13/viper"
)

// Config struct holds the configuration params for a Migrator
type Config struct {
	ReadDB          postgres.Config
	WriteDB         postgres.Config
	WorkersPerTable int
}

// NewConfig returns a new Config
func NewConfig() *Config {
	viper.BindEnv(TOML_MIGRATION_WORKERS_PER_TABLE, MIGRATION_WORKERS_PER_TABLE)

	viper.BindEnv(TOML_OLD_DATABASE_NAME, OLD_DATABASE_NAME)
	viper.BindEnv(TOML_OLD_DATABASE_PASSWORD, OLD_DATABASE_PASSWORD)
	viper.BindEnv(TOML_OLD_DATABASE_PORT, OLD_DATABASE_PORT)
	viper.BindEnv(TOML_OLD_DATABASE_USER, OLD_DATABASE_USER)
	viper.BindEnv(TOML_OLD_DATABASE_HOSTNAME, OLD_DATABASE_HOSTNAME)
	viper.BindEnv(TOML_OLD_DATABASE_MAX_CONN_LIFETIME, OLD_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_OLD_DATABASE_MAX_OPEN_CONNECTIONS, OLD_DATABASE_MAX_OPEN_CONNECTIONS)
	viper.BindEnv(TOML_OLD_DATABASE_MAX_CONN_LIFETIME, OLD_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_OLD_DATABASE_MAX_IDLE_CONNECTIONS, OLD_DATABASE_MAX_IDLE_CONNECTIONS)

	viper.BindEnv(TOML_NEW_DATABASE_NAME, NEW_DATABASE_NAME)
	viper.BindEnv(TOML_NEW_DATABASE_PASSWORD, NEW_DATABASE_PASSWORD)
	viper.BindEnv(TOML_NEW_DATABASE_PORT, NEW_DATABASE_PORT)
	viper.BindEnv(TOML_NEW_DATABASE_USER, NEW_DATABASE_USER)
	viper.BindEnv(TOML_NEW_DATABASE_HOSTNAME, NEW_DATABASE_HOSTNAME)
	viper.BindEnv(TOML_NEW_DATABASE_MAX_CONN_LIFETIME, NEW_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_NEW_DATABASE_MAX_OPEN_CONNECTIONS, NEW_DATABASE_MAX_OPEN_CONNECTIONS)
	viper.BindEnv(TOML_NEW_DATABASE_MAX_CONN_LIFETIME, NEW_DATABASE_MAX_CONN_LIFETIME)
	viper.BindEnv(TOML_NEW_DATABASE_MAX_IDLE_CONNECTIONS, NEW_DATABASE_MAX_IDLE_CONNECTIONS)

	return &Config{
		WorkersPerTable: viper.GetInt(TOML_MIGRATION_WORKERS_PER_TABLE),
		ReadDB: postgres.Config{
			Username:        viper.GetString(TOML_OLD_DATABASE_USER),
			Password:        viper.GetString(TOML_OLD_DATABASE_PASSWORD),
			Hostname:        viper.GetString(TOML_OLD_DATABASE_HOSTNAME),
			DatabaseName:    viper.GetString(TOML_OLD_DATABASE_NAME),
			Port:            viper.GetInt(TOML_OLD_DATABASE_PORT),
			MaxConns:        viper.GetInt(TOML_OLD_DATABASE_MAX_OPEN_CONNECTIONS),
			MaxIdle:         viper.GetInt(TOML_OLD_DATABASE_MAX_IDLE_CONNECTIONS),
			MaxConnLifetime: viper.GetDuration(TOML_OLD_DATABASE_MAX_CONN_LIFETIME),
		},
		WriteDB: postgres.Config{
			Username:        viper.GetString(TOML_NEW_DATABASE_USER),
			Password:        viper.GetString(TOML_NEW_DATABASE_PASSWORD),
			Hostname:        viper.GetString(TOML_NEW_DATABASE_HOSTNAME),
			DatabaseName:    viper.GetString(TOML_NEW_DATABASE_NAME),
			Port:            viper.GetInt(TOML_NEW_DATABASE_PORT),
			MaxConns:        viper.GetInt(TOML_NEW_DATABASE_MAX_OPEN_CONNECTIONS),
			MaxIdle:         viper.GetInt(TOML_NEW_DATABASE_MAX_IDLE_CONNECTIONS),
			MaxConnLifetime: viper.GetDuration(TOML_NEW_DATABASE_MAX_CONN_LIFETIME),
		},
	}
}
