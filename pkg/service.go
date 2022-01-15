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

	"github.com/sirupsen/logrus"
)

const defaultNumWorkersPerTable = 1

// Migrator interface for migrating from v2 DB to v3 DB
type Migrator interface {
	Migrate(wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan struct{}, chan error)
	io.Closer
}

// Service struct underpinning the Migrator interface
type Service struct {
	reader *Reader
	writer *Writer

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
		writer:             NewWriter(writeDB),
		closeChan:          make(chan struct{}),
		numWorkersPerTable: numWorkers,
	}, nil
}

// Migrate satisfies Migrator
// Migrate spins up a goroutine to process the block ranges provided through the blockRanges work chan for the specified tables
// Migrate returns a channel for emitting read gaps and failed write ranges, a quitChan for closing the process once we
// are sending it ranges to process, a channel for signalling successful shutdown of the process, and a channel for writing out errors
func (s *Service) Migrate(wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan struct{}, chan error) {
	quitChan := make(chan struct{})
	doneChan := make(chan struct{})
	transformer := NewTableTransformer(tableName)
	readPgStr := tableReaderStrMappings[tableName]
	writePgStr := tableWriterStrMappings[tableName]
	readGapChan := make(chan [2]uint64)
	writeGapChan := make(chan [2]uint64)
	errChan := make(chan error)

	for workerNum := 1; workerNum <= s.numWorkersPerTable; workerNum++ {
		wg.Add(1)
		go func(workerNum int, tableName TableName) {
			logrus.Infof("starting migration worker %d for table %s", workerNum, tableName)
			defer wg.Done()
			defer close(doneChan)
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
					if reflect.Indirect(reflect.ValueOf(oldModels)).Len() == 0 {
						if tableName == EthHeaders || tableName == EthState || tableName == EthAccounts {
							// all other tables can, at least in theory, be empty within a range
							// e.g. a block that has no txs or uncles will only
							// have a header and an updated state account for the miner's reward
							readGapChan <- rng
						} else {
							logrus.Infof("table %s worker %d finished range (%d, %d)", tableName, workerNum, rng[0], rng[1])
						}
						continue
					}
					newModels, gaps, err := transformer.Transform(oldModels, rng)
					if err != nil {
						errChan <- fmt.Errorf("table %s worker %d transform error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						writeGapChan <- rng
						continue
					}
					if err := s.writer.Write(writePgStr, newModels); err != nil {
						errChan <- fmt.Errorf("table %s worker %d write error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						writeGapChan <- rng
						continue
					}
					for _, gap := range gaps {
						readGapChan <- gap
					}
					logrus.Infof("table %s worker %d finished range (%d, %d)", tableName, workerNum, rng[0], rng[1])
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

	return readGapChan, writeGapChan, doneChan, quitChan, errChan
}

// Close satisfied io.Closer
// Close shuts down the Migrator, it quits all Migrate goroutines that are currently running
// whereas closing the chan returned by Migrate only closes the goroutines spun up by that method call
func (s *Service) Close() error {
	close(s.closeChan)
	return nil
}
