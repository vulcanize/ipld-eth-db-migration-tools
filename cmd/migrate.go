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
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Tool for migrating old DB to new DB",
	Long: `Tool for reading data from a old database, transforming them into new DB models,
and writing them to a new database.

Can be configured to work over a subset of the tables, and over specific block ranges.

While processing headers, state, or state accounts it checks for gaps in the data and writes these out to a file`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *logrus.WithField("SubCommand", subCommand)
		migrate()
	},
}

func migrate() {
	logWithCommand.Info("----- running migration -----")
	conf := migration_tools.NewConfig()
	logWithCommand.Infof("initializing a new Migrator with config params: %+v", conf)
	migrator, err := migration_tools.NewMigrator(context.Background(), conf)
	if err != nil {
		logWithCommand.Fatalf("failed to initialize a new Migrator: %v", err)
	}

	tables, err := getTableNames()
	if err != nil {
		logWithCommand.Fatalf("failed to generate set of TableNames for processing: %v", err)
	}

	if err := getGapDirs(); err != nil {
		logWithCommand.Fatalf("failed to open directories for writing read and write gaps: %v", err)
	}

	ranges, err := getRanges(conf.ReadDB)
	if err != nil {
		logWithCommand.Fatalf("failed to load block ranges for processing: %v", err)
	}

	wg := new(sync.WaitGroup)
	go func() {
		for _, table := range tables {
			migrateTable(wg, migrator, table, ranges)
		}
	}()

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		<-shutdown
		migrator.Close()
	}()
	wg.Wait()
}

var (
	readGapsDir  string
	writeGapsDir string
)

func getGapDirs() error {
	viper.BindEnv(migration_tools.TOML_LOG_READ_GAPS_DIR, migration_tools.LOG_READ_GAPS_DIR)
	viper.BindEnv(migration_tools.TOML_LOG_WRITE_GAPS_DIR, migration_tools.LOG_WRITE_GAPS_DIR)
	readGapsDir = viper.GetString(migration_tools.TOML_LOG_READ_GAPS_DIR)
	writeGapsDir = viper.GetString(migration_tools.TOML_LOG_WRITE_GAPS_DIR)
	if _, err := os.Stat(readGapsDir); os.IsNotExist(err) {
		if err := os.Mkdir(readGapsDir, 0777); err != nil {
			return err
		}
	}
	if _, err := os.Stat(writeGapsDir); os.IsNotExist(err) {
		if err := os.Mkdir(writeGapsDir, 0777); err != nil {
			return err
		}
	}
	return nil
}

func migrateTable(wg *sync.WaitGroup, migrator migration_tools.Migrator,
	tableName migration_tools.TableName, blockRanges [][2]uint64) {

	now := time.Now().Unix()
	readGapFilePath := filepath.Join(readGapsDir, string(tableName)+"_"+strconv.Itoa(int(now)))
	writeGapFilePath := filepath.Join(writeGapsDir, string(tableName)+"_"+strconv.Itoa(int(now)))
	readGapFile, err := os.OpenFile(readGapFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		migrator.Close()
		logWithCommand.Fatalf("unable to open read gap file at %s", readGapFilePath)
	}
	writeGapFile, err := os.OpenFile(writeGapFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		migrator.Close()
		logWithCommand.Fatalf("unable to open write gap file at %s", writeGapFilePath)
	}

	rangeChan := make(chan [2]uint64)
	readGapsChan, writeGapsChan, doneChan, quitChan, errChan := migrator.Migrate(wg, tableName, rangeChan)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i, blockRange := range blockRanges {
			select {
			case <-doneChan:
				logWithCommand.Infof("closing sendRanges subprocess\r\nunsent ranges: %+v", blockRanges[i:])
				return
			default:
				rangeChan <- blockRange
			}
			if tableName == migration_tools.PublicNodes {
				// public nodes will be migrated in one batch, since it is not segmented by block height
				break
			}
		}
		logWithCommand.Infof("finished sending block ranges for table %s\r\nshutting down migration process for table %s", tableName, tableName)
		close(quitChan)
	}()

	// handle writing out gaps
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer readGapFile.Close()
		defer writeGapFile.Close()
		for {
			select {
			case readGap := <-readGapsChan:
				logWithCommand.Infof("Migrator %s table read gap: %v", tableName, readGap)
				if _, err := readGapFile.WriteString(fmt.Sprintf("%d, %d\r\n", readGap[0], readGap[1])); err != nil {
					logWithCommand.Errorf("error writing read gap to file at %s; err: %s", readGapFilePath, err.Error())
				}
			case writeGap := <-writeGapsChan:
				logWithCommand.Infof("Migrator %s table write gap: %v", tableName, writeGap)
				if _, err := writeGapFile.WriteString(fmt.Sprintf("%d, %d\r\n", writeGap[0], writeGap[1])); err != nil {
					logWithCommand.Errorf("error writing write gap to file at %s; err: %s", writeGapFilePath, err.Error())
				}
			case <-doneChan:
				return
			}
		}
	}()

	// handle writing out errors
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-errChan:
				logWithCommand.Errorf("Migrator %s table migration error: %v", tableName, err)
			case <-doneChan:
				return
			}
		}
	}()
}

func getTableNames() ([]migration_tools.TableName, error) {
	viper.BindEnv(migration_tools.TOML_MIGRATION_TABLE_NAMES, migration_tools.MIGRATION_TABLE_NAMES)
	tableNameStrs := viper.GetStringSlice(migration_tools.TOML_MIGRATION_TABLE_NAMES)
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

func getRanges(readConf postgres.Config) ([][2]uint64, error) {
	viper.BindEnv(migration_tools.TOML_MIGRATION_AUTO_RANGE, migration_tools.MIGRATION_AUTO_RANGE)
	viper.BindEnv(migration_tools.TOML_MIGRATION_AUTO_RANGE_SEGMENT_SIZE, migration_tools.MIGRATION_AUTO_RANGE_SEGMENT_SIZE)
	if viper.GetBool(migration_tools.TOML_MIGRATION_AUTO_RANGE) && viper.IsSet(migration_tools.TOML_MIGRATION_AUTO_RANGE_SEGMENT_SIZE) {
		segmentSize := viper.GetUint64(migration_tools.TOML_MIGRATION_AUTO_RANGE_SEGMENT_SIZE)
		if segmentSize == 0 {
			return nil, errors.New("auto range detection and segmenting is on, but segment size is set to 0")
		}
		logWithCommand.Infof("auto range detection and segmenting is on, with segment size of %d", segmentSize)
		return migration_tools.DetectAndSegmentRangeByChunkSize(readConf, segmentSize)
	}
	viper.BindEnv(migration_tools.TOML_MIGRATION_START, migration_tools.MIGRATION_START)
	viper.BindEnv(migration_tools.TOML_MIGRATION_STOP, migration_tools.MIGRATION_STOP)
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

	// migrator flags
	migrateCmd.PersistentFlags().Uint64(migration_tools.CLI_MIGRATION_START, 0, "start height")
	migrateCmd.PersistentFlags().Uint64(migration_tools.CLI_MIGRATION_STOP, 0, "stop height")
	migrateCmd.PersistentFlags().StringArray(migration_tools.CLI_MIGRATION_TABLE_NAMES, nil, "list of table names to migrate")
	migrateCmd.PersistentFlags().Int(migration_tools.CLI_MIGRATION_WORKERS_PER_TABLE, 1, "number of workers per table")
	migrateCmd.PersistentFlags().Bool(migration_tools.CLI_MIGRATION_AUTO_RANGE, false, "turn on or off auto range detection and chunking")
	migrateCmd.PersistentFlags().Uint64(migration_tools.CLI_MIGRATION_AUTO_RANGE_SEGMENT_SIZE, 0, "segment size for auto range detection and chunking")

	// migrator TOML bindings
	viper.BindPFlag(migration_tools.TOML_MIGRATION_START, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_START))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_STOP, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_STOP))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_TABLE_NAMES, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_TABLE_NAMES))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_WORKERS_PER_TABLE, migrateCmd.PersistentFlags().Lookup(migration_tools.TOML_MIGRATION_WORKERS_PER_TABLE))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_AUTO_RANGE, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_AUTO_RANGE))
	viper.BindPFlag(migration_tools.TOML_MIGRATION_AUTO_RANGE_SEGMENT_SIZE, migrateCmd.PersistentFlags().Lookup(migration_tools.CLI_MIGRATION_AUTO_RANGE_SEGMENT_SIZE))
}
