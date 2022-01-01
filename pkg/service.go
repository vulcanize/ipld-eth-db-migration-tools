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
)

type Migrator interface {
	Migrate(wg *sync.WaitGroup, tables []TableName, blockHeights <-chan uint64) (chan error, error)
	io.Closer
}

type Service struct {
	Readers     map[TableName]Reader
	Writers     map[TableName]Writer
	Transformer map[TableName]Transformer

	wg       *sync.WaitGroup
	errChan  chan error
	quitChan chan struct{}
}

func NewMigrator(conf *Config) Migrator {
	return &Service{
		Readers:     make(map[TableName]Reader, numTables),
		Writers:     make(map[TableName]Writer, numTables),
		Transformer: make(map[TableName]Transformer, numTables),
		quitChan:    make(chan struct{}),
		errChan:     make(chan error),
	}
}

func (s *Service) Migrate(wg *sync.WaitGroup, tables []TableName, blockHeights <-chan uint64) (chan error, error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

		}
	}()
	return s.errChan, nil
}

func (s *Service) Close() error {
	close(s.quitChan)
	return nil
}
