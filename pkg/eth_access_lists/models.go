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

package public_blocks

import "github.com/lib/pq"

// AccessListElementModelV2 is the db model for eth.access_list_entry for v2 DB
type AccessListElementModelV2 struct {
	ID          int64          `db:"id"`
	Index       int64          `db:"index"`
	TxID        int64          `db:"tx_id"`
	Address     string         `db:"address"`
	StorageKeys pq.StringArray `db:"storage_keys"`
}

// AccessListElementModelV3 is the db model for eth.access_list_entry for v3 DB
type AccessListElementModelV3 struct {
	Index       int64          `db:"index"`
	TxID        string         `db:"tx_id"`
	Address     string         `db:"address"`
	StorageKeys pq.StringArray `db:"storage_keys"`
}
