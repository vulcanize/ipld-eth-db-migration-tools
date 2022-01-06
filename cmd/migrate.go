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

package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/les/vflux/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	migration_tools "github.com/vulcanize/migration-tools/pkg"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *logrus.WithField("SubCommand", subCommand)
		migrate()
	},
}

func migrate() {
	conf := migration_tools.NewConfig()
	logWithCommand.Infof("initializing a new Migrator with config params: %+v", conf)
	migrator, err := migration_tools.NewMigrator(context.Background(), conf)
	if err != nil {
		logWithCommand.Fatalf("failed to initialize a new Migrator: %v", err)
	}
	wg := new(sync.WaitGroup)
	tables := getTablesNames()
	rangeChan := make(chan [2]uint64)
	readGapsChan, writeGapsChan, _,  errChan := migrator.Migrate(wg, tables, rangeChan)
	quitChan := make(chan struct{})
	go func() {
		sendRanges(wg, rangeChan, quitChan)
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case readGap := <- readGapsChan:
				logWithCommand.Infof("Migrator read gap: %v", readGap)
			case writeGap := <- writeGapsChan:
				logWithCommand.Infof("Migrator write gap: %v", writeGap)
			case err := <- errChan:
				logWithCommand.Errorf("Migrator error: %v", err)
			case <- quitChan:
				return
			}
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	migrator.Close()
	wg.Wait()
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// groupcache flags
	migrateCmd.PersistentFlags().Bool("find-gaps", false, "turn on the gap finding")
	migrateCmd.PersistentFlags().Uint64("start-height", 0, "start height")
	migrateCmd.PersistentFlags().Uint64("stop-height", 0, "stop height")
	migrateCmd.PersistentFlags().StringArray("table-names", nil, "list of table names to migrate")
	migrateCmd.PersistentFlags().Int("workers-per-table", 1, "number of workers per table")

	migrateCmd.PersistentFlags().Bool("v2-db-name", false, "turn on the gap finding")
	migrateCmd.PersistentFlags().Uint64("v2-db-", 0, "start height")
	migrateCmd.PersistentFlags().Uint64("stop-height", 0, "stop height")
	migrateCmd.PersistentFlags().StringArray("table-names", nil, "list of table names to migrate")
	migrateCmd.PersistentFlags().Int("workers-per-table", 1, "number of workers per table")

	migrateCmd.PersistentFlags().Bool("find-gaps", false, "turn on the gap finding")
	migrateCmd.PersistentFlags().Uint64("start-height", 0, "start height")
	migrateCmd.PersistentFlags().Uint64("stop-height", 0, "stop height")
	migrateCmd.PersistentFlags().StringArray("table-names", nil, "list of table names to migrate")
	migrateCmd.PersistentFlags().Int("workers-per-table", 1, "number of workers per table")

	// state validator flags
	migrateCmd.PersistentFlags().Bool("validator-enabled", false, "turn on the state validator")
	migrateCmd.PersistentFlags().Uint("validator-every-nth-block", 1500, "only validate every Nth block")

	// groupcache flags
	viper.BindPFlag("groupcache.pool.enabled", migrateCmd.PersistentFlags().Lookup("gcache-pool-enabled"))
	viper.BindPFlag("groupcache.pool.httpEndpoint", migrateCmd.PersistentFlags().Lookup("gcache-pool-http-path"))
	viper.BindPFlag("groupcache.pool.peerHttpEndpoints", migrateCmd.PersistentFlags().Lookup("gcache-pool-http-peers"))
	viper.BindPFlag("groupcache.statedb.cacheSizeInMB", migrateCmd.PersistentFlags().Lookup("gcache-statedb-cache-size"))
	viper.BindPFlag("groupcache.statedb.cacheExpiryInMins", migrateCmd.PersistentFlags().Lookup("gcache-statedb-cache-expiry"))

	// state validator flags
	viper.BindPFlag("validator.enabled", migrateCmd.PersistentFlags().Lookup("validator-enabled"))
	viper.BindPFlag("validator.everyNthBlock", migrateCmd.PersistentFlags().Lookup("validator-every-nth-block"))
}

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
