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

package eth_access_lists

const (
	PgReadAccessListElementsStr = `SELECT eth.transaction_cids.tx_hash, eth.access_list_elements.*
									FROM eth.access_list_elements
									INNER JOIN eth.transaction_cids ON (access_list_elements.tx_id = transaction_cids.id)
									INNER JOIN eth.header_cids ON (transaction_cids.header_id = header_cids.id)
									WHERE block_number = $1`

	PgWriteAccessListElementsStr = `INSERT INTO eth.access_list_elements (tx_id, index, address, storage_keys)
									VALUES (unnest($1::VARCHAR(66)[]), unnest($2::INTEGER[]), unnest($3::VARCHAR(66)[]),
									unnest($4::VARCHAR(66)[][]))`
)
