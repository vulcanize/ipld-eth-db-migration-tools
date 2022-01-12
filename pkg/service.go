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
	"sync"

	"github.com/sirupsen/logrus"
)

const defaultNumWorkersPerTable = 1

// Migrator interface for migrating from v2 DB to v3 DB
type Migrator interface {
	Migrate(wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan error)
	io.Closer
}

// Service struct underpinning the Migrator interface
type Service struct {
	reader *Reader
	writer *Writer

	wg                 *sync.WaitGroup
	readGapsChan       chan [2]uint64
	writeFailuresChan  chan [2]uint64
	errChan            chan error
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
		readGapsChan:       make(chan [2]uint64),
		writeFailuresChan:  make(chan [2]uint64),
		closeChan:          make(chan struct{}),
		errChan:            make(chan error),
		numWorkersPerTable: numWorkers,
	}, nil
}

// Migrate satisfies Migrator
// Migrate spins up a goroutine to process the block ranges provided through the blockRanges work chan for the specified tables
// Migrate returns a channel for emitting read gaps and failed write ranges, a doneChan for closing the process once we
// are sending it ranges to process, and a channel for writing out errors
func (s *Service) Migrate(wg *sync.WaitGroup, tableName TableName, blockRanges <-chan [2]uint64) (chan [2]uint64, chan [2]uint64, chan struct{}, chan error) {
	doneChan := make(chan struct{})
	transformer := NewTableTransformer(tableName)
	readPgStr := tableReaderStrMappings[tableName]
	writePgStr := tableWriterStrMappings[tableName]

	for workerNum := 1; workerNum <= s.numWorkersPerTable; workerNum++ {
		wg.Add(1)
		go func(workerNum int, tableName TableName) {
			logrus.Infof("starting migration worker %d for table %s", workerNum, tableName)
			defer wg.Done()
			for {
				select {
				case rng := <-blockRanges: // TODO: we don't want to trigger the public.nodes iterator more than once
					oldModels, err := NewTableReadModels(tableName)
					if err := s.reader.Read(rng, readPgStr, oldModels); err != nil {
						s.errChan <- fmt.Errorf("table %s worker %d read error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						if tableName == EthHeaders || tableName == EthState || tableName == EthAccounts {
							// all other tables can, at least in theory, be empty within a range
							// e.g. a block that has no txs or uncles will only
							// have a header and an updated state account for the miner's reward
							s.readGapsChan <- rng
						}
						continue
					}
					newModels, gaps, err := transformer.Transform(oldModels, rng)
					if err != nil {
						s.errChan <- fmt.Errorf("table %s worker %d transform error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						s.writeFailuresChan <- rng
						continue
					}
					if err := s.writer.Write(writePgStr, newModels); err != nil {
						s.errChan <- fmt.Errorf("table %s worker %d write error (%v) in range (%d, %d)", tableName, workerNum, err, rng[0], rng[1])
						s.writeFailuresChan <- rng
						continue
					}
					for _, gap := range gaps { // TODO: this will send duplicate ranges if running header, state, or account processing simultaneously
						s.readGapsChan <- gap
					}
				case <-s.closeChan:
					return
				default:
					select {
					case <-doneChan:
						logrus.Infof("quitting migration worker %d for table %s", workerNum, tableName)
						return
					default:
					}
				}
			}
		}(workerNum, tableName)
	}

	return s.readGapsChan, s.writeFailuresChan, doneChan, s.errChan
}

// Close satisfied io.Closer
// Close shuts down the Migrator, it quits all Migrate goroutines that are currently running
// whereas closing the chan returned by Migrate only closes the goroutines spun up by that method call
func (s *Service) Close() error {
	close(s.closeChan)
	return nil
}
