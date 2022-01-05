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

package interfaces

import "github.com/jmoiron/sqlx"

// Reader interface for reading v2 DB models for a specific table in a database
type Reader interface {
	Read(blockHeights []uint64) ([][]interface{}, []uint64, error)
}

// ReaderConstructor func sig for constructing a Reader for a specific table
type ReaderConstructor func(db *sqlx.DB) Reader

// Writer interface for writing v3 DB models for a specific table in a database
type Writer interface {
	Write(models [][]interface{}) error
}

// WriterConstructor func sig for constructing a Writer for a specific table
type WriterConstructor func(db *sqlx.DB) Writer

// Transformer interface for transforming v2 DB models into v3 DB models for a specific table
type Transformer interface {
	Transform(models [][]interface{}) ([][]interface{}, error)
}

// TransformerConstructor func sig for constructing a Transformer for a specific table
type TransformerConstructor func() Transformer
