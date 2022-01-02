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
	"io"
	"sync"

	"github.com/vulcanize/migration-tools/pkg/interfaces"
)

// Migrator interface for migrating from v2 DB to v3 DB
type Migrator interface {
	Migrate(wg *sync.WaitGroup, tables []TableName, blockRanges <-chan []uint64) (chan []uint64, chan struct{}, chan error)
	io.Closer
}

// Service struct underpinning the Migrator interface
type Service struct {
	Readers     map[TableName]interfaces.Reader
	Writers     map[TableName]interfaces.Writer
	Transformer map[TableName]interfaces.Transformer

	wg       *sync.WaitGroup
	gapsChan chan []uint64
	errChan  chan error
	stopChan chan struct{}
}

// NewMigrator returns a new Migrator from the given Config
func NewMigrator(conf *Config) Migrator {
	return &Service{
		Readers:     make(map[TableName]interfaces.Reader, numTables),
		Writers:     make(map[TableName]interfaces.Writer, numTables),
		Transformer: make(map[TableName]interfaces.Transformer, numTables),
		gapsChan:    make(chan []uint64),
		stopChan:    make(chan struct{}),
		errChan:     make(chan error),
	}
}

// Migrate satisfies Migrator
// Migrate spins up a goroutine to process the block ranges provided through the blockRanges work chan for the specified tables
// Migrate returns a channel for emitting gaps and failed ranges, a quitChannel for its goroutine, and a channel for writing out errors
func (s *Service) Migrate(wg *sync.WaitGroup, tables []TableName, blockRanges <-chan []uint64) (chan []uint64, chan struct{}, chan error) {
	wg.Add(1)
	quitChan := make(chan struct{})
	go func() {
		defer wg.Done()
		for {
			select {
			case rng := <-blockRanges:
				for _, tableName := range tables {
					oldModels, gaps, err := s.Readers[tableName].Read(rng)
					if err != nil {
						s.errChan <- err
						s.gapsChan <- rng
						break
					}
					newModels, err := s.Transformer[tableName].Transform(oldModels)
					if err != nil {
						s.errChan <- err
						s.gapsChan <- rng
					}
					if err := s.Writers[tableName].Write(newModels); err != nil {
						s.errChan <- err
						s.gapsChan <- rng
					}
					s.gapsChan <- gaps
				}
			case <-s.stopChan:
				close(quitChan)
			case <-quitChan:
				return
			}
		}
	}()
	return s.gapsChan, quitChan, s.errChan
}

// Close satisfied io.Closer
// Close shuts down the Migrator, it quits any/all Migrate goroutines that are currently running
func (s *Service) Close() error {
	close(s.stopChan)
	return nil
}
