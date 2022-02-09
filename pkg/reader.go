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
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/migration-tools/pkg/sql"
)

// Reader struct for reading v2 DB eth.log_cids models
type Reader struct {
	db *sqlx.DB
}

// NewReader satisfies interfaces.ReaderConstructor for eth.log_cids
func NewReader(db *sqlx.DB) *Reader {
	return &Reader{db: db}
}

// Read satisfies interfaces.Reader for eth.log_cids
// Read is safe for concurrent use, as the only shared state is the concurrent safe *sqlx.DB
func (r *Reader) Read(blockRange [2]uint64, pgStr sql.ReadPgStr, models interface{}) error {
	return r.db.Select(models, string(pgStr), blockRange[0], blockRange[1])
}

// Close satisfies io.Closer
func (r *Reader) Close() error {
	return r.db.Close()
}
