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
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Tool for migrating v2 DB to v3 DB",
	Long: `Tool for reading data from a v2 database, transforming them into v3 DB models,
and writing them to a v3 database.

Can be configured to work over a subset of the tables, and over specific block ranges.

While processing headers, state, or state accounts it checks for gaps in the data and writes these out to a file`,
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

	tables, err := getTableNames()
	if err != nil {
		logWithCommand.Fatalf("failed to generate set of TableNames to process: %v", err)
	}

	ranges, err := getRanges()
	if err != nil {
		logWithCommand.Fatalf("failed to load block ranges to process: %v", err)
	}

	wg := new(sync.WaitGroup)
	quitChan := make(chan struct{})
	for _, table := range tables {
		migrateTable(wg, migrator, table, ranges, quitChan)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	close(quitChan)
	migrator.Close()
	wg.Wait()
}

func migrateTable(wg *sync.WaitGroup, migrator migration_tools.Migrator,
	tableName migration_tools.TableName, blockRanges [][2]uint64, quitChan chan struct{}) {

	rangeChan := make(chan [2]uint64)
	readGapsChan, writeGapsChan, doneChan, errChan := migrator.Migrate(wg, tableName, rangeChan)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i, blockRange := range blockRanges {
			select {
			case <-quitChan:
				logWithCommand.Infof("closing sendRanges subprocess\r\nunsent ranges: %+v", blockRanges[i:])
				return
			default:
				rangeChan <- blockRange
			}
		}
		logWithCommand.Infof("sendRanges subprocess has finished sending all of its ranges for table %s", tableName)
		close(doneChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case readGap := <-readGapsChan:
				logWithCommand.Infof("Migrator %s table read gap: %v", tableName, readGap)
			case writeGap := <-writeGapsChan:
				logWithCommand.Infof("Migrator %s table write gap: %v", tableName, writeGap)
			case err := <-errChan:
				logWithCommand.Errorf("Migrator %s table error: %v", tableName, err)
			case <-quitChan:
				return
			case <-doneChan:
				return
			}
		}
	}()
}

func getTableNames() ([]migration_tools.TableName, error) {
	tableNameStrs := viper.GetStringSlice(migration_tools.TOML_MIGRATION_STOP)
	tableNames := make([]migration_tools.TableName, 0, len(tableNameStrs))
	for _, tableNameStr := range tableNameStrs {
		tableName, err := migration_tools.NewTableNameFromString(tableNameStr)
		if err != nil {
			logWithCommand.Warnf("unable to convert table name string to TableName: %v", err)
			continue
		}
		tableNames = append(tableNames, tableName)
	}
	if len(tableNames) == 0 {
		return nil, fmt.Errorf("migrator needs to be configured with a set of table names to process")
	}
	return tableNames, nil
}

func getRanges() ([][2]uint64, error) {
	var blockRanges [][2]uint64
	viper.UnmarshalKey(migration_tools.TOML_MIGRATION_RANGES, &blockRanges)
	if viper.IsSet(migration_tools.TOML_MIGRATION_START) && viper.IsSet(migration_tools.TOML_MIGRATION_STOP) {
		hardStart := viper.GetUint64(migration_tools.TOML_MIGRATION_START)
		hardStop := viper.GetUint64(migration_tools.TOML_MIGRATION_STOP)
		blockRanges = append(blockRanges, [2]uint64{hardStart, hardStop})
	}
	if len(blockRanges) == 0 {
		return nil, errors.New("migrator needs to be configured with a set of block ranges to process")
	}
	return blockRanges, nil
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// process flags
	migrateCmd.PersistentFlags().Uint64(migration_tools.CLI_MIGRATION_START, 0, "start height")
	migrateCmd.PersistentFlags().Uint64(migration_tools.CLI_MIGRATION_STOP, 0, "stop height")
	migrateCmd.PersistentFlags().StringArray(migration_tools.CLI_MIGRATION_TABLE_NAMES, nil, "list of table names to migrate")
	migrateCmd.PersistentFlags().Int(migration_tools.CLI_MIGRATION_WORKERS_PER_TABLE, 1, "number of workers per table")

	// v2 db flags
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V2_DATABASE_NAME, "vulcanize_v2", "name for the v2 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V2_DATABASE_HOSTNAME, "localhost", "hostname for the v2 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V2_DATABASE_PORT, "5432", "port for the v2 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V2_DATABASE_USER, "postgres", "username to use with the v2 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V2_DATABASE_PASSWORD, "", "password to use for the v2 database")
	migrateCmd.PersistentFlags().Int(migration_tools.CLI_V2_DATABASE_MAX_IDLE_CONNECTIONS, 0, "max idle connections for the v2 database")
	migrateCmd.PersistentFlags().Int(migration_tools.CLI_V2_DATABASE_MAX_OPEN_CONNECTIONS, 0, "max open connections for the v2 database")
	migrateCmd.PersistentFlags().Duration(migration_tools.CLI_V2_DATABASE_MAX_CONN_LIFETIME, 0, "max connection lifetime for the v2 database")

	// v3 db flags
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V3_DATABASE_NAME, "vulcanize_v3", "name for the v3 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V3_DATABASE_HOSTNAME, "localhost", "hostname for the v3 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V3_DATABASE_PORT, "5432", "port for the v3 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V3_DATABASE_USER, "postgres", "username to use with the v3 database")
	migrateCmd.PersistentFlags().String(migration_tools.CLI_V3_DATABASE_PASSWORD, "", "password to use for the v3 database")
	migrateCmd.PersistentFlags().Int(migration_tools.CLI_V3_DATABASE_MAX_IDLE_CONNECTIONS, 0, "max idle connections for the v3 database")
	migrateCmd.PersistentFlags().Int(migration_tools.CLI_V3_DATABASE_MAX_OPEN_CONNECTIONS, 0, "max open connections for the v3 database")
	migrateCmd.PersistentFlags().Duration(migration_tools.CLI_V3_DATABASE_MAX_CONN_LIFETIME, 0, "max connection lifetime for the v3 database")

	// process TOML bindings
	viper.BindPFlag(migration_tools.TOML_MIGRATION_START, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_START))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_STOP, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_STOP))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_TABLE_NAMES, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_TABLE_NAMES))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_WORKERS_PER_TABLE, migrateCmd.PersistentFlags().Lookup(migration_tools.TOML_MIGRATION_WORKERS_PER_TABLE))

	// v2 db TOML bindings
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_NAME, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V2_DATABASE_NAME))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_HOSTNAME, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V2_DATABASE_HOSTNAME))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_PORT, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_PORT))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_USER, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_USER))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_PASSWORD, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V2_DATABASE_PASSWORD))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_MAX_IDLE_CONNECTIONS, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V2_DATABASE_MAX_IDLE_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_MAX_OPEN_CONNECTIONS, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V2_DATABASE_MAX_OPEN_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_V2_DATABASE_MAX_CONN_LIFETIME, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V2_DATABASE_MAX_CONN_LIFETIME))

	// v3 db TOML bindings
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_NAME, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_NAME))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_HOSTNAME, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_HOSTNAME))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_PORT, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_PORT))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_USER, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_USER))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_PASSWORD, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_PASSWORD))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_MAX_IDLE_CONNECTIONS, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_MAX_IDLE_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_MAX_OPEN_CONNECTIONS, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_MAX_OPEN_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_V3_DATABASE_MAX_CONN_LIFETIME, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_V3_DATABASE_MAX_CONN_LIFETIME))
}
