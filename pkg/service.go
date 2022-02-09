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
	"context"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/migration-tools/pkg/csv"
	"github.com/vulcanize/migration-tools/pkg/public_blocks"
	"github.com/vulcanize/migration-tools/pkg/sql"
)

const defaultNumWorkersPerTable = 1

// Migrator interface for migrating from v2 DB to v3 DB
type Migrator interface {
	Migrate(wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan struct{}, chan error)
	Transfer(wg *sync.WaitGroup, fdwTableName string, segmentSize uint64) (chan [2]uint64, chan struct{}, chan error, error)
	TransformToCSV(csvWriter csv.Writer, wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan struct{}, chan error)
	io.Closer
}

// Service struct underpinning the Migrator interface
type Service struct {
	reader       *Reader
	writer       *sql.Writer
	oldDB, newDB *sqlx.DB

	wg                 *sync.WaitGroup
	closeChan          chan struct{}
	numWorkersPerTable int
}

// NewMigrator returns a new Migrator from the given Config
func NewMigrator(ctx context.Context, conf *Config) (Migrator, error) {
	readDB, err := NewDB(ctx, conf.ReadDB)
	if err != nil {
		return nil, err
	}
	writeDB, err := NewDB(ctx, conf.WriteDB)
	if err != nil {
		return nil, err
	}
	numWorkers := defaultNumWorkersPerTable
	if conf.WorkersPerTable != 0 {
		numWorkers = conf.WorkersPerTable
	}
	return &Service{
		reader:             NewReader(readDB),
		writer:             sql.NewWriter(writeDB),
		oldDB:              readDB,
		newDB:              writeDB,
		closeChan:          make(chan struct{}),
		numWorkersPerTable: numWorkers,
	}, nil
}

// TransformToCSV satisfies Migrator
// TransformToCSV spins up a goroutine to process the block ranges provided through the blockRanges work chan for the specified tables
// TransformToCSV returns a channel for emitting read gaps and failed write ranges, a channel for signaling completion
// of the process, a quitChan for closing the single process, and a channel for writing out errors
func (s *Service) TransformToCSV(csvWriter csv.Writer, wg *sync.WaitGroup, tableName TableName,
	blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan struct{}, chan error) {
	quitChan := make(chan struct{})
	doneChan := make(chan struct{})
	transformer := NewTableTransformer(tableName)
	readPgStr := tableReaderStrMappings[tableName]
	writeCSVStr := csvWriterStrMappings[tableName]
	readGapChan := make(chan [2]uint64)
	writeGapChan := make(chan [2]uint64)
	errChan := make(chan error)
	innerWg := new(sync.WaitGroup)

	for workerNum := 1; workerNum <= s.numWorkersPerTable; workerNum++ {
		innerWg.Add(1)
		go func(workerNum int, tableName TableName) {
			logrus.Infof("starting migration worker %d for table %s", workerNum, tableName)
			defer innerWg.Done()
			for {
				select {
				case rng := <-blockRanges:
					logrus.Debugf("table %s worker %d received block range (%d, %d)", tableName, workerNum, rng[0], rng[1])
					oldModels, err := NewTableReadModels(tableName)
					if err != nil {
						errChan <- fmt.Errorf("table %s worker %d unable to create tabel models for range (%d, %d): %v", tableName, workerNum, rng[0], rng[1], err)
						readGapChan <- rng
						continue
					}
					if err := s.reader.Read(rng, readPgStr, oldModels); err != nil {
						errChan <- fmt.Errorf("table %s worker %d read error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						readGapChan <- rng
						continue
					}
					numReadRecords := reflect.Indirect(reflect.ValueOf(oldModels)).Len()
					if numReadRecords == 0 {
						if tableName == EthHeaders || tableName == EthState || tableName == EthAccounts {
							// all other tables can, at least in theory, be empty within a range
							// e.g. a block that has no txs or uncles will only
							// have a header and an updated state account for the miner's reward
							readGapChan <- rng
						} else {
							logrus.Infof("table %s worker %d finished range (%d, %d)- no read records found in range", tableName, workerNum, rng[0], rng[1])
						}
						continue
					}
					logrus.Debugf("table %s worker %d block range (%d, %d) read models count: %d", tableName, workerNum, rng[0], rng[1], numReadRecords)
					newModels, gaps, err := transformer.Transform(oldModels, rng)
					if err != nil {
						errChan <- fmt.Errorf("table %s worker %d transform error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						writeGapChan <- rng
						continue
					}
					logrus.Debugf("table %s worker %d block range (%d, %d) write models count: %d", tableName, workerNum, rng[0], rng[1], reflect.ValueOf(newModels).Len())
					if err := csvWriter.Write(writeCSVStr, newModels); err != nil {
						errChan <- fmt.Errorf("table %s worker %d write error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						writeGapChan <- rng
						continue
					}
					for _, gap := range gaps {
						readGapChan <- gap
					}
					logrus.Infof("table %s worker %d finished range (%d, %d)- %d records processed", tableName, workerNum, rng[0], rng[1], numReadRecords)
				case <-s.closeChan:
					logrus.Infof("quitting migration worker %d for table %s", workerNum, tableName)
					return
				default:
					select {
					case <-quitChan:
						logrus.Infof("quitting migration worker %d for table %s", workerNum, tableName)
						return
					default:
					}
				}
			}
		}(workerNum, tableName)
	}

	wg.Add(1)
	go func() {
		innerWg.Wait()
		wg.Done()
		close(doneChan)
	}()

	return readGapChan, writeGapChan, doneChan, quitChan, errChan
}

// Migrate satisfies Migrator
// Migrate spins up a goroutine to process the block ranges provided through the blockRanges work chan for the specified tables
// Migrate returns a channel for emitting read gaps and failed write ranges, a channel for signaling
// completion of the process, a quitChan for closing the single process, and a channel for writing out errors
func (s *Service) Migrate(wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64,
	chan [2]uint64, chan struct{}, chan struct{}, chan error) {
	quitChan := make(chan struct{})
	doneChan := make(chan struct{})
	transformer := NewTableTransformer(tableName)
	readPgStr := tableReaderStrMappings[tableName]
	writePgStr := tableWriterStrMappings[tableName]
	readGapChan := make(chan [2]uint64)
	writeGapChan := make(chan [2]uint64)
	errChan := make(chan error)
	innerWg := new(sync.WaitGroup)

	for workerNum := 1; workerNum <= s.numWorkersPerTable; workerNum++ {
		innerWg.Add(1)
		go func(workerNum int, tableName TableName) {
			logrus.Infof("starting migration worker %d for table %s", workerNum, tableName)
			defer innerWg.Done()
			for {
				select {
				case rng := <-blockRanges:
					logrus.Debugf("table %s worker %d received block range (%d, %d)", tableName, workerNum, rng[0], rng[1])
					oldModels, err := NewTableReadModels(tableName)
					if err != nil {
						errChan <- fmt.Errorf("table %s worker %d unable to create tabel models for range (%d, %d): %v", tableName, workerNum, rng[0], rng[1], err)
						readGapChan <- rng
						continue
					}
					if err := s.reader.Read(rng, readPgStr, oldModels); err != nil {
						errChan <- fmt.Errorf("table %s worker %d read error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						readGapChan <- rng
						continue
					}
					numReadRecords := reflect.Indirect(reflect.ValueOf(oldModels)).Len()
					if numReadRecords == 0 {
						if tableName == EthHeaders || tableName == EthState || tableName == EthAccounts {
							// all other tables can, at least in theory, be empty within a range
							// e.g. a block that has no txs or uncles will only
							// have a header and an updated state account for the miner's reward
							readGapChan <- rng
						} else {
							logrus.Infof("table %s worker %d finished range (%d, %d)- no read records found in range", tableName, workerNum, rng[0], rng[1])
						}
						continue
					}
					logrus.Debugf("table %s worker %d block range (%d, %d) read models count: %d", tableName, workerNum, rng[0], rng[1], numReadRecords)
					newModels, gaps, err := transformer.Transform(oldModels, rng)
					if err != nil {
						errChan <- fmt.Errorf("table %s worker %d transform error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						writeGapChan <- rng
						continue
					}
					logrus.Debugf("table %s worker %d block range (%d, %d) write models count: %d", tableName, workerNum, rng[0], rng[1], reflect.ValueOf(newModels).Len())
					if err := s.writer.Write(writePgStr, newModels); err != nil {
						errChan <- fmt.Errorf("table %s worker %d write error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						writeGapChan <- rng
						continue
					}
					for _, gap := range gaps {
						readGapChan <- gap
					}
					logrus.Infof("table %s worker %d finished range (%d, %d)- %d records processed", tableName, workerNum, rng[0], rng[1], numReadRecords)
				case <-s.closeChan:
					logrus.Infof("quitting migration worker %d for table %s", workerNum, tableName)
					return
				default:
					select {
					case <-quitChan:
						logrus.Infof("quitting migration worker %d for table %s", workerNum, tableName)
						return
					default:
					}
				}
			}
		}(workerNum, tableName)
	}

	wg.Add(1)
	go func() {
		innerWg.Wait()
		wg.Done()
		close(doneChan)
	}()

	return readGapChan, writeGapChan, doneChan, quitChan, errChan
}

// Transfer for transferring public.blocks to a new DB page-by-page
// Transfer assumes the targeted postgres_fdw is already in the db
// returns a chan for logging failed transfer page ranges, a chan for the errors that caused them,
// a chan for signalling success, and any error during initialization
func (s *Service) Transfer(wg *sync.WaitGroup, fdwTableName string, segmentSize uint64) (chan [2]uint64,
	chan struct{}, chan error, error) {
	db := s.newDB
	if fdwTableName == "" {
		fdwTableName = public_blocks.DefaultV2FDWTableName
	}

	maxPage, err := public_blocks.GetMaxPage(db, fdwTableName)
	if err != nil {
		return nil, nil, nil, err
	}

	doneChan := make(chan struct{})
	segments := public_blocks.GetPageSegments(maxPage, segmentSize)
	errChan := make(chan error)
	transferFailChan := make(chan [2]uint64)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(doneChan)
		for _, segment := range segments {
			select {
			case <-s.closeChan:
				logrus.Infof("quitting transfer process for table %s", fdwTableName)
				return
			default:
			}
			logrus.Infof("transfer %s page range (%d, %d) from old DB to new DB", fdwTableName, segment[0], segment[1])
			if err := public_blocks.TransferPages(db, fdwTableName, segment[0], segment[1]); err != nil {
				errChan <- fmt.Errorf("failed to transfer %s page range (%d, %d): %v", fdwTableName, segment[0], segment[1], err)
				transferFailChan <- segment
			}
		}
	}()
	return transferFailChan, doneChan, errChan, nil
}

// Close satisfied io.Closer
// Close shuts down the Migrator, it quits all Migrate goroutines that are currently running
// whereas closing the chan returned by Migrate only closes the goroutines spun up by that method call
func (s *Service) Close() error {
	close(s.closeChan)
	if err := s.reader.Close(); err != nil {
		return err
	}
	return s.writer.Close()
}
