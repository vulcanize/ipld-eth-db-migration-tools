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

package eth_accounts

import (
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/migration-tools/pkg/interfaces"
)

// Reader struct for reading v2 DB eth.state_accounts models
type Reader struct {
	db *sqlx.DB
}

// NewReader satisfies interfaces.ReaderConstructor for eth.state_accounts
func NewReader(db *sqlx.DB) interfaces.Reader {
	return &Reader{db: db}
}

// Read satisfies interfaces.Reader for eth.state_accounts
func (r *Reader) Read(blockHeights []uint64) ([][]interface{}, []uint64, error) {

}
