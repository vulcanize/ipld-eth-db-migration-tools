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
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
	"github.com/vulcanize/migration-tools/pkg/public_blocks"
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Tool for transfer old DB data to new DB",
	Long: `Tool for reading data from a old database in segments, and writing it to a new database
in a manner that is able to handle potential unique constraints conflict
(e.g. calling out to a postgres_fdw procedure).

Enables batching and proper logging of errors and progress during the process.`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *logrus.WithField("SubCommand", subCommand)
		transfer()
	},
}

func transfer() {
	logWithCommand.Info("----- running transfer -----")
	conf := migration_tools.NewConfig()
	logWithCommand.Infof("initializing a new Transferor with config params: %+v", conf)
	transferor, err := migration_tools.NewMigrator(context.Background(), conf)
	if err != nil {
		logWithCommand.Fatalf("failed to initialize a new Transferor: %v", err)
	}

	if err := getTransferGapDir(); err != nil {
		logWithCommand.Fatalf("failed to open directory for writing transfer gaps: %v", err)
	}

	wg := new(sync.WaitGroup)
	viper.BindEnv(migration_tools.TOML_TRANSFER_TABLE_NAME, migration_tools.TRANSFER_TABLE_NAME)
	viper.BindEnv(migration_tools.TOML_TRANSFER_SEGMENT_SIZE, migration_tools.TRANSFER_SEGMENT_SIZE)
	viper.BindEnv(migration_tools.TOML_TRANSFER_SEGMENT_OFFSET, migration_tools.TRANSFER_SEGMENT_OFFSET)
	viper.BindEnv(migration_tools.TOML_TRANSFER_MAX_PAGE, migration_tools.TRANSFER_MAX_PAGE)
	transferTable(wg, transferor,
		viper.GetString(migration_tools.TOML_TRANSFER_TABLE_NAME),
		viper.GetUint64(migration_tools.TOML_TRANSFER_SEGMENT_SIZE),
		viper.GetUint64(migration_tools.TOML_TRANSFER_SEGMENT_OFFSET),
		viper.GetUint64(migration_tools.TOML_TRANSFER_MAX_PAGE))

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		<-shutdown
		transferor.Close()
	}()
	wg.Wait()
}

var transferGapDir string

func getTransferGapDir() error {
	viper.BindEnv(migration_tools.TOML_LOG_TRANSFER_GAPS_DIR, migration_tools.LOG_TRANSFER_GAPS_DIR)
	transferGapDir = viper.GetString(migration_tools.TOML_LOG_TRANSFER_GAPS_DIR)
	if _, err := os.Stat(transferGapDir); os.IsNotExist(err) {
		if err := os.Mkdir(transferGapDir, 0777); err != nil {
			return err
		}
	}
	return nil
}

func transferTable(wg *sync.WaitGroup, transferor migration_tools.Migrator, tableName string,
	segmentSize, segmentOffset, maxPage uint64) {
	now := time.Now().Unix()
	transferGapFilePath := filepath.Join(transferGapDir, string(tableName)+"_"+strconv.Itoa(int(now)))
	transferGapFile, err := os.OpenFile(transferGapFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		transferor.Close()
		logWithCommand.Fatalf("unable to open transfer gap file at %s", transferGapFilePath)
	}

	gapChan, doneChan, errChan, err := transferor.Transfer(wg, tableName, segmentSize, segmentOffset, maxPage)
	if err != nil {
		logWithCommand.Fatalf("transfer initialization failed: %v", err)
	}

	// handle writing out gaps
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer transferGapFile.Close()
		for {
			select {
			case readGap := <-gapChan:
				logWithCommand.Infof("Migrator %s table read gap: %v", tableName, readGap)
				if _, err := transferGapFile.WriteString(fmt.Sprintf("%d, %d\r\n", readGap[0], readGap[1])); err != nil {
					logWithCommand.Errorf("error writing read gap to file at %s; err: %s", transferGapFilePath, err.Error())
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
				logWithCommand.Errorf("Migrator %s table transfer error: %v", tableName, err)
			case <-doneChan:
				return
			}
		}
	}()
}

func init() {
	rootCmd.AddCommand(transferCmd)

	// transferor flags
	transferCmd.PersistentFlags().Uint64(migration_tools.CLI_TRANSFER_SEGMENT_SIZE, 1000, "number of pages transferred per tx")
	transferCmd.PersistentFlags().String(migration_tools.CLI_TRANSFER_TABLE_NAME, public_blocks.DefaultV2FDWTableName, "postgres_fdw table name in the new database")
	transferCmd.PersistentFlags().Uint64(migration_tools.CLI_TRANSFER_MAX_PAGE, 0, "configure the max page; if left 0 then MAX(ctid) is queried from the DB (which can take a long time)")
	transferCmd.PersistentFlags().Uint64(migration_tools.CLI_TRANSFER_SEGMENT_OFFSET, 0, "starting offset for the number of segments we process (for picking up where a previous process stopped)")

	// transferor TOML bindings
	viper.BindPFlag(migration_tools.TOML_TRANSFER_SEGMENT_SIZE, transferCmd.PersistentFlags().Lookup(migration_tools.CLI_TRANSFER_SEGMENT_SIZE))
	viper.BindPFlag(migration_tools.TOML_TRANSFER_TABLE_NAME, transferCmd.PersistentFlags().Lookup(migration_tools.CLI_TRANSFER_TABLE_NAME))
	viper.BindPFlag(migration_tools.TOML_TRANSFER_MAX_PAGE, transferCmd.PersistentFlags().Lookup(migration_tools.CLI_TRANSFER_MAX_PAGE))
	viper.BindPFlag(migration_tools.TOML_TRANSFER_SEGMENT_OFFSET, transferCmd.PersistentFlags().Lookup(migration_tools.CLI_TRANSFER_SEGMENT_OFFSET))
}
